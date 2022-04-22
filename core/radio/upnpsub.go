package radio

import (
	"encoding/xml"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type eventStateXML struct {
	Metadata string `xml:"playback-details>stream>metadata"`
	State    string `xml:"playback-details>state"`
	Title    string `xml:"playback-details>stream>title"`
	URL      string `xml:"playback-details>stream>url"`
}

func parseEvent(event *upnpsub.Event, frag *state.Fragment) {
	for _, v := range event.Properties {
		if v.Name == "PowerState" {
			// Power change
			newPower := v.Value == "On"
			frag.Power = &newPower
		} else if v.Name == "PlaybackXML" {
			if v.Value == "" {
				continue
			}

			esXML := eventStateXML{}
			if err := xml.Unmarshal([]byte(v.Value), &esXML); err != nil {
				log.Println("radio.parseEvent:", err)
				continue
			}

			// Status change
			newStatus := state.ParseStatus(esXML.State)
			frag.Status = &newStatus

			// Title change
			newTitle := esXML.Title
			frag.Title = &newTitle

			// Url change
			newURL := esXML.URL
			frag.URL = &newURL

			// Metadata change
			newMetadata := esXML.Metadata
			frag.Metadata = &newMetadata
		} else if v.Name == "IsMuted" {
			// IsMuted change
			newIsMuted := v.Value == "TRUE"
			frag.IsMuted = &newIsMuted
		} else if v.Name == "AudioSource" {
			// AudioSource change
			newAudioSource := v.Value
			frag.AudioSource = &newAudioSource
		}
	}
}
