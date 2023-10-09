package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// HiveMQ Public Broker
	broker := "mqtt://3a3bd399e18940d5a6ce7b2b9d8476d0.s1.eu.hivemq.cloud:8883/mqtt"
	clientID := "testDev"
	topic := "Topictest"

	opts := createClientOptions(clientID, broker)
	client := mqtt.NewClient(opts)

	// Connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to broker: %s", token.Error())
		os.Exit(1)
	}

	// Subscribe to a topic
	token := client.Subscribe(topic, 1, messageHandler)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)

	// Publish a message to the topic
	token = client.Publish(topic, 1, false, "Hello from Go!")
	token.Wait()

	// Keep the application running to receive messages
	select {}
}

func messageHandler(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message on topic %s: %s\n", message.Topic(), message.Payload())
}

func createClientOptions(clientID, uri string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(uri)
	opts.SetClientID(clientID)
	opts.SetConnectionLostHandler(connectionLostHandler)
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	return opts
}

func connectionLostHandler(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %s\n", err)
}