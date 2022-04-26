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

func parseEvent(event *upnpsub.Event, s *state.State) (changed state.Changed) {
	for _, v := range event.Properties {
		if v.Name == "PowerState" {
			// Power change
			power := v.Value == "On"
			changed = changed.Merge(s.SetPower(power))
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
			changed = changed.Merge(s.SetStatus(status))

			// Title change
			title := esXML.Title
			changed = changed.Merge(s.SetTitle(title))

			// URL change
			url := esXML.URL
			changed = changed.Merge(s.SetURL(url))

			// Metadata change
			metadata := esXML.Metadata
			changed = changed.Merge(s.SetMetadata(metadata))
		} else if v.Name == "IsMuted" {
			// IsMuted change
			isMuted := v.Value == "TRUE"
			changed = changed.Merge(s.SetIsMuted(isMuted))
		} else if v.Name == "AudioSource" {
			// AudioSource change
			audioSource := v.Value
			c, err := s.SetAudioSource(audioSource)
			if err != nil {
				log.Println("radio.parseEvent:", err)
				continue
			}

			changed = changed.Merge(c)
		}
	}

	return changed
}
