package store

import (
	"context"
	"encoding/json"
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

	// Init store settings from disk
	if sg, err := s.readSettings(); err == nil {
		s.sg = sg
	} else {
		cfg.EnablePresets = false
		return nil, err
	}

	// Merge presets from cfg to store
	if len(cfg.Presets) > 0 {
		p := make([]Preset, 0, len(cfg.Presets))
		for _, v := range cfg.Presets {
			p = append(p, Preset{URI: v})
		}
		s.sg.mergePresets(p)
	}

	// Merge presets from store to cfg
	if len(s.sg.Presets) > 0 {
		cfg.EnablePresets = true
		var u []string
		for _, v := range s.sg.Presets {
			u = append(u, v.URI)
		}
		cfg.Presets = u
	}

	// Prioritize cfg port over settings port unless it is default port
	if cfg.Port != s.sg.Port && cfg.Port == config.DefaultPort {
		cfg.Port = s.sg.Port
	}

	// Prioritize cfg cport over settings cport unless it is default cport
	if cfg.CPort != s.sg.CPort && cfg.CPort == goupnpsub.DefaultPort {
		cfg.CPort = s.sg.CPort
	}

	// Save and start loop
	s.writeSettings(s.sg)
	go s.storeLoop()

	return s, nil
}

func (s *Store) writeSettings(sg *Settings) error {
	b, err := json.MarshalIndent(sg, "", "	")
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

	sg := Settings{}
	if err = json.Unmarshal(b, &sg); err != nil {
		return nil, err
	}

	return &sg, nil
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
		s.sgMutex.Lock()
		sg := *s.sg
		s.sgMutex.Unlock()
		return s.writeSettings(&sg)
	}
	readSettings := func() error {
		sg, err := s.readSettings()
		if err != nil {
			return err
		}
		s.sgMutex.Lock()
		s.sg = sg
		s.sgMutex.Unlock()
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

func (s *Store) DeleteStream(sid int) int {
	// Delete stream
	deleted := 0
	s.sgMutex.Lock()
	newStreams := make([]Stream, 0, len(s.sg.Streams))
	for i := range s.sg.Streams {
		if s.sg.Streams[i].SID != sid {
			newStreams[i] = s.sg.Streams[i]
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

// clearStream sets SID to 0 for all presets that have a SID of sid.
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

func (s *Store) ClearStream(sid int) bool {
	s.sgMutex.Lock()
	ok := s.clearStream(sid)
	s.sgMutex.Unlock()
	if ok {
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

func (s *Store) UpdatePreset(pt *Preset) bool {
	if pt.SID == 0 {
		return s.ClearPreset(pt.URI)
	}

	s.sgMutex.Lock()
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

func (s *Store) GetStream(id int) (*Stream, bool) {
	s.sgMutex.Lock()
	for i := range s.sg.Streams {
		if s.sg.Streams[i].SID == id {
			st := s.sg.Streams[i]
			s.sgMutex.Unlock()
			return &st, true
		}
	}
	s.sgMutex.Unlock()
	return nil, false
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

func (s *Store) GetPresets() []Preset {
	s.sgMutex.Lock()
	pt := make([]Preset, len(s.sg.Presets))
	copy(pt, s.sg.Presets)
	s.sgMutex.Unlock()
	return pt
}

func (s *Store) GetStreams() []Stream {
	s.sgMutex.Lock()
	st := make([]Stream, len(s.sg.Streams))
	copy(st, s.sg.Streams)
	s.sgMutex.Unlock()
	return st
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
