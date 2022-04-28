package dto

import "github.com/ItsNotGoodName/reciva-web-remote/core/state"

type SlimState struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

func NewSlimState(s *state.State) SlimState {
	return SlimState{Name: s.Name, UUID: s.UUID}
}

func NewSlimStates(ss []state.State) []SlimState {
	states := make([]SlimState, len(ss))
	for i, s := range ss {
		states[i] = NewSlimState(&s)
	}

	return states
}
