package state

func (s *State) SetAudioSource(audioSource string) (Changed, error) {
	if err := ValidAudioSource(s, audioSource); err != nil {
		return 0, err
	}

	if s.AudioSource == audioSource {
		return 0, nil
	}

	s.AudioSource = audioSource
	return ChangedAudioSource, nil
}

func (s *State) SetAudioSources(audioSources []string) {
	s.AudioSources = audioSources
}

func (s *State) SetIsMuted(isMuted bool) Changed {
	if s.IsMuted == isMuted {
		return 0
	}

	s.IsMuted = isMuted
	return ChangedIsMuted
}

func (s *State) SetMetadata(metadata string) Changed {
	if s.Metadata == metadata {
		return 0
	}

	s.Metadata = metadata
	return ChangedMetadata
}

func (s *State) SetTitleNew(titleNew string) Changed {
	if s.TitleNew == titleNew {
		return 0
	}

	s.TitleNew = titleNew
	return ChangedTitleNew
}

func (s *State) SetURLNew(urlNew string) Changed {
	if s.URLNew == urlNew {
		return 0
	}

	s.URLNew = urlNew
	return ChangedURLNew
}

func (s *State) SetPower(power bool) Changed {
	if s.Power == power {
		return 0
	}

	s.Power = power
	return ChangedPower
}

func (s *State) SetPresets(presets []Preset) Changed {
	s.Presets = presets
	return ChangedPresets.Merge(s.calculatePresetNumber())
}

func (s *State) SetStatus(status Status) Changed {
	if s.Status == status {
		return 0
	}

	s.Status = status
	return ChangedStatus
}

func (s *State) SetTitle(title string) Changed {
	if s.Title == title {
		return 0
	}

	s.Title = title
	return ChangedTitle.Merge(s.calculatePresetNumber())
}

func (s *State) calculatePresetNumber() Changed {
	s.PresetNumber = 0
	for _, preset := range s.Presets {
		if preset.Title == s.Title {
			s.PresetNumber = preset.Number
			return ChangedPresetNumber
		}
	}
	return 0
}

func (s *State) SetURL(url string) Changed {
	if s.URL == url {
		return 0
	}

	s.URL = url
	return ChangedURL
}

func (s *State) SetVolume(volume int) Changed {
	volume = NormalizeVolume(volume)
	if s.Volume == volume {
		return 0
	}

	s.Volume = volume
	return ChangedVolume
}
