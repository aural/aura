package aural

type Track struct {
	Location string

	source AudioSource
}

func (track *Track) Open() error {
	track.source = NewAudioSource(track.Location)
	err := track.source.Open(track.Location)

	if err != nil {
		return err
	}

	return nil
}

func (track *Track) Update(playstate *Playstate) (bool, error) {
	frames, err := track.source.ReadFrames(playstate.out)
	done := frames == 0

	if err != nil {
		return false, err
	}

	return done, playstate.stream.Write()
}

func (track *Track) Close() {
	track.source.Close()
}
