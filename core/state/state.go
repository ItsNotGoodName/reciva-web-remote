package state

func (s *State) SetAudioSource(audioSource string) (int, error) {
	if err := ValidAudioSource(s, audioSource); err != nil {
		return 0, err
	}

	s.AudioSource = audioSource

	return ChangedAudioSource, nil
}

func (s *State) SetAudioSources(audioSources []string) {
	s.AudioSources = audioSources
}

func (s *State) SetIsMuted(isMuted bool) int {
	s.IsMuted = isMuted
	return ChangedIsMuted
}

func (s *State) SetMetadata(metadata string) int {
	s.Metadata = metadata
	return ChangedMetadata
}

func (s *State) SetTitleNew(titleNew string) int {
	s.TitleNew = titleNew
	return ChangedTitleNew
}

func (s *State) SetURLNew(urlNew string) int {
	s.URLNew = urlNew
	return ChangedURLNew
}

func (s *State) SetPower(power bool) int {
	s.Power = power
	return ChangedPower
}

func (s *State) SetPresets(presets []Preset) int {
	s.Presets = presets
	return ChangedPresets | s.setPresetNumber()
}

func (s *State) SetStatus(status Status) int {
	s.Status = status
	return ChangedStatus
}

func (s *State) SetTitle(title string) int {
	s.Title = title
	return ChangedTitle | s.setPresetNumber()
}

func (s *State) setPresetNumber() int {
	s.PresetNumber = 0
	for _, preset := range s.Presets {
		if preset.Title == s.Title {
			s.PresetNumber = preset.Number
			return ChangedPresetNumber
		}
	}
	return 0
}

func (s *State) SetURL(url string) int {
	s.URL = url
	return ChangedURL
}

func (s *State) SetVolume(volume int) int {
	s.Volume = NormalizeVolume(volume)
	return ChangedVolume
}
