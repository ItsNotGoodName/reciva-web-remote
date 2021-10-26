package store

// AddStream creates stream. Returns false when name is not unique.
func (s *Store) AddStream(name string, content string) (*Stream, bool) {
	s.sgMutex.Lock()
	// Make sure no duplicate name and find new id for stream
	id := 1
	for i := range s.sg.Streams {
		if s.sg.Streams[i].Name == name {
			s.sgMutex.Unlock()
			return nil, false
		}
		if s.sg.Streams[i].SID >= id {
			id = s.sg.Streams[i].SID + 1
		}
	}

	// Create stream and add to settings
	st := Stream{SID: id, Name: name, Content: content}
	s.sg.Streams = append(s.sg.Streams, st)
	s.sgMutex.Unlock()
	s.queueWrite()

	return &st, true
}

// GetStreams returns all streams.
func (s *Store) GetStreams() []Stream {
	s.sgMutex.Lock()
	st := make([]Stream, len(s.sg.Streams))
	copy(st, s.sg.Streams)
	s.sgMutex.Unlock()
	return st
}

func (s *Store) getStream(sid int) (*Stream, bool) {
	for i := range s.sg.Streams {
		if s.sg.Streams[i].SID == sid {
			st := s.sg.Streams[i]
			return &st, true
		}
	}
	return nil, false
}

// GetStream finds a stream by the given sid.
func (s *Store) GetStream(sid int) (*Stream, bool) {
	s.sgMutex.Lock()
	st, ok := s.getStream(sid)
	s.sgMutex.Unlock()
	return st, ok
}

// UpdateStream changes stream's name and content.
func (s *Store) UpdateStream(st *Stream) bool {
	idx := -1
	s.sgMutex.Lock()
	for i := range s.sg.Streams {
		if s.sg.Streams[i].SID == st.SID {
			idx = i
		} else if s.sg.Streams[i].Name == st.Name {
			s.sgMutex.Unlock()
			return false
		}
	}
	if idx == -1 {
		s.sgMutex.Unlock()
		return false
	}
	s.sg.Streams[idx] = *st
	s.sgMutex.Unlock()
	s.queueWrite()
	return true
}

// DeleteStream remove a stream and it's preset bindings.
func (s *Store) DeleteStream(sid int) int {
	deleted := 0
	s.sgMutex.Lock()

	// Delete stream
	newStreams := make([]Stream, 0, len(s.sg.Streams))
	for i := range s.sg.Streams {
		if s.sg.Streams[i].SID != sid {
			newStreams = append(newStreams, s.sg.Streams[i])
		} else {
			deleted += 1
		}
	}
	s.sg.Streams = newStreams

	s.clearStream(sid)
	s.sgMutex.Unlock()
	s.queueWrite()

	return deleted
}

func (s *Store) clearStream(sid int) bool {
	changed := false
	for i := range s.sg.Presets {
		if s.sg.Presets[i].SID == sid {
			s.sg.Presets[i].SID = 0
			changed = true
		}
	}
	return changed
}

// clearStream sets sid to 0 for all presets that have the given sid.
func (s *Store) ClearStream(sid int) bool {
	s.sgMutex.Lock()
	ok := s.clearStream(sid)
	s.sgMutex.Unlock()
	if ok {
		s.queueWrite()
	}
	return ok
}
