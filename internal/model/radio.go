package model

import "github.com/ItsNotGoodName/reciva-web-remote/internal/hub"

type Radio struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func NewRadio(radio hub.Radio) Radio {
	return Radio{
		UUID: radio.UUID,
		Name: radio.Name,
	}
}
