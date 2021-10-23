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

	if cfg.Port != s.st.Port && cfg.Port == config.DefaultPort {
		cfg.Port = s.st.Port
	}

	if cfg.CPort != s.st.CPort && cfg.CPort == goupnpsub.DefaultPort {
		cfg.CPort = s.st.CPort
	}

	go s.storeLoop()

	return s, nil
}

func (s *Store) writeSettings(st *Settings) error {
	b, err := json.MarshalIndent(st, "", "")
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

func (s *Store) GetSettings() *Settings {
	s.stMutex.Lock()
	st := *s.st
	s.stMutex.Unlock()
	return &st
}

func (s *Store) WriteSettings() error {
	errChan := make(chan error)
	select {
	case s.writeChan <- errChan:
		return <-errChan
	case <-s.dctx.Done():
		return s.dctx.Err()
	}
}

func (s *Store) ReadSettings() error {
	errChan := make(chan error)
	select {
	case s.readChan <- errChan:
		return <-errChan
	case <-s.dctx.Done():
		return s.dctx.Err()
	}
}

func (s *Store) DeleteStream(id int) int {
	deleted := 0
	s.stMutex.Lock()
	newStreams := make([]Stream, 0, len(s.st.Streams))
	for i := range s.st.Streams {
		if s.st.Streams[i].ID != id {
			newStreams[i] = s.st.Streams[i]
		} else {
			deleted += 1
		}
	}
	s.st.Streams = newStreams
	s.stMutex.Unlock()
	s.queueWrite()
	return deleted
}

func (s *Store) GetStream(id int) *Stream {
	s.stMutex.Lock()
	for i := range s.st.Streams {
		if s.st.Streams[i].ID == id {
			s.stMutex.Unlock()
			return &s.st.Streams[i]
		}
	}
	s.stMutex.Unlock()
	return nil
}

func (s *Store) UpdateStream(stream *Stream) {
	s.stMutex.Lock()
	for i := range s.st.Streams {
		if stream.ID == s.st.Streams[i].ID {
			s.st.Streams[i] = *stream
			s.stMutex.Unlock()
			s.queueWrite()
			return
		}
	}
	s.st.Streams = append(s.st.Streams, *stream)
	s.stMutex.Unlock()
	s.queueWrite()
}

func (s *Store) GetPreset(uri string) *Preset {
	s.stMutex.Lock()
	for i := range s.st.Presets {
		if s.st.Presets[i].URI == uri {
			s.stMutex.Unlock()
			return &s.st.Presets[i]
		}
	}
	s.stMutex.Unlock()
	return nil
}

func (s *Store) UpdatePreset(newPreset *Preset) {
	s.stMutex.Lock()
	for i := range s.st.Presets {
		if newPreset.URI == s.st.Presets[i].URI {
			s.st.Presets[i] = *newPreset
			s.stMutex.Unlock()
			s.queueWrite()
			return
		}
	}
	s.stMutex.Unlock()
}
