package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/api/youtube/v3"
)

// ServiceClients -
type ServiceClients struct {
	Sheets  sheets.Service
	YouTube youtube.Service
}

var Services ServiceClients

var tokenFile string

func InitClientsWithAPIKey(apiKey string) {
	sheetsCtx := context.Background()
	youtubeCtx := context.Background()

	sheetsClient, err := sheets.NewService(sheetsCtx, option.WithAPIKey(apiKey), option.WithScopes(sheets.SpreadsheetsReadonlyScope))
	if err != nil {
		log.Fatalf("Could not get sheets client %v\n", err)
	}
	youtubeClient, err := youtube.NewService(youtubeCtx, option.WithAPIKey(apiKey), option.WithScopes(youtube.YoutubeReadonlyScope))
	if err != nil {
		log.Fatalf("Could not get youtube client %v\n", err)
	}

	Services.Sheets = *sheetsClient
	Services.YouTube = *youtubeClient
}

func InitClientsWithSecretJSONFile(filename string) {
	dir, _ := filepath.Split(filename)
	tokenFile = filepath.Join(dir, "token.json")

	Services.Sheets = *getSheetsClientOAuth(filename)
	Services.YouTube = *getYoutubeClientOAuth(filename)
}

func getSheetsClientOAuth(credsFilename string) *sheets.Service {
	b, err := ioutil.ReadFile(credsFilename)
	if err != nil {
		log.Fatalf("Unable to read client secret file (%s): %v", credsFilename, err)
	}

	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClientWithToken(config)
	sheetsClient, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Could not get sheets client%v\n", err)
	}

	return sheetsClient
}

// getYoutubeClientOAuth - reads a credentials file, calls a config and creates and
// returns a youtube client
func getYoutubeClientOAuth(credsFilename string) *youtube.Service {
	b, err := ioutil.ReadFile(credsFilename)
	if err != nil {
		log.Fatalf("Unable to read client secret file (%s): %v", credsFilename, err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClientWithToken(config)
	youtubeClient, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error obtaining client: %v\n", err.Error())
	}

	return youtubeClient
}

// OAuth Boilerplate Below

// Retrieve a token, saves the token, then returns the generated client.
func getClientWithToken(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
