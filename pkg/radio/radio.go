package radio

import (
	"context"
	"encoding/xml"
	"log"

	"github.com/avast/retry-go"
)

func (rd *Radio) radioLoop(dctx context.Context) {
	log.Println("radioLoop: started")

	rd.initState(dctx)

	for {
		select {
		case <-dctx.Done():
			log.Println("radioLoop: dctx is done, exiting")
			return
		case rd.GetStateChan <- *rd.state:
		case rd.state.Volume = <-rd.UpdateVolumeChan:
			rd.stateChanged()
		case newEvent := <-rd.subscription.EventChan:
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
						log.Println(err)
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

func (rd *Radio) initState(dctx context.Context) {
	// Set name of radio
	rd.state.Name = rd.Client.RootDevice.Device.FriendlyName

	// Get number of presets
	var presets int
	if err := retry.Do(func() error {
		if p, e := rd.GetNumberOfPresets(dctx); e != nil {
			return e
		} else {
			presets = p
			return nil
		}
	}, retry.Context(dctx)); err != nil {
		log.Println(err)
	} else {
		presets = presets - 2
		if presets < 1 {
			log.Println("radioLoop(ERROR): invalid number of presets were given from radio, ", rd.state.Presets)
		} else {
			rd.state.Presets = presets
		}
	}

	var volume int
	// Get volume
	if err := retry.Do(func() error {
		if v, e := rd.GetVolume(dctx); e != nil {
			return e
		} else {
			volume = v
			return nil
		}
	}, retry.Context(dctx)); err != nil {
		log.Println(err)
	} else {
		if !IsValidVolume(volume) {
			log.Println("radioLoop(ERROR): invalid volume was given from radio, ", rd.state.Volume)
		} else {
			rd.state.Volume = volume
		}
	}
}

func (rd *Radio) stateChanged() {
	rd.allStateChan <- *rd.state
}

func (rd *Radio) IsPresetValid(preset int) bool {
	return !(preset < 1 || preset > rd.state.Presets)
}
