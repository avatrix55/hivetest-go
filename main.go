package main
/*

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
	clientID := "testCred"
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
*/

import (
  "fmt"
  "time"

  mqtt "github.com/eclipse/paho.mqtt.golang"
)

// this callback triggers when a message is received, it then prints the message (in the payload) and topic
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
  fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

// upon connection to the client, this is called
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
  fmt.Println("Connected")
}

// this is called when the connection to the client is lost, it prints "Connection lost" and the corresponding error
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
  fmt.Printf("Connection lost: %v", err)
}

func main() {
  var broker = "tls://3a3bd399e18940d5a6ce7b2b9d8476d0.s1.eu.hivemq.cloud" // find the host name in the Overview of your cluster (see readme)
  var port = 8883            // find the port right under the host name, standard is 8883
  opts := mqtt.NewClientOptions()
  opts.AddBroker(fmt.Sprintf("tls://%s:%d", broker, port))
  opts.SetClientID("testCred") // set a name as you desire
  opts.SetUsername("avatrix")    // these are the credentials that you declare for your cluster
  opts.SetPassword("@Testcred1")
  // (optionally) configure callback handlers that get called on certain events
  opts.SetDefaultPublishHandler(messagePubHandler)
  opts.OnConnect = connectHandler
  opts.OnConnectionLost = connectLostHandler
  // create the client using the options above
  client := mqtt.NewClient(opts)
  // throw an error if the connection isn't successfull
  if token := client.Connect(); token.Wait() && token.Error() != nil {
    panic(token.Error())
  }

  subscribe(client)
  publish(client)

  client.Disconnect(250)
}

func subscribe(client mqtt.Client) {
  // subscribe to the same topic, that was published to, to receive the messages
  topic := "topic/test"
  token := client.Subscribe(topic, 1, nil)
  token.Wait()
  // Check for errors during subscribe (More on error reporting https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang#readme-error-handling)
  if token.Error() != nil {
    fmt.Printf("Failed to subscribe to topic")
    panic(token.Error())
  }
  fmt.Printf("Subscribed to topic: %s", topic)
}

func publish(client mqtt.Client) {
  // publish the message "Message" to the topic "topic/test" 10 times in a for loop
  num := 10
  for i := 0; i < num; i++ {
    text := fmt.Sprintf("Message %d", i)
    token := client.Publish("topic/test", 0, false, text)
    token.Wait()
    // Check for errors during publishing (More on error reporting https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang#readme-error-handling)
    if token.Error() != nil {
      fmt.Printf("Failed to publish to topic")
      panic(token.Error())
    }
    time.Sleep(time.Second)
  }
}