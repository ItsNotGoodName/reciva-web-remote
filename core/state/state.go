package state

import (
	"fmt"
	"log"
)

func (s *State) ValidPresetNumber(preset int) error {
	if preset < 1 || preset > len(s.Presets) {
		return fmt.Errorf("invalid preset number: %d", preset)
	}

	return nil
}

func (s *State) SetPresets(presets []Preset) {
	s.Presets = presets
	s.SetTitle(s.Title)
}

func (s *State) SetVolume(volume int) {
	s.Volume = NormalizeVolume(volume)
}

func (s *State) SetTitle(title string) {
	s.Title = title
	s.PresetNumber = 0
	for _, preset := range s.Presets {
		if preset.Title == title {
			s.PresetNumber = preset.Number
			return
		}
	}
}

func (s *State) SetAudioSources(audioSources []string) {
	s.AudioSources = audioSources
}

func (s *State) SetAudioSource(audioSource string) error {
	for _, source := range s.AudioSources {
		if source == audioSource {
			s.AudioSource = audioSource
			return nil
		}
	}

	return fmt.Errorf("invalid audio source: %s", audioSource)
}

// Merge fragment into state and return a fragment that has the merged changes.
func (s *State) Merge(f Fragment) (Fragment, bool) {
	changed := false

	if f.AudioSource != nil {
		if *f.AudioSource != s.AudioSource {
			oldAudioSource := s.AudioSource

			if err := s.SetAudioSource(*f.AudioSource); err != nil {
				log.Println("state.State.Merge:", err)
			} else if oldAudioSource != s.AudioSource {
				changed = true

				newAudioSource := s.AudioSource
				f.AudioSource = &newAudioSource
			}
		} else {
			f.AudioSource = nil
		}
	}

	if f.IsMuted != nil {
		if *f.IsMuted != s.IsMuted {
			s.IsMuted = *f.IsMuted
			changed = true
		} else {
			f.IsMuted = nil
		}
	}

	if f.Metadata != nil {
		if *f.Metadata != s.Metadata {
			s.Metadata = *f.Metadata
			changed = true
		} else {
			f.Metadata = nil
		}
	}

	if f.NewTitle != nil {
		if *f.NewTitle != s.NewTitle {
			s.NewTitle = *f.NewTitle
			changed = true
		} else {
			f.NewTitle = nil
		}
	}

	if f.NewURL != nil {
		if *f.NewURL != s.NewURL {
			s.NewURL = *f.NewURL
			changed = true
		} else {
			f.NewURL = nil
		}
	}
	if f.Power != nil {
		if *f.Power != s.Power {
			s.Power = *f.Power
			changed = true
		} else {
			f.Power = nil
		}
	}

	if f.Presets != nil {
		oldPresetNumber := s.PresetNumber
		oldTitle := s.Title

		s.SetPresets(f.Presets)
		changed = true

		newPresets := s.Presets
		f.Presets = newPresets
		if oldPresetNumber != s.PresetNumber {
			newPresetNumber := s.PresetNumber
			f.PresetNumber = &newPresetNumber
		}
		if oldTitle != s.Title {
			newTitle := s.Title
			f.Title = &newTitle
		}
	}

	if f.Status != nil {
		if *f.Status != s.Status {
			s.Status = *f.Status
			changed = true
		} else {
			f.Status = nil
		}
	}

	if f.Title != nil {
		if *f.Title != s.Title {
			oldPresetNumber := s.PresetNumber

			s.SetTitle(*f.Title)
			changed = true

			newTitle := s.Title
			f.Title = &newTitle
			if oldPresetNumber != s.PresetNumber {
				newPresetNumber := s.PresetNumber
				f.PresetNumber = &newPresetNumber
			}
		} else {
			f.Title = nil
		}
	}

	if f.URL != nil {
		if *f.URL != s.URL {
			s.URL = *f.URL
			changed = true
		} else {
			f.URL = nil
		}
	}

	if f.Volume != nil {
		if *f.Volume != s.Volume {
			oldVolume := s.Volume

			s.SetVolume(*f.Volume)
			if oldVolume != s.Volume {
				changed = true

				newVolume := s.Volume
				f.Volume = &newVolume
			}
		} else {
			f.Volume = nil
		}
	}

	return f, changed
}
