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
	ServerHost = "localhost"
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
	fmt.Println("Aguardando conexões dos cliente em " + ServerHost + ":" + ServerPort)
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

func processRequestBytes(conn net.Conn) {
	var fromClient = make([]byte, 1024)

	for {
		// recebe dados
		mLen, err := conn.Read(fromClient)
		if err != nil {
			fmt.Println("Erro na leitura dos dados do cliente:", err.Error())
		}
		//fmt.Println("Dado recebido: ", string(fromClient[:mLen]))

		// envia resposta
		_, err = conn.Write([]byte(string(fromClient[:mLen])))

		if string(fromClient[:mLen]) == EndMessage {
			break
		}
	}

	// fecha conexão
	conn.Close()
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
