package aural

import "github.com/mkb218/gosndfile/sndfile"

type AudioSource interface {
	ReadFrames(out interface{}) (int64, error)

	Channels() int32
	SampleRate() int32

	Open(string) error
	Close()
}

type LibSndFileAudioSource struct {
	isOpen bool

	info *sndfile.Info
	io   *sndfile.File
}

func (source *LibSndFileAudioSource) ReadFrames(out interface{}) (int64, error) {
	return source.io.ReadFrames(out)
}

func (source *LibSndFileAudioSource) Close() {
	source.isOpen = false
	source.io.Close()
}

func (source *LibSndFileAudioSource) Channels() int32 {
	return source.info.Channels
}

func (source *LibSndFileAudioSource) SampleRate() int32 {
	return source.info.Samplerate
}

func (source *LibSndFileAudioSource) Open(identifier string) error {
	source.info = &sndfile.Info{}
	io, err := sndfile.Open(identifier, sndfile.Read, source.info)

	if err != nil {
		return err
	}

	source.io = io
	source.isOpen = true

	return nil
}

func NewAudioSource() AudioSource {
	return new(LibSndFileAudioSource)
}
