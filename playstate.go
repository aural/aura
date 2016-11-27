package aural

import (
	"log"

	"github.com/gordonklaus/portaudio"
)

const (
	FRAMES_PER_BUFFER = 8196
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

	stream *portaudio.Stream
	out    []int32

	isStarted bool
}

func NewPlaystate() (*Playstate, error) {
	playstate := Playstate{
		Playlist:  NewPlaylist([]*Track{}),
		out:       make([]int32, FRAMES_PER_BUFFER),
		isStarted: false,
	}

	stream, err := portaudio.OpenDefaultStream(
		0, int(2),
		float64(44100),
		FRAMES_PER_BUFFER, &playstate.out)

	if err != nil {
		return nil, err
	}

	playstate.stream = stream

	return &playstate, nil
}

func (playstate *Playstate) Queue(playlist *Playlist) *Playlist {
	previous := playstate.Playlist
	playstate.Playlist = playlist
	return previous
}

func (playstate *Playstate) Clear() {
	playstate.Playlist = NewPlaylist([]*Track{})
}

func (playstate *Playstate) updateStreamState() bool {
	if playstate.Playlist.Length() == 0 {
		if playstate.isStarted {
			log.Println("Playlist is now empty.")

			playstate.stream.Stop()
			playstate.isStarted = false
		}

		return true
	}

	if !playstate.isStarted {
		if err := playstate.stream.Start(); err != nil {
			playstate.stream.Close()
			log.Fatalln(err)
		}

		playstate.isStarted = true
	}

	return false
}

func (playstate *Playstate) Update() *Playstate {
	if skipUpdate := playstate.updateStreamState(); skipUpdate == true {
		return playstate
	}

	playstate.isStarted = true
	track := playstate.Playlist.Current()

	if track != playstate.current {
		if err := track.Open(); err != nil {
			log.Fatalln(err)
		}

		log.Println("Now playing", track.Location)
		playstate.current = track
	}

	done, err := track.Update(playstate)

	if done || err != nil {
		playstate.Playlist.Pop()
	}

	if err != nil {
		log.Println("Skipping due to error:", err)
	}

	if err = playstate.stream.Write(); err != nil {
		log.Println("Error writing data to audio hardware.")
	}

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
