#!/bin/bash

# Check if username file exists
if [ ! -f username ] && [ ! -f username.txt ]; then
  echo "Error: username file not found. Please create a file named 'username' or 'username.txt' with your Instagram username."
  exit 1
fi

# Check if access_token file exists
if [ ! -f access_token ] && [ ! -f access_token.txt ]; then
  echo "Error: access_token file not found. Please create a file named 'access_token' or 'access_token.txt' with your Instagram access token."
  exit 1
fi

# Fetch IG username
if [ -f username ]; then
  USERNAME=$(cat username)
elif [ -f username.txt ]; then
  USERNAME=$(cat username.txt)
fi

# Fetch IG access_token
if [ -f access_token ]; then
  ACCESS_TOKEN=$(cat access_token)
elif [ -f access_token.txt ]; then
  ACCESS_TOKEN=$(cat access_token.txt)
fi

go run main.go -username="$USERNAME" -access_token="$ACCESS_TOKEN"