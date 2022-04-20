package radio

import (
	"encoding/xml"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type RunServiceImpl struct {
	fragmentPub state.FragmentPub
}

func NewRunService(fragmentPub state.FragmentPub) *RunServiceImpl {
	return &RunServiceImpl{
		fragmentPub: fragmentPub,
	}
}

func (rs *RunServiceImpl) Run(radio Radio, s state.State) {
	handle := func(f state.Fragment) {
		if frag, changed := s.Merge(f); changed {
			rs.fragmentPub.Publish(frag)
		}
	}

	for {
		select {
		case <-radio.Done():
			return
		case radio.readC <- s:
		case fragment := <-radio.updateC:
			handle(fragment)
		case event := <-radio.subscription.Events():
			fragment := state.NewFragment(radio.UUID)
			parseEvent(event, &fragment)
			handle(fragment)
		}
	}
}

type stateXML struct {
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

			sXML := stateXML{}
			if err := xml.Unmarshal([]byte(v.Value), &sXML); err != nil {
				log.Println("radio.parseEvent:", err)
				continue
			}

			// Status change
			newStatus := state.ParseStatus(sXML.State)
			frag.Status = &newStatus

			// Title change
			newTitle := sXML.Title
			frag.Title = &newTitle

			// Url change
			newURL := sXML.URL
			frag.URL = &newURL

			// Metadata change
			newMetadata := sXML.Metadata
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
