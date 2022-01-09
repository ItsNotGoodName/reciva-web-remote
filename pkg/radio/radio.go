package radio

import (
	"context"
	"encoding/xml"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/huin/goupnp"
)

type Radio struct {
	UUID       string                // UUID of the radio.
	Done       chan struct{}         // Channel to signal that the radio is done.
	client     goupnp.ServiceClient  // client is the SOAP client.
	getStateC  chan State            // getStateC is used to read the state.
	setVolumeC chan int              // setVolumeC is used to set the volume.
	pub        *Pub                  // pub is where state changes events are sent.
	sub        *upnpsub.Subscription // subscription that belongs to this Radio.
}

func newRadio(uuid string, client goupnp.ServiceClient, sub *upnpsub.Subscription, pub *Pub) *Radio {
	return &Radio{
		UUID:       uuid,
		Done:       make(chan struct{}),
		client:     client,
		getStateC:  make(chan State),
		setVolumeC: make(chan int),
		pub:        pub,
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

func (rd *Radio) publish(state *State) {
	state.UUID = rd.UUID
	rd.pub.publish(state)
}

func (rd *Radio) start(ctx context.Context, state State) {
	log.Println("Radio.start: started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Radio.start: ctx is done, exiting")
			close(rd.Done)
			return
		case rd.getStateC <- state:
		case newVolume := <-rd.setVolumeC:
			// Volume change
			if *state.Volume != newVolume {
				state.Volume = &newVolume
				rd.publish(&State{Volume: &newVolume})
			}
		case event := <-rd.sub.Event:
			newState := State{}
			changed := false

			for _, v := range event.Properties {
				if v.Name == "PowerState" {
					newPower := v.Value == "On"

					// Power change
					if newPower != *state.Power {
						state.Power = &newPower
						newState.Power = &newPower
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
						newState.State = sXML.State
						changed = true
					}

					// Title change
					if sXML.Title != state.Title {
						state.Title = sXML.Title
						newState.Title = sXML.Title

						// Preset Change
						newPreset := -1
						for i := range state.Presets {
							if state.Presets[i].Title == sXML.Title {
								newPreset = i + 1
								state.Title = state.Presets[i].Name
								newState.Title = state.Presets[i].Name
								break
							}
						}
						state.Preset = newPreset
						newState.Preset = newPreset

						changed = true
					}

					// Url change
					if sXML.URL != state.URL {
						state.URL = sXML.URL
						newState.URL = sXML.URL
						changed = true
					}

					// Metadata change
					if sXML.Metadata != *state.Metadata {
						newMetadata := sXML.Metadata
						state.Metadata = &newMetadata
						newState.Metadata = &newMetadata
						changed = true
					}
				} else if v.Name == "IsMuted" {
					newIsMuted := v.Value == "TRUE"

					if newIsMuted != *state.IsMuted {
						state.IsMuted = &newIsMuted
						newState.IsMuted = &newIsMuted
						changed = true
					}
				}
			}

			if changed {
				rd.publish(&newState)
			}
		}
	}
}
