package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
)

func TestStore(t *testing.T) {
	s := testStore(t)

	testStream(t, s)

	// Test canceling context right after writing settings
	s.WriteSettings()
	s.Cancel()
}

func testStore(t *testing.T) *Store {
	// Create temp directory
	p, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	cfg := config.NewConfig()
	cfg.ConfigPath = filepath.Join(p, "settings.json")

	// Create store
	s, err := NewService(cfg)
	if err != nil {
		t.Error(err)
	}

	// Settings file exists
	if _, err := os.Stat(s.file); err != nil {
		t.Error(err)
	}

	return s
}

func testStream(t *testing.T, s *Store) {
	// Test adding a stream and getting the stream back
	addST, err := s.AddStream("Name", "Content")
	if err != nil {
		t.Error("could not add stream,", err)
	}
	if _, err = s.AddStream("Name", "Content"); err == nil {
		t.Error("duplicate add stream name")
	}
	testGetStream := func() {
		if getStream, ok := s.GetStream(addST.ID); !ok {
			t.Error("stream should exist")
		} else {
			if !reflect.DeepEqual(*addST, *getStream) {
				t.Error("saved stream is not equal", addST, *getStream)
			}
		}
	}
	testGetStream()

	// Test getting the stream back after writing and reading
	s.WriteSettings()
	s.ReadSettings()
	testGetStream()

	// Test deleting stream
	if c := s.DeleteStream(addST.ID); c != 1 {
		t.Error("deleting stream should return 1, got", c)
	}
	if _, ok := s.GetStream(addST.ID); ok {
		t.Error("stream should be deleted")
	}
}
