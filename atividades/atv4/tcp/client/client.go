// socket-client project main.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
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

func main() {
	// estabelece conexão
	conn, err := net.Dial(ServerType, ServerHost+":"+ServerPort)
	fmt.Printf("Conectado ao servidor %s:%s\n", ServerHost, ServerPort)
	if err != nil {
		panic(err)
	}

	// envia dado/recebe resposta
	t1 := time.Now()
	// comServerBytes(conn)
	comServerJson(conn)
	fmt.Println(time.Now().Sub(t1).Milliseconds())

	// fecha conexão
	defer conn.Close()
}

type Request struct {
	Number int `json:"number"`
}

type Response struct {
	Fibonacci int `json:"fibonacci"`
}

func openCSVFile(basePath string, uniqueID int64) (*csv.Writer, *os.File, error) {
	timestamp := time.Now().Format("20060102150405")
	filePath := fmt.Sprintf("%s_%d_%s.csv", basePath, uniqueID, timestamp)

	// Open the file in append mode, create if doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	// Initialize csv writer
	writer := csv.NewWriter(file)
	return writer, file, nil
}

func writeToCSV(writer *csv.Writer, fibonacciIn, fibonacciOut int, timeTaken int64) error {
	return writer.Write([]string{fmt.Sprintf("%d", fibonacciIn), fmt.Sprintf("%d", fibonacciOut), fmt.Sprintf("%d", timeTaken)})
}

func writeCSVHeader(writer *csv.Writer) error {
	return writer.Write([]string{"Input", "Output", "timeTaken"})
}

func clearCSVContent(path string) error {
	// Open the file in append mode, create if doesn't exist
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Initialize csv writer
	writer := csv.NewWriter(file)
	return writer.Write([]string{})
}

func comServerJson(conn net.Conn) {
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)
	request := Request{}
	response := Response{}

	uniqueID := time.Now().UnixNano()

	// Open CSV file
	writer, file, err := openCSVFile(DataFilePath, uniqueID)

	if err != nil {
		fmt.Println("Error opening the file:", err.Error())
		return
	}
	defer file.Close()
	defer writer.Flush()

	// Clear CSV content
	err = clearCSVContent(DataFilePath)

	// Write CSV header
	err = writeCSVHeader(writer)
	if err != nil {
		fmt.Println("Error writing CSV header:", err.Error())
		return
	}

	for i := 0; i < NumberRequests; i++ {
		// Prepare the request
		request.Number = i

		// Time the request
		t1 := time.Now()

		// Sends the request to the server
		err := enc.Encode(&request)
		if err != nil {
			fmt.Println("Error sending data to the server:", err.Error())
		}

		// Receives response from the server
		err = dec.Decode(&response)
		if err != nil {
			fmt.Println("Error receiving data from the server:", err.Error())
		}

		// Calculate the request time
		timeTaken := time.Now().Sub(t1).Nanoseconds()

		// Write the Fibonacci number and the request time to CSV
		err = writeToCSV(writer, i, response.Fibonacci, timeTaken)
		if err != nil {
			fmt.Println("Error writing to CSV:", err.Error())
		}
	}
}
