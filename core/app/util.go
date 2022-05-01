package app

import (
	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
	"github.com/ItsNotGoodName/reciva-web-remote/core/preset"
	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
)

func newPreset(p *preset.Preset) dto.Preset {
	return dto.Preset{
		TitleNew: p.TitleNew,
		URL:      p.URL,
		URLNew:   p.URLNew,
	}
}

func newPresets(ps []preset.Preset) []dto.Preset {
	presets := make([]dto.Preset, len(ps))
	for i, p := range ps {
		presets[i] = dto.Preset{TitleNew: p.TitleNew, URL: p.URL, URLNew: p.URLNew}
	}

	return presets
}

func newRadio(r *radio.Radio) dto.Radio {
	return dto.Radio{UUID: r.UUID, Name: r.Name}
}

func newRadios(rds []radio.Radio) []dto.Radio {
	radios := make([]dto.Radio, len(rds))
	for i, r := range rds {
		radios[i] = newRadio(&r)
	}

	return radios
}
