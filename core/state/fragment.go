package state

import "log"

func NewFragment(uuid string) Fragment {
	return Fragment{
		UUID: uuid,
	}
}

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

// Merge fragment into state and return a fragment that has the merged changes.
func (s *State) Merge(f Fragment) bool {
	changed := false

	// AudioSource
	if f.AudioSource != nil && *f.AudioSource != s.AudioSource {
		if err := s.SetAudioSource(*f.AudioSource); err != nil {
			log.Println("state.State.Merge:", err)
		} else {
			changed = true
		}
	}

	// IsMuted
	if f.IsMuted != nil && *f.IsMuted != s.IsMuted {
		s.SetIsMuted(*f.IsMuted)
		changed = true
	}

	// Metadata
	if f.Metadata != nil && *f.Metadata != s.Metadata {
		s.SetMetadata(*f.Metadata)
		changed = true
	}

	// TitleNew
	if f.TitleNew != nil && *f.TitleNew != s.TitleNew {
		s.SetTitleNew(*f.TitleNew)
		changed = true
	}

	// URLNew
	if f.URLNew != nil && *f.URLNew != s.URLNew {
		s.SetURLNew(*f.URLNew)
		changed = true
	}

	// Power
	if f.Power != nil && *f.Power != s.Power {
		s.SetPower(*f.Power)
		changed = true
	}

	// Presets
	if f.Presets != nil {
		s.SetPresets(f.Presets)
		changed = true
	}

	// Status
	if f.Status != nil && *f.Status != s.Status {
		s.SetStatus(*f.Status)
		changed = true
	}

	// Title
	if f.Title != nil && *f.Title != s.Title {
		s.SetTitle(*f.Title)
		changed = true
	}

	// URL
	if f.URL != nil && *f.URL != s.URL {
		s.SetURL(*f.URL)
		changed = true
	}

	// Volume
	if f.Volume != nil && *f.Volume != s.Volume {
		s.SetVolume(*f.Volume)
		changed = true
	}

	return changed
}
