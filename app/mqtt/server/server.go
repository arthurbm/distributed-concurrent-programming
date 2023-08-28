package main

import (
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	BrokerAddress     = "tcp://mosquitto:1883"
	RequestTopicBase  = "fibonacci/request/"
	ResponseTopicBase = "fibonacci/response/"
	QoS               = 1
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

	// Listen to all client requests
	client.Subscribe(RequestTopicBase+"#", QoS, func(client mqtt.Client, msg mqtt.Message) {
		go func(m mqtt.Message) {
			numberStr := string(m.Payload())
			number, err := strconv.Atoi(numberStr)
			if err != nil {
				fmt.Println("Error converting payload to integer:", err)
				return
			}

			// Calculate the Fibonacci number
			result := fibonacci(number)

			// Extract client ID from the topic
			clientID := strings.Split(m.Topic(), "/")[2]

			// Construct the response topic based on the client ID and the number requested
			respTopic := ResponseTopicBase + clientID + "/" + numberStr

			// Publish the response
			client.Publish(respTopic, QoS, false, fmt.Sprintf("%d", result))
		}(msg)
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
