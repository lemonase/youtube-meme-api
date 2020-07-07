package server

import (
	"log"
	"net/http"
	"os"

	"github.com/lemonase/youtube-meme-api/handlers"
	"github.com/lemonase/youtube-meme-api/wrappers/sheets"
	"github.com/lemonase/youtube-meme-api/wrappers/youtube"
)

// InitServer - Sets all routes and initializes the server
func InitServer() {
	// get port from environment
	var port string = "8000"
	var ok bool
	portEnv, ok := os.LookupEnv("PORT")
	if ok {
		port = portEnv
	}
	port = ":" + port

	// TODO use flags and add configuration options
	// to use an API key instead of default OAuth,
	// which looks for credentials.json, performs auth,
	// and writes an auto-refreshing token.json file.
	// While this is robust, a simple API key would
	// be much more lightweight

	mux := http.NewServeMux()

	// all
	mux.HandleFunc("/", handlers.CatchAll)

	// random
	mux.HandleFunc("/api/v1/random/video", handlers.RandomVideo)
	mux.HandleFunc("/api/v1/random/playlist", handlers.RandomPlaylist)
	mux.HandleFunc("/api/v1/random/channel", handlers.RandomChannel)

	// all
	mux.HandleFunc("/api/v1/all/videos", handlers.AllVideos)
	mux.HandleFunc("/api/v1/all/playlists", handlers.AllPlaylists)
	mux.HandleFunc("/api/v1/all/channels", handlers.AllChannels)

	// updates
	mux.HandleFunc("/api/v1/update/all", handlers.UpdateAllValuesFromSheet)
	mux.HandleFunc("/api/v1/update/videos", handlers.UpdateAllVideosFromSheet)
	mux.HandleFunc("/api/v1/update/playlists", handlers.UpdateAllPlaylistsFromSheet)
	mux.HandleFunc("/api/v1/update/channels", handlers.UpdateAllChannelsFromSheet)

	server := http.Server{Addr: port, Handler: mux}
	log.Printf("Server listenting on *%s", port)
	log.Fatal(server.ListenAndServe())
}

// FetchInitResources - Calls sheets and youtube APIs for data
func FetchInitResources() {
	log.Printf("::Fetching initial data::\n")
	sheets.FetchAllValues()
	youtube.FetchAllListsFromSheet()
}
