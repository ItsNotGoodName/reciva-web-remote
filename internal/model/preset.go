package model

import (
	"fmt"
	"net/url"
)

type Preset struct {
	TitleNew string
	URL      string
	URLNew   string
}

func NewPreset(urL, titleNew, urlNew string) (*Preset, error) {
	urlObj, err := url.ParseRequestURI(urL)
	if err != nil {
		return nil, err
	}

	if urlObj.Scheme == "" || urlObj.Host == "" {
		return nil, fmt.Errorf("invalid URL: %s", urL)
	}

	return &Preset{
		URL:      urL,
		TitleNew: titleNew,
		URLNew:   urlNew,
	}, nil
}
