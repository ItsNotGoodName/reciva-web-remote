package dto

import (
	"context"
)

type App interface {
	PresetGet(context.Context, *PresetGetRequest) (*PresetGetResponse, error)
	PresetList(context.Context) (*PresetListResponse, error)
	PresetUpdate(context.Context, *PresetUpdateRequest) error
	RadioDiscover(*RadioDiscoverRequest) (*RadioDiscoverResponse, error)
	RadioGet(*RadioRequest) (*RadioGetResponse, error)
	RadioList() (*RadioListResponse, error)
	RadioRefreshSubscription(context.Context, *RadioRequest) error
	RadioRefreshVolume(context.Context, *RadioRequest) error
	StateGet(context.Context, *StateRequest) (*StateGetResponse, error)
	StateList(context.Context) (*StateListResponse, error)
	StatePatch(context.Context, *StatePatchRequest) error
	Build() Build
}
