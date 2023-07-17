// socket-server project main.go
// developer.com/languages/intro-socket-programming-go/
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const (
	ServerHost = "0.0.0.0" // docker
	// ServerHost = "localhost" // local
	ServerPort = "1313"
	ServerType = "tcp"
	EndMessage = "END"
)

func main() {

	// cria listener
	fmt.Println("Servidor em execução...")
	server, err := net.Listen(ServerType, ServerHost+":"+ServerPort)
	if err != nil {
		fmt.Println("Erro na escuta por conexões:", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	// aguarda conexões
	fmt.Println("Aguardando conexões dos clientes em " + ServerHost + ":" + ServerPort)
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Cliente conectado")

		// cria thread para o cliente
		// go processRequestBytes(conn)
		go processRequestJson(conn)
	}
}

type Request struct {
	Number int `json:"number"`
}

type Response struct {
	Fibonacci int `json:"fibonacci"`
}

// Fibonacci function
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func processRequestJson(conn net.Conn) {
	var request Request
	var response Response
	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	for {
		// Receives data
		err := dec.Decode(&request)
		if err != nil {
			fmt.Println("Error reading data from the client:", err.Error())
			break
		}

		// Process the request (calculates Fibonacci)
		response.Fibonacci = fibonacci(request.Number)

		// Sends the response
		err = enc.Encode(&response)
		if err != nil {
			fmt.Println("Error sending data to the client:", err.Error())
			break
		}
	}
	// Closes the connection
	conn.Close()
}
