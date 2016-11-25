package aural

import "github.com/gordonklaus/portaudio"

const (
	FRAMES_PER_BUFFER = 8196
)

type Track struct {
	Location string

	source AudioSource

	out    []int32
	stream *portaudio.Stream
}

func (track *Track) Open() error {
	track.source = NewAudioSource()
	err := track.source.Open(track.Location)

	if err != nil {
		return err
	}

	track.out = make([]int32, FRAMES_PER_BUFFER)

	stream, err := portaudio.OpenDefaultStream(
		0, int(track.source.Channels()),
		float64(track.source.SampleRate()),
		FRAMES_PER_BUFFER, &track.out)

	if err != nil {
		track.source.Close()
		return err
	}

	track.stream = stream
	return nil
}

func (track *Track) Update() (bool, error) {
	frames, err := track.source.ReadFrames(track.out)
	done := frames == 0

	if err != nil {
		return false, err
	}

	return done, track.stream.Write()
}

func (track *Track) Start() error {
	return track.stream.Start()
}

func (track *Track) Close() {
	track.stream.Close()
	track.Close()
}
