package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lemonase/youtube-meme-api/wrappers/sheets"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

/*
 * Code Samples and Docs:
 * https://developers.google.com/youtube/v3/docs/playlists
 * https://developers.google.com/youtube/v3/code_samples/go#search_by_keyword
 */

// creds and token files for auth
var credentialsFile string = "auth/youtube/credentials.json"
var tokenFile string = "auth/youtube/token.json"

// Client - youtube client for auth and API methods
var Client = getYoutubeClient()

// PageSize - the number of items that will be returned in a single API call
var PageSize int64 = 50

// VideoResponses - holds responses from videos
var VideoResponses []*youtube.VideoListResponse

// PlaylistResponses - holds responses from playlists
var PlaylistResponses []*youtube.PlaylistListResponse

// PlaylistItemResponses - holds responses for items of a playlist
var PlaylistItemResponses []*youtube.PlaylistItemListResponse

// ChannelResponses - holds responses from channels
var ChannelResponses []*youtube.ChannelListResponse

/*
 * Fetching
 */

// FetchAllListsFromSheet - Fetches data for all the values in the sheet ranges
func FetchAllListsFromSheet() {
	log.Println("::Fetching Youtube Data::")
	FetchAllVideos()
	FetchAllPlaylists()
	FetchAllChannels()
}

// FetchAllChannels - Fetches youtube data for all channel values on the sheet
func FetchAllChannels() {
	for _, url := range sheets.ChannelValues {
		channelURL := string(url[0].(string))
		ChannelResponses = append(ChannelResponses, GetChannelResponseFromURL(channelURL))
	}
	for _, channelRes := range ChannelResponses {
		uploadPl := GetPlaylistResponseFromID(channelRes.Items[0].ContentDetails.RelatedPlaylists.Uploads)
		PlaylistResponses = append(PlaylistResponses, uploadPl)
	}
	log.Printf("Number of Channels: %d\n", len(ChannelResponses))
}

// FetchAllPlaylists - Fetches youtube data for all playlist values on the sheet
func FetchAllPlaylists() {
	for _, url := range sheets.PlaylistValues {
		playlistURL := string(url[0].(string))
		PlaylistResponses = append(PlaylistResponses, GetPlaylistRepsonseFromURL(playlistURL))
	}
	log.Printf("Number of Playlists: %d\n", len(PlaylistResponses))
}

// FetchAllVideos - Fetches youtube data for all the videos on the sheet
func FetchAllVideos() {
	for _, url := range sheets.VideoValues {
		videoURL := string(url[0].(string))
		VideoResponses = append(VideoResponses, GetVideoResponseFromURL(videoURL))
	}
	log.Printf("Number of Videos: %d\n", len(VideoResponses))
}

/*
 * Randomizers
 */

// GetRandomVideo - Returns a random video response
func GetRandomVideo() *youtube.VideoListResponse {
	rand.Seed(time.Now().UnixNano())
	randIndex := int64(rand.Intn(len(VideoResponses)))
	return VideoResponses[randIndex]
}

// GetRandomPlaylist - Returns a random playlist response
func GetRandomPlaylist() *youtube.PlaylistListResponse {
	rand.Seed(time.Now().UnixNano())
	randIndex := int64(rand.Intn(len(PlaylistResponses)))
	return PlaylistResponses[randIndex]
}

// GetRandomPlaylistItem - Returns a random playlist video response
func GetRandomPlaylistItem() *youtube.PlaylistItem {
	randomPlaylist := GetRandomPlaylist()
	id := randomPlaylist.Items[0].Id
	videoCount := randomPlaylist.Items[0].ContentDetails.ItemCount

	rand.Seed(time.Now().UnixNano())
	randIndex := int64(rand.Intn(int(videoCount)))

	res := GetPlaylistItemsResponseFromIDAtIndex(id, randIndex)
	return res.Items[randIndex%PageSize]
}

// GetRandomChannel - Returns a random playlist channel response
func GetRandomChannel() *youtube.ChannelListResponse {
	rand.Seed(time.Now().UnixNano())
	randIndex := int64(rand.Intn(len(ChannelResponses)))
	return ChannelResponses[randIndex]
}

/*
 * Video Utils
 */

// GetVideoIDFromURL - Get the video id from a given url
func GetVideoIDFromURL(url string) string {
	param := "v="
	if strings.Contains(url, param) {
		return url[strings.LastIndex(url, param)+len(param):]
	}
	log.Fatalf("Could not retrive Video ID from URL: %s", url)
	return ""
}

// GetVideoResponseFromID - Returns a video response from video ID
func GetVideoResponseFromID(id string) *youtube.VideoListResponse {
	part := []string{"snippet,contentDetails"}

	Call := Client.Videos.List(part)
	Call = Call.Id(id)

	res, err := Call.Do()
	if err != nil {
		log.Fatalf("Error fetching youtube video %v\n", err)
	}

	return res
}

// GetVideoResponseFromURL - Returns a video response from a video URL
func GetVideoResponseFromURL(url string) *youtube.VideoListResponse {
	return GetVideoResponseFromID(GetVideoIDFromURL(url))
}

