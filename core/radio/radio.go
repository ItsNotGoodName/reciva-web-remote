package radio

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
	"github.com/ItsNotGoodName/reciva-web-remote/core/upnp"
)

type RadioServiceImpl struct{}

func NewRadioService() *RadioServiceImpl {
	return &RadioServiceImpl{}
}

func (rs *RadioServiceImpl) SetVolume(ctx context.Context, radio Radio, volume int) error {
	volume = state.NormalizeVolume(volume)

	if err := upnp.SetVolume(ctx, radio.client, volume); err != nil {
		return err
	}

	return radio.update(context.Background(), func(s *state.State) state.Changed {
		return s.SetVolume(volume)
	})
}

func (rs *RadioServiceImpl) PlayPreset(ctx context.Context, radio Radio, preset int) error {
	s, err := rs.GetState(ctx, radio)
	if err != nil {
		return err
	}

	if err := state.ValidAudioSource(s, state.AudioSourceInternetRadio); err != nil {
		log.Printf("radio.RadioService.PlayPreset: %s is not a valid audio source", state.AudioSourceInternetRadio)
	}

	if err := state.ValidPresetNumber(s, preset); err != nil {
		return err
	}

	if !s.Power {
		if err := upnp.SetPowerState(ctx, radio.client, true); err != nil {
			return err
		}
	}

	return upnp.PlayPreset(ctx, radio.client, preset)
}

func (rs *RadioServiceImpl) SetPower(ctx context.Context, radio Radio, power bool) error {
	return upnp.SetPowerState(ctx, radio.client, power)
}

func (rs *RadioServiceImpl) RefreshVolume(ctx context.Context, radio Radio) error {
	volume, err := upnp.GetVolume(ctx, radio.client)
	if err != nil {
		return err
	}

	return radio.update(ctx, func(s *state.State) state.Changed {
		return s.SetVolume(volume)
	})
}

func (rs *RadioServiceImpl) RefreshSubscription(ctx context.Context, radio Radio) error {
	radio.subscription.Renew()
	return nil
}

func (rs *RadioServiceImpl) GetState(ctx context.Context, radio Radio) (*state.State, error) {
	return radio.state(ctx)
}

func (rs *RadioServiceImpl) SetAudioSource(ctx context.Context, radio Radio, audioSource string) error {
	s, err := rs.GetState(ctx, radio)
	if err != nil {
		return err
	}

	if err := state.ValidAudioSource(s, audioSource); err != nil {
		return err
	}

	return upnp.SetAudioSource(ctx, radio.client, audioSource)
}
