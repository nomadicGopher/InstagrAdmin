package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
    // Get user's Instagram username
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter your Instagram username: ")
    username, _ := reader.ReadString('\n')

    // Replace with your Instagram access token
    accessToken := "YOUR_ACCESS_TOKEN"

    // Get user's ID from username
    userIdUrl := fmt.Sprintf("https://graph.instagram.com/?fields=id&access_token=%s&username=%s", accessToken, username)
    userIdResponse, err := http.Get(userIdUrl)
    if err != nil {
        fmt.Println("Error fetching user ID:", err)
        return
    }
    defer userIdResponse.Body.Close()

    // Parse user ID
    userIdData := map[string]interface{}{}
    err = json.NewDecoder(userIdResponse.Body).Decode(&userIdData)
    if err != nil {
        fmt.Println("Error parsing user ID data:", err)
        return
    }
    userId := userIdData["id"].(string)

    // Get user's following list
    followingUrl := fmt.Sprintf("https://graph.instagram.com/%s/following?access_token=%s", userId, accessToken)
    followingResponse, err := http.Get(followingUrl)
    if err != nil {
        fmt.Println("Error fetching following list:", err)
        return
    }
    defer followingResponse.Body.Close()

    // Parse following list
    followingData := map[string]interface{}{}
    err = json.NewDecoder(followingResponse.Body).Decode(&followingData)
    if err != nil {
        fmt.Println("Error parsing following list data:", err)
        return
    }
    followingList := followingData["data"].([]interface{})

    // Check if users follow back
    for _, following := range followingList {
        followingUser := following.(map[string]interface{})
        followingUserId := followingUser["id"].(string)

        // Get following user's profile
        profileUrl := fmt.Sprintf("https://graph.instagram.com/%s?fields=username,follows_count&access_token=%s", followingUserId, accessToken)
        profileResponse, err := http.Get(profileUrl)
        if err != nil {
            fmt.Println("Error fetching profile:", err)
            continue
        }
        defer profileResponse.Body.Close()

        // Parse profile data
        profileData := map[string]interface{}{}
        err = json.NewDecoder(profileResponse.Body).Decode(&profileData)
        if err != nil {
            fmt.Println("Error parsing profile data:", err)
            continue
        }

        // Check if user follows back
        followsCount := profileData["follows_count"].(float64)
        if followsCount > 0 {
            username := profileData["username"].(string)
            fmt.Println("User", username, "follows you back!")
        }
    }
}