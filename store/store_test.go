package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
)

var PresetURIs = []string{"/01.m3u", "/02.m3u", "/03.m3u"}

func TestStore(t *testing.T) {
	s := testStore(t)

	// Test write and read settings
	if err := s.WriteSettings(); err != nil {
		t.Fatal("TestStore(WriteSettings)", err)
	}
	if err := s.ReadSettings(); err != nil {
		t.Fatal("TestStore(ReadSettings)", err)
	}
	s.Cancel()

	testReadSettings(t, "TestStore", s)
}

func TestReadSettings(t *testing.T) {
	testReadSettings(t, "TestReadSettings", testStore(t))
}

func testReadSettings(t *testing.T, name string, s *Store) {
	// Compare disk settings to current settings
	if sg, err := s.readSettings(); err != nil {
		t.Fatalf("%s(err) = %s, want nil", name, err)
	} else if !reflect.DeepEqual(*sg, *s.sg) {
		t.Fatalf("%s(sg) = %+v, want %+v", name, *sg, *s.sg)
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
		Presets:       PresetURIs,
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
