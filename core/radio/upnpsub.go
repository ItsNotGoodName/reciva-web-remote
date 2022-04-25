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
			power := v.Value == "On"
			frag.Power = &power
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
			status := state.ParseStatus(esXML.State)
			frag.Status = &status

			// Title change
			title := esXML.Title
			frag.Title = &title

			// URL change
			url := esXML.URL
			frag.URL = &url

			// Metadata change
			metadata := esXML.Metadata
			frag.Metadata = &metadata
		} else if v.Name == "IsMuted" {
			// IsMuted change
			isMuted := v.Value == "TRUE"
			frag.IsMuted = &isMuted
		} else if v.Name == "AudioSource" {
			// AudioSource change
			audioSource := v.Value
			frag.AudioSource = &audioSource
		}
	}
}
