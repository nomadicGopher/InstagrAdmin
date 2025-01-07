package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	userName        *string = flag.String("username", "", "Your Instagram user handle. *Required")
	accessToken     *string = flag.String("access_token", "", "Your Instagram user access token. *Required")
	outDir          *string = flag.String("outDir", "", "Output directory path for your report.")
	includeVerified *bool   = flag.Bool("includeVerified", false, "Boolean to include verified accounts in report.")
	debug           *bool   = flag.Bool("debug", false, "Enable developer logging.")
)

const baseURL = "https://graph.instagram.com/"

type Config struct {
	Username        string `json:"username"`
	AccessToken     string `json:"access_token"`
	OutDir          string `json:"outDir"`
	IncludeVerified bool   `json:"includeVerified"`
}

func checkRequiredConfigs(config Config, method string) {
	if config.Username == "" {
		log.Fatalf("username %s must not be empty.\n", method)
	}

	if config.AccessToken == "" {
		log.Fatalf("access_token %s must not be empty.\n", method)
	}
}

func loadConfig() Config {
	config := Config{
		Username:        *userName,
		AccessToken:     *accessToken,
		OutDir:          *outDir,
		IncludeVerified: *includeVerified,
	}

	configFile, err := os.Open("config.json")
	if err != nil {
		log.Println("No config file was found; command line arguments will be used.")

		checkRequiredConfigs(config, "commannd line argument")

		return config
	}
	defer configFile.Close()

	log.Println("A config file was found. Specified key/value pairs will override command line arguments & their potential defaults.")

	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalln("Error decoding config file: ", err)
	}

	checkRequiredConfigs(config, "config file value")

	return config
}

func fetchData(userName, accessToken string) ([]byte, error) {
	url := fmt.Sprintf("%s%s?access_token=%s", baseURL, userName, accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if *debug {
		log.Printf("Raw response for user %s: %s", userName, string(body))
	}

	return body, nil
}

func main() {
	flag.Parse()

	logFile, err := os.Create(*outDir + "UnmutualConnections_" + time.Now().Format("2006-01-02_15-04-05") + ".log")
	if err != nil {
		log.Println("Error creating log file:", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	config := loadConfig()

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
