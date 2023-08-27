package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const (
	ServerAddr = "0.0.0.0" // docker
	// ServerAddr = "localhost" // local
	ServerPort = "1313"
	ServerType = "udp"
	EndMessage = "END"
)

func main() {

	// creates UDP listener
	fmt.Println("Server running...")
	serverAddr, _ := net.ResolveUDPAddr(ServerType, ServerAddr+":"+ServerPort)
	serverConn, err := net.ListenUDP(ServerType, serverAddr)
	if err != nil {
		fmt.Println("Error listening for connections:", err.Error())
		os.Exit(1)
	}
	defer serverConn.Close()

	// waits for connections
	fmt.Println("Waiting for client connections at " + ServerAddr + ":" + ServerPort)
	for {
		buffer := make([]byte, 1024)
		n, addr, err := serverConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading data from the client:", err.Error())
			continue
		}

		go processRequestJson(serverConn, buffer[:n], addr)
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

func processRequestJson(serverConn *net.UDPConn, data []byte, addr *net.UDPAddr) {
	var request Request
	json.Unmarshal(data, &request)

	var response Response

	// Process the request (calculates Fibonacci)
	response.Fibonacci = fibonacci(request.Number)

	// Sends the response
	respData, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error encoding response:", err.Error())
		return
	}

	serverConn.WriteToUDP(respData, addr)
}
