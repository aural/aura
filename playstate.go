package aural

import (
	"log"

	"code.google.com/p/portaudio-go/portaudio"
)

func init() {
	if err := portaudio.Initialize(); err != nil {
		log.Fatalln(err)
	}
}

func Terminate() {
	portaudio.Terminate()
}

type Playable interface {
	Play()
	Queue()
}

type PlayState struct {
	CurrentTrack Track

	Tracks      []Track
	Description string
	Shuffled    bool
}

func (playState *PlayState) Update(channel chan *PlayState) {
	if len(playState.Tracks) == 0 {
		channel <- playState
		return
	}

	playState.CurrentTrack = playState.Tracks[0]
	track := playState.CurrentTrack

	channel <- playState

	if err := track.Open(); err != nil {
		log.Fatalln(err)
	}

	if err := track.stream.Start(); err != nil {
		log.Fatalln(err)
	}

	for {
		remaining, err := track.file.ReadFrames(track.out)

		if err != nil {
			log.Fatalln(err)
		}

		if remaining == 0 {
			var newTracks []Track
			for i := 1; i < len(playState.Tracks); i++ {
				newTracks = append(newTracks, playState.Tracks[i])
			}
			playState.Tracks = newTracks
			break
		}

		track.stream.Write()
	}
}
