package main

import (
	"encoding/csv"
	"fmt"
	"net/rpc"
	"os"
	"time"
)

const (
	// ServerHost = "localhost" // local
	// DataFilePath = "../data/"  // local
	ServerHost     = "server"     // docker
	DataFilePath   = "/app/data/" // docker
	ServerPort     = "1313"
	ServerType     = "tcp"
	NumberRequests = 40
	EndMessage     = "END"
)

type Request struct {
	Number int
}

type Response struct {
	Fibonacci int
}

func main() {
	client, err := rpc.Dial("tcp", ServerHost+":"+ServerPort)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}

	uniqueID := time.Now().UnixNano()
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

	request := Request{}
	response := Response{}

	for i := 0; i < NumberRequests; i++ {
		request.Number = i

		t1 := time.Now()

		err := client.Call("Fibonacci.Calc", request, &response)
		if err != nil {
			fmt.Println("Error during the function call:", err)
		}

		timeTaken := time.Now().Sub(t1).Nanoseconds()
		err = writeToCSV(writer, i, response.Fibonacci, timeTaken)
		if err != nil {
			fmt.Println("Error writing to CSV:", err.Error())
		}
	}
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
	return writer.Write([]string{"Input", "Output", "timeTaken"})
}
