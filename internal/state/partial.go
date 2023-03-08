package state

func GetPartial(s *State, c Changed) Partial {
	p := NewPartial(s.UUID)

	// AudioSource
	if c&ChangedAudioSource == ChangedAudioSource {
		audioSource := s.AudioSource
		p.AudioSource = &audioSource
	}

	// IsMuted
	if c&ChangedIsMuted == ChangedIsMuted {
		isMuted := s.IsMuted
		p.IsMuted = &isMuted
	}

	// Metadata
	if c&ChangedMetadata == ChangedMetadata {
		metadata := s.Metadata
		p.Metadata = &metadata
	}

	// Power
	if c&ChangedPower == ChangedPower {
		power := s.Power
		p.Power = &power
	}

	// PresetNumber
	if c&ChangedPresetNumber == ChangedPresetNumber {
		presetNumber := s.PresetNumber
		p.PresetNumber = &presetNumber
	}

	// Presets
	if c&ChangedPresets == ChangedPresets {
		p.Presets = s.Presets
	}

	// Status
	if c&ChangedStatus == ChangedStatus {
		status := s.Status
		p.Status = &status
	}

	// Title
	if c&ChangedTitle == ChangedTitle {
		title := s.Title
		p.Title = &title
	}

	// TitleNew
	if c&ChangedTitleNew == ChangedTitleNew {
		titleNew := s.TitleNew
		p.TitleNew = &titleNew
	}

	// URL
	if c&ChangedURL == ChangedURL {
		url := s.URL
		p.URL = &url
	}

	// URLNew
	if c&ChangedURLNew == ChangedURLNew {
		urlNew := s.URLNew
		p.URLNew = &urlNew
	}

	// Volume
	if c&ChangedVolume == ChangedVolume {
		volume := s.Volume
		p.Volume = &volume
	}

	return p
}
