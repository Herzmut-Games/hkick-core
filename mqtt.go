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
	client.Subscribe("score/increase", 0, func(client mqtt.Client, msg mqtt.Message) {
		increaseScore(string(msg.Payload()))
	})
	client.Subscribe("score/undo", 0, func(client mqtt.Client, msg mqtt.Message) {
		undoScore()
	})
	client.Subscribe("score/reset", 0, func(client mqtt.Client, msg mqtt.Message) {
		resetScore()
	})
	client.Subscribe("game/start", 0, func(client mqtt.Client, msg mqtt.Message) {
		startGame()
	})
	client.Subscribe("game/stop", 0, func(client mqtt.Client, msg mqtt.Message) {
		stopGame()
	})
}

func publish(topic string, message string, retain bool) {
	client.Publish(topic, 0, retain, message)
}
