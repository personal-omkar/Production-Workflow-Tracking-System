#!/bin/bash

# Define an array of directory paths
DIRECTORIES=(
  "/RUBBER"
  "/RUBBER/qrimg"
  "/RUBBER/pdf"
  "/RUBBER/xlsx"
  "/RUBBER/scripts"
)

# Function to check and create directories
check_and_create_dirs() {
  for DIR in "${DIRECTORIES[@]}"; do
    if [ ! -d "$DIR" ]; then
      echo "Directory $DIR does not exist. Creating..."
      mkdir -p "$DIR"
      chmod -R 777 "$DIR"
      if [ $? -eq 0 ]; then
        echo "Directory $DIR created successfully."
      else
        echo "Failed to create directory $DIR."
        exit 1
      fi
    else
      echo "Directory $DIR already exists."
    fi
  done
}

# Run the functions
check_and_create_dirs