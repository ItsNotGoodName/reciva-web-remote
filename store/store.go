package store

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
)

func NewService(cfg *config.Config) (*Store, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	s := &Store{
		ctx:               ctx,
		Cancel:            cancel,
		file:              cfg.ConfigPath,
		getSettingsChan:   make(chan Settings),
		writeSettingsChan: make(chan chan error),
		readSettingsChan:  make(chan chan error),
		deleteStreamChan:  make(chan int),
		updateStreamChan:  make(chan Stream),
		updatePresetChan:  make(chan Preset),
		setPresetsChan:    make(chan []Preset),
	}

	st, err := s.readSettings()
	if err != nil {
		cfg.EnablePresets = false
		return nil, err
	}

	if len(cfg.Presets) > 0 {
		p := make([]Preset, len(cfg.Presets))
		for _, v := range cfg.Presets {
			p = append(p, Preset{URI: v})
		}
		st.mergePresets(p)
	}

	if len(st.Presets) > 0 {
		cfg.EnablePresets = true
		var u []string
		for _, v := range st.Presets {
			u = append(u, v.URI)
		}
		cfg.Presets = u
	}

	go s.storeLoop(st)

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

func (s *Store) storeLoop(st *Settings) {
	ticker := time.NewTicker(15 * time.Second)
	shouldSave := false

	for {
		select {
		case s.getSettingsChan <- *st:
		case res := <-s.writeSettingsChan:
			err := s.writeSettings(st)
			res <- err
			shouldSave = false
		case res := <-s.readSettingsChan:
			var err error
			st, err = s.readSettings()
			res <- err
			shouldSave = false
		case id := <-s.deleteStreamChan:
			newStreams := make([]Stream, len(st.Streams))
			for i := range st.Streams {
				if st.Streams[i].ID != id {
					newStreams[i] = st.Streams[i]
				}
			}
			shouldSave = true
		case newStream := <-s.updateStreamChan:
			idx := -1
			for i := range st.Streams {
				if newStream.ID == st.Streams[i].ID {
					idx = i
					break
				}
			}
			if idx >= 0 {
				st.Streams[idx] = newStream
			} else {
				st.Streams = append(st.Streams, newStream)
			}
			shouldSave = true
		case newPreset := <-s.updatePresetChan:
			idx := -1
			for i := range st.Presets {
				if newPreset.URI == st.Presets[i].URI {
					idx = i
					break
				}
			}
			if idx >= 0 {
				st.Presets[idx] = newPreset
				shouldSave = true
			}
		case st.Presets = <-s.setPresetsChan:
		case <-ticker.C:
			if shouldSave {
				log.Println("Store.storeLoop: settings saved")
				s.writeSettings(st)
				shouldSave = false
			}
		case <-s.ctx.Done():
		}
	}
}

func (s *Store) GetSettings() Settings {
	return <-s.getSettingsChan
}

func (s *Store) WriteSettings() error {
	errChan := make(chan error)
	s.writeSettingsChan <- errChan
	return <-errChan
}

func (s *Store) ReadSettings() error {
	errChan := make(chan error)
	s.readSettingsChan <- errChan
	return <-errChan
}

func (s *Store) DeleteStream(id int) {
	s.deleteStreamChan <- id
}

func (s *Store) UpdateStream(stream *Stream) {
	s.updateStreamChan <- *stream
}

func (s *Store) UpdatePreset(preset *Preset) {
	s.updatePresetChan <- *preset
}

func (s *Store) SetPresets(presets []Preset) {
	s.setPresetsChan <- presets
}
