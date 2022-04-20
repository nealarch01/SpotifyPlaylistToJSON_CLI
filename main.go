package main

import (
	base64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	url "net/url"
	"os"
	"strings"
)

type APIConfig struct {
	SpotifyClientID     string `json:"SpotifyClientID"`
	SpotifyClientSecret string `json:"SpotifyClientSecret"`
}

type TokenConfig struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func jsonOpenError() {
	fmt.Println("There was an error opening api-config.json")
	fmt.Println("Please make sure you api-config.json file exists with your spotify API key")
	fmt.Println("If your file has a different name, change the os.Open argument")
}

func getAccessToken(clientID string, clientSecret string) TokenConfig {
	var tokenData TokenConfig
	tokenData.AccessToken = ""
	tokenData.TokenType = ""
	tokenData.ExpiresIn = -1
	urlEndpoint := "https://accounts.spotify.com/api/token"

	b64Data := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	var authorization string = "Basic " + b64Data

	// url.Values {"grant_type": "client_credentials"} is the Form which is encoded
	callReq, err := http.NewRequest("POST", urlEndpoint, strings.NewReader(url.Values{"grant_type": {"client_credentials"}}.Encode())) // Make the request

	callReq.Header.Add("Authorization", authorization) // Add authorization string
	callReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{}
	response, err := httpClient.Do(callReq) // do request

	if err != nil {
		fmt.Println("Error making call request\nTerminating program")
		return tokenData
	}

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error reading response body\nTerminate program")
		return tokenData
	}

	// Unmarshal / parse data into object
	json.Unmarshal(data, &tokenData)

	response.Body.Close() // Close response body
	return tokenData
}

func main() {
	fmt.Println("Spotify public playlist to JSON")
	jsonFile, err := os.Open("api-config.json")

	if err != nil { // An error exists
		jsonOpenError()
		return
	}

	var apiConfig APIConfig                     // Instantiate as Type APIConfig
	jsonByteData, _ := ioutil.ReadAll(jsonFile) // Read the JSON file into ByeData
	json.Unmarshal(jsonByteData, &apiConfig)    // Unmarshal data

	var clientID string = apiConfig.SpotifyClientID //
	var clientSecret string = apiConfig.SpotifyClientSecret
	if clientID == "" {
		fmt.Println("No Client ID read. Check object name\nTerminating process")
		return
	}
	if clientSecret == "" {
		fmt.Println("No Client Secret read. Check object name\nTerminate process")
		return
	}

	var tokenData TokenConfig
	tokenData = getAccessToken(clientID, clientSecret)
	accessToken := tokenData.AccessToken

	if accessToken == "" {
		return
	}

	var userInput string
	fmt.Println("Enter spotify playlist id: ")
	fmt.Scanln(&userInput)

	// Similar to C printf
	// %s denotes string of characters
	playlistURLEndpoint := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s?market=ES", userInput)

	callReq, err := http.NewRequest("GET", playlistURLEndpoint, nil)
	if err != nil {
		fmt.Println("There was an error creating a call request")
	}
	callReq.Header.Add("Authorization", "Bearer "+accessToken)
	callReq.Header.Add("Content-Type", "application/json")
	httpClient := &http.Client{}
	response, err := httpClient.Do(callReq)

	if err != nil {
		fmt.Println("There was an error making call request to get Playlist Data")
	}
	// fmt.Println(response)
	defer response.Body.Close() // Close after processing body
	responseData, err := ioutil.ReadAll(response.Body)
	// fmt.Println(string(responseData)) // Prints

	writeErr := ioutil.WriteFile("./PlaylistData.json", []byte(responseData), 0644) // 0644 permission: "readable by all the user groups, but writable by the user only"

	if writeErr != nil {
		fmt.Println("Error writing to file")
		return
	}

	fmt.Println("Program terminated with no errors")
	jsonFile.Close() // Close the file
}
