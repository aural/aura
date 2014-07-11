package aural

type Playlist struct {
	Tracks []Track
}

func (playlist Playlist) Play(playState *PlayState) {
	playState.Tracks = playlist.Tracks
}

func (playlist Playlist) Queue(playState *PlayState) {
	for _, track := range playlist.Tracks {
		track.Queue(playState)
	}
}
