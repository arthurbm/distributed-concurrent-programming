// socket-client project main.go
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	ServerHost     = "localhost"
	ServerPort     = "1313"
	ServerType     = "tcp"
	SampleSize     = 30
	NumberRequests = 1000
	EndMessage     = "END"
)

func main() {

	for i := 0; i < SampleSize; i++ {

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
}

func comServerBytes(conn net.Conn) {
	fromServer := make([]byte, 1024)
	toServer := ""

	for i := 0; i < NumberRequests; i++ {

		// envia mensagem
		toServer = "Mensagem #" + strconv.Itoa(i)
		_, err := conn.Write([]byte(toServer))
		if err != nil {
			fmt.Println("Erro no envio dos dados do servidor:", err.Error())
		}

		// recebe resposta do servidor
		//mLen, err := conn.Read(fromServer)
		_, err = conn.Read(fromServer)
		if err != nil {
			fmt.Println("Erro no recebimento dos dados do servidor:", err.Error())
		}
		//fmt.Println("Dado: ", string(fromServer[:mLen]))
	}
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
