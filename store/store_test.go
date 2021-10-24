package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
)

func TestStore(t *testing.T) {
	s := testStore(t)

	// Test canceling context right after writing settings
	s.WriteSettings()
	s.Cancel()
}

func testStore(t *testing.T) *Store {
	// Temp directory should be created
	p, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	cfg := config.NewConfig()
	cfg.ConfigPath = filepath.Join(p, "settings.json")

	// Store should be created
	s, err := NewService(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Settings file should exists
	if _, err := os.Stat(s.file); err != nil {
		t.Fatal(err)
	}

	return s
}
