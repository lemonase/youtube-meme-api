package server

import (
	"log"
	"net/http"

	"github.com/lemonase/youtube-meme-api/handlers"
)

// InitServer - Sets all routes and initializes the server
func InitServer(port string) {

	mux := http.NewServeMux()

	// serves webpage
	mux.HandleFunc("/", handlers.Home)

	// api
	mux.HandleFunc("/api/", handlers.APIHelper)

	// random
	mux.HandleFunc("/api/v1/random/video", handlers.RandomVideo)
	mux.HandleFunc("/api/v1/random/playlist", handlers.RandomPlaylist)
	mux.HandleFunc("/api/v1/random/playlist/item", handlers.RandomPlaylistItem)
	mux.HandleFunc("/api/v1/random/channel", handlers.RandomChannel)

	// all
	mux.HandleFunc("/api/v1/all/video", handlers.AllVideos)
	mux.HandleFunc("/api/v1/all/playlist", handlers.AllPlaylists)
	mux.HandleFunc("/api/v1/all/playlist/item", handlers.AllPlaylistsWithItems)
	mux.HandleFunc("/api/v1/all/channel", handlers.AllChannels)

	// updates
	mux.HandleFunc("/api/v1/update/all", handlers.UpdateAllValuesFromSheet)
	mux.HandleFunc("/api/v1/update/video", handlers.UpdateAllVideosFromSheet)
	mux.HandleFunc("/api/v1/update/playlist", handlers.UpdateAllPlaylistsFromSheet)
	mux.HandleFunc("/api/v1/update/channel", handlers.UpdateAllChannelsFromSheet)

	server := http.Server{Addr: port, Handler: mux}
	log.Printf("Server listenting on *%s", port)
	log.Fatal(server.ListenAndServe())
}

// FetchInitResources - Calls sheets and youtube APIs for data
func FetchInitResources() {
	handlers.FetchAllValuesFromSheet()
}
