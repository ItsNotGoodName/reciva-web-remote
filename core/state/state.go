package state

func (s *State) SetAudioSource(audioSource string) error {
	if err := ValidAudioSource(s, audioSource); err != nil {
		return err
	}

	s.AudioSource = audioSource

	return nil
}

func (s *State) SetAudioSources(audioSources []string) {
	s.AudioSources = audioSources
}

func (s *State) SetIsMuted(isMuted bool) {
	s.IsMuted = isMuted
}

func (s *State) SetMetadata(metadata string) {
	s.Metadata = metadata
}

func (s *State) SetNewTitle(newTitle string) {
	s.NewTitle = newTitle
}

func (s *State) SetNewURL(newURL string) {
	s.NewURL = newURL
}

func (s *State) SetPower(power bool) {
	s.Power = power
}

func (s *State) SetPresets(presets []Preset) {
	s.Presets = presets
	s.SetTitle(s.Title)
}

func (s *State) SetStatus(status Status) {
	s.Status = status
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

func (s *State) SetURL(url string) {
	s.URL = url
}

func (s *State) SetVolume(volume int) {
	s.Volume = NormalizeVolume(volume)
}
