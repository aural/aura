package aural

import (
	"io/ioutil"
	"path"

	"github.com/mkb218/gosndfile/sndfile"
	"gopkg.in/h2non/filetype.v0"

	mpg "github.com/bobertlo/go-mpg123/mpg123"
)

type AudioSourceFactory func() AudioSource

type AudioSource interface {
	ReadFrames(out []int32) (int64, error)

	Channels() int32
	SampleRate() int32

	Open(string) error
	Close()
}

var sourceTypes map[string]AudioSourceFactory

func init() {
	sourceTypes = make(map[string]AudioSourceFactory)
	sourceTypes["mp3"] = NewMP3AudioSource
}

func GetExtensionFor(identifier string) string {
	buf, err := ioutil.ReadFile(identifier)
	extension := path.Ext(identifier)

	if kind, unknown := filetype.Match(buf); err == nil {
		if unknown == nil {
			extension = kind.Extension
		}
	}

	return extension
}

func NewAudioSource(identifier string) AudioSource {
	// TODO: Support URLs, not just file paths

	extension := GetExtensionFor(identifier)
	factory, ok := sourceTypes[extension]

	if ok != true {
		factory = NewLibSndFileAudioSource
	}

	return factory()
}

type LibSndFileAudioSource struct {
	isOpen bool

	info *sndfile.Info
	file *sndfile.File
}

func NewLibSndFileAudioSource() AudioSource {
	return &LibSndFileAudioSource{}
}

func (this *LibSndFileAudioSource) Open(identifier string) error {
	this.info = &sndfile.Info{}
	file, err := sndfile.Open(identifier, sndfile.Read, this.info)

	if err != nil {
		return err
	}

	this.file = file
	this.isOpen = true

	return nil
}

func (this *LibSndFileAudioSource) ReadFrames(out []int32) (int64, error) {
	return this.file.ReadFrames(out)
}

func (this *LibSndFileAudioSource) Close() {
	this.isOpen = false
	this.file.Close()
}

func (this *LibSndFileAudioSource) Channels() int32 {
	return this.info.Channels
}

func (this *LibSndFileAudioSource) SampleRate() int32 {
	return this.info.Samplerate
}

type MP3AudioSource struct {
	sampleRate int64
	channels   int

	decoder *mpg.Decoder
}

func (this *MP3AudioSource) Open(identifier string) error {
	decoder, err := mpg.NewDecoder("")

	if err != nil {
		return err
	}

	err = decoder.Open(identifier)

	if err != nil {
		return err
	}

	sampleRate, channels, encoding := decoder.GetFormat()

	decoder.FormatNone()
	decoder.Format(sampleRate, channels, encoding)

	this.sampleRate = sampleRate
	this.channels = channels
	this.decoder = decoder

	return nil
}

func NewMP3AudioSource() AudioSource {
	return &MP3AudioSource{}
}

func (this *MP3AudioSource) ReadFrames(out []int32) (int64, error) {
	buffer := make([]byte, len(out))
	length, err := this.decoder.Read(buffer)

	if err != nil {
		return 0, err
	}

	for index := 0; index < len(out); index++ {
		out[index] = int32(buffer[index])
	}

	return int64(length), nil
}

func (this *MP3AudioSource) Close() {
	this.sampleRate = 0
	this.channels = 0
	this.decoder.Close()
}

func (this *MP3AudioSource) Channels() int32 {
	return int32(this.channels)
}

func (this *MP3AudioSource) SampleRate() int32 {
	return int32(this.sampleRate)
}
