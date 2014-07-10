package aural

import "github.com/mkb218/gosndfile/sndfile"

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
