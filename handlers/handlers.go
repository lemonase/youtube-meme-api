package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/lemonase/youtube-meme-api/wrappers/sheets"
	"github.com/lemonase/youtube-meme-api/wrappers/youtube"
)

// TemplateData - The data the goes into the served html page
type TemplateData struct {
	SiteTitle string `json:"siteTitle"`
	Title     string `json:"title"`
	VideoID   string `json:"videoID"`
}

// Home - Displays the home page
func Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	tmpl := template.Must(template.ParseFiles("html/index.html"))
	data := &TemplateData{
		SiteTitle: "🔀 YouTube Meme Shuffle 🔀",
		Title:     "🔀 YouTube Meme Shuffle 🔀",
		VideoID:   youtube.GetRandomPlaylistItem().ContentDetails.VideoId,
	}

	tmpl.Execute(w, data)
}

// APIHelper - Prints a helpful error message
func APIHelper(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404: URL Not Found")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Endpoints are: ")
	fmt.Fprintln(w, "	Random:")
	fmt.Fprintln(w, "GET	/api/v1/random/video")
	fmt.Fprintln(w, "GET	/api/v1/random/playlist")
	fmt.Fprintln(w, "GET	/api/v1/random/playlist/item")
	fmt.Fprintln(w, "GET	/api/v1/random/channel")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "	All:")
	fmt.Fprintln(w, "GET   	/api/v1/all/video")
	fmt.Fprintln(w, "GET   	/api/v1/all/playlist")
	fmt.Fprintln(w, "GET   	/api/v1/all/playlist/item")
	fmt.Fprintln(w, "GET   	/api/v1/all/channel")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "	Update:")
	fmt.Fprintln(w, "GET   	/api/v1/update/all")
	fmt.Fprintln(w, "GET   	/api/v1/update/video")
	fmt.Fprintln(w, "GET   	/api/v1/update/playlist")
	fmt.Fprintln(w, "GET   	/api/v1/update/channel")
}

// Videos

// AllVideos - Get all singular videos responses
func AllVideos(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(youtube.VideoResponses, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// RandomVideo - Get a random playlist item from a random playlist
func RandomVideo(w http.ResponseWriter, r *http.Request) {
	item := youtube.GetRandomVideo()

	j, err := json.MarshalIndent(item, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// Playlists

// AllPlaylists - Get all playlist responses
func AllPlaylists(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(youtube.PlaylistResponses, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// AllPlaylistsWithItems - Get all playlist responses
func AllPlaylistsWithItems(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(youtube.PlaylistItemResponses, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// RandomPlaylist - Get a random playlist response
func RandomPlaylist(w http.ResponseWriter, r *http.Request) {
	randomPlaylist := youtube.GetRandomPlaylist()
	j, err := json.MarshalIndent(randomPlaylist, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Could not unmarshal %v", err)
	}
	fmt.Fprintf(w, string(j))
}

// RandomPlaylistItem - Get a random playlist response
func RandomPlaylistItem(w http.ResponseWriter, r *http.Request) {
	item := youtube.GetRandomPlaylistItem()
	j, err := json.MarshalIndent(item, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Could not unmarshal %v", err)
	}
	fmt.Fprintf(w, string(j))
}

// Channels

// AllChannels - Get all youtube channel responses
func AllChannels(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(youtube.ChannelResponses, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

// RandomChannel - Get a random channel from youtube responses
func RandomChannel(w http.ResponseWriter, r *http.Request) {
	randomChannel := youtube.GetRandomChannel()
	j, err := json.MarshalIndent(randomChannel, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Could not unmarshal %v", err)
	}
	fmt.Fprintf(w, string(j))
}

// FetchAllYoutubeInfoFromSheet - Gets sheet values, resets responses and fetches youtube data
func FetchAllYoutubeInfoFromSheet(forceRefresh bool) {
	sheets.FetchAllValues()

	youtube.ChannelResponses = nil
	youtube.PlaylistResponses = nil
	youtube.VideoResponses = nil

	youtube.FetchOrReadAll(forceRefresh)
}

// UpdateAllValuesFromSheet - Updates json files by enforcing refresh
func UpdateAllValuesFromSheet(w http.ResponseWriter, r *http.Request) {
	FetchAllYoutubeInfoFromSheet(true)
}

// UpdateAllChannelsFromSheet - Refetches channel responses and forces refresh
func UpdateAllChannelsFromSheet(w http.ResponseWriter, r *http.Request) {
	oldLen := sheets.ChannelLength

	sheets.FetchChannelValues()
	if oldLen != sheets.ChannelLength {
		youtube.ChannelResponses = nil
		youtube.FetchOrRead("channel", true)
	}
}

// UpdateAllPlaylistsFromSheet - Refetches playlist responses and forces refresh
func UpdateAllPlaylistsFromSheet(w http.ResponseWriter, r *http.Request) {
	oldLen := sheets.PlaylistLength

	sheets.FetchPlaylistValues()
	if oldLen != sheets.PlaylistLength {
		youtube.PlaylistResponses = nil
		youtube.PlaylistItemResponses = nil
		youtube.FetchOrRead("playlist", true)
		youtube.FetchOrRead("playlistItem", true)
	}
}

// UpdateAllVideosFromSheet - Refetches video responses and forces refresh
func UpdateAllVideosFromSheet(w http.ResponseWriter, r *http.Request) {
	oldLen := sheets.VideoLength

	sheets.FetchVideoValues()
	if oldLen != sheets.VideoLength {
		youtube.VideoResponses = nil
		youtube.FetchOrRead("video", true)
	}
}
