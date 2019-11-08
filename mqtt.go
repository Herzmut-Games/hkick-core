package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	client    mqtt.Client
	goalTopic = "goals"
)

func connect(clientID string, uri *url.URL) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))

	client = mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}

	return client
}

func subscribe(uri *url.URL) {
	client.Subscribe("goals", 0, func(client mqtt.Client, msg mqtt.Message) {
		increaseScore(string(msg.Payload()))
	})
	client.Subscribe("score/decrease", 0, func(client mqtt.Client, msg mqtt.Message) {
		decreaseScore(string(msg.Payload()))
	})
	client.Subscribe("score/increase", 0, func(client mqtt.Client, msg mqtt.Message) {
		increaseScore(string(msg.Payload()))
	})
	client.Subscribe("score/reset", 0, func(client mqtt.Client, msg mqtt.Message) {
		resetScore()
	})
}

func publish(topic string, message string, retain bool) {
	client.Publish(topic, 0, retain, message)
}
