package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/lemonase/youtube-meme-api/wrappers/sheets"
	"github.com/lemonase/youtube-meme-api/wrappers/youtube"
)

// Local struct wrappers for Youtube API
type playlist struct {
	URL        string `json:"url"`
	ID         string `json:"id"`
	Title      string `json:"title"`
	VideoCount int64  `json:"videocount"`
	Date       string `json:"date"`
}

type video struct {
	URL       string `json:"url"`
	ID        string `json:"id"`
	Position  int64  `json:"pos"`
	Title     string `json:"title"`
	VideoDate string `json:"videodate"`
	Date      string `json:"date"`
}

var playlists []playlist
var videos []video

// Init - Sets all routes and initializes the server
func Init() {
	// get port from environment
	var port string = "8000"
	var ok bool
	portEnv, ok := os.LookupEnv("PORT")
	if ok {
		port = portEnv
	}
	port = ":" + port

	playlists = fetchAllPlaylistsFromSheet()

	// TODO add middleware for logging and checking/settings http headers for a JSON response
	http.HandleFunc("/playlist/all", getAllPlaylists)
	http.HandleFunc("/playlist/random", getRandomPlaylist)
	http.HandleFunc("/video/all", getAllVideos)
	http.HandleFunc("/video/random", getRandomVideo)
	http.HandleFunc("/update", updatePlaylistValues)

	log.Printf("Server listenting on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

/*
 * API Router Functions
 */

func updatePlaylistValues(w http.ResponseWriter, r *http.Request) {
	// TODO set up a function on the sheet that calls this endpoint when the sheet changes
	sheets.PlaylistLength, sheets.PlaylistValues = sheets.GetSheetValues(sheets.SheetID, sheets.PlaylistRange)
	playlists = nil
	playlists = fetchAllPlaylistsFromSheet()
}

func getAllPlaylists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	j, err := json.MarshalIndent(playlists, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal data %v", err)
	}

	fmt.Fprintf(w, string(j))
}

func getAllVideos(w http.ResponseWriter, r *http.Request) {
	videoList := fetchAllVideosFromSheet()
	j, err := json.MarshalIndent(videoList, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal %v", err)
	}

	fmt.Fprintf(w, string(j))
}

func getRandomPlaylist(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(fetchRandomPlaylist(), "", "  ")
	if err != nil {
		log.Fatalf("Could not unmarshal %v", err)
	}
	fmt.Fprintf(w, string(j))
}

func getRandomVideo(w http.ResponseWriter, r *http.Request) {
	j, err := json.MarshalIndent(fetchRandomVideo(), "", "  ")
	if err != nil {
		log.Fatalf("Could not unmarshal %v", err)
	}
	fmt.Fprintf(w, string(j))
}

/*
 * API Fetch Functions
 */

func fetchPlaylistDataFromURL(url string) playlist {
	// TODO add goroutine here or in the API wrapper
	id := youtube.GetPlaylistIDFromURL(url)
	resp := youtube.GetPlaylistResponse(id)
	if len(resp.Items) < 1 {
		log.Fatal("No Items In Playlist")
	}
	info := resp.Items[0]

	pl := playlist{
		URL:        url,
		ID:         id,
		Title:      info.Snippet.Title,
		VideoCount: info.ContentDetails.ItemCount,
		Date:       info.Snippet.PublishedAt}

	return pl
}

func fetchAllPlaylistsFromSheet() []playlist {
	for _, playlistURL := range sheets.PlaylistValues {
		url := string(playlistURL[0].(string))
		pl := fetchPlaylistDataFromURL(url)
		playlists = append(playlists, pl)
	}

	return playlists
}

func fetchVideoAtIndex(playlistID string, index int64) video {
	res := youtube.GetPlaylistItemsResponse(playlistID, index)
	item := res.Items[index%youtube.PageSize]

	vid := video{
		URL:       "https://www.youtube.com/watch?v=" + item.ContentDetails.VideoId,
		ID:        item.ContentDetails.VideoId,
		Position:  item.Snippet.Position,
		Title:     item.Snippet.Title,
		VideoDate: item.ContentDetails.VideoPublishedAt,
		Date:      item.Snippet.PublishedAt}

	return vid
}

func fetchAllVideosFromPlaylist(pl playlist) []video {
	var vids []video
	plResp := youtube.GetPlaylistItemsResponse(pl.ID, pl.VideoCount)
	plItems := plResp.Items

	for _, item := range plItems {
		vid := video{
			URL:       "https://www.youtube.com/watch?v=" + item.ContentDetails.VideoId,
			ID:        item.ContentDetails.VideoId,
			Position:  item.Snippet.Position,
			Title:     item.Snippet.Title,
			VideoDate: item.ContentDetails.VideoPublishedAt,
			Date:      item.Snippet.PublishedAt}
		vids = append(vids, vid)
	}
	return vids
}

func fetchAllVideosFromSheet() []video {
	var allVids []video
	playlists := fetchAllPlaylistsFromSheet()
	for _, pl := range playlists {
		plVids := fetchAllVideosFromPlaylist(pl)
		allVids = append(allVids, plVids...)
	}

	return allVids
}

func fetchRandomPlaylist() playlist {
	return fetchPlaylistDataFromURL(sheets.GetRandomPlaylist())
}

func fetchRandomVideo() video {
	pl := fetchRandomPlaylist()
	rand.Seed(time.Now().UnixNano())
	randIndex := int64(rand.Intn(int(pl.VideoCount)))
	return fetchVideoAtIndex(pl.ID, randIndex)
}
