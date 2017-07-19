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
)

type Specification struct {
	Debug bool    `default:"false"`
	Port  int     `default:"8080"`
	UUID  string
}

func main() {
	var s Specification
	err := envconfig.Process("webhook", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	hook := bitbucket.New(&bitbucket.Config{UUID: s.UUID})
	hook.RegisterEvents(HandleMultiple, bitbucket.RepoPushEvent) // Add as many as you want

	r := mux.NewRouter()
	r.Path("/bitbucket").
		Handler(webhooks.Handler(hook)).
		Headers("X-Hook-UUID", s.UUID)

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + strconv.Itoa(s.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//log.Printf("Listening on addr: %s path: %s", addr, s.Path)
	log.Printf("Listening on addr: %s path: %s", ":" + strconv.Itoa(s.Port), "/bitbucket")
	log.Fatal(srv.ListenAndServe())
}

// HandleMultiple handles multiple GitHub events
func HandleMultiple(payload interface{}, header webhooks.Header) {

	log.Println("Handling Payload..")

	switch payload.(type) {

	case bitbucket.RepoPushPayload:
		release := payload.(bitbucket.RepoPushPayload)
		// Do whatever you want from here...
		// TODO: post payload to Jenkins
		log.Printf("%+v", release)
	}
}
