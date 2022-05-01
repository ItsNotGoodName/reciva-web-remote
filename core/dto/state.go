package dto

import "github.com/ItsNotGoodName/reciva-web-remote/core/state"

type StateRequest struct {
	UUID string `json:"uuid"`
}

type StateGetResponse struct {
	State state.State `json:"state"`
}

type StateListResponse struct {
	States []state.State `json:"states"`
}

type StatePatchRequest struct {
	UUID        string  `json:"uuid"`
	Power       *bool   `json:"power,omitempty"`
	AudioSource *string `json:"audio_source,omitempty"`
	Preset      *int    `json:"preset,omitempty"`
	Volume      *int    `json:"volume,omitempty"`
}
