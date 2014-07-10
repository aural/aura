package aural

import (
	"log"

	"code.google.com/p/portaudio-go/portaudio"
)

var FRAMES_PER_BUFFER = 8916

func init() {
	err := portaudio.Initialize()

	if err != nil {
		log.Fatalln(err)
	}
}

func (track *Track) Play() error {
	out := make([]int32, FRAMES_PER_BUFFER)

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
			break
		}

		stream.Write()
	}

	return nil
}

func Terminate() {
	portaudio.Terminate()
}