/*
 * Playlist Utils
 */

// GetPlaylistIDFromURL - Takes a URL string and gets everything to the right of playlist param
func GetPlaylistIDFromURL(url string) string {
	possibleParams := []string{"list=", "p="}
	for _, p := range possibleParams {
		if strings.Contains(url, p) {
			return url[strings.LastIndex(url, p)+len(p):]
		}
	}
	log.Fatalf("Could not retrieve Playlist ID from URL: %s", url)
	return ""
}

// GetPlaylistResponseFromID - Takes a playlist id and executes API call to playlists service
func GetPlaylistResponseFromID(id string) *youtube.PlaylistListResponse {
	part := []string{"snippet,contentDetails"}

	Call := Client.Playlists.List(part)
	Call = Call.Id(id)

	res, err := Call.Do()
	if err != nil {
		log.Fatalf("Error fetching playlist %v\n", err)
	}
	if len(res.Items) < 1 {
		log.Fatalf("No items in playlist %s\n", id)
	}

	return res
}

// GetPlaylistRepsonseFromURL - Takes a URL string and returns an playlist response
func GetPlaylistRepsonseFromURL(url string) *youtube.PlaylistListResponse {
	return GetPlaylistResponseFromID(GetPlaylistIDFromURL(url))
}

// Playlist *Items* (different from playlist response)

// GetPlaylistItemsResponseFromIDAtIndex - Takes an id and position of a video in a playlist and returns a response
func GetPlaylistItemsResponseFromIDAtIndex(id string, videoIndex int64) *youtube.PlaylistItemListResponse {
	var correctPageRes *youtube.PlaylistItemListResponse

	part := []string{"snippet,contentDetails"}
	Call := Client.PlaylistItems.List(part)

	Call = Call.PlaylistId(id)
	Call = Call.MaxResults(PageSize)

	// pagination occurs in the API with tokens, so we iterate through
	// pages until the index of the requested video is within the page
	for pageIndex := int64(0); pageIndex <= videoIndex; pageIndex += PageSize {
		res, err := Call.Do()
		if err != nil {
			log.Fatalf("Error fetching playlist %v\n", err)
		}
		correctPageRes = res
		Call.PageToken(res.NextPageToken)
	}

	return correctPageRes
}

// GetPlaylistItemsResponseFromURLAtIndex - Takes a URL string and index, returns playlist items response
func GetPlaylistItemsResponseFromURLAtIndex(url string, videoIndex int64) *youtube.PlaylistItemListResponse {
	return GetPlaylistItemsResponseFromIDAtIndex(GetPlaylistIDFromURL(url), videoIndex)
}

/*
 * Channels
 */

// GetChannelIDFromURL - Takes a string, splits it on "/" and gets the last field
func GetChannelIDFromURL(url string) string {
	paramList := []string{"/channel/", "/c/", "/user/"}
	for _, p := range paramList {
		if strings.Contains(url, p) {
			slice := strings.Split(url, "/")
			return string(slice[len(slice)-1:][0])
		}
	}
	log.Fatalf("Could not retrieve Channel ID from URL: %s", url)
	return ""
}

// GetChannelResponseFromID - Returns a channel response given an ID
func GetChannelResponseFromID(id string) *youtube.ChannelListResponse {
	part := []string{"snippet,contentDetails"}

	Call := Client.Channels.List(part)
	Call.MaxResults(PageSize)
	Call.Id(id)

	res, err := Call.Do()
	if err != nil {
		log.Fatalf("Error fetching channel details %v\n", err)
	}
	if len(res.Items) < 1 {
		newCall := Client.Channels.List(part)
		newCall.ForUsername(id)
		newRes, err := newCall.Do()
		if err != nil {
			log.Fatalf("Error fetching channel details %v\n", err)
		}
		if len(newRes.Items) < 1 {
			log.Fatalf("No items in channel response for username: %s", id)
		}

		return newRes
	}

	return res
}

// GetChannelResponseFromURL - Returns a channel response from a URL
func GetChannelResponseFromURL(url string) *youtube.ChannelListResponse {
	return GetChannelResponseFromID(GetChannelIDFromURL(url))
}

// ChannelsListByUsername - example function from docs
func ChannelsListByUsername(username string) {
	call := Client.Channels.List(strings.Split("snippet,contentDetails", ","))
	call = call.ForUsername(username)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error calling API: %v\n", err.Error())
	}
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}

/*
 * Client Functions
 */

// getYoutubeClient - reads a credentials file, calls a config and creates and
// returns a youtube client
func getYoutubeClient() *youtube.Service {
	ctx := context.Background()

	b, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(ctx, config)
	youtubeClient, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error obtaining client: %v\n", err.Error())
	}

	return youtubeClient
}

// getClient - takes http context, and oauth config type and returns a http
// client with an auto-refreshing token
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(ctx, tok)
}

/*
 * Auth Functions
 */

// getTokenFromWeb - default auth flow to get token
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v\n", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
