package main

import (
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	BrokerAddress     = "tcp://mosquitto:1883"
	RequestTopic      = "fibonacci/request"
	ResponseTopicBase = "fibonacci/response/"
)

// Fibonacci function
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	opts := createClientOptions("server", BrokerAddress)
	client := connect("server", opts)

	client.Subscribe(RequestTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		number, _ := strconv.Atoi(string(msg.Payload()))
		result := fibonacci(number)

		respTopic := ResponseTopicBase + string(msg.Payload())
		client.Publish(respTopic, 0, false, fmt.Sprintf("%d", result))
	})

	fmt.Println("Server is running...")

	select {} // This keeps the program running indefinitely to process incoming messages.
}

func createClientOptions(clientID, uri string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(uri)
	opts.SetClientID(clientID)
	return opts
}

func connect(clientID string, opts *mqtt.ClientOptions) mqtt.Client {
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
