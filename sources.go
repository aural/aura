package aural

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/badgerodon/mp3"
	"github.com/mkb218/gosndfile/sndfile"

	"gopkg.in/h2non/filetype.v0"
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

func (source *LibSndFileAudioSource) Open(identifier string) error {
	source.info = &sndfile.Info{}
	file, err := sndfile.Open(identifier, sndfile.Read, source.info)

	if err != nil {
		return err
	}

	source.file = file
	source.isOpen = true

	return nil
}

func (source *LibSndFileAudioSource) ReadFrames(out []int32) (int64, error) {
	return source.file.ReadFrames(out)
}

func (source *LibSndFileAudioSource) Close() {
	source.isOpen = false
	source.file.Close()
}

func (source *LibSndFileAudioSource) Channels() int32 {
	return source.info.Channels
}

func (source *LibSndFileAudioSource) SampleRate() int32 {
	return source.info.Samplerate
}

type MP3AudioSource struct {
	file   *os.File
	frames *mp3.Frames
}

func (source *MP3AudioSource) Open(identifier string) error {
	var fileSeeker io.ReadSeeker
	file, err := os.Open(identifier)

	if err != nil {
		return err
	}

	fileSeeker = file
	frames, err := mp3.GetFrames(fileSeeker)

	if err != nil {
		file.Close()
		return err
	}

	source.file = file
	source.frames = frames

	return nil
}

func NewMP3AudioSource() AudioSource {
	log.Fatalln("Sorry, but MP3 support is not yet fully functional.")
	return &MP3AudioSource{}
}

func (source *MP3AudioSource) ReadFrames(out []int32) (totalSize int64, err error) {
	for i := 0; i < len(out); i++ {
		hasFrame := source.frames.Next()

		if hasFrame == false {
			log.Println()
			// TODO: Should we continue or return an error here?...
			return totalSize, source.frames.Error()
		}

		header := source.frames.Header()
		out[i] = int32(header.Samples)
		totalSize += header.Size
	}

	return totalSize, nil
}

func (source *MP3AudioSource) Close() {
	source.file.Close()
}

func (source *MP3AudioSource) Channels() int32 {
	switch source.frames.Header().ChannelMode {
	case mp3.SingleChannel:
		return 1
	default:
		return 2
	}
}

func (source *MP3AudioSource) SampleRate() int32 {
	return int32(source.frames.Header().SampleRate)
}
