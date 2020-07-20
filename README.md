# YouTube Meme API

W.I.P

An API to get meme videos on YouTube.

Playlist URLs are stored in [this Google Sheet](https://docs.google.com/spreadsheets/d/1MuvC8JpJte1wzAS0m9qR0rr2-gxzL8aaX6lvlKeAqvs/edit?usp=sharing)
and fetched through the Sheets API

Additional video/playlist data is then retrieved from the YouTube Data API.

The endpoints (so far) are:

## HTTP Endpoints

## GET

### Random Endpoints

- `/api/v1/random/video` - Gets a random video
- `/api/v1/random/playlist` - Gets a random playlist
- `/api/v1/random/channel` - Gets a random channel

### Aggregate Endpoints

- `/api/v1/all/videos` - Gets all videos
- `/api/v1/all/playlists` - Gets all playlists
- `/api/v1/all/channels` - Gets all channels

## Client usage

### Video

```shell
$ xdg-open "https://www.youtube.com/watch?v=$(curl -sS localhost:8000/api/v1/random/video | jq .contentDetails.videoId | sed 's/"//g')"
```

### Playlist

```shell
xdg-open "https://www.youtube.com/playlist?list=$(curl -sS localhost:8000/api/v1/random/playlist | jq .items[0].id | sed 's/"//g')"
```

### Channel

```shell
$ xdg-open "https://www.youtube.com/channel/$(curl -sS localhost:8000/api/v1/random/channel | jq .items[0].id | sed 's/"//g')"
```
