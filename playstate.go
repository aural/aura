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

func (playstate *Playstate) Update() *Playstate {
	if playstate.Playlist.Length() == 0 {
		return playstate
	}

	track := playstate.Playlist.Current()

	if track != playstate.current {
		if err := track.Open(); err != nil {
			log.Fatalln(err)
		}

		if err := track.Start(); err != nil {
			log.Fatalln(err)
		}

		log.Println("Now playing", track.Location)
	}

	done, err := track.Update()

	if done || err != nil {
		playstate.Playlist.Pop()
	}

	if err != nil {
		log.Println("Skipping due to error:", err)
	}

	playstate.current = track
	return playstate
}

func (playstate *Playstate) MainLoop() chan *Playstate {
	channel := make(chan *Playstate)

	go func() {
		for {
			channel <- playstate.Update()
		}
	}()

	return channel
}
