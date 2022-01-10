package core

import (
	"context"
	"fmt"
)

var (
	ErrPresetNotFound = fmt.Errorf("preset not found")
)

type (
	Preset struct {
		URL     string `json:"url"`
		NewName string `json:"newName"`
		NewURL  string `json:"newUrl"`
	}

	PresetStore interface {
		ListPresets(ctx context.Context) ([]Preset, error)
		GetPreset(ctx context.Context, url string) (*Preset, error)
		UpdatePreset(ctx context.Context, preset *Preset) error
		PresetChanged() <-chan struct{}
	}
)
