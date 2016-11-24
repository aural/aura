package aural

import (
	"github.com/gordonklaus/portaudio"
	"github.com/mkb218/gosndfile/sndfile"
)

const (
	FRAMES_PER_BUFFER = 8196
)

type Track struct {
	Location string

	info *sndfile.Info
	io   *sndfile.File

	isOpen bool
	out    []int32
	stream *portaudio.Stream
}

func (track *Track) Open() error {
	track.info = &sndfile.Info{}
	io, err := sndfile.Open(track.Location, sndfile.Read, track.info)

	if err != nil {
		return err
	}

	track.io = io
	track.out = make([]int32, FRAMES_PER_BUFFER)

	stream, err := portaudio.OpenDefaultStream(
		0, int(track.info.Channels),
		float64(track.info.Samplerate),
		FRAMES_PER_BUFFER, &track.out)

	if err != nil {
		track.io.Close()
		return err
	}

	track.stream = stream
	track.isOpen = true

	return nil
}

func (track *Track) Update() error {
	_, err := track.io.ReadFrames(track.out)
	if err != nil {
		return err
	}
	track.stream.Write()
	return nil
}

func (track *Track) Start() error {
	return track.stream.Start()
}

func (track *Track) Close() error {
	track.isOpen = false
	track.stream.Close()
	return track.io.Close()
}
