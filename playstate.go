package aural

import "code.google.com/p/portaudio-go/portaudio"

var FRAMES_PER_BUFFER = 8916

type Playable interface {
	Play()
	Queue()
}

var PlayState struct {
	CurrentTrack Track

	Tracks      []Track
	Description string
	Shuffled    bool
	Started     bool
}

func Continue() error {
	if len(PlayState.Tracks) == 0 {
		return nil
	}

	out := make([]int32, FRAMES_PER_BUFFER)

	PlayState.CurrentTrack = PlayState.Tracks[0]
	track := PlayState.CurrentTrack

	if track.info == nil {
		if err := track.Load(); err != nil {
			return err
		}
	}

	stream, err := portaudio.OpenDefaultStream(
		0, int(track.info.Channels),
		float64(track.info.Samplerate),
		FRAMES_PER_BUFFER, &out)

	if err != nil {
		return err
	}

	defer stream.Close()

	err = stream.Start()

	if err != nil {
		return err
	}

	defer stream.Stop()

	for {
		remaining, err := track.file.ReadFrames(out)

		if err != nil {
			return err
		}

		if remaining == 0 {
			var newTracks []Track
			for i := 1; i < len(PlayState.Tracks); i++ {
				newTracks = append(newTracks, PlayState.Tracks[i])
			}
			PlayState.Tracks = newTracks
			break
		}

		stream.Write()
	}

	return nil
}

func Start() error {
	err := portaudio.Initialize()

	if err != nil {
		return err
	}

	defer portaudio.Terminate()

	PlayState.Started = true

	for PlayState.Started == true {
		if err := Continue(); err != nil {
			// TODO: Handle error case here.
			return err
		}

		if len(PlayState.Tracks) == 0 {
			PlayState.Started = false
		}
	}

	return nil
}
