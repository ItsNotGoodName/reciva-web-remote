package radio

import (
	"context"
	"encoding/xml"
	"log"

	"github.com/avast/retry-go"
)

func (rd *Radio) radioLoop() {
	log.Println("Radio.radioLoop: started")

	rd.initState()

	for {
		select {
		case <-rd.dctx.Done():
			log.Println("Radio.radioLoop: dctx is done, exiting")
			return
		case rd.getStateChan <- *rd.state:
		case newVolume := <-rd.updateVolumeChan:
			// Volume change
			if *rd.state.Volume != newVolume {
				rd.state.Volume = &newVolume
				rd.emitState(&State{Volume: &newVolume})
			}
		case newEvent := <-rd.Subscription.EventChan:
			newState := State{}
			changed := false

			for _, v := range newEvent.Properties {
				if v.Name == "PowerState" {
					newPower := v.Value == "On"

					// Power change
					if newPower != *rd.state.Power {
						rd.state.Power = &newPower
						newState.Power = &newPower
						changed = true
					}
				} else if v.Name == "PlaybackXML" {
					if v.Value == "" {
						continue
					}

					sXML := stateXML{}
					if err := xml.Unmarshal([]byte(v.Value), &sXML); err != nil {
						log.Println("Radio.radioLoop:", err)
						continue
					}

					// State change
					if sXML.State != rd.state.State {
						rd.state.State = sXML.State
						newState.State = sXML.State
						changed = true
					}

					// Title change
					if sXML.Title != rd.state.Title {
						rd.state.Title = sXML.Title
						newState.Title = sXML.Title
						changed = true
					}

					// Url change
					if sXML.Url != rd.state.Url {
						rd.state.Url = sXML.Url
						newState.Url = sXML.Url
						changed = true
					}

					// Metadata change
					if sXML.Metadata != *rd.state.Metadata {
						newMetadata := sXML.Metadata
						rd.state.Metadata = &newMetadata
						newState.Metadata = &newMetadata
						changed = true
					}
				} else if v.Name == "IsMuted" {
					newIsMuted := v.Value == "TRUE"

					if newIsMuted != *rd.state.IsMuted {
						rd.state.IsMuted = &newIsMuted
						newState.IsMuted = &newIsMuted
						changed = true
					}
				}
			}
			if changed {
				rd.emitState(&newState)
			}
		}
	}
}

func (rd *Radio) emitState(state *State) {
	state.UUID = rd.state.UUID
	rd.allStateChan <- *state
}

func (rd *Radio) initState() {
	// Set name of radio
	rd.state.Name = rd.Client.RootDevice.Device.FriendlyName

	// Get number of presets
	var numPresets int
	if err := retry.Do(func() error {
		if p, e := rd.GetNumberOfPresets(rd.dctx); e != nil {
			return e
		} else {
			numPresets = p
			return nil
		}
	}, retry.Context(rd.dctx)); err != nil {
		log.Println("Radio.initState:", err)
	} else {
		numPresets = numPresets - 2
		if numPresets < 1 {
			log.Println("Radio.initState(ERROR): invalid number of presets were given from radio,", numPresets)
		} else {
			rd.state.NumPresets = numPresets
		}
	}

	// Get volume
	var volume int
	if err := retry.Do(func() error {
		if v, e := rd.GetVolume(rd.dctx); e != nil {
			return e
		} else {
			volume = v
			return nil
		}
	}, retry.Context(rd.dctx)); err != nil {
		log.Println("Radio.initState:", err)
	} else {
		if !IsValidVolume(volume) {
			log.Println("Radio.initState(ERROR): invalid volume was given from radio,", volume)
		} else {
			rd.state.Volume = &volume
		}
	}

	// Get presets
	var presets []Preset
	if err := retry.Do(func() error {
		if p, e := rd.GetPresets(rd.dctx); e != nil {
			return e
		} else {
			presets = p
			return nil
		}
	}, retry.Context(rd.dctx)); err != nil {
		log.Println("Radio.initState:", err)
	} else {
		rd.state.Presets = presets
	}
}

func (rd *Radio) RefreshVolume(ctx context.Context) error {
	// Get volume
	volume, err := rd.GetVolume(ctx)
	if err != nil {
		return err
	}

	// Update volume
	if err = rd.UpdateVolume(volume); err != nil {
		return err
	}

	return nil
}

func (rd *Radio) UpdateVolume(volume int) error {
	select {
	case <-rd.dctx.Done():
		return rd.dctx.Err()
	case rd.updateVolumeChan <- volume:
		return nil
	}
}

func (rd *Radio) RefreshPresets(ctx context.Context) error {
	// Get presets
	presets, err := rd.GetPresets(ctx)
	if err != nil {
		return err
	}

	// Update presets
	return rd.UpdatePresets(presets)
}

func (rd *Radio) UpdatePresets(presets []Preset) error {
	select {
	case <-rd.dctx.Done():
		return rd.dctx.Err()
	case rd.updatePresetsChan <- presets:
		return nil
	}
}

func (rd *Radio) GetState(ctx context.Context) (*State, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-rd.dctx.Done():
		return nil, rd.dctx.Err()
	case state := <-rd.getStateChan:
		return &state, nil
	}
}

func (rd *Radio) IsPresetValid(preset int) bool {
	return !(preset < 1 || preset > rd.state.NumPresets)
}
