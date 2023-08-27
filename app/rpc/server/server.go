package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
)

const (
	ServerHost = "0.0.0.0"
	ServerPort = "1313"
)

type Fibonacci struct{}

type Request struct {
	Number int
}

type Response struct {
	Fibonacci int
}

// Fibonacci function
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func (f *Fibonacci) Calc(req Request, res *Response) error {
	res.Fibonacci = fibonacci(req.Number)
	return nil
}

func main() {
	fib := new(Fibonacci)
	rpc.Register(fib)
	l, err := net.Listen("tcp", ServerHost+":"+ServerPort)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	fmt.Println("Server is running...")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
