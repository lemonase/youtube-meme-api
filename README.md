# YouTube Meme API

W.I.P

An API to get meme videos on YouTube.

Playlist URLs are stored in [this Google Sheet](https://docs.google.com/spreadsheets/d/1MuvC8JpJte1wzAS0m9qR0rr2-gxzL8aaX6lvlKeAqvs/edit?usp=sharing)
and fetched through the Sheets API

Additional video/playlist data is then retrieved from the YouTube Data API.

The endpoints (so far) are:

## GET

- /playlists/all - Get all playlists
- /playlists/random - Get a random playlist
- /video/random - Get a random video from a random playlist
- /video/all - Get all videos from all playlists
- /update - Refresh playlist data from the Google Sheet
