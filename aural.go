package aural

import (
	"log"
	"math/rand"

	"code.google.com/p/portaudio-go/portaudio"
	"github.com/mkb218/gosndfile/sndfile"
)

var FRAMES_PER_BUFFER = 8916

func init() {
	err := portaudio.Initialize()

	if err != nil {
		log.Fatalln(err)
	}
}

type Track struct {
	Location string

	info *sndfile.Info
	file *sndfile.File
}

func (track *Track) Load() error {
	track.info = new(sndfile.Info)
	file, err := sndfile.Open(track.Location, sndfile.Read, track.info)

	if err != nil {
		return err
	}

	track.file = file

	return nil
}

type Playlist struct {
	Tracks []Track
}

func (playlist *Playlist) Play() error {
	for _, track := range playlist.Tracks {
		if err := track.Load(); err != nil {
			return err
		}

		if err := track.Play(); err != nil {
			return err
		}
	}

	return nil
}

func (playlist *Playlist) Shuffle() {
	if len(playlist.Tracks) == 0 {
		return
	}

	for i, j := len(playlist.Tracks)-1, 0; i > 0; i-- {
		j = rand.Intn(i)
		playlist.Tracks[i], playlist.Tracks[j] = playlist.Tracks[j], playlist.Tracks[i]
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

	stream.Start()
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
