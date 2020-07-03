package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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

// GetPlaylistIDFromURL - takes a URL string and gets everything to the right
// of "list=", which should be the playlistID
func GetPlaylistIDFromURL(url string) string {
	return url[strings.LastIndex(url, "list=")+5:]
}

// GetPlaylistRepsonseFromURL - takes a URL string and returns an playlist response
func GetPlaylistRepsonseFromURL(url string) *youtube.PlaylistListResponse {
	return GetPlaylistResponse(GetPlaylistIDFromURL(url))
}

// GetPlaylistItemsResponseFromURL - takes a URL string and returns playlist items response
func GetPlaylistItemsResponseFromURL(url string) *youtube.PlaylistItemListResponse {
	return GetPlaylistItemsResponse(GetPlaylistIDFromURL(url), 0)
}

// GetPlaylistResponse - takes a playlist id and executes API call to playlists service
// https://developers.google.com/youtube/v3/docs/playlists/list#response
func GetPlaylistResponse(id string) *youtube.PlaylistListResponse {
	// part is a required parameter that describes the resource properties in the response
	// https://developers.google.com/youtube/v3/docs/playlists/list#parameters
	part := []string{"snippet,contentDetails"}

	// TODO make api call concurrently
	Call := Client.Playlists.List(part)
	Call = Call.Id(id)

	res, err := Call.Do()
	if err != nil {
		log.Fatalf("Error fetching playlist %v\n", err)
	}

	return res
}

// GetPlaylistItemsResponse - takes an id and position of a video in a playlist
// returns a response with a list of items (videos) and page info
// https://developers.google.com/youtube/v3/docs/playlistItems/list
func GetPlaylistItemsResponse(id string, videoPosition int64) *youtube.PlaylistItemListResponse {
	var correctPageRes *youtube.PlaylistItemListResponse

	part := []string{"snippet,contentDetails"}
	Call := Client.PlaylistItems.List(part)

	Call = Call.PlaylistId(id)
	Call = Call.MaxResults(PageSize)

	// pagination occurs in the API with tokens, so we iterate through
	// pages until the index of the requested video is within the page
	// TODO make this faster with goroutines
	for pageIndex := int64(0); pageIndex <= videoPosition; pageIndex += PageSize {
		res, err := Call.Do()
		if err != nil {
			log.Fatalf("Error fetching playlist %v\n", err)
		}
		correctPageRes = res
		Call.PageToken(res.NextPageToken)
	}

	return correctPageRes
}

// ChannelsListByUsername - example function
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
