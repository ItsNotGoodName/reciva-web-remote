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

	// Test canceling context right after writing settings
	s.WriteSettings()
	s.Cancel()

	// Compare saved settings to current settings
	sg, err := s.readSettings()
	if err != nil {
		t.Fatalf("TestStore(err) = %s, want nil", err)
	}
	if !reflect.DeepEqual(*sg, *s.sg) {
		t.Fatalf("TestStore(sg) = %+v, want %+v", *sg, *s.sg)
	}
}

func testStore(t *testing.T) *Store {
	// Temp directory should be created
	p, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	cfg := config.Config{
		Port:          0,
		CPort:         0,
		ConfigPath:    filepath.Join(p, "settings.json"),
		EnablePresets: true,
		Presets:       []string{"/01.m3u", "/02.m3u"},
		APIURI:        "",
	}

	// Store should be created
	s, err := NewService(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Settings file should exists
	if _, err := os.Stat(s.file); err != nil {
		t.Fatal(err)
	}

	return s
}
