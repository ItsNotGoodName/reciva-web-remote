package store

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestStore(t *testing.T) {
	// Create temp directory
	p, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}

	// Create store
	s, err := NewServiceWithPath(p)
	if err != nil {
		t.Error(err)
	}

	// Settings file exists
	if _, err := os.Stat(s.file); err != nil {
		t.Error(err)
	}

	testReadSettings(t, s)

	// Add stream
	s.updateStreamChan <- Stream{ID: 23, Name: "Name", Content: "Content"}

	// Write to disk
	errChan := make(chan error)
	s.writeSettingsChan <- errChan
	<-errChan

	testReadSettings(t, s)

	// Stop store
	s.Cancel()
}

func testReadSettings(t *testing.T, s *Store) {
	// Get settings from loop
	loopSettings := <-s.getSettingsChan

	// Read settings
	if readSettings, err := s.readSettings(); err != nil {
		t.Error(err)
	} else {
		if !reflect.DeepEqual(loopSettings, *readSettings) {
			t.Error("current settings are not equal to settings on disk, ", loopSettings, *readSettings)
		}
	}
}
