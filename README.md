# YouTube Meme API

The intention of this API is to fetch random meme videos/playlist/channels from a Google Sheet and
learn a little bit about Go while doing it.

The API uses the Go standard library to route HTTP paths to functions.

All "backend" data for this project is stored in [this sheet](https://docs.google.com/spreadsheets/d/1MuvC8JpJte1wzAS0m9qR0rr2-gxzL8aaX6lvlKeAqvs/)
for ease of editing and sharing.
Recommendations can be added using [this form](https://docs.google.com/forms/d/1j62PxUnAuFop-o7z0C0PKfBOAYMMyjmom8u_7l2qUDY)

The "frontend" is a single `index.html` template with embedded iframes for the YT video and form (this is served by Go's builtin http server).

## How it works

The data is fetched from the [Google Sheets API](https://developers.google.com/sheets/api/reference/rest) with the [Go client library](https://pkg.go.dev/google.golang.org/api/sheets/v4).
Additional video/playlist info is retrieved from the [YouTube Data API](https://developers.google.com/youtube/v3/docs) with the associated [Go library](https://developers.google.com/youtube/v3/quickstart/go).

## HTTP Endpoints

Hosted on Heroku:
<https://youtube-meme-api.herokuapp.com/>

## GET

- `/` - Home page

### API "Random" Endpoints

- `/api/v1/random/video` - Gets a random video
- `/api/v1/random/playlist` - Gets a random playlist
- `/api/v1/random/playlist/item` - Gets a random playlist item (playlist video)
- `/api/v1/random/channel` - Gets a random channel

### API "Aggregate" Endpoints

- `/api/v1/all/video` - Gets all videos
- `/api/v1/all/playlist` - Gets all playlists
- `/api/v1/all/playlist/item` - Gets all playlists items/videos
- `/api/v1/all/channel` - Gets all channels

## Client usage examples

### Web browser

Go to https://youtube-meme-api.herokuapp.com/api to see all endpoints.

### Bash

```shell
endpoint="https://youtube-meme-api.herokuapp.com/api/v1/random/playlist/item"
vid_id="$(curl -sSL $endpoint | jq .contentDetails.videoId | sed 's/"//g')"
vid_url="https://www.youtube.com/watch?v=$vid_id"

# on linux
xdg-open $vid_url

# on mac
open $vid_url
```

### PowerShell

```powershell
$endpoint = "https://youtube-meme-api.herokuapp.com/api/v1/random/playlist/item"
$vid_id = ((iwr "$endpoint").Content | ConvertFrom-Json).contentDetails.videoId
$vid_url = "https://www.youtube.com/watch?v=$vid_id"

start chrome $vid_url
```
