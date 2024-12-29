#!/bin/bash

# Function to check if the executable exists in either format
check_executable() {
    local exe_name=$1
    # Check for Linux (Debian) and Windows formats
    if [ -f "$exe_name" ] || [ -f "${exe_name}.exe" ]; then
        return 0
    else
        return 1
    fi
}

# Function to build the executable if it does not exist
build_executable() {
    go build main.go
}

if check_executable "InstagrAdmin"; then
    # Run the executable with parameters
    ./InstagrAdmin -username="$(cat username)" -access_token="$(cat access_token)"
else
    # Build the executable first if it does not exist
    build_executable
    # Run the executable with parameters after building
    ./InstagrAdmin -username="$(cat username)" -access_token="$(cat access_token)"
fi