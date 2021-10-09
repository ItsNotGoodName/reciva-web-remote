package radio

import (
	"context"
	"encoding/xml"
	"log"
)

func (rd Radio) radioLoop(dctx context.Context) {
	// TODO: Refactor this function
	log.Println("radioLoop: started")
	emit := func() {
		rd.stateChan <- *rd.state
	}

	// Set name of radio
	rd.state.Name = rd.Client.ServiceClient.RootDevice.Device.FriendlyName
	// Get number of presets
	if presets, err := rd.Client.GetNumberOfPresets(dctx); err != nil {
		log.Println(err)
	} else {
		presets = presets - 2
		if presets < 1 {
			log.Println("radioLoop(ERROR): invalid number of presets were given from radio, ", rd.state.Presets)
		} else {
			rd.state.Presets = presets
		}
	}
	// Get volume
	if volume, err := rd.Client.GetVolume(dctx); err != nil {
		log.Println(err)
	} else {
		if !IsValidVolume(volume) {
			log.Println("radioLoop(ERROR): invalid volume was given from radio, ", rd.state.Volume)
		} else {
			rd.state.Volume = volume
			emit()
		}
	}

	for {
		select {
		case <-dctx.Done():
			log.Println("radioLoop: dctx is done, exiting")
			return
		case rd.GetStateChan <- *rd.state:
		case rd.state.Volume = <-rd.UpdateVolumeChan:
			emit()
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
				emit()
			}
		}
	}
}

func (rd Radio) IsPresetValid(preset int) bool {
	return !(preset < 1 || preset > rd.state.Presets)
}
