package aural

import (
	"os"

	"github.com/mkb218/gosndfile/sndfile"
	"github.com/tcolgate/mp3"
)

type AudioSource interface {
	ReadFrames(out interface{}) (int64, error)

	Channels() int32
	SampleRate() int32

	Open(string) error
	Close()
}

func NewAudioSource() AudioSource {
	return new(LibSndFileAudioSource)
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

type MP3AudioSource struct {
	file    *os.File
	decoder *mp3.Decoder

	lastFrame mp3.Frame
}

func (source *MP3AudioSource) Channels() int32 {
	switch source.lastFrame.Header().ChannelMode() {
	case mp3.Stereo, mp3.JointStereo, mp3.DualChannel:
		return 2
	default:
		return 1
	}
}

func (source *MP3AudioSource) SampleRate() int32 {
	return int32(source.lastFrame.Header().SampleRate())
}

func (source *MP3AudioSource) ReadFrames(out interface{}) (int64, error) {
	err := source.decoder.Decode(&source.lastFrame)
	return int64(source.lastFrame.Size()), err
}

func (source *MP3AudioSource) Open(identifier string) error {
	file, err := os.Open(identifier)

	if err != nil {
		return err
	}

	source.file = file
	source.decoder = mp3.NewDecoder(file)

	return nil
}

func (source *MP3AudioSource) Close() {
	source.file.Close()
}
