package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	DataFilePath      = "/app/data/"
	BrokerAddress     = "tcp://mosquitto:1883"
	RequestTopic      = "fibonacci/request"
	ResponseTopicBase = "fibonacci/response/"
	NumberRequests    = 40
)

func main() {
	uniqueID := time.Now().UnixNano()

	// convert uniqueId to string
	uniqueIDStr := strconv.FormatInt(uniqueID, 10)

	opts := createClientOptions(uniqueIDStr, BrokerAddress)
	client := connect(uniqueIDStr, opts)

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

	token := client.Subscribe(ResponseTopicBase+"#", 0, func(client mqtt.Client, msg mqtt.Message) {
		number, _ := strconv.Atoi(string(msg.Payload()))
		err := writeToCSV(writer, int(number), 2, int64(time.Now().Nanosecond()))
		if err != nil {
			fmt.Println("Error writing to CSV:", err.Error())
		}
	})
	token.Wait()

	for i := 0; i < NumberRequests; i++ {
		text := fmt.Sprintf("%d", i)
		token := client.Publish(RequestTopic, 0, false, text)
		token.Wait()
	}

	// Sleep for some time to allow all responses to arrive.
	time.Sleep(10 * time.Second)
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

func openCSVFile(basePath string, uniqueID int64) (*csv.Writer, *os.File, error) {
	timestamp := time.Now().Format("20060102150405")
	filePath := fmt.Sprintf("%s_%d_%s.csv", basePath, uniqueID, timestamp)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	writer := csv.NewWriter(file)
	return writer, file, nil
}

func writeToCSV(writer *csv.Writer, fibonacciIn, fibonacciOut int, timeTaken int64) error {
	return writer.Write([]string{fmt.Sprintf("%d", fibonacciIn), fmt.Sprintf("%d", fibonacciOut), fmt.Sprintf("%d", timeTaken)})
}

func writeCSVHeader(writer *csv.Writer) error {
	return writer.Write([]string{"Input", "Output", "TimeTaken"})
}
