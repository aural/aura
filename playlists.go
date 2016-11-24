package aural

type Repeat int

const (
	REPEAT_ONE Repeat = iota
	REPEAT_ALL
)

type Playlist struct {
	tracks []*Track
}

func NewPlaylist(tracks []*Track) *Playlist {
	playlist := new(Playlist)
	playlist.tracks = tracks
	return playlist
}

func (playlist *Playlist) Queue(location string) {
	track := Track{Location: location}
	playlist.tracks = append(playlist.tracks, &track)
}

func (playlist *Playlist) Current() *Track {
	return playlist.tracks[0]
}

func (playlist *Playlist) Length() int {
	return len(playlist.tracks)
}

func (playlist *Playlist) Pop() *Track {
	poppedTrack := playlist.Current()
	playlist.tracks = playlist.tracks[1:]
	return poppedTrack
}
