package store

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
)

func NewService(cfg *config.Config) (*Store, error) {
	dctx := context.Background()
	dctx, cancel := context.WithCancel(dctx)

	s := &Store{
		dctx:           dctx,
		Cancel:         cancel,
		file:           cfg.ConfigPath,
		writeChan:      make(chan chan error),
		readChan:       make(chan chan error),
		queueWriteChan: make(chan bool),
	}

	if st, err := s.readSettings(); err != nil {
		cfg.EnablePresets = false
		return nil, err
	} else {
		s.st = st
	}

	if len(cfg.Presets) > 0 {
		p := make([]Preset, 0, len(cfg.Presets))
		for _, v := range cfg.Presets {
			p = append(p, Preset{URI: v})
		}
		s.st.mergePresets(p)
	}

	if len(s.st.Presets) > 0 {
		cfg.EnablePresets = true
		var u []string
		for _, v := range s.st.Presets {
			u = append(u, v.URI)
		}
		cfg.Presets = u
	}

	// Prioritize flag port over settings port unless default
	if cfg.Port != s.st.Port && cfg.Port == config.DefaultPort {
		cfg.Port = s.st.Port
	}

	// Prioritize flag cport over settings cport unless default
	if cfg.CPort != s.st.CPort && cfg.CPort == goupnpsub.DefaultPort {
		cfg.CPort = s.st.CPort
	}

	// Save and start loop
	s.writeSettings(s.st)
	go s.storeLoop()

	return s, nil
}

func (s *Store) writeSettings(st *Settings) error {
	b, err := json.MarshalIndent(st, "", "	")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.file, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) readSettings() (*Settings, error) {
	b, err := ioutil.ReadFile(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Store.readSettings: creating new settings file", s.file)
			st := NewSettings()
			err = s.writeSettings(st)
			return st, err
		}
		return nil, err
	}

	st := Settings{}
	if err = json.Unmarshal(b, &st); err != nil {
		return nil, err
	}

	return &st, nil
}

func (s *Store) queueWrite() {
	select {
	case s.queueWriteChan <- true:
	case <-s.dctx.Done():
	}
}

func (s *Store) storeLoop() {
	ticker := time.NewTicker(15 * time.Second)
	save := false
	writeSettings := func() error {
		s.stMutex.Lock()
		st := *s.st
		s.stMutex.Unlock()
		return s.writeSettings(&st)
	}
	readSettings := func() error {
		st, err := s.readSettings()
		if err != nil {
			return err
		}
		s.stMutex.Lock()
		s.st = st
		s.stMutex.Unlock()
		return nil
	}

	for {
		select {
		case save = <-s.queueWriteChan:
		case res := <-s.readChan:
			err := readSettings()
			select {
			case res <- err:
			case <-s.dctx.Done():
			}
		case res := <-s.writeChan:
			err := writeSettings()
			select {
			case res <- err:
			case <-s.dctx.Done():
			}
		case <-ticker.C:
			if save {
				if err := writeSettings(); err != nil {
					log.Println("Store.storeLoop(ERROR):", err)
				} else {
					log.Println("Store.storeLoop: settings saved")
					save = false
				}
			}
		case <-s.dctx.Done():
			if save {
				log.Println("Store.storeLoop: dctx is done, saving")
				if err := writeSettings(); err != nil {
					log.Println("Store.storeLoop(ERROR):", err)
				}
				return
			}
			log.Println("Store.storeLoop: dctx is done")
			return
		}
	}
}

func (s *Store) AddStream(name string, content string) (*Stream, error) {
	s.stMutex.Lock()
	// Make sure no duplicate name and find new id for stream
	id := 1
	for i := range s.st.Streams {
		if s.st.Streams[i].Name == name {
			s.stMutex.Unlock()
			return nil, errors.New("duplicate stream name")
		}
		if s.st.Streams[i].SID >= id {
			id = s.st.Streams[i].SID + 1
		}
	}

	// Create stream and add to settings
	st := Stream{SID: id, Name: name, Content: content}
	s.st.Streams = append(s.st.Streams, st)
	s.stMutex.Unlock()
	s.queueWrite()

	return &st, nil
}

