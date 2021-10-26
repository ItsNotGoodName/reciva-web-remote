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
		Cancel:         cancel,
		dctx:           dctx,
		file:           cfg.ConfigPath,
		queueWriteChan: make(chan bool),
		readChan:       make(chan chan error),
		writeChan:      make(chan chan error),
	}

	// Read settings from disk
	if sg, err := s.readSettings(); err != nil {
		cfg.PresetsEnabled = false
		return nil, err
	} else {
		s.sg = sg
	}

	// Sync settings with config
	s.config(cfg)

	// Write settings to disk
	if err := s.writeSettings(s.sg); err != nil {
		log.Fatal(err)
	}

	// Start loop
	go s.storeLoop()

	return s, nil
}

func (s *Store) config(cfg *config.Config) {
	// port
	if cfg.PortFlag {
		s.sg.Port = cfg.Port
	} else {
		cfg.Port = s.sg.Port
	}

	// cport
	if cfg.CPortFlag {
		s.sg.CPort = cfg.CPort
	} else {
		cfg.CPort = s.sg.CPort
	}

	// presets
	for _, u := range cfg.Presets {
		if i := getPresetIndex(s.sg.Presets, u); i < 0 {
			s.sg.Presets = append(s.sg.Presets, Preset{URI: u})
		}
	}
	cfg.Presets = make([]string, 0, len(s.sg.Presets))
	for _, p := range s.sg.Presets {
		cfg.Presets = append(cfg.Presets, p.URI)
	}
	if len(cfg.Presets) > 0 {
		cfg.PresetsEnabled = true
	}
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

// queueWrite tells store to that settings have changed and it should write settings to disk at some point.
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

func NewSettings() *Settings {
	return &Settings{
		Port:    config.DefaultPort,
		CPort:   goupnpsub.DefaultPort,
		Streams: make([]Stream, 0),
		Presets: make([]Preset, 0),
	}
}

// ReadSettings from disk, do not use this function as it may discard current settings that have pending writes.
func (s *Store) ReadSettings() error {
	errChan := make(chan error)
	select {
	case s.readChan <- errChan:
		return <-errChan
	case <-s.dctx.Done():
		return s.dctx.Err()
	}
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
