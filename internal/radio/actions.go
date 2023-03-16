package radio

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
)

func SetVolume(ctx context.Context, r hub.Radio, volume int) error {
	volume = state.NormalizeVolume(volume)

	if err := r.Reciva.SetVolume(ctx, volume); err != nil {
		return err
	}

	return r.Update(context.Background(), func(s *state.State) state.Changed {
		return s.SetVolume(volume)
	})
}

func PlayPreset(ctx context.Context, r hub.Radio, preset int) error {
	s, err := GetState(ctx, r)
	if err != nil {
		return err
	}

	if err := state.ValidAudioSource(s, state.AudioSourceInternetRadio); err != nil {
		log.Printf("radio.PlayPreset: %s is not a valid audio source", state.AudioSourceInternetRadio)
	}

	if err := state.ValidPresetNumber(s, preset); err != nil {
		return err
	}

	if !s.Power {
		if err := r.Reciva.SetPowerState(ctx, true); err != nil {
			return err
		}
	}

	return r.Reciva.PlayPreset(ctx, preset)
}

func SetAudioSource(ctx context.Context, r hub.Radio, audioSource string) error {
	s, err := GetState(ctx, r)
	if err != nil {
		return err
	}

	if err := state.ValidAudioSource(s, audioSource); err != nil {
		return err
	}

	if !s.Power {
		if err := r.Reciva.SetPowerState(ctx, true); err != nil {
			return err
		}
	}

	return r.Reciva.SetAudioSource(ctx, audioSource)
}

func SetPower(ctx context.Context, r hub.Radio, power bool) error {
	return r.Reciva.SetPowerState(ctx, power)
}

func RefreshVolume(ctx context.Context, r hub.Radio) error {
	volume, err := r.Reciva.GetVolume(ctx)
	if err != nil {
		return err
	}

	return r.Update(context.Background(), func(s *state.State) state.Changed {
		return s.SetVolume(volume)
	})
}

func RefreshSubscription(ctx context.Context, r hub.Radio) error {
	r.Subscription.Renew()
	return nil
}

func GetState(ctx context.Context, r hub.Radio) (*state.State, error) {
	return r.State(ctx)
}
