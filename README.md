# YouTube Meme API

Frontend grabs a random video from a list of playlists.

Everything is stored in [this sheet](https://docs.google.com/spreadsheets/d/1MuvC8JpJte1wzAS0m9qR0rr2-gxzL8aaX6lvlKeAqvs/)
and fetched from the Google Sheets API. Additional video/playlist info is retrieved from the YouTube Data API.

Recommendations can be added using [this form](https://docs.google.com/forms/d/1j62PxUnAuFop-o7z0C0PKfBOAYMMyjmom8u_7l2qUDY)

## HTTP Endpoints

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

### Bash

```shell
endpoint="http://localhost:8000/api/v1/random/playlist/item"
vid_id="$(curl -sSL $endpoint | jq .contentDetails.videoId | sed 's/"//g')"
vid_url="https://www.youtube.com/watch?v=$vid_id"

# on linux
xdg-open $vid_url

# on mac
open $vid_url
```

### PowerShell

```powershell
$endpoint = "http://localhost:8000/api/v1/random/playlist/item"
$vid_id = ((iwr "$endpoint").Content | ConvertFrom-Json).contentDetails.videoId
$vid_url = "https://www.youtube.com/watch?v=$vid_id"

start chrome $vid_url
```
