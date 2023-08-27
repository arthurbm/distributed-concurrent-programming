#!/bin/bash

# Define the sequence of client numbers
client_numbers="1 5 10 20 40 80"

# Define the methods
methods="rpc tcp udp"

# Iterate through the methods
for method in $methods; do
  echo "Running experiments for $method..."

  # Change directory to where the docker compose.yml file is located for the current method
  cd "$method"

  # Run the experiments with different numbers of clients
  for clients in $client_numbers; do
    echo "Running experiment with $clients clients..."

    # Stop any running containers
    docker compose down

    # Start the server
    docker compose up -d server

    # Wait for the server to initialize (you may adjust the time as needed)
    sleep 10

    # Start the clients
    docker compose up client -d --scale client=$clients

    sleep 30
    
    # Stop the client and server containers
    docker compose down

    # Create the destination directory if it doesn't exist
    mkdir -p ./samples/data${clients}/

    # Move the CSV files to the corresponding folder
    mv ./data/*.csv ./samples/data${clients}/

    # Wait before the next experiment (you may adjust the time as needed)
    sleep 5
  done

  # Change back to the parent directory
  cd ..
done
