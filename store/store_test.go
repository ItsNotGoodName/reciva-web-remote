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

	// Add stream
	stream, err := s.AddStream("Name", "Content")
	if err != nil {
		t.Error("could not create stream,", err)
	}

	// Get stream
	if getStream, ok := s.GetStream(stream.ID); !ok {
		t.Error("stream should exists")
	} else {
		if !reflect.DeepEqual(*stream, *getStream) {
			t.Error("saved stream is not equal", stream, *getStream)
		}
	}

	// Delete stream
	if c := s.DeleteStream(stream.ID); c != 1 {
		t.Error("streams deleted should be 1, got", c)
	}
	if _, ok := s.GetStream(stream.ID); ok {
		t.Error("stream should not exists")
	}

	s.WriteSettings()

	// Stop store
	s.Cancel()
}
