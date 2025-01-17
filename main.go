package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const apiBaseURL = "https://graph.instagram.com/" // https://developer.microsoft.com/en-us/graph/graph-explorer

type Config struct {
	// Instagram
	Username        string `json:"username"`
	AccessToken     string `json:"access_token"`
	// Report
	OutDir          string `json:"output_directory"`
	IncludeVerified bool   `json:"include_verified"`
}

func promptUserForConfig(config *Config) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter your Instagram username: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading username: %v", err)
		}
		config.Username = strings.TrimSpace(username)
		if config.Username != "" && !strings.Contains(config.Username, " ") {
			break
		}
		fmt.Println("username cannot be empty or contain spaces.")
	}

	for {
		fmt.Print("Enter your Instagram access_token: ")
		accessToken, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading access_token: %v", err)
		}
		config.AccessToken = strings.TrimSpace(accessToken)
		if config.AccessToken != "" && !strings.Contains(config.AccessToken, " ") {
			break
		}
		fmt.Println("access_token cannot be empty or contain spaces.")
	}

	for {
		fmt.Print("Enter the output directory path for your report (leave empty for the same directory as the program): ")
		outDir, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading output directory: %v", err)
		}
		config.OutDir = strings.TrimSpace(outDir)
		if config.OutDir == "" {
			config.OutDir = "./"
			break
		}
		if !strings.Contains(config.OutDir, " ") {
			break
		}
		fmt.Println("output_directory cannot contain spaces.")
	}

	for {
		fmt.Print("Include verified accounts in report? (true/false): ")
		includeVerified, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading include_verified: %v", err)
		}
		config.IncludeVerified = strings.TrimSpace(includeVerified) == "true"
		if !strings.Contains(includeVerified, " ") {
			break
		}
		fmt.Println("include_verified input cannot contain spaces.")
	}
}

func loadConfig() (Config, bool) {
	config := Config{}

	configFile, err := os.Open("config.json")
	if err != nil {
		log.Println("No config file was found; prompting user for configuration.")
		promptUserForConfig(&config)
		return config, false
	}
	defer configFile.Close()

	log.Println("A config file was found.")

	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalln("Error decoding config file: ", err)
	}

	config.Username = strings.TrimSpace(config.Username)
	if config.Username == "" || strings.Contains(config.Username, " ") {
		log.Fatalln("Config file must contain a valid 'username' field without spaces.")
	}

	config.AccessToken = strings.TrimSpace(config.AccessToken)
	if config.AccessToken == "" || strings.Contains(config.AccessToken, " ") {
		log.Fatalln("Config file must contain a valid 'access_token' field without spaces.")
	}

	config.OutDir = strings.TrimSpace(config.OutDir)
	if config.OutDir == "" {
		config.OutDir = "./" // Default to current directory if empty
	} else if strings.Contains(config.OutDir, " ") {
		log.Fatalln("Config file must contain a valid 'output_directory' field without spaces.")
	}

	config.IncludeVerified = strings.TrimSpace(fmt.Sprintf("%v", config.IncludeVerified)) == "true" // Default to false if not == 'true'

	return config, true
}

func fetchData(userName, accessToken string) ([]byte, error) {
	url := fmt.Sprintf("%s%s?access_token=%s", apiBaseURL, userName, accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	log.Printf("Raw response for user %s: %s", userName, string(body))

	return body, nil
}

func main() {
	config, _ := loadConfig()

	logFileName := config.OutDir + "UnmutualConnections_" + time.Now().Format("2006-01-02_15-04-05") + ".log"

	logFile, err := os.Create(logFileName)
	if err != nil {
		log.Println("Error creating log file:", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	not := ""
	if !config.IncludeVerified {
		not = " not"
	}
	log.Printf("Verified accounts are%s included.", not)

	// Fetch data for the user
	data, err := fetchData(config.Username, config.AccessToken)
	if err != nil {
		log.Fatalf("Error fetching data for user %s: %v", config.Username, err)
	}

	var formattedData map[string]interface{}
	if err := json.Unmarshal(data, &formattedData); err != nil {
		log.Fatalf("Error parsing %s's data: %v", config.Username, err)
	}

	// Process the result...
}