func (s *Store) DeleteStream(sid int) int {
	// Delete stream
	deleted := 0
	s.stMutex.Lock()
	newStreams := make([]Stream, 0, len(s.st.Streams))
	for i := range s.st.Streams {
		if s.st.Streams[i].SID != sid {
			newStreams[i] = s.st.Streams[i]
		} else {
			deleted += 1
		}
	}
	s.st.Streams = newStreams

	s.clearStream(sid)
	s.stMutex.Unlock()
	s.queueWrite()

	return deleted
}

// clearStream sets StreamID to 0 for all presets that have a StreamID of sid.
func (s *Store) clearStream(sid int) bool {
	changed := false
	for i := range s.st.Presets {
		if s.st.Presets[i].StreamID == sid {
			s.st.Presets[i].StreamID = 0
			changed = true
		}
	}
	return changed
}

func (s *Store) ClearStream(sid int) bool {
	s.stMutex.Lock()
	ok := s.clearStream(sid)
	s.stMutex.Unlock()
	if ok {
		s.queueWrite()
	}
	return ok
}

// ClearPreset sets preset's sid to 0.
func (s *Store) ClearPreset(uri string) bool {
	s.stMutex.Lock()
	for i := range s.st.Presets {
		if s.st.Presets[i].URI != uri {
			s.st.Presets[i].StreamID = 0
			s.queueWrite()
			s.stMutex.Unlock()
			return true
		}
	}
	s.stMutex.Unlock()
	return false
}

func (s *Store) UpdatePreset(p *Preset) bool {
	if p.StreamID == 0 {
		return s.ClearPreset(p.URI)
	}

	s.stMutex.Lock()
	changed := false
	ok := false
	for i := range s.st.Presets {
		if s.st.Presets[i].URI == p.URI {
			ok = true
			if s.st.Presets[i].StreamID != p.StreamID {
				s.st.Presets[i].StreamID = p.StreamID
				changed = true
			}
		} else if s.st.Presets[i].StreamID == p.StreamID {
			// Clear preset if it is equal to preset stream id
			s.st.Presets[i].StreamID = 0
			changed = true
		}
	}
	s.stMutex.Unlock()
	if changed {
		s.queueWrite()
	}
	return ok
}

func (s *Store) UpdateStream(stream *Stream) bool {
	idx := -1
	s.stMutex.Lock()
	for i := range s.st.Streams {
		if s.st.Streams[i].SID == stream.SID {
			idx = i
		} else if s.st.Streams[i].Name == stream.Name {
			s.stMutex.Unlock()
			return false
		}
	}
	if idx == -1 {
		s.stMutex.Unlock()
		return false
	}
	s.st.Streams[idx] = *stream
	s.stMutex.Unlock()
	s.queueWrite()
	return true
}

func (s *Store) GetStream(id int) (*Stream, bool) {
	s.stMutex.Lock()
	for i := range s.st.Streams {
		if s.st.Streams[i].SID == id {
			newST := s.st.Streams[i]
			s.stMutex.Unlock()
			return &newST, true
		}
	}
	s.stMutex.Unlock()
	return nil, false
}

func (s *Store) GetPreset(uri string) (*Preset, bool) {
	s.stMutex.Lock()
	for i := range s.st.Presets {
		if s.st.Presets[i].URI == uri {
			newP := s.st.Presets[i]
			s.stMutex.Unlock()
			return &newP, true
		}
	}
	s.stMutex.Unlock()
	return nil, false
}

func (s *Store) GetPresets() []Preset {
	s.stMutex.Lock()
	p := make([]Preset, len(s.st.Presets))
	copy(p, s.st.Presets)
	s.stMutex.Unlock()
	return p
}

func (s *Store) GetStreams() []Stream {
	s.stMutex.Lock()
	ss := make([]Stream, len(s.st.Streams))
	copy(ss, s.st.Streams)
	s.stMutex.Unlock()
	return ss
}

// WriteSettings to disk.
func (s *Store) WriteSettings() error {
	errChan := make(chan error)
	select {
	case s.writeChan <- errChan:
		return <-errChan
	case <-s.dctx.Done():
		return s.dctx.Err()
	}
}

// ReadSettings from disk, do not use this function as it may discard current settings that have pending saves.
func (s *Store) ReadSettings() error {
	errChan := make(chan error)
	select {
	case s.readChan <- errChan:
		return <-errChan
	case <-s.dctx.Done():
		return s.dctx.Err()
	}
}
