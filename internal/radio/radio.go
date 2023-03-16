package radio

import (
	"context"
	"encoding/xml"
	"log"
	"time"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
)

type StateHook interface {
	// OnStart is called when when state is first created.
	OnStart(context.Context, *state.State, state.Changed) state.Changed
	// OnChanged is called when state changes.
	OnChanged(context.Context, *state.State, state.Changed) state.Changed
}

func run(ctx context.Context, radio hub.Radio, s state.State, stateC hub.RadioStateC, updateFnC hub.RadioUpdateFnC, stateHook StateHook) {
	handle := func(c state.Changed) {
		c = stateHook.OnChanged(ctx, &s, c)
		if c != state.ChangedNone {
			pubsub.DefaultPub.Publish(pubsub.StateTopic, pubsub.StateMessage{State: s, Changed: c})
		}
	}
	if c := stateHook.OnStart(ctx, &s, state.ChangedAll); c != state.ChangedNone {
		pubsub.DefaultPub.Publish(pubsub.StateTopic, pubsub.StateMessage{State: s, Changed: c})
	}

	msgC, unsub := pubsub.DefaultPub.Subscribe([]pubsub.Topic{pubsub.StateHookStaleTopic})
	defer unsub()
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-radio.Done():
			return
		case <-ticker.C:
			go func() {
				if err := RefreshVolume(ctx, radio); err != nil {
					log.Println("radio.run: failed to refresh volume:", err)
				}
			}()
		case msg := <-msgC:
			data := msg.Data.(pubsub.StateHookStaleMessage)
			handle(data.Changed)
		case stateC <- s:
		case fn := <-updateFnC:
			handle(fn(&s))
		case event := <-radio.Subscription.Events():
			handle(parseEvent(event, &s))
		}
	}
}

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
