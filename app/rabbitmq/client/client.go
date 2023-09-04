package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// ServerHost = "localhost" // local
	// DataFilePath = "../data/"  // local
	DataFilePath   = "/app/data/" // docker
	NumberRequests = 40
	BrokerAddress  = "amqp://guest:guest@localhost:5672/"
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err.Error())
	}
}

func fibonacciRPC(n int) (res int, err error) {
	// estabelece conexão
	amqpServerURL := BrokerAddress
	conn, err := amqp.Dial(amqpServerURL)
	fmt.Printf("Connected to server %s\n", amqpServerURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// abrindo um cahnnel para a instância de rabbitmq que estabelecemos
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // queue name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	corrId := randomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			failOnError(err, "Failed to convert body to integer")
		}
	}
	return
}

func main() {
	uniqueID := time.Now().UnixNano()
	// file configs
	writer, file, err := openCSVFile(DataFilePath, uniqueID)
	failOnError(err, "Error opening the file")
	defer file.Close()
	defer writer.Flush()

	err = writeCSVHeader(writer)
	failOnError(err, "Error writing CSV header")

	for i := 0; i < NumberRequests; i++ {
		log.Printf(" [x] Requesting fib(%d)\n", i)
		start := time.Now()
		res, _ := fibonacciRPC(i)
		err = writeToCSV(writer, i, 2, int64(time.Since(start)))
		failOnError(err, "Failed to write to csv file")
		log.Printf(" [.] Got %d\n", res)
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
	return writer.Write([]string{"Input", "Output", "TimeTaken"})
}
