package model

import "github.com/ItsNotGoodName/reciva-web-remote/internal/hub"

type Radio struct {
	UUID string `json:"uuid" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func NewRadio(radio hub.Radio) Radio {
	return Radio{
		UUID: radio.UUID,
		Name: radio.Name,
	}
}
