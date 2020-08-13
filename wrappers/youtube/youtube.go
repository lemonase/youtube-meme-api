package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lemonase/youtube-meme-api/client"
	"github.com/lemonase/youtube-meme-api/wrappers/sheets"
	"google.golang.org/api/youtube/v3"
)

/*
 * Code Samples and Docs:
 * https://developers.google.com/youtube/v3/docs/playlists
 * https://developers.google.com/youtube/v3/code_samples/go#search_by_keyword
 */

// Client - youtube client for auth and API methods
var Client = &client.Services.YouTube

// PageSize - the number of items that will be returned in a single API call
var PageSize int64 = 50

// TODO write responses to one or more JSON files
// instead of storing in memory.

var dataBaseDir = "data"

// VideoResponses - holds responses from videos
var VideoResponses []*youtube.VideoListResponse
var videoJSONFile = filepath.Join(dataBaseDir, "video.json")

// PlaylistResponses - holds responses from playlists
var PlaylistResponses []*youtube.PlaylistListResponse
var playlistJSONFile = filepath.Join(dataBaseDir, "playlist.json")

// PlaylistItemResponses - holds responses for items of a playlist
var PlaylistItemResponses []*youtube.PlaylistItemListResponse
var playlistItemJSONFile = filepath.Join(dataBaseDir, "playlist_item.json")
var playlistItemCount int

var PlaylistItem *youtube.PlaylistItem

// ChannelResponses - holds responses from channels
var ChannelResponses []*youtube.ChannelListResponse
var channelJSONFile = filepath.Join(dataBaseDir, "channel.json")

// SearchResponses - holds response from a search call
var SearchResponses []*youtube.SearchListResponse
var searchJSONFile = filepath.Join(dataBaseDir, "search.json")

// Files

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createDataDir() {
	_, statErr := os.Stat(dataBaseDir)
	if statErr != nil {
		fErr := os.Mkdir(dataBaseDir, 0644)
		if fErr != nil {
			log.Fatal(fErr)
		}
	}
}

