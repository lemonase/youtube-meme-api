package youtube

import (
	"fmt"
	"log"
	"math/rand"
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

// VideoResponses - holds responses from videos
var VideoResponses []*youtube.VideoListResponse

// PlaylistResponses - holds responses from playlists
var PlaylistResponses []*youtube.PlaylistListResponse

// PlaylistItemResponses - holds responses for items of a playlist
var PlaylistItemResponses []*youtube.PlaylistItemListResponse

// ChannelResponses - holds responses from channels
var ChannelResponses []*youtube.ChannelListResponse

// SearchResponses - holds response from a search call
var SearchResponses []*youtube.SearchListResponse

// Fetching

// FetchAllListsFromSheet - Fetches data for all the values in the sheet ranges
func FetchAllListsFromSheet() {
	log.Println("::Fetching Youtube Data::")
	FetchAllChannels()
	FetchAllPlaylists()
	FetchAllVideos()
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
	log.Printf("	Number of Channel Playlists: %d\n", len(ChannelResponses))
}

// FetchAllPlaylists - Fetches youtube data for all playlist values on the sheet
func FetchAllPlaylists() {
	for _, url := range sheets.PlaylistValues {
		playlistURL := string(url[0].(string))
		PlaylistResponses = append(PlaylistResponses, GetPlaylistRepsonseFromURL(playlistURL))
	}
	log.Printf("	Number of Regular Playlists: %d\n", len(PlaylistResponses))
}

// FetchAllVideos - Fetches youtube data for all the videos on the sheet
func FetchAllVideos() {
	for _, url := range sheets.VideoValues {
		videoURL := string(url[0].(string))
		VideoResponses = append(VideoResponses, GetVideoResponseFromURL(videoURL))
	}

	// 	if len(PlaylistResponses) < 1 {
	// 		FetchAllPlaylists()
	// 	}
	// 	for _, pl := range PlaylistResponses {
	// 		plVideos := GetAllVideoItemsFromPlaylistId(pl.Items[0].Id)
	// 		VideoResponses = append(VideoResponses, plVideos...)
	// 	}

	log.Printf("	Number of Total Videos: %d\n", len(VideoResponses))
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

// Playlist *Items* (different from playlist response)

// TODO these playlist API calls can be very expensive (in quota terms)
// it would be better to just get video urls from PlaylistItemsList
// instead of doing 2x+ the amount of api calls

// GetAllVideoItemsFromPlaylistId - Retruns a list of videos from playlist
func GetAllVideoItemsFromPlaylistId(id string) []*youtube.VideoListResponse {
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
