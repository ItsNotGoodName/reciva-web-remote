package file

import (
	"context"
	"sync"

	"github.com/ItsNotGoodName/reciva-web-remote/core/preset"
)

type PresetStore struct {
	file string

	presetsMapMu sync.RWMutex
	presetsMap   map[string]preset.Preset
}

func NewPresetStore(file string) (*PresetStore, error) {
	presetsMap, err := readConfig(file)
	if err != nil {
		return nil, err
	}

	return &PresetStore{
		file:       file,
		presetsMap: presetsMap,
	}, nil
}

func (ps *PresetStore) List(ctx context.Context) ([]preset.Preset, error) {
	ps.presetsMapMu.RLock()
	pts := make([]preset.Preset, 0, len(ps.presetsMap))
	for _, preset := range ps.presetsMap {
		pts = append(pts, preset)
	}
	ps.presetsMapMu.RUnlock()

	return pts, nil
}

func (ps *PresetStore) Get(ctx context.Context, url string) (*preset.Preset, error) {
	ps.presetsMapMu.RLock()
	p, ok := ps.presetsMap[url]
	ps.presetsMapMu.RUnlock()

	if !ok {
		return nil, preset.ErrPresetNotFound
	}

	return &p, nil
}

func (ps *PresetStore) Update(ctx context.Context, p *preset.Preset) error {
	ps.presetsMapMu.Lock()
	defer ps.presetsMapMu.Unlock()

	old, ok := ps.presetsMap[p.URL]
	if !ok {
		return preset.ErrPresetNotFound
	}

	ps.presetsMap[p.URL] = *p

	err := writeConfig(ps.file, ps.presetsMap)
	if err != nil {
		ps.presetsMap[p.URL] = old
		return err
	}

	return nil
}
