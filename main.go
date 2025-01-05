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
	userName        *string   = flag.String("username", "", "Your Instagram user handle. *Required")
	accessToken     *string   = flag.String("access_token", "", "Your Instagram user access token. *Required")
	outDir          *string   = flag.String("outDir", "", "Output directory path for your report.")
	includeVerified *bool     = flag.Bool("includeVerified", false, "Boolean to include verified accounts in report.")
	now             time.Time = time.Now()
)

const baseURL = "https://graph.instagram.com/"

type Config struct {
	Username        string `json:"username"`
	AccessToken     string `json:"access_token"`
	OutDir          string `json:"outDir"`
	IncludeVerified bool   `json:"includeVerified"`
}

func config() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Println("No config file was found; command line arguments will be used.")

		// Verify required arguments are met
		if *userName == "" {
			log.Fatalln("username commannd line argument must not be empty.")
		}

		if *accessToken == "" {
			log.Fatalln("access_token commannd line argument must not be empty.")
		}

		return
	}
	defer file.Close()

	log.Println("A config file was found; key/values that are specified will override command line arguments & defaults.")

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}

	if config.Username == "" {
		log.Fatalln("username config file value cannot be empty. Please set to your Instagram user handle.")
	} else if config.Username != "" {
		*userName = config.Username
	}

	if config.AccessToken == "" {
		log.Fatalln("access_token config file value cannot be empty. Please set to your Instagram user access token.")
	} else if config.AccessToken != "" {
		*accessToken = config.AccessToken
	}

	if config.OutDir != "" {
		*outDir = config.OutDir
	}

	if config.IncludeVerified {
		*includeVerified = true
	}
}

func fetchData(userName string, accessToken string) ([]byte, error) {
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

	// Log the raw response for debugging
	log.Printf("Raw response for user %s: %s", userName, string(body))

	return body, nil
}

func main() {
	flag.Parse()

	logFile, err := os.Create(*outDir + "UnmutualConnections_" + now.Format("2006-01-02_15-04-05") + ".log")
	if err != nil {
		log.Println("Error creating log file:", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	config()

	not := ""
	if !*includeVerified {
		not = " not"
	}
	log.Printf("Verified accounts are%s included.", not)

	// Fetch data for the user
	rawData, err := fetchData(*userName, *accessToken)
	if err != nil {
		log.Fatalf("Error fetching data for user %s: %v", *userName, err)
	}

	// Parse the data (assuming JSON format)
	var formattedData map[string]interface{}
	if err := json.Unmarshal(rawData, &formattedData); err != nil {
		log.Fatalf("Error parsing %s's data: %v", *userName, err)
	}

	// Process the result...
}
