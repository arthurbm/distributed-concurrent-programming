# Go Middleware Performance Benchmark

This project implements a performance benchmarking system for various middleware technologies in Go. The primary benchmark uses a Fibonacci calculation function as the test workload, which provides an excellent test case as its computational complexity increases with larger input numbers.

## Project Structure

- `app/`: Main application directory containing all middleware implementations and benchmarks
  - `tcp/`: TCP implementation of client-server architecture
  - `udp/`: UDP implementation of client-server architecture
  - `rpc/`: RPC implementation using Go's built-in RPC package
  - `mqtt/`: MQTT implementation using the Paho MQTT client
  - `rabbitmq/`: RabbitMQ implementation using AMQP
  - `utils/`: Shared utility functions
  - `benchmark/`: Contains Jupyter notebooks with benchmark analysis
  - `run_all_experiments.sh`: Script to run all experiments across different middleware implementations

## Benchmark Function

The core function used for benchmarking is the recursive Fibonacci calculation:

```go
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}
```

This function is ideal for benchmarking because:
1. It has exponentially increasing complexity as input size grows
2. It's CPU-intensive, which helps measure processing capabilities
3. The function is deterministic, making results reproducible
4. Different middleware technologies can be compared using the same workload

## How to Run Experiments

1. Make sure you have Docker and Docker Compose installed on your system
2. To run all experiments at once:
   ```bash
   cd app
   chmod +x run_all_experiments.sh
   ./run_all_experiments.sh
   ```
3. To run experiments for a specific middleware:
   ```bash
   cd app/<middleware-name>
   docker compose build
   docker compose up -d server
   docker compose up client -d --scale client=<number-of-clients>
   ```

## Viewing Results

The benchmark results are stored in CSV files in the `data/` directory within each middleware implementation folder. Each CSV file contains:
- Input: The Fibonacci number to calculate
- Output: The calculated Fibonacci value
- timeTaken: Time taken for the request/response cycle in milliseconds

## Analyzing Results

The `app/benchmark/` directory contains Jupyter notebooks that can be used to analyze the benchmark results. These notebooks provide visualizations and comparative analysis of the different middleware technologies.

To view the analysis:
1. Install Jupyter Notebook on your system
2. Navigate to the benchmark directory:
   ```bash
   cd app/benchmark
   jupyter notebook
   ```
3. Open the `main.ipynb` file to see the analysis

## Middleware Implementations

Each middleware implementation follows the same pattern:
1. A server that listens for client requests
2. A client that sends numbers to calculate Fibonacci values
3. The server calculates and returns results
4. The client measures and records the time taken

The project currently supports these middleware technologies:
- TCP: Direct socket communication using TCP
- UDP: Connectionless communication using UDP
- RPC: Remote Procedure Call using Go's built-in rpc package
- MQTT: Message Queuing Telemetry Transport protocol
- RabbitMQ: Advanced Message Queuing Protocol implementation

## Prerequisites

- Go 1.18+
- Docker and Docker Compose
- For analytics: Python with Jupyter Notebook

## License

[Insert your license information here] 