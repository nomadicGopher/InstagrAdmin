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
	userName        *string   = flag.String("username", "", "Create a file named username or username.txt with only your IG user handle as the file contents.")
	accessToken     *string   = flag.String("access_token", "", "Create a file named access_token or access_token.txt with only your IG access token as the file contents.")
	outDir          *string   = flag.String("outDir", "", "Output directory of your results.")
	includeVerified *bool     = flag.Bool("includeVerified", false, "Boolean to include verified accounts in report.")
	now             time.Time = time.Now()
)

const baseURL = "https://graph.instagram.com/"

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

	not := ""
	if !*includeVerified {
		not = " not"
	}
	log.Printf("Verified accounts are%s included.", not)

	// Fetch data for the user
	data, err := fetchData(*userName, *accessToken)
	if err != nil {
		log.Fatalf("Error fetching data for user %s: %v", *userName, err)
	}

	// Parse the data (assuming JSON format)
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatalf("Error parsing %s's data: %v", *userName, err)
	}

	// Process the result...
}
