package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	DataFilePath      = "/app/data/"
	BrokerAddress     = "tcp://mosquitto:1883"
	RequestTopicBase  = "fibonacci/request/"
	ResponseTopicBase = "fibonacci/response/"
	NumberRequests    = 40
	QoS               = 1
)

func main() {
	startTime := time.Now() // record the start time
	uniqueID := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Modify the client options to use the unique ID
	opts := createClientOptions("client-"+uniqueID, BrokerAddress)
	client := connect("client-"+uniqueID, opts)

	writer, file, err := openCSVFile(DataFilePath, uniqueID)
	if err != nil {
		fmt.Println("Error opening the file:", err.Error())
		return
	}
	defer file.Close()
	defer writer.Flush()

	err = writeCSVHeader(writer)
	if err != nil {
		fmt.Println("Error writing CSV header:", err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(NumberRequests) // Set the counter

	responseTopic := ResponseTopicBase + uniqueID + "/+"
	client.Subscribe(responseTopic, QoS, func(client mqtt.Client, msg mqtt.Message) {
		numberStr := strings.Split(msg.Topic(), "/")[3]
		number, _ := strconv.Atoi(numberStr)
		response, _ := strconv.Atoi(string(msg.Payload()))
		err := writeToCSV(writer, number, response)
		if err != nil {
			fmt.Println("Error writing to CSV:", err.Error())
		}
		wg.Done() // Decrement the counter
	})

	for i := 0; i < NumberRequests; i++ {
		text := fmt.Sprintf("%d", i)
		requestTopic := RequestTopicBase + uniqueID
		client.Publish(requestTopic, QoS, false, text)
	}

	// Instead of sleeping, wait until all responses have been processed
	wg.Wait()

	var totalTime = time.Now().Sub(startTime).Nanoseconds()
	err = addTotalTimeCSV(writer, totalTime)
	if err != nil {
		fmt.Println("Error writing to CSV:", err.Error())
	}

	// disconnect from the broker
	defer client.Disconnect(250)
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

func openCSVFile(basePath string, uniqueID string) (*csv.Writer, *os.File, error) {
	timestamp := time.Now().Format("20060102150405")
	filePath := fmt.Sprintf("%s_%s_%s.csv", basePath, uniqueID, timestamp)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	writer := csv.NewWriter(file)
	return writer, file, nil
}

func addTotalTimeCSV(writer *csv.Writer, totalTime int64) error {
	return writer.Write([]string{"0", "0", fmt.Sprintf("%d", totalTime)})
}

func writeToCSV(writer *csv.Writer, fibonacciIn, fibonacciOut int) error {
	return writer.Write([]string{fmt.Sprintf("%d", fibonacciIn), fmt.Sprintf("%d", fibonacciOut)})
}

func writeCSVHeader(writer *csv.Writer) error {
	return writer.Write([]string{"Input", "Output", "timeTaken"})
}
