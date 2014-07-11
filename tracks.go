package aural

import (
	"code.google.com/p/portaudio-go/portaudio"
	"github.com/mkb218/gosndfile/sndfile"
)

const (
	FRAMES_PER_BUFFER = 8196
)

type Track struct {
	Location string

	info *sndfile.Info
	file *sndfile.File

	isOpen bool
	out    []int32
	stream *portaudio.Stream
}

func (track *Track) Open() error {
	track.info = new(sndfile.Info)
	file, err := sndfile.Open(track.Location, sndfile.Read, track.info)

	if err != nil {
		return err
	}

	track.file = file

	track.out = make([]int32, FRAMES_PER_BUFFER)

	stream, err := portaudio.OpenDefaultStream(
		0, int(track.info.Channels),
		float64(track.info.Samplerate),
		FRAMES_PER_BUFFER, &track.out)

	if err != nil {
		track.file.Close()
		return err
	}

	track.stream = stream
	track.isOpen = true

	return nil
}

func (track *Track) Close() error {
	track.isOpen = false
	track.stream.Close()
	return track.file.Close()
}

func (track Track) Play(playState *PlayState) {
	var tracks []Track
	playState.Tracks = append(tracks, track)
}

func (track Track) Queue(playState *PlayState) {
	playState.Tracks = append(playState.Tracks, track)
}
