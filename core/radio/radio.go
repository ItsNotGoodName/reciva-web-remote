package radio

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type RadioServiceImpl struct{}

func NewRadioService() *RadioServiceImpl {
	return &RadioServiceImpl{}
}

func (rs *RadioServiceImpl) SetVolume(ctx context.Context, r Radio, volume int) error {
	volume = state.NormalizeVolume(volume)

	if err := r.reciva.SetVolume(ctx, volume); err != nil {
		return err
	}

	return r.update(context.Background(), func(s *state.State) state.Changed {
		return s.SetVolume(volume)
	})
}

func (rs *RadioServiceImpl) PlayPreset(ctx context.Context, r Radio, preset int) error {
	s, err := rs.GetState(ctx, r)
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
		if err := r.reciva.SetPowerState(ctx, true); err != nil {
			return err
		}
	}

	return r.reciva.PlayPreset(ctx, preset)
}

func (rs *RadioServiceImpl) SetPower(ctx context.Context, r Radio, power bool) error {
	return r.reciva.SetPowerState(ctx, power)
}

func (rs *RadioServiceImpl) RefreshVolume(ctx context.Context, r Radio) error {
	volume, err := r.reciva.GetVolume(ctx)
	if err != nil {
		return err
	}

	return r.update(ctx, func(s *state.State) state.Changed {
		return s.SetVolume(volume)
	})
}

func (rs *RadioServiceImpl) RefreshSubscription(ctx context.Context, r Radio) error {
	r.subscription.Renew()
	return nil
}

func (rs *RadioServiceImpl) GetState(ctx context.Context, r Radio) (*state.State, error) {
	return r.state(ctx)
}

func (rs *RadioServiceImpl) SetAudioSource(ctx context.Context, r Radio, audioSource string) error {
	s, err := rs.GetState(ctx, r)
	if err != nil {
		return err
	}

	if err := state.ValidAudioSource(s, audioSource); err != nil {
		return err
	}

	return r.reciva.SetAudioSource(ctx, audioSource)
}
