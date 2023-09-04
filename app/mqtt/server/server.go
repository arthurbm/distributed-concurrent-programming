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
			splitData := strings.Split(string(m.Payload()), ":")
			numberStr := splitData[0]
			sentTime := splitData[1]

			number, err := strconv.Atoi(numberStr)
			if err != nil {
				fmt.Println("Error converting payload to integer:", err)
				return
			}

			result := fibonacci(number)
			clientID := strings.Split(m.Topic(), "/")[2]
			respTopic := ResponseTopicBase + clientID + "/" + numberStr

			responsePayload := fmt.Sprintf("%d:%s", result, sentTime)
			client.Publish(respTopic, QoS, false, responsePayload)
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
