package main

import (
	"fmt"
	"strconv"
	"github.com/namsral/flag"
	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/bitbucket"
)

func main() {
	var (
		port int
		uuid string
	)

	flag.IntVar(&port, "port", 8080, "TCP listener port")
	flag.StringVar(&uuid, "uuid", "", "Webhook UUID")
	flag.Parse()

	hook := bitbucket.New(&bitbucket.Config{UUID: uuid})
	hook.RegisterEvents(HandleMultiple, bitbucket.RepoPushEvent) // Add as many as you want

	err := webhooks.Run(hook, ":"+strconv.Itoa(port), "/bitbucket")
	if err != nil {
		fmt.Println(err)
	}
}

// HandleMultiple handles multiple GitHub events
func HandleMultiple(payload interface{}, header webhooks.Header) {

	fmt.Println("Handling Payload..")

	switch payload.(type) {

	case bitbucket.RepoPushPayload:
		release := payload.(bitbucket.RepoPushPayload)
		// Do whatever you want from here...
		// TODO: post payload to Jenkins
		fmt.Printf("%+v", release)
	}
}
