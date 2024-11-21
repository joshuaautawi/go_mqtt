package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// MQTT broker address
	broker := "tcp://localhost:1883" // Replace with your broker address
	topic := "test/topic"            // The topic you want to subscribe to

	// MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("go-mqtt-client")
	opts.SetKeepAlive(60 * time.Second)
	opts.SetCleanSession(true)

	// Create MQTT client
	client := mqtt.NewClient(opts)

	// Connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	defer client.Disconnect(250)

	// Subscribe to the topic
	if token := client.Subscribe(topic, 0, messageHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Gracefully shutting down...")
}

// Callback function to handle incoming MQTT messages
func messageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), string(msg.Payload()))
}
