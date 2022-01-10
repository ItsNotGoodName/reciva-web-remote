package store

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
)

type Store struct {
	DoneC          chan struct{}
	PresetChangedC chan struct{}
	configFile     string
	op             chan func(map[string]core.Preset)
	presets        map[string]core.Preset
	readonly       bool
}

func New(configFile string) *Store {
	var (
		readonly bool
		presets  map[string]core.Preset
		err      error
	)

	if presets, err = readPresets(configFile); err != nil {
		log.Println("store.New(WARNING): store is readonly:", err)
		presets = make(map[string]core.Preset)
		readonly = true
	}

	return &Store{
		DoneC:          make(chan struct{}),
		PresetChangedC: make(chan struct{}, 1),
		configFile:     configFile,
		op:             make(chan func(map[string]core.Preset)),
		presets:        presets,
		readonly:       readonly,
	}
}

func (s *Store) PresetChanged() <-chan struct{} {
	return s.PresetChangedC
}

func (s *Store) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(s.DoneC)
			return
		case op := <-s.op:
			op(s.presets)
		}
	}
}

func (s *Store) ListPresets(ctx context.Context) ([]core.Preset, error) {
	var presets []core.Preset
	errC := make(chan error)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case s.op <- func(m map[string]core.Preset) {
		for _, p := range m {
			presets = append(presets, p)
		}

		errC <- nil
	}:
	}

	return presets, <-errC
}

func (s *Store) GetPreset(ctx context.Context, url string) (*core.Preset, error) {
	var preset *core.Preset
	errC := make(chan error)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case s.op <- func(m map[string]core.Preset) {
		p, ok := m[url]
		if !ok {
			errC <- core.ErrPresetNotFound
			return
		}

		preset = &p
		errC <- nil
	}:
	}

	return preset, <-errC
}

func (s *Store) UpdatePreset(ctx context.Context, preset *core.Preset) error {
	if s.readonly {
		return core.ErrConfigReadonly
	}

	errC := make(chan error)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.op <- func(m map[string]core.Preset) {
		_, ok := m[preset.URL]
		if !ok {
			errC <- core.ErrPresetNotFound
			return
		}

		m[preset.URL] = *preset

		select {
		case s.PresetChangedC <- struct{}{}:
		default:
		}

		errC <- saveConfig(s.configFile, m)
	}:
	}

	return <-errC
}
