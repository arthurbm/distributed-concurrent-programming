#!/bin/bash

# Check if the argument is provided
if [ -z "$1" ]; then
  echo "Please provide the data folder number as an argument."
  exit 1
fi

# Create the samples directory if it doesn't exist
mkdir -p ../samples

# Create the specific data folder inside samples
mkdir -p ../samples/data$1

# Check if there are any CSV files in the data folder
if [ -z "$(ls -A ../data/*.csv 2>/dev/null)" ]; then
  echo "No CSV files found in the data folder."
  exit 1
fi

# Move the CSV files to the respective data folder inside samples
mv ../data/*.csv ../samples/data$1/

echo "Files moved to samples/data$1/"
