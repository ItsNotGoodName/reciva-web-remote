package state

func (s *State) Fragment() Fragment {
	frag := NewFragment(s.UUID)

	// AudioSource
	audioSource := s.AudioSource
	frag.AudioSource = &audioSource

	// IsMuted
	isMuted := s.IsMuted
	frag.IsMuted = &isMuted

	// Metadata
	metadata := s.Metadata
	frag.Metadata = &metadata

	// TitleNew
	titleNew := s.TitleNew
	frag.TitleNew = &titleNew

	// URLNew
	urlNew := s.URLNew
	frag.URLNew = &urlNew

	// Power
	power := s.Power
	frag.Power = &power

	// Presets
	frag.Presets = s.Presets

	// Status
	status := s.Status
	frag.Status = &status

	// Title
	title := s.Title
	frag.Title = &title

	// URL
	url := s.URL
	frag.URL = &url

	// Volume
	volume := s.Volume
	frag.Volume = &volume

	return frag
}
