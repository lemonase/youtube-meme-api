package sheets

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/lemonase/youtube-meme-api/client"
)

/*
 * Code Samples and Docs:
 * https://developers.google.com/sheets/api/quickstart/go
 * https://docs.google.com/spreadsheets/d/1MuvC8JpJte1wzAS0m9qR0rr2-gxzL8aaX6lvlKeAqvs/edit?usp=sharing
 */

// SheetID - The main sheet ID that we are working with
const SheetID = "1MuvC8JpJte1wzAS0m9qR0rr2-gxzL8aaX6lvlKeAqvs"

// Client - The authroized youtube service client (either with a key or a token)
var Client = &client.Services.Sheets

// VideoRange -
var VideoRange = "Sheet1!A2:A1000"

// PlaylistRange -
var PlaylistRange = "Sheet1!C2:C1000"

// ChannelRange -
var ChannelRange = "Sheet1!E2:E1000"

// ChannelValues -
var ChannelValues [][]interface{}

// PlaylistValues -
var PlaylistValues [][]interface{}

// VideoValues -
var VideoValues [][]interface{}

// ChannelLength -
var ChannelLength int

// PlaylistLength -
var PlaylistLength int

// VideoLength -
var VideoLength int

/*
 * Fetch Functions
 */

// FetchAllValues -
func FetchAllValues() {
	log.Printf(":: Fetching Values From Google Sheet ::\n")
	log.Printf("https://docs.google.com/spreadsheets/d/%s\n", SheetID)

	FetchChannelValues()
	FetchPlaylistValues()
	FetchVideoValues()
}

// FetchSheetValues - wrapper to SheetsAPI
// Params - takes a sheetID for a spreadsheet, and a range of values to get
// Returns - the length of the values and the values themselves
func FetchSheetValues(sheetID string, playlistRange string) (int, [][]interface{}) {
	resp, err := Client.Spreadsheets.Values.Get(SheetID, playlistRange).Do()
	if err != nil {
		log.Fatalf("Error fetching sheet%v\n", err)
	} else if len(resp.Values) < 1 {
		log.Fatalf("No values on sheet %s in range %s!", sheetID, playlistRange)
	}
	return len(resp.Values), resp.Values
}

// FetchChannelValues -
func FetchChannelValues() {
	ChannelLength, ChannelValues = FetchSheetValues(SheetID, ChannelRange)
	log.Printf("	Fetching Channels From Range %s\n", ChannelRange)
	log.Printf("		Number of Channel URLs: %d\n", ChannelLength)
}

// FetchPlaylistValues -
func FetchPlaylistValues() {
	PlaylistLength, PlaylistValues = FetchSheetValues(SheetID, PlaylistRange)
	log.Printf("	Fetching Playlists From Range %s\n", PlaylistRange)
	log.Printf("		Number of Playlist URLs: %d\n", PlaylistLength)
}

// FetchVideoValues -
func FetchVideoValues() {
	VideoLength, VideoValues = FetchSheetValues(SheetID, VideoRange)
	log.Printf("	Fetching Videos From Range %s\n", VideoRange)
	log.Printf("		Number of Video URLs: %d\n", VideoLength)
}

/*
 * Random Functions
 */

// GetRandomVideo -
func GetRandomVideo() string {
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(VideoLength)
	return fmt.Sprintf("%s", VideoValues[randIndex][0])
}

// GetRandomPlaylist - returns a random value from PlaylistValues array
func GetRandomPlaylist() string {
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(PlaylistLength)
	return fmt.Sprintf("%s", PlaylistValues[randIndex][0])
}

// GetRandomChannel -
func GetRandomChannel() string {
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(ChannelLength)
	return fmt.Sprintf("%s", ChannelValues[randIndex][0])
}
