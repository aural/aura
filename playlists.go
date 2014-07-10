package aural

type Playlist struct {
	Tracks []Track
}

func (playlist Playlist) Play() {
	PlayState.Tracks = playlist.Tracks
}

func (playlist Playlist) Queue() {
	for _, track := range playlist.Tracks {
		track.Queue()
	}
}
