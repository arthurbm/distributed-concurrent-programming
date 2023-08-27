#!/bin/bash
 
# Directory paths
SCRIPT_DIR="$(dirname "$0")"
DATA_DIR="$SCRIPT_DIR/../data"
SAMPLES_DIR="$SCRIPT_DIR/../samples/async"
 
# Check if the argument is provided
if [ -z "$1" ]; then
  echo "Please provide the data folder number as an argument."
  exit 1
fi
 
# Create the samples directory if it doesn't exist
mkdir -p "$SAMPLES_DIR"
 
# Create the specific data folder inside samples
mkdir -p "$SAMPLES_DIR/data$1"
 
# Check if there are any CSV files in the data folder
if [ -z "$(ls -A "$DATA_DIR"/*.csv 2>/dev/null)" ]; then
  echo "No CSV files found in the data folder."
  exit 1
fi
 
# Move the CSV files to the respective data folder inside samples
mv "$DATA_DIR"/*.csv "$SAMPLES_DIR/data$1/"
 
echo "Files moved to samples/async/data$1/"