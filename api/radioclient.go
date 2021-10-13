package api

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
)

func (a *API) SetRadioPower(ctx context.Context, rd *radio.Radio, power bool) error {
	return rd.SetPowerState(ctx, power)
}

func (a *API) PlayRadioPreset(ctx context.Context, rd *radio.Radio, preset int) error {
	select {
	case state := <-rd.GetStateChan:
		// Turn on radio if it is not already on
		if !state.Power {
			if err := rd.SetPowerState(ctx, true); err != nil {
				return err
			}
		}

		// Play preset
		return rd.PlayPreset(ctx, preset)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *API) SetRadioVolume(ctx context.Context, rd *radio.Radio, volume int) error {
	// Set volume
	if err := rd.SetVolume(ctx, volume); err != nil {
		return nil
	}

	// Get volume
	vol, err := rd.GetVolume(ctx)
	if err != nil {
		return err
	}

	// Update state volume
	select {
	case rd.UpdateVolumeChan <- vol:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
