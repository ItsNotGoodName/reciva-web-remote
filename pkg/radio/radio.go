package radio

import (
	"context"
	"encoding/xml"
	"log"

	"github.com/avast/retry-go"
)

func (rd *Radio) GetState(ctx context.Context) (*State, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-rd.ctx.Done():
		return nil, rd.ctx.Err()
	case state := <-rd.getStateChan:
		return &state, nil
	}
}

func (rd *Radio) SetPower(ctx context.Context, power bool) error {
	return rd.setPowerState(ctx, power)
}

func (rd *Radio) PlayPreset(ctx context.Context, preset int) error {
	if preset < 1 || preset > rd.state.NumPresets {
		return ErrInvalidPreset
	}

	state, err := rd.GetState(ctx)
	if err != nil {
		return err
	}

	if !*state.Power && rd.setPowerState(ctx, true) != nil {
		return err
	}

	return rd.playPreset(ctx, preset)
}

func (rd *Radio) SetVolume(volume int) error {
	volume = NormalizeVolume(volume)
	if err := rd.setVolume(rd.ctx, volume); err != nil {
		return err
	}
	return rd.updateVolume(volume)
}

func (rd *Radio) RefreshVolume(ctx context.Context) error {
	vol, err := rd.getVolume(ctx)
	if err != nil {
		return err
	}

	return rd.updateVolume(vol)
}

func (rd *Radio) Renew() {
	rd.Subscription.Renew()
}

func (rd *Radio) radioLoop() {
	log.Println("Radio.radioLoop: started")

	rd.initState()

	for {
		select {
		case <-rd.ctx.Done():
			log.Println("Radio.radioLoop: ctx is done, exiting")
			return
		case rd.getStateChan <- *rd.state:
		case newVolume := <-rd.updateVolumeChan:
			// Volume change
			if *rd.state.Volume != newVolume {
				rd.state.Volume = &newVolume
				rd.sendState(&State{Volume: &newVolume})
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

						// PresetNum Change
						newPresetNum := -1
						for i := range rd.state.Presets {
							if rd.state.Presets[i].Name == sXML.Title {
								newPresetNum = i + 1
							}
						}
						rd.state.PresetNum = newPresetNum
						newState.PresetNum = newPresetNum

						changed = true
					}

					// Url change
					if sXML.URL != rd.state.URL {
						rd.state.URL = sXML.URL
						newState.URL = sXML.URL
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
				rd.sendState(&newState)
			}
		}
	}
}

func (rd *Radio) updateVolume(volume int) error {
	select {
	case <-rd.ctx.Done():
		return rd.ctx.Err()
	case rd.updateVolumeChan <- volume:
		return nil
	}
}

func (rd *Radio) sendState(state *State) {
	state.UUID = rd.state.UUID
	rd.emitState(state)
}

func (rd *Radio) initState() {
	// Set name of radio
	rd.state.Name = rd.Client.RootDevice.Device.FriendlyName

	// Get number of presets
	var numPresets int
	if err := retry.Do(func() error {
		if p, e := rd.getNumberOfPresets(rd.ctx); e != nil {
			return e
		} else {
			numPresets = p
			return nil
		}
	}, retry.Context(rd.ctx)); err != nil {
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
		if v, e := rd.getVolume(rd.ctx); e != nil {
			return e
		} else {
			volume = v
			return nil
		}
	}, retry.Context(rd.ctx)); err != nil {
		log.Println("Radio.initState:", err)
	} else {
		rd.state.Volume = &volume
	}

	// Get presets
	var presets []Preset
	if err := retry.Do(func() error {
		if p, e := rd.getPresets(rd.ctx); e != nil {
			return e
		} else {
			presets = p
			return nil
		}
	}, retry.Context(rd.ctx)); err != nil {
		log.Println("Radio.initState:", err)
	} else {
		rd.state.Presets = presets
	}
}
