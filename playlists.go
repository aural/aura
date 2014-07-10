package aural

import (
	"math/rand"
	"time"
)

type Playlist struct {
	Tracks []Track
}

func (playlist *Playlist) Play() error {
	for _, track := range playlist.Tracks {
		if err := track.Load(); err != nil {
			return err
		}

		if err := track.Play(); err != nil {
			return err
		}
	}

	return nil
}

func (playlist *Playlist) Shuffle() {
	if len(playlist.Tracks) == 0 {
		return
	}

	rand.Seed(time.Now().UTC().UnixNano())

	for i, j := len(playlist.Tracks)-1, 0; i > 0; i-- {
		j = rand.Intn(i)
		playlist.Tracks[i], playlist.Tracks[j] = playlist.Tracks[j], playlist.Tracks[i]
	}
}
