package radio

import (
	"context"
	"fmt"
)

const radioServiceType = "urn:reciva-com:service:RecivaRadio:0.0"

// const controlServiceType = "urn:schemas-upnp-org:service:RenderingControl:1"

func (rd *Radio) SetPowerState(ctx context.Context, power bool) error {
	// Create request
	request := &struct {
		NewPowerStateValue string
	}{}
	response := &struct {
		RetPowerStateValue string
	}{}
	if power {
		request.NewPowerStateValue = "On"
	} else {
		request.NewPowerStateValue = "Off"
	}

	// Send request
	return rd.Client.SOAPClient.PerformActionCtx(ctx, rd.Client.Service.ServiceType, "SetPowerState", request, response)
}

func (rd *Radio) PlayPreset(ctx context.Context, preset int) error {
	// Create request
	request := &struct {
		NewPresetNumberValue string
	}{}
	response := interface{}(nil)

	request.NewPresetNumberValue = fmt.Sprint(preset)

	select {
	case state := <-rd.getStateChan:
		// Turn on radio if it is not already on
		if !state.Power {
			if err := rd.SetPowerState(ctx, true); err != nil {
				return err
			}
		}

		// Play preset
		return rd.Client.SOAPClient.PerformActionCtx(ctx, rd.Client.Service.ServiceType, "PlayPreset", request, response)
	case <-ctx.Done():
		return ctx.Err()
	case <-rd.dctx.Done():
		return rd.dctx.Err()
	}
}

func (rd *Radio) GetNumberOfPresets(ctx context.Context) (int, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetNumberOfPresetsValue int
	}{}

	// Send request
	err := rd.Client.SOAPClient.PerformActionCtx(ctx, rd.Client.Service.ServiceType, "GetNumberOfPresets", request, response)

	// Return number of presets
	if err != nil {
		return 0, err
	}
	return response.RetNumberOfPresetsValue, nil
}

func (rd *Radio) SetVolume(volume int) error {
	// Create request
	request := &struct {
		NewVolumeValue string
	}{}
	response := interface{}(nil)

	request.NewVolumeValue = fmt.Sprint(volume)

	// Send request
	if err := rd.Client.SOAPClient.PerformActionCtx(rd.dctx, rd.Client.Service.ServiceType, "SetVolume", request, response); err != nil {
		return err
	}

	// Update state volume
	return rd.UpdateVolume(volume)
}

func (rd *Radio) GetVolume(ctx context.Context) (int, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetVolumeValue int
	}{}

	// Send request
	err := rd.Client.SOAPClient.PerformActionCtx(ctx, rd.Client.Service.ServiceType, "GetVolume", request, response)

	// Return volume
	if err != nil {
		return 0, err
	}
	return response.RetVolumeValue, nil
}
