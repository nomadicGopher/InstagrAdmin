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
	userName        *string   = flag.String("username", "", "YOUR_IG_USERNAME")
	accessToken     *string   = flag.String("access_token", "", "YOUR_IG_USER_ACCESS_TOKEN")
	outDir          *string   = flag.String("outDir", "", "Output directory of log file.")
	includeVerified *bool     = flag.Bool("includeVerified", false, "Boolean to include verified accounts in report.")
	now             time.Time = time.Now()
)

const baseURL = "https://graph.instagram.com/"

func main() {
	flag.Parse()

	logFile, err := os.Create(*outDir + "UnmutualConnections_" + now.Format("2006-01-02_15-04-05") + ".log")
	if err != nil {
		log.Println("Error creating log file:", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.Println("Include verified accounts in report: ", *includeVerified)

	// Get user data
	userProfileURL := fmt.Sprintf("%s?fields=id&access_token=%s&username=%s", baseURL, *accessToken, *userName)
	userProfileResponse, err := http.Get(userProfileURL)
	if err != nil {
		log.Fatalln("Error fetching ", *userName, "'s ID:", err)
	}
	defer userProfileResponse.Body.Close()

	// Parse userID
	userData := map[string]interface{}{}
	err = json.NewDecoder(userProfileResponse.Body).Decode(&userData)
	if err != nil {
		log.Fatalln("Error parsing ", *userName, "'s data:", err)
	}
	userID := userData["id"].(string)

	// Get user's following list
	usersFollowingDataURL := fmt.Sprintf("%s%s/following?access_token=%s", baseURL, userID, *accessToken)
	usersFollowingDataResponse, err := http.Get(usersFollowingDataURL)
	if err != nil {
		log.Fatalln("Error fetching ", *userName, "'s following list:", err)
	}
	defer usersFollowingDataResponse.Body.Close()

	// Parse following list
	usersFollowingData := map[string]interface{}{}
	err = json.NewDecoder(usersFollowingDataResponse.Body).Decode(&usersFollowingData)
	if err != nil {
		log.Fatalf("Error parsing %s's following list data: %v\n", *userName, err)
	}
	usersFollowingList := usersFollowingData["data"].([]interface{})

	// Check if user follows back
	for _, followee := range usersFollowingList {
		followsBack := false
		followeesInfo := followee.(map[string]interface{})
		followeesID := followeesInfo["id"].(string)
		followeesUserName := followeesInfo["username"].(string)
		followeesFullName := followeesInfo["full_name"].(string)
		followeesAccountStatus := followeesInfo["account_status"].(map[string]interface{})
		followeeIsVerified := followeesAccountStatus["is_verified"].(bool)

		if !*includeVerified && followeeIsVerified {
			continue
		}

		// Get followee's following list
		followeesFollowingDataURL := fmt.Sprintf("%s%s/following?access_token=%s", baseURL, followeesID, *accessToken)
		followeesFollowingDataResponse, err := http.Get(followeesFollowingDataURL)
		if err != nil {
			log.Fatalln("Error fetching ", followeesUserName, "'s following list:", err)
			continue
		}
		defer followeesFollowingDataResponse.Body.Close()

		// Parse followee's following list
		followeesFollowingData := map[string]interface{}{}
		err = json.NewDecoder(followeesFollowingDataResponse.Body).Decode(&followeesFollowingData)
		if err != nil {
			log.Fatalln("Error parsing ", followeesUserName, "'s following list data:", err)
			continue
		}
		followeesFollowingList := followeesFollowingData["data"].([]interface{})

		// Check if userID is included in the list of followee's following list
		for _, followeesFollowee := range followeesFollowingList {
			followeesFolloweeInfo := followeesFollowee.(map[string]interface{})
			if followeesFolloweeInfo["id"].(string) == userID {
				followsBack = true
			}
		}

		if !followsBack {
			log.Println(followeesUserName, " ( ", followeesFullName, " ) does not follow ", userID, " back.")
		}
	}
}
