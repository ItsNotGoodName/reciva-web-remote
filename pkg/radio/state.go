package radio

import (
	"context"
	"fmt"

	"github.com/avast/retry-go"
	"github.com/huin/goupnp"
)

// TODO: json camelCase -> snake_case

type State struct {
	IsMuted    *bool    `json:"isMuted,omitempty"`    // IsMuted represents if the Radio's volume is muted.
	Metadata   *string  `json:"metadata,omitempty"`   // Metadata that is received from the stream url.
	Name       string   `json:"name,omitempty"`       // Name of the radio.
	NumPresets int      `json:"numPresets,omitempty"` // NumPresets on the radio, it will not change after it is set.
	Power      *bool    `json:"power,omitempty"`      // Power represents if the radio is not in standby.
	Preset     int      `json:"preset,omitempty"`     // Preset is the current preset that is playing, -1 means it is unknown.
	Presets    []Preset `json:"presets,omitempty"`    // Presets on the radio, it's length should not change after it is set.
	State      string   `json:"state,omitempty"`      // State is either playing, connecting, or stopped.
	Title      string   `json:"title,omitempty"`      // Title of the preset.
	URL        string   `json:"url,omitempty"`        // URL of the stream that is currently selected.
	UUID       string   `json:"uuid"`                 // UUID will not change once it has been set.
	Volume     *int     `json:"volume,omitempty"`     // Volume of the Radio.
}

type Preset struct {
	Name   string `json:"name"`   // Name of the preset that overrides Title.
	Number int    `json:"number"` //
	Title  string `json:"-"`      // Title given by the radio.
	URL    string `json:"url"`    //
}

type stateXML struct {
	Metadata string `xml:"playback-details>stream>metadata"`
	State    string `xml:"playback-details>state"`
	Title    string `xml:"playback-details>stream>title"`
	URL      string `xml:"playback-details>stream>url"`
}

// newState creates a State.
func newState(uuid, name string) *State {
	boolDefault := false
	intDefault := 0
	stringDefault := ""
	return &State{
		IsMuted:  &boolDefault,
		Metadata: &stringDefault,
		Name:     name,
		Power:    &boolDefault,
		Preset:   -1,
		UUID:     uuid,
		Volume:   &intDefault,
	}
}

// newStateFromClient creates a State from the given ServiceClient.
func newStateFromClient(ctx context.Context, client goupnp.ServiceClient) (*State, error) {
	uuid, err := getServiceClientUUID(client)
	if err != nil {
		return nil, err
	}

	state := newState(uuid, client.RootDevice.Device.FriendlyName)

	// Get number of presets
	var numPresets int
	if err := retry.Do(func() error {
		if p, e := getNumberOfPresets(ctx, client); e != nil {
			return e
		} else {
			numPresets = p
			return nil
		}
	}, retry.Context(ctx)); err != nil {
		return nil, err
	} else {
		numPresets = numPresets - 2
		if numPresets < 1 {
			return nil, fmt.Errorf("invalid number of presets from radio: %d", numPresets)
		} else {
			state.NumPresets = numPresets
		}
	}

	// Get presets
	var presets []Preset
	if err := retry.Do(func() error {
		if p, e := getPresets(ctx, client, numPresets); e != nil {
			return e
		} else {
			presets = p
			return nil
		}
	}, retry.Context(ctx)); err != nil {
		return nil, err
	} else {
		// TODO: Mutate preset somewhere else
		// for i := range presets {
		// 	 rd.h.PresetMutator(rd.ctx, &presets[i])
		// }
		state.Presets = presets
	}

	// Get volume
	var volume int
	if err := retry.Do(func() error {
		if v, e := getVolume(ctx, client); e != nil {
			return e
		} else {
			volume = v
			return nil
		}
	}, retry.Context(ctx)); err != nil {
		return nil, err
	} else {
		state.Volume = &volume
	}

	return state, nil
}

func (s *State) Merge(ss *State) {
	if ss.IsMuted != nil {
		s.IsMuted = ss.IsMuted
	}
	if ss.Metadata != nil {
		s.Metadata = ss.Metadata
	}
	if ss.Name != "" {
		s.Name = ss.Name
	}
	if ss.NumPresets != 0 {
		s.NumPresets = ss.NumPresets
	}
	if ss.Power != nil {
		s.Power = ss.Power
	}
	if len(ss.Presets) != 0 {
		s.Presets = ss.Presets
	}
	if ss.Preset != 0 {
		s.Preset = ss.Preset
	}
	if ss.State != "" {
		s.State = ss.State
	}
	if ss.Title != "" {
		s.Title = ss.Title
	}
	if ss.URL != "" {
		s.URL = ss.URL
	}
	if ss.Volume != nil {
		s.Volume = ss.Volume
	}
}

func (s *State) ValidPreset(preset int) error {
	if preset < 1 || preset > s.NumPresets {
		return ErrPresetInvalid
	}

	return nil
}

func (s *State) On() bool {
	return s.Power != nil && *s.Power
}
