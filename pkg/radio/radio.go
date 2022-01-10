package radio

import (
	"context"
	"encoding/xml"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/huin/goupnp"
)

type Radio struct {
	Done       chan struct{}         // Channel to signal that the radio is done.
	UUID       string                // UUID of the radio.
	client     goupnp.ServiceClient  // client is the SOAP client.
	getStateC  chan State            // getStateC is used to read the state.
	mutateC    chan struct{}         // mutateC is used to signal that mutator has changed.
	pub        *Pub                  // pub is where state change events are sent.
	setVolumeC chan int              // setVolumeC is used to set the volume.
	sub        *upnpsub.Subscription // sub that belongs to this Radio.
}

func newRadio(uuid string, client goupnp.ServiceClient, sub *upnpsub.Subscription, pub *Pub) *Radio {
	return &Radio{
		UUID:       uuid,
		Done:       make(chan struct{}),
		client:     client,
		getStateC:  make(chan State),
		mutateC:    make(chan struct{}, 1),
		pub:        pub,
		setVolumeC: make(chan int),
		sub:        sub,
	}
}

// GetState returns the current state of the radio.
func (rd *Radio) GetState(ctx context.Context) (*State, error) {
	select {
	case s := <-rd.getStateC:
		return &s, nil
	case <-rd.Done:
		return nil, ErrRadioNotFound
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// SetPowerState sets the power state of the radio.
func (rd *Radio) SetPower(ctx context.Context, power bool) error {
	return setPowerState(ctx, rd.client, power)
}

// PlayPreset plays the given preset.
func (rd *Radio) PlayPreset(ctx context.Context, preset int) error {
	state, err := rd.GetState(ctx)
	if err != nil {
		return err
	}

	if err := state.ValidPreset(preset); err != nil {
		return err
	}

	if !state.On() {
		if err := setPowerState(ctx, rd.client, true); err != nil {
			return err
		}
	}

	return playPreset(ctx, rd.client, preset)
}

// SetVolume sets the volume of the radio.
func (rd *Radio) SetVolume(ctx context.Context, volume int) error {
	volume = normalizeVolume(volume)
	if err := setVolume(ctx, rd.client, volume); err != nil {
		return err
	}

	select {
	case rd.setVolumeC <- volume:
		return nil
	case <-rd.Done:
		return ErrRadioNotFound
	case <-ctx.Done():
		return ctx.Err()
	}
}

// RefreshVolume fetches volume from the radio.
func (rd *Radio) RefreshVolume(ctx context.Context) error {
	volume, err := getVolume(ctx, rd.client)
	if err != nil {
		return err
	}

	select {
	case rd.setVolumeC <- volume:
		return nil
	case <-rd.Done:
		return ErrRadioNotFound
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Refresh renews subscription to the radio.
func (rd *Radio) Refresh() {
	rd.sub.Renew()
}

func (rd *Radio) Mutate() {
	select {
	case rd.mutateC <- struct{}{}:
	default:
	}
}

func (rd *Radio) publish(state *State) {
	state.UUID = rd.UUID
	rd.pub.publish(state)
}

func (rd *Radio) start(ctx context.Context, state State, mutator MutatorPort) {
	log.Println("Radio.start: started radio", rd.UUID)

	mutator.Mutate(ctx, &state)

	for {
		select {
		case <-ctx.Done():
			<-rd.sub.Done
			close(rd.Done)
			log.Println("Radio.start: stopped radio", rd.UUID)
			return
		case <-rd.mutateC:
			rd.publish(mutator.Mutate(ctx, &state))
		case rd.getStateC <- state:
		case newVolume := <-rd.setVolumeC:
			// Volume change
			if *state.Volume != newVolume {
				state.Volume = &newVolume
				rd.publish(&State{Volume: &newVolume})
			}
		case event := <-rd.sub.Event:
			diffState := State{}
			changed := false

			for _, v := range event.Properties {
				if v.Name == "PowerState" {
					newPower := v.Value == "On"

					// Power change
					if newPower != *state.Power {
						state.Power = &newPower
						diffState.Power = &newPower
						changed = true
					}
				} else if v.Name == "PlaybackXML" {
					if v.Value == "" {
						continue
					}

					sXML := stateXML{}
					if err := xml.Unmarshal([]byte(v.Value), &sXML); err != nil {
						log.Println("Radio.start(ERROR):", err)
						continue
					}

					// State change
					if sXML.State != state.State {
						state.State = sXML.State
						diffState.State = sXML.State
						changed = true
					}

					// Title change
					if sXML.Title != state.Title {
						state.Title = sXML.Title
						diffState.Title = sXML.Title

						// Preset Change
						newPreset := -1
						for i := range state.Presets {
							if state.Presets[i].Title == sXML.Title {
								newPreset = i + 1
								state.Title = state.Presets[i].Name
								diffState.Title = state.Presets[i].Name
								break
							}
						}
						state.Preset = newPreset
						diffState.Preset = newPreset

						changed = true
					}

					// Url change
					if sXML.URL != state.URL {
						state.URL = sXML.URL
						diffState.URL = sXML.URL

						mutator.MutateNewURL(ctx, &state)
						diffState.NewURL = state.NewURL

						changed = true
					}

					// Metadata change
					if sXML.Metadata != *state.Metadata {
						newMetadata := sXML.Metadata
						state.Metadata = &newMetadata
						diffState.Metadata = &newMetadata
						changed = true
					}
				} else if v.Name == "IsMuted" {
					newIsMuted := v.Value == "TRUE"

					if newIsMuted != *state.IsMuted {
						state.IsMuted = &newIsMuted
						diffState.IsMuted = &newIsMuted
						changed = true
					}
				}
			}

			if changed {
				rd.publish(&diffState)
			}
		}
	}
}
