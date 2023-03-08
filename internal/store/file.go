package store

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sort"
	"sync"

	"github.com/ItsNotGoodName/reciva-web-remote/internal"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/model"
)

type FileConfig struct {
	Presets []FilePreset `json:"presets"`
}

type FilePreset struct {
	URL      string `json:"url"`
	TitleNew string `json:"newName"`
	URLNew   string `json:"newUrl"`
}

type File struct {
	file string

	presetsMapMu sync.RWMutex
	presetsMap   map[string]model.Preset
}

func NewFile(file string) (*File, error) {
	presetsMap, err := readConfig(file)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}

		presetsMap = make(map[string]model.Preset)
	}

	return &File{
		file:       file,
		presetsMap: presetsMap,
	}, nil
}

func (f *File) ListPresets(ctx context.Context) ([]model.Preset, error) {
	f.presetsMapMu.RLock()
	pts := make([]model.Preset, 0, len(f.presetsMap))
	for _, preset := range f.presetsMap {
		pts = append(pts, preset)
	}
	f.presetsMapMu.RUnlock()

	// Sort presets
	sort.Slice(pts, func(i, j int) bool {
		return pts[i].URL < pts[j].URL
	})

	return pts, nil
}

func (f *File) GetPreset(ctx context.Context, url string) (*model.Preset, error) {
	f.presetsMapMu.RLock()
	p, ok := f.presetsMap[url]
	f.presetsMapMu.RUnlock()

	if !ok {
		return nil, internal.ErrPresetNotFound
	}

	return &p, nil
}

func (f *File) UpdatePreset(ctx context.Context, p *model.Preset) error {
	f.presetsMapMu.Lock()
	defer f.presetsMapMu.Unlock()

	old, ok := f.presetsMap[p.URL]
	if !ok {
		return internal.ErrPresetNotFound
	}

	f.presetsMap[p.URL] = *p

	err := writeConfig(f.file, f.presetsMap)
	if err != nil {
		f.presetsMap[p.URL] = old
		return err
	}

	return nil
}

func readConfig(file string) (map[string]model.Preset, error) {
	// Read config file
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal cfg file
	cfg := FileConfig{}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	// Create map from config struct
	presets := make(map[string]model.Preset, len(cfg.Presets))
	for _, p := range cfg.Presets {
		pp, err := model.NewPreset(p.URL, p.TitleNew, p.URLNew)
		if err != nil {
			return nil, err
		}

		presets[p.URL] = *pp
	}

	return presets, nil
}

func writeConfig(file string, m map[string]model.Preset) error {
	// Create presets slice from map
	presets := make([]FilePreset, 0, len(m))
	for _, p := range m {
		presets = append(presets, FilePreset{
			URL:      p.URL,
			TitleNew: p.TitleNew,
			URLNew:   p.URLNew,
		})
	}

	// Sort presets
	sort.Slice(presets, func(i, j int) bool {
		return presets[i].URL < presets[j].URL
	})

	// Create config struct
	cfg := FileConfig{
		Presets: presets,
	}

	// Marshal config struct
	b, err := json.MarshalIndent(cfg, "", "	")
	if err != nil {
		return err
	}

	// Write cfg file
	return os.WriteFile(file, b, 0644)
}
