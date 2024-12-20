#!/bin/bash

# Check if the executable exists
if [ -f "InstagrAdmin" ]; then
  # Run the executable with parameters
  ./InstagrAdmin -username="$(cat username)" -access_token="$(cat access_token)"
else
  # Build the executable first
  go build main.go
  # Run the executable with parameters
  ./InstagrAdmin -username="$(cat username)" -access_token="$(cat access_token)"
fi