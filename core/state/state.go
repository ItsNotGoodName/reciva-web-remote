package state

import (
	"fmt"
	"log"
)

func New(uuid, name, modelName, modelNumber string) State {
	return State{
		ModelName:   modelName,
		ModelNumber: modelNumber,
		Name:        name,
		Status:      StatusUnknown,
		UUID:        uuid,
	}
}

func NormalizeVolume(volume int) int {
	if volume < 0 {
		return 0
	}
	if volume > 100 {
		return 100
	}

	return volume
}

func ParsePresetsCount(presetsCount int) (int, error) {
	presetsCount = presetsCount - 2
	if presetsCount < 1 {
		return 0, fmt.Errorf("invalid presets count: %d", presetsCount)
	}

	return presetsCount, nil
}

func ParseStatus(status string) Status {
	switch {
	case status == StatusConnecting:
		return StatusConnecting
	case status == StatusPlaying:
		return StatusPlaying
	case status == StatusStopped:
		return StatusStopped
	default:
		return StatusUnknown
	}
}

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
	for _, preset := range s.Presets {
		if preset.Title == title {
			s.PresetNumber = preset.Number
			return
		}
	}
	s.PresetNumber = 0
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

// Merge fragment into state and return a fragment of what changed.
func (s *State) Merge(f Fragment) (Fragment, bool) {
	changed := false
	if f.AudioSource != nil {
		oldAudioSource := s.AudioSource

		if err := s.SetAudioSource(*f.AudioSource); err != nil {
			log.Println("state.State.Merge:", err)
		} else if oldAudioSource != s.AudioSource {
			changed = true

			newAudioSource := s.AudioSource
			f.AudioSource = &newAudioSource
		}
	}
	if f.IsMuted != nil && *f.IsMuted != s.IsMuted {
		s.IsMuted = *f.IsMuted
		changed = true
	}
	if f.Metadata != nil && *f.Metadata != s.Metadata {
		s.Metadata = *f.Metadata
		changed = true
	}
	if f.Power != nil && *f.Power != s.Power {
		s.Power = *f.Power
		changed = true
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
	if f.Status != nil && *f.Status != s.Status {
		s.Status = *f.Status
		changed = true
	}
	if f.Title != nil && *f.Title != s.Title {
		oldPresetNumber := s.PresetNumber

		s.SetTitle(*f.Title)
		changed = true

		newTitle := s.Title
		f.Title = &newTitle
		if oldPresetNumber != s.PresetNumber {
			newPresetNumber := s.PresetNumber
			f.PresetNumber = &newPresetNumber
		}
	}
	if f.URL != nil && *f.URL != s.URL {
		s.URL = *f.URL
		changed = true
	}
	if f.Volume != nil {
		oldVolume := s.Volume

		s.SetVolume(*f.Volume)
		if oldVolume != s.Volume {
			changed = true

			newVolume := s.Volume
			f.Volume = &newVolume
		}
	}

	return f, changed
}
