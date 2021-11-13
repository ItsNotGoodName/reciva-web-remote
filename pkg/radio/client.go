package radio

import (
	"context"
	"strconv"
)

const radioServiceType = "urn:reciva-com:service:RecivaRadio:0.0"

// const controlServiceType = "urn:schemas-upnp-org:service:RenderingControl:1"

func (rd *Radio) setPowerState(ctx context.Context, power bool) error {
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

func (rd *Radio) playPreset(ctx context.Context, preset int) error {
	// Create request
	request := &struct {
		NewPresetNumberValue string
	}{}
	response := interface{}(nil)

	request.NewPresetNumberValue = strconv.Itoa(preset)

	// Play preset
	return rd.Client.SOAPClient.PerformActionCtx(ctx, rd.Client.Service.ServiceType, "PlayPreset", request, response)
}

func (rd *Radio) getNumberOfPresets(ctx context.Context) (int, error) {
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

func (rd *Radio) getPreset(ctx context.Context, num int) (*Preset, error) {
	// Create request
	request := &struct {
		RetPresetNumberValue string
	}{}
	response := &struct {
		RetPresetName string
		RetPresetURL  string
	}{}

	request.RetPresetNumberValue = strconv.Itoa(num)

	// Send request
	err := rd.Client.SOAPClient.PerformActionCtx(ctx, rd.Client.Service.ServiceType, "GetPreset", request, response)

	if err != nil {
		return nil, err
	}

	return &Preset{Number: num, Title: response.RetPresetName, URL: response.RetPresetURL}, nil
}

func (rd *Radio) getPresets(ctx context.Context) ([]Preset, error) {
	var pts []Preset

	for i := 1; i <= rd.state.NumPresets; i++ {
		p, err := rd.getPreset(ctx, i)
		if err != nil {
			return pts, err
		}
		pts = append(pts, *p)
	}

	return pts, nil
}

func (rd *Radio) setVolume(ctx context.Context, volume int) error {
	// Create request
	request := &struct {
		NewVolumeValue string
	}{}
	response := interface{}(nil)

	request.NewVolumeValue = strconv.Itoa(volume)

	// Send request
	return rd.Client.SOAPClient.PerformActionCtx(ctx, rd.Client.Service.ServiceType, "SetVolume", request, response)
}

func (rd *Radio) getVolume(ctx context.Context) (int, error) {
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
