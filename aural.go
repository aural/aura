package aural

import (
	"log"

	"code.google.com/p/portaudio-go/portaudio"
	"github.com/mkb218/gosndfile/sndfile"
)

func init() {
	err := portaudio.Initialize()

	if err != nil {
		log.Fatalln(err)
	}
}

var FRAMES_PER_BUFFER = 8916

func Play(fileName string) (err error) {
	var info sndfile.Info

	out := make([]int32, FRAMES_PER_BUFFER)
	file, err := sndfile.Open(fileName, sndfile.Read, &info)

	if err != nil {
		return err
	}

	stream, err := portaudio.OpenDefaultStream(
		0, int(info.Channels),
		float64(info.Samplerate),
		FRAMES_PER_BUFFER, &out)

	if err != nil {
		return err
	}

	defer stream.Close()

	stream.Start()
	defer stream.Stop()

	log.Printf("Playing %vKhz stream with %v channels.\n", info.Samplerate, info.Channels)

	remaining := int64(1)

	for remaining > 0 {
		remaining, err = file.ReadFrames(out)

		if err != nil {
			return err
		}

		stream.Write()
	}

	return nil
}

func Terminate() {
	portaudio.Terminate()
}
