# Go TCP Client-Server Fibonacci Calculator

This project consists of a server application that calculates the Fibonacci sequence given a number and a client application that requests the calculations. It uses Go, Docker, and Docker Compose.

The server is a TCP server that listens for connections from multiple clients. The client sends a series of numbers to the server, which then calculates the Fibonacci number for each and sends it back. The client then writes the received Fibonacci number and the time taken for the request to a unique CSV file.

## Prerequisites

You need to have the following installed on your machine:

- [Go (version 1.18 or later)](https://golang.org/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)


## How to run

1. Clone the repository to your local machine.
2. Navigate to the atividades/atv4 directory of the project.
3. Run `docker compose build` to build the Docker images for the client and the server.
4. Run `docker compose up` to start the server and the clients.

The server will start and wait for connections. The clients will connect to the server, send numbers, receive the corresponding Fibonacci numbers, and write them along with the request time to unique CSV files in the `./data` directory.

To stop the execution, use `CTRL+C` in the terminal where `docker compose up` is running.

If you want to run the server and the clients in separate terminals, you can run `docker compose up server` in one terminal and `docker compose up client` in another.

If you want to run multiple clients, you can run `docker compose up --scale client=3` to run 3 clients, for example.
## Clean up

When you're done, you can remove the Docker containers, networks, and volumes by running `docker-compose down`. If you want to rebuild the Docker images for any reason (e.g., after modifying the Go code), run `docker-compose build` again before `docker-compose up`.
