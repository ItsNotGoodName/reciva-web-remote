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

	stream := Stream{ID: 23, Name: "Name", Content: "Content"}

	// Add stream
	s.UpdateStream(&stream)
	if getStream := s.GetStream(stream.ID); getStream == nil {
		t.Error("stream should not be nil")
	} else {
		if !reflect.DeepEqual(stream, *getStream) {
			t.Error("saved stream is not equal", stream, *getStream)
		}
	}

	// Delete stream
	if c := s.DeleteStream(stream.ID); c != 1 {
		t.Error("streams deleted should be 1, got", c)
	}
	if stream := s.GetStream(stream.ID); stream != nil {
		t.Error("stream should be nil")
	}

	s.WriteSettings()

	// Stop store
	s.Cancel()
}
