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

const (
	apiBaseURL   = "https://graph.instagram.com/" // https://developer.microsoft.com/en-us/graph/graph-explorer
	clientID     = "your-app-id"
	clientSecret = "your-app-secret"
	redirectURI  = "http://localhost:8080/callback"
)

type Config struct {
	Username        string `json:"username"`
	AccessToken     string `json:"access_token"`
	OutDir          string `json:"output_directory"`
	IncludeVerified bool   `json:"include_verified"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

var (
	config   Config
	logFile  *os.File
	userData map[string]interface{}
)

func main() {
	loadConfig()

	initLogFile()
	defer logFile.Close()

	http.ListenAndServe(":8080", nil)

	fetchAccessToken()

	fetchUserData()

	buildResults()
}

func promptUserForConfig() {
	var (
		err             error
		username        string
		outDir          string
		includeVerified string
	)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter your Instagram username: ")
		if username, err = reader.ReadString('\n'); err != nil {
			log.Fatalln("Error reading username: ", err)
		}
		config.Username = strings.TrimSpace(username)
		if config.Username != "" && !strings.Contains(config.Username, " ") {
			break
		}
		fmt.Println("username cannot be empty or contain spaces.")
	}

	for {
		fmt.Print("Enter the output directory path for your report (leave empty for the same directory as the program): ")
		if outDir, err = reader.ReadString('\n'); err != nil {
			log.Fatalln("Error reading output directory: ", err)
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
		if includeVerified, err = reader.ReadString('\n'); err != nil {
			log.Fatalf("Error reading include_verified: %v", err)
		}
		config.IncludeVerified = strings.TrimSpace(includeVerified) == "true"
		if !strings.Contains(includeVerified, " ") {
			break
		}
		fmt.Println("include_verified input cannot contain spaces.")
	}
}

func loadConfig() {
	var (
		err        error
		configFile *os.File
	)

	if configFile, err = os.Open("config.json"); err != nil {
		log.Println("No config file was found; prompting user for configuration.")
		promptUserForConfig()
		return
	}
	defer configFile.Close()

	log.Println("A config file was found.")

	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&config); err != nil {
		log.Fatalln("Error decoding config file: ", err)
	}

	config.Username = strings.TrimSpace(config.Username)
	if config.Username == "" || strings.Contains(config.Username, " ") {
		log.Fatalln("Config file must contain a valid 'username' field without spaces.")
	}

	config.OutDir = strings.TrimSpace(config.OutDir)
	if config.OutDir == "" {
		config.OutDir = "./" // Default to current directory if empty
	} else if strings.Contains(config.OutDir, " ") {
		log.Fatalln("Config file must contain a valid 'output_directory' field without spaces.")
	}

	config.IncludeVerified = strings.TrimSpace(fmt.Sprint(config.IncludeVerified)) == "true" // Default to false if not == 'true'
}

func initLogFile() {
	var err error

	logFileName := config.OutDir + "UnmutualConnections_" + time.Now().Format("2006-01-02_15-04-05") + ".log"

	if logFile, err = os.Create(logFileName); err != nil {
		log.Println("Error creating log file:", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
}

func fetchAccessToken() {
	var (
		err       error
		resp      *http.Response
		tokenResp AccessTokenResponse
	)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprint("https://www.facebook.com/v12.0/dialog/oauth?client_id=", clientID, "&redirect_uri=", redirectURI, "&scope=email")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			log.Fatalln("HTTP Error 400 (code not found)") // http.Error(w, "Code not found", http.StatusBadRequest) return
		}

		tokenURL := fmt.Sprint("https://graph.facebook.com/v12.0/oauth/access_token?client_id=", clientID, "&redirect_uri=", redirectURI, "&client_secret=", clientSecret, "&code=", code)
		if resp, err = http.Get(tokenURL); err != nil {
			log.Fatalln("HTTP Error 500 (error getting access token)") // http.Error(w, "Error getting access token", http.StatusInternalServerError) return
		}
		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			log.Fatalln("HTTP Error 500 (error decoding access token response)") // http.Error(w, "Error decoding access token response", http.StatusInternalServerError) return
		}

		config.AccessToken = tokenResp.AccessToken
		log.Println("Set Access Token: ", tokenResp.AccessToken)
	})
}

func fetchUserData() {
	var (
		err  error
		resp *http.Response
		body []byte
	)

	url := fmt.Sprint(apiBaseURL, config.Username, "?access_token=", config.AccessToken)
	if resp, err = http.Get(url); err != nil {
		log.Fatalf("error fetching %s's user data: %v", config.Username, err)
	}
	defer resp.Body.Close()

	if body, err = io.ReadAll(resp.Body); err != nil {
		log.Fatalf("error reading %s's response body: %v", config.Username, err)
	}

	log.Println("Raw response for user ", config.Username, ": ", string(body))

	if err = json.Unmarshal(body, &userData); err != nil {
		log.Fatalf("Error parsing %s's data: %v", config.Username, err)
	}
}

func buildResults() {
	not := ""
	if !config.IncludeVerified {
		not = " not"
	}
	log.Print("Verified accounts are", not, " included.")

	log.Println(userData) // TODO
}
