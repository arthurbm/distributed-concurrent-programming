// socket-client project main.go
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const (
	ServerHost     = "localhost"
	ServerPort     = "1313"
	ServerType     = "tcp"
	NumberRequests = 40
	EndMessage     = "END"
)

func main() {
	// estabelece conexão
	conn, err := net.Dial(ServerType, ServerHost+":"+ServerPort)
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

func comServerJson(conn net.Conn) {
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)
	request := Request{}
	response := Response{}

	for i := 0; i < NumberRequests; i++ {
		// Prepare the request
		request.Number = i

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

		fmt.Printf("Fibonacci of %d is %d\n", i, response.Fibonacci)
	}
}
