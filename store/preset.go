package store

func (s *Store) GetPresets() []Preset {
	s.sgMutex.Lock()
	pt := make([]Preset, len(s.sg.Presets))
	copy(pt, s.sg.Presets)
	s.sgMutex.Unlock()
	return pt
}

func (s *Store) GetPreset(uri string) (*Preset, bool) {
	s.sgMutex.Lock()
	for i := range s.sg.Presets {
		if s.sg.Presets[i].URI == uri {
			pt := s.sg.Presets[i]
			s.sgMutex.Unlock()
			return &pt, true
		}
	}
	s.sgMutex.Unlock()
	return nil, false
}

func (s *Store) UpdatePreset(pt *Preset) bool {
	if pt.SID == 0 {
		return s.ClearPreset(pt.URI)
	}

	s.sgMutex.Lock()

	if _, ok := s.getStream(pt.SID); !ok {
		s.sgMutex.Unlock()
		return false
	}

	changed := false
	ok := false
	for i := range s.sg.Presets {
		if s.sg.Presets[i].URI == pt.URI {
			ok = true
			if s.sg.Presets[i].SID != pt.SID {
				s.sg.Presets[i].SID = pt.SID
				changed = true
			}
		} else if s.sg.Presets[i].SID == pt.SID {
			// Clear duplicate SID mappings
			s.sg.Presets[i].SID = 0
			changed = true
		}
	}
	s.sgMutex.Unlock()
	if changed {
		s.queueWrite()
	}
	return ok
}

// ClearPreset sets preset's SID to 0.
func (s *Store) ClearPreset(uri string) bool {
	s.sgMutex.Lock()
	for i := range s.sg.Presets {
		if s.sg.Presets[i].URI != uri {
			s.sg.Presets[i].SID = 0
			s.queueWrite()
			s.sgMutex.Unlock()
			return true
		}
	}
	s.sgMutex.Unlock()
	return false
}
