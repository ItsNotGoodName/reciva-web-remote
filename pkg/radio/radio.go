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
			if rd.state.Volume != newVolume {
				rd.state.Volume = newVolume
				rd.stateChanged()
			}
		case newEvent := <-rd.Subscription.EventChan:
			changed := false
			for _, v := range newEvent.Properties {
				if v.Name == "PowerState" {
					rd.state.Power = v.Value == "On"
					changed = true
				} else if v.Name == "PlaybackXML" {
					if v.Value == "" {
						continue
					}

					oldMetadata := rd.state.Metadata
					rd.state.Metadata = ""
					if err := xml.Unmarshal([]byte(v.Value), rd.state); err != nil {
						log.Println("Radio.radioLoop:", err)
						rd.state.Metadata = oldMetadata
						continue
					}

					changed = true
				} else if v.Name == "IsMuted" {
					rd.state.IsMuted = v.Value == "TRUE"
					changed = true
				}
			}
			if changed {
				rd.stateChanged()
			}
		}
	}
}

func (rd *Radio) initState() {
	// Set name of radio
	rd.state.Name = rd.Client.RootDevice.Device.FriendlyName

	// Get number of presets
	var presets int
	if err := retry.Do(func() error {
		if p, e := rd.GetNumberOfPresets(rd.dctx); e != nil {
			return e
		} else {
			presets = p
			return nil
		}
	}, retry.Context(rd.dctx)); err != nil {
		log.Println("Radio.initState:", err)
	} else {
		presets = presets - 2
		if presets < 1 {
			log.Println("Radio.initState(ERROR): invalid number of presets were given from radio,", rd.state.Presets)
		} else {
			rd.state.Presets = presets
		}
	}

	var volume int
	// Get volume
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
			log.Println("Radio.initState(ERROR): invalid volume was given from radio,", rd.state.Volume)
		} else {
			rd.state.Volume = volume
		}
	}
}

func (rd *Radio) UpdateVolume(volume int) error {
	select {
	case <-rd.dctx.Done():
		return rd.dctx.Err()
	case rd.updateVolumeChan <- volume:
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

func (rd *Radio) stateChanged() {
	rd.allStateChan <- *rd.state
}

func (rd *Radio) IsPresetValid(preset int) bool {
	return !(preset < 1 || preset > rd.state.Presets)
}
