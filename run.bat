@echo off

:: Check if username file exists
if not exist username if not exist username.txt (
  echo Error: username file not found. Please create a file named 'username' or 'username.txt' with your Instagram username.
  exit /b 1
)

:: Check if access_token file exists
if not exist access_token if not exist access_token.txt (
  echo Error: access_token file not found. Please create a file named 'access_token' or 'access_token.txt' with your Instagram access token.
  exit /b 1
)

:: Fetch IG username
if exist username (
  set /p USERNAME=<username
) else if exist username.txt (
  set /p USERNAME=<username.txt
)

:: Fetch IG access_token
if exist access_token (
  set /p ACCESS_TOKEN=<access_token
) else if exist access_token.txt (
  set /p ACCESS_TOKEN=<access_token.txt
)

go run main.go -username="%USERNAME%" -access_token="%ACCESS_TOKEN%"