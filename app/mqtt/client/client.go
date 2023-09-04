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
		splitData := strings.Split(string(msg.Payload()), ":")
		number, _ := strconv.Atoi(splitData[0])
		sentTime, _ := strconv.ParseInt(splitData[1], 10, 64)
		responseTime := time.Now().UnixNano()
		timeTaken := responseTime - sentTime

		err := writeToCSV(writer, number, 0, timeTaken)
		if err != nil {
			fmt.Println("Error writing to CSV:", err.Error())
		}
		wg.Done()
	})

	for i := 0; i < NumberRequests; i++ {
		currentTime := time.Now().UnixNano() // get current time in nanoseconds
		payload := fmt.Sprintf("%d:%d", i, currentTime)
		requestTopic := RequestTopicBase + uniqueID
		client.Publish(requestTopic, QoS, false, payload)
	}

	// Instead of sleeping, wait until all responses have been processed
	wg.Wait()

	if err != nil {
		fmt.Println("Error writing to CSV:", err.Error())
	}

	// disconnect from the broker
	defer client.Disconnect(250)
}

func clientRoutine(wg *sync.WaitGroup) {
	defer wg.Done()
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

func writeToCSV(writer *csv.Writer, fibonacciIn, fibonacciOut int, timeTaken int64) error {
	return writer.Write([]string{fmt.Sprintf("%d", fibonacciIn), fmt.Sprintf("%d", fibonacciOut), fmt.Sprintf("%d", timeTaken)})
}

func writeCSVHeader(writer *csv.Writer) error {
	return writer.Write([]string{"Input", "Output", "timeTaken"})
}
