package aural

import (
	"log"

	"github.com/gordonklaus/portaudio"
)

func init() {
	if err := portaudio.Initialize(); err != nil {
		log.Fatalln(err)
	}
}

func Terminate() {
	portaudio.Terminate()
}

type Playstate struct {
	current  *Track
	Playlist *Playlist
}

func NewPlaystate() *Playstate {
	p := new(Playstate)
	p.Playlist = NewPlaylist([]*Track{})
	return p
}

func (playstate *Playstate) Queue(playlist *Playlist) *Playlist {
	previous := playstate.Playlist
	playstate.Playlist = playlist
	return previous
}

func (playstate *Playstate) Clear() {
	playstate.Playlist = NewPlaylist([]*Track{})
}

func (playstate *Playstate) Update() {
	if playstate.Playlist.Length() == 0 {
		return
	}

	track := playstate.Playlist.Current()

	if track != playstate.current {
		if err := track.Open(); err != nil {
			log.Fatalln(err)
		}

		if err := track.Start(); err != nil {
			log.Fatalln(err)
		}
	}

	if err := track.Update(); err != nil {
		log.Fatalln(err)
	}

	playstate.current = track
}

func (playstate *Playstate) MainLoop() chan *Playstate {
	channel := make(chan *Playstate)

	go func() {
		for {
			playstate.Update()
			channel <- playstate
		}
	}()

	return channel
}
