# YouTube Meme API

W.I.P

An API to get meme videos on YouTube.

Playlist URLs are stored in [this Google Sheet](https://docs.google.com/spreadsheets/d/1MuvC8JpJte1wzAS0m9qR0rr2-gxzL8aaX6lvlKeAqvs/edit?usp=sharing)
and fetched through the Sheets API

Additional video/playlist data is then retrieved from the YouTube Data API.

The endpoints (so far) are:

## GET

### Random Endpoints

- `/api/v1/random/video` - Gets a random video
- `/api/v1/random/playlist` - Gets a random playlist
- `/api/v1/random/channel` - Gets a random channel

### Aggregate Endpoints

- `/api/v1/playlist/all` - Gets all playlists
- `/api/v1/channel/all` - Gets all channels

