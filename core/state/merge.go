package state

import "log"

// Merge fragment into state and return a fragment that has the merged changes.
func (s *State) Merge(f Fragment) (changed int) {
	// AudioSource
	if f.AudioSource != nil && *f.AudioSource != s.AudioSource {
		if num, err := s.SetAudioSource(*f.AudioSource); err != nil {
			log.Println("state.State.Merge:", err)
		} else {
			changed = changed | num
		}
	}

	// IsMuted
	if f.IsMuted != nil && *f.IsMuted != s.IsMuted {
		changed = s.SetIsMuted(*f.IsMuted) | changed
	}

	// Metadata
	if f.Metadata != nil && *f.Metadata != s.Metadata {
		changed = s.SetMetadata(*f.Metadata) | changed
	}

	// TitleNew
	if f.TitleNew != nil && *f.TitleNew != s.TitleNew {
		changed = s.SetTitleNew(*f.TitleNew) | changed
	}

	// URLNew
	if f.URLNew != nil && *f.URLNew != s.URLNew {
		changed = s.SetURLNew(*f.URLNew) | changed
	}

	// Power
	if f.Power != nil && *f.Power != s.Power {
		changed = s.SetPower(*f.Power) | changed
	}

	// Presets
	if f.Presets != nil {
		changed = s.SetPresets(f.Presets) | changed
	}

	// Status
	if f.Status != nil && *f.Status != s.Status {
		changed = s.SetStatus(*f.Status) | changed
	}

	// Title
	if f.Title != nil && *f.Title != s.Title {
		changed = s.SetTitle(*f.Title) | changed
	}

	// URL
	if f.URL != nil && *f.URL != s.URL {
		changed = s.SetURL(*f.URL) | changed
	}

	// Volume
	if f.Volume != nil && *f.Volume != s.Volume {
		changed = s.SetVolume(*f.Volume) | changed
	}

	return changed
}
