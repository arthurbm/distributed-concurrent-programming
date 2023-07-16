package main

import (
	"log"
	"os/exec"
)

const (
	NumClients = 5
)

func main() {
	for i := 0; i < NumClients; i++ {
		cmd := exec.Command("go", "run", "../client/client.go")
		err := cmd.Start()
		if err != nil {
			log.Fatalf("Failed to start client %d: %s", i, err)
		}
		log.Printf("Started client %d", i)
	}
}
