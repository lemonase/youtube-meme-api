package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

/*
 * Code Samples and Docs:
 * https://developers.google.com/sheets/api/quickstart/go
 * https://docs.google.com/spreadsheets/d/1MuvC8JpJte1wzAS0m9qR0rr2-gxzL8aaX6lvlKeAqvs/edit?usp=sharing
 */

// credentials and tokens for auth (can also use an API key)
var credFile string = "auth/sheets/credentials.json"
var tokenFile string = "auth/sheets/token.json"

// SheetID - The main sheet ID that we are working with
const SheetID = "1MuvC8JpJte1wzAS0m9qR0rr2-gxzL8aaX6lvlKeAqvs"

// PlaylistRange - The range of values in the sheet containing the playlist urls
var PlaylistRange = "Sheet1!A2:A1000"

// CountRange - The range of values correlated to the times a playlist is selected
// TODO get more than write access to the sheet to implement this
var CountRange = "Sheet1!B2:B10000"

// Client - The authroized youtube service client (either with a key or a token)
var Client = getSheetsClientAuth()

// PlaylistLength, PlaylistValues - The values returned from the API call to sheets
var PlaylistLength, PlaylistValues = GetSheetValues(SheetID, PlaylistRange)

// GetSheetValues - wrapper to SheetsAPI
// Params - takes a sheetID for a spreadsheet, and a range of values to get
// Returns - the length of the values and the values themselves
func GetSheetValues(sheetID string, playlistRange string) (int, [][]interface{}) {
	// TODO fetch data concurrently
	log.Printf("Started fetching playlists from Google Sheet")
	resp, err := Client.Spreadsheets.Values.Get(SheetID, playlistRange).Do()
	if err != nil {
		log.Fatalf("Error fetching sheet%v\n", err)
	} else if len(resp.Values) < 1 {
		log.Printf("No values on sheet %s in range %s!", sheetID, playlistRange)
	}
	log.Printf("Finished fetching playlists from Google Sheet")
	return len(resp.Values), resp.Values
}

// GetRandomPlaylist - returns a random value from PlaylistValues array
func GetRandomPlaylist() string {
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(PlaylistLength)
	return fmt.Sprintf("%s", PlaylistValues[randIndex][0])
}

// Client Functions

func getSheetsClientUnauth() *sheets.Service {
	httpClient := &http.Client{}
	sheetsClient, err := sheets.New(httpClient)
	if err != nil {
		log.Fatalf("Could not get sheets client%v\n", err)
	}

	return sheetsClient
}

// Reads API creds from a json file
func getSheetsClientAuth() *sheets.Service {
	b, err := ioutil.ReadFile(credFile)
	if err != nil {
		log.Fatalf("Could not read credentials.json %v", err)
	}
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")

	httpClient := getClient(config)
	sheetsClient, err := sheets.New(httpClient)
	if err != nil {
		log.Fatalf("Could not get sheets client%v\n", err)
	}

	return sheetsClient
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

/*
 * OAuth Functions
 */

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
