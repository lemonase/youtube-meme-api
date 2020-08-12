package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lemonase/youtube-meme-api/wrappers/sheets"
	"github.com/lemonase/youtube-meme-api/wrappers/youtube"
)

// CatchAll -
func CatchAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404: URL Not Found")

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Endpoints are: ")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "/api/v1/random/video")
	fmt.Fprintln(w, "/api/v1/random/playlist")
	fmt.Fprintln(w, "/api/v1/random/channel")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "/api/v1/all/video")
	fmt.Fprintln(w, "/api/v1/all/playlist")
	fmt.Fprintln(w, "/api/v1/all/channel")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "/api/v1/update/all")
	fmt.Fprintln(w, "/api/v1/update/video")
	fmt.Fprintln(w, "/api/v1/update/playlist")
	fmt.Fprintln(w, "/api/v1/update/channel")
}

// Videos

// AllVideos - Get all singular videos responses
func AllVideos(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(youtube.VideoResponses, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// RandomVideo - Get a random playlist item from a random playlist
func RandomVideo(w http.ResponseWriter, r *http.Request) {
	item := youtube.GetRandomVideo()

	j, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// Playlists

// AllPlaylists - Get all playlist responses
func AllPlaylists(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(youtube.PlaylistResponses, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// AllPlaylistsWithItems - Get all playlist responses
func AllPlaylistsWithItems(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(youtube.PlaylistItemResponses, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// RandomPlaylist - Get a random playlist response
func RandomPlaylist(w http.ResponseWriter, r *http.Request) {
	randomPlaylist := youtube.GetRandomPlaylist()
	j, err := json.MarshalIndent(randomPlaylist, "", "  ")
	if err != nil {
		log.Fatalf("Could not unmarshal %v", err)
	}
	fmt.Fprintf(w, string(j))
}

// RandomPlaylistItem - Get a random playlist response
func RandomPlaylistItem(w http.ResponseWriter, r *http.Request) {
	item := youtube.GetRandomPlaylistItem()
	j, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		log.Fatalf("Could not unmarshal %v", err)
	}
	fmt.Fprintf(w, string(j))
}

// Channels

// AllChannels - Get all youtube channel responses
func AllChannels(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(youtube.ChannelResponses, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// RandomChannel - Get a random channel from youtube responses
func RandomChannel(w http.ResponseWriter, r *http.Request) {
	randomChannel := youtube.GetRandomChannel()
	j, err := json.MarshalIndent(randomChannel, "", "  ")
	if err != nil {
		log.Fatalf("Could not unmarshal %v", err)
	}
	fmt.Fprintf(w, string(j))
}

// Updates
// TODO set up a function on the sheet that calls these endpoint when the sheet changes

// FetchAllValuesFromSheet -
func FetchAllValuesFromSheet() {
	sheets.FetchAllValues()

	youtube.ChannelResponses = nil
	youtube.PlaylistResponses = nil
	youtube.VideoResponses = nil
	youtube.FetchAllListsFromSheet()
}

// UpdateAllValuesFromSheet -
func UpdateAllValuesFromSheet(w http.ResponseWriter, r *http.Request) {
	FetchAllValuesFromSheet()
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
