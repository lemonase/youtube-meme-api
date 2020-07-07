package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/lemonase/youtube-meme-api/wrappers/sheets"
	"github.com/lemonase/youtube-meme-api/wrappers/youtube"
)

// All

// HandleAll -
func HandleAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "Hello\nHandler")
}

// Videos

// HandleVideosAll -
func HandleVideosAll(w http.ResponseWriter, r *http.Request) {
	videos := youtube.VideoResponses

	j, err := json.MarshalIndent(videos, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// HandleVideosRandom -
func HandleVideosRandom(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())

	item := youtube.GetRandomPlaylistItem()
	// url = "https://www.youtube.com/watch?v=" + item.ContentDetails.VideoId
	// date := item.ContentDetails.VideoPublishedAt

	j, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// Playlists

// HandlePlaylistsAll -
func HandlePlaylistsAll(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(youtube.PlaylistResponses, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// HandlePlaylistsRandom -
func HandlePlaylistsRandom(w http.ResponseWriter, r *http.Request) {
	randomPlaylist := youtube.GetRandomPlaylist()
	j, err := json.MarshalIndent(randomPlaylist, "", "  ")
	if err != nil {
		log.Fatalf("Could not unmarshal %v", err)
	}
	fmt.Fprintf(w, string(j))
}

// Channels

// HandleChannelsAll -
func HandleChannelsAll(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(youtube.ChannelResponses, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// HandleChannelsRandom -
func HandleChannelsRandom(w http.ResponseWriter, r *http.Request) {
	randomChannel := youtube.GetRandomChannel()
	j, err := json.MarshalIndent(randomChannel, "", "  ")
	if err != nil {
		log.Fatalf("Could not unmarshal %v", err)
	}
	fmt.Fprintf(w, string(j))
}

// Updates
// TODO set up a function on the sheet that calls this endpoint when the sheet changes

// UpdateAllValuesFromSheet -
func UpdateAllValuesFromSheet(w http.ResponseWriter, r *http.Request) {
	sheets.FetchAllValues()

	youtube.ChannelResponses = nil
	youtube.PlaylistResponses = nil
	youtube.VideoResponses = nil
	youtube.FetchAllListsFromSheet()
}

// UpdateAllChannelsFromSheet -
func UpdateAllChannelsFromSheet(w http.ResponseWriter, r *http.Request) {
	sheets.FetchChannelValues()
	youtube.ChannelResponses = nil
	youtube.FetchAllChannels()
}

// UpdateAllPlaylistsFromSheet -
func UpdateAllPlaylistsFromSheet(w http.ResponseWriter, r *http.Request) {
	sheets.FetchPlaylistValues()
	youtube.PlaylistResponses = nil
	youtube.FetchAllPlaylists()
}

// UpdateAllVideosFromSheet -
func UpdateAllVideosFromSheet(w http.ResponseWriter, r *http.Request) {
	sheets.FetchVideoValues()
	youtube.VideoResponses = nil
	youtube.FetchAllVideos()
}
