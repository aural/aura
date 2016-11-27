package aural

type Track struct {
	Location string

	source AudioSource
}

func (this *Track) Open() error {
	this.source = NewAudioSource(this.Location)
	err := this.source.Open(this.Location)

	if err != nil {
		return err
	}

	return nil
}

func (this *Track) Update(playstate *Playstate) (bool, error) {
	size, err := this.source.ReadFrames(playstate.out)

	if err != nil {
		return false, err
	}

	return size == 0, nil
}

func (this *Track) Close() {
	this.source.Close()
}