func FetchOrRead(pageType string, forceRefresh bool) {
	if pageType == "channel" {
		if fileExists(channelJSONFile) && !forceRefresh {
			log.Printf("	Fetching Channel Info From %s", channelJSONFile)
			data, err := ioutil.ReadFile(channelJSONFile)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(data, &ChannelResponses)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("	Fetching Channel Info From API\n")
			FetchAllChannels()
			j, err := json.Marshal(ChannelResponses)
			if err != nil {
				log.Fatalf("Error marshalling json")
			}
			err = ioutil.WriteFile(channelJSONFile, j, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("		Number of Channel Playlists: %d\n", len(ChannelResponses))

	} else if pageType == "playlist" {
		if fileExists(playlistJSONFile) && !forceRefresh {
			log.Printf("	Fetching Playlist Info From %s\n", playlistJSONFile)
			data, err := ioutil.ReadFile(playlistJSONFile)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(data, &PlaylistResponses)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("	Fetching Playlist Info From API\n")
			FetchAllPlaylists()
			j, err := json.Marshal(PlaylistResponses)
			if err != nil {
				log.Fatalf("Error marshalling json")
			}
			err = ioutil.WriteFile(playlistJSONFile, j, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("		Number of Playlists: %d\n", len(PlaylistResponses))

	} else if pageType == "playlistItem" {
		if fileExists(playlistItemJSONFile) && !forceRefresh {
			log.Printf("	Fetching Playlist Items From %s\n", playlistItemJSONFile)
			data, err := ioutil.ReadFile(playlistItemJSONFile)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(data, &PlaylistItemResponses)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("	Fetching Playlist Items From API\n")
			FetchAllPlaylistItems()
			j, err := json.Marshal(PlaylistItemResponses)
			if err != nil {
				log.Fatalf("Error marshalling json")
			}
			err = ioutil.WriteFile(playlistItemJSONFile, j, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("		Number of Playlist Pages: %d\n", len(PlaylistItemResponses))
		for _, pl := range PlaylistItemResponses {
			playlistItemCount += len(pl.Items)
		}
		log.Printf("		Number of Playlist Items: %d\n", playlistItemCount)

	} else if pageType == "video" {
		if fileExists(videoJSONFile) && !forceRefresh {
			log.Printf("	Fetching Videos From %s\n", videoJSONFile)
			data, err := ioutil.ReadFile(videoJSONFile)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(data, &VideoResponses)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("	Fetching Videos From API\n")
			FetchAllVideos()
			j, err := json.Marshal(VideoResponses)
			if err != nil {
				log.Fatalf("Error marshalling json")
			}
			err = ioutil.WriteFile(videoJSONFile, j, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("		Number of Videos: %d\n", len(VideoResponses))

	} else {
		log.Fatalf("Unknown page type to fetch\n")
	}
}

func FetchOrReadAll(forceRefresh bool) {
	FetchOrRead("channel", forceRefresh)
	FetchOrRead("playlist", forceRefresh)
	FetchOrRead("playlistItem", forceRefresh)
	FetchOrRead("video", forceRefresh)
}

// Fetching

// FetchAllListsFromSheet - Fetches data for all the values in the sheet ranges
func FetchAllListsFromSheet() {
	createDataDir()
	log.Println(":: Fetching YouTube Data ::")
	FetchOrReadAll(false)
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
}

// FetchAllPlaylists - Fetches youtube data for all playlist values on the sheet
func FetchAllPlaylists() {
	for _, url := range sheets.PlaylistValues {
		playlistURL := string(url[0].(string))
		PlaylistResponses = append(PlaylistResponses, GetPlaylistRepsonseFromURL(playlistURL))
	}
}

// FetchAllPlaylistItems - Fetch all playlist items from playlist responses
func FetchAllPlaylistItems() {
	for _, pl := range PlaylistResponses {
		PlaylistItemResponses = append(PlaylistItemResponses, GetAllPlaylistItemResponsesFromPlaylistID(pl.Items[0].Id)...)
		playlistItemCount += len(pl.Items)
	}
}

// FetchAllVideos - Fetches youtube data for all the videos on the sheet
func FetchAllVideos() {
	for _, url := range sheets.VideoValues {
		videoURL := string(url[0].(string))
		VideoResponses = append(VideoResponses, GetVideoResponseFromURL(videoURL))
	}
}

// Randomizers

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
	// randomPlaylist := GetRandomPlaylist()
	// id := randomPlaylist.Items[0].Id
	// videoCount := randomPlaylist.Items[0].ContentDetails.ItemCount
	rand.Seed(time.Now().UnixNano())
	randIndex := int64(rand.Intn(int(len(PlaylistItemResponses))))

	randPl := PlaylistItemResponses[randIndex]
	item := randPl.Items[rand.Intn(len(randPl.Items))]

	return item
}

// GetRandomChannel - Returns a random playlist channel response
func GetRandomChannel() *youtube.ChannelListResponse {
	rand.Seed(time.Now().UnixNano())
	randIndex := int64(rand.Intn(len(ChannelResponses)))
	return ChannelResponses[randIndex]
}

// Video Utils

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

// Playlist Utils

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

// GetAllVideoItemsFromPlaylistID - Retruns a list of videos from playlist
func GetAllVideoItemsFromPlaylistID(id string) []*youtube.VideoListResponse {
	var playlistVideos []*youtube.VideoListResponse

	part := []string{"contentDetails"}
	Call := Client.PlaylistItems.List(part)

	Call = Call.PlaylistId(id)
	Call = Call.MaxResults(PageSize)

	res, err := Call.Do()
	if err != nil {
		log.Fatalf("Error fetching playlist %v\n", err)
	}
	if len(res.Items) < 1 {
		log.Fatalf("No items in playlist %s\n", id)
	}

	for pageIndex := int64(0); pageIndex <= int64(len(res.Items)); pageIndex += PageSize {
		res, err := Call.Do()
		if err != nil {
			log.Fatalf("Error fetching playlist %v\n", err)
		}
		for _, item := range res.Items {
			playlistVideos = append(playlistVideos, GetVideoResponseFromID(item.ContentDetails.VideoId))
		}

		Call.PageToken(res.NextPageToken)
	}

	return playlistVideos
}

// GetAllPlaylistItemResponsesFromPlaylistID - Appends all playlist item responses from an ID to the main list
func GetAllPlaylistItemResponsesFromPlaylistID(id string) []*youtube.PlaylistItemListResponse {
	var playlistItemResponses []*youtube.PlaylistItemListResponse

	part := []string{"contentDetails"}
	Call := Client.PlaylistItems.List(part)

	Call = Call.PlaylistId(id)
	Call = Call.MaxResults(PageSize)

	res, err := Call.Do()
	if err != nil {
		log.Fatalf("Error fetching playlist %v\n", err)
	}
	if len(res.Items) < 1 {
		log.Fatalf("No items in playlist %s\n", id)
	}

	for pageIndex := int64(0); pageIndex <= int64(len(res.Items)); pageIndex += PageSize {
		res, err := Call.Do()
		if err != nil {
			log.Fatalf("Error fetching playlist %v\n", err)
		}
		playlistItemResponses = append(playlistItemResponses, res)
	}

	return playlistItemResponses
}

// GetPlaylistItemsResponseFromIDAtIndex - Takes an id and position of a video in a playlist and returns a response
func GetPlaylistItemsResponseFromIDAtIndex(id string, videoIndex int64) *youtube.PlaylistItemListResponse {
	var correctPageRes *youtube.PlaylistItemListResponse

	part := []string{"contentDetails"}
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

	if len(correctPageRes.Items) < 1 {
		log.Fatalf("No items returned in response: Check playlist https:/www.yotube.com/playlist?list=%v at index %v", id, videoIndex)
	}

	return correctPageRes
}

// GetPlaylistItemsResponseFromURLAtIndex - Takes a URL string and index, returns playlist items response
func GetPlaylistItemsResponseFromURLAtIndex(url string, videoIndex int64) *youtube.PlaylistItemListResponse {
	return GetPlaylistItemsResponseFromIDAtIndex(GetPlaylistIDFromURL(url), videoIndex)
}

// Channels

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
