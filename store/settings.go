package store

import (
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
)

func NewSettings() *Settings {
	return &Settings{
		Port:    config.DefaultPort,
		CPort:   goupnpsub.DefaultPort,
		Streams: make([]Stream, 0),
		Presets: make([]Preset, 0),
	}
}
