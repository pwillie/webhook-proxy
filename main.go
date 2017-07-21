package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/bitbucket"
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

	r.HandleFunc("/status", healthcheck_handler)

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

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println()
		log.Printf("Received signal: %v", sig)
		log.Println("Shutting down...")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		srv.Shutdown(ctx)
	}()

	log.Printf("Listening on addr: %v path: %v", ":"+strconv.Itoa(c.Port), c.Path)
	log.Fatal(srv.ListenAndServe())
}

// status handler for kubernetes health checks
func healthcheck_handler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
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
