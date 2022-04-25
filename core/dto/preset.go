package dto

import (
	"github.com/ItsNotGoodName/reciva-web-remote/core/preset"
)

type Preset struct {
	TitleNew string `json:"title_new"`
	URL      string `json:"url"`
	URLNew   string `json:"url_new"`
}

func ConvertPreset(p *Preset) (*preset.Preset, error) {
	return preset.ParsePreset(p.URL, p.TitleNew, p.URLNew)
}

func NewPreset(p *preset.Preset) Preset {
	return Preset{
		TitleNew: p.TitleNew,
		URL:      p.URL,
		URLNew:   p.URLNew,
	}
}

func NewPresets(ps []preset.Preset) []Preset {
	presets := make([]Preset, len(ps))
	for i, p := range ps {
		presets[i] = Preset{TitleNew: p.TitleNew, URL: p.URL, URLNew: p.URLNew}
	}

	return presets
}
