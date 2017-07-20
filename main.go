package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/bitbucket"
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Debug      bool     `default:"false"`
	Port       int      `default:"8080"`
	Path       string   `default:"/bitbucket"`
	Uuid       []string `required:"true"`
	JenkinsUrl string   `required:"true"`
}

var c Configuration

func main() {
	err := envconfig.Process("webhook", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := mux.NewRouter()

	for _, uuid := range c.Uuid {
		hook := bitbucket.New(&bitbucket.Config{UUID: uuid})
		hook.RegisterEvents(HandleMultiple, bitbucket.RepoPushEvent) // Add as many as you want

		r.Path(c.Path).
			Handler(webhooks.Handler(hook)).
			Headers("X-Hook-Uuid", uuid)
	}

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + strconv.Itoa(c.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on addr: %v path: %v", ":"+strconv.Itoa(c.Port), c.Path)
	log.Fatal(srv.ListenAndServe())
}

// HandleMultiple handles multiple GitHub events
func HandleMultiple(payload interface{}, header webhooks.Header) {

	log.Println("Handling Payload..")

	switch payload.(type) {

	case bitbucket.RepoPushPayload:
		release := payload.(bitbucket.RepoPushPayload)
		// post payload to Jenkins
		log.Printf("%+v", release)
		jsonValue, _ := json.Marshal(release)
		resp, err := http.Post(c.JenkinsUrl, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Printf("ERROR: Jenkins: %s", err.Error())
		}

		defer resp.Body.Close()

		log.Println("Jenkins response Status:", resp.Status)
		log.Println("Jenkins response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Jenkins response Body:", string(body))
	}
}
