package radio

import (
	"context"
	"fmt"
	"strconv"

	"github.com/huin/goupnp"
)

const radioServiceType = "urn:reciva-com:service:RecivaRadio:0.0"

// const controlServiceType = "urn:schemas-upnp-org:service:RenderingControl:1"

func setPowerState(ctx context.Context, client goupnp.ServiceClient, power bool) error {
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
	err := client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "SetPowerState", request, response)
	if err != nil {
		return fmt.Errorf("could not set power to %t: %v", power, err)
	}

	return nil
}

func playPreset(ctx context.Context, client goupnp.ServiceClient, preset int) error {
	// Create request
	request := &struct {
		NewPresetNumberValue string
	}{}
	response := interface{}(nil)

	request.NewPresetNumberValue = strconv.Itoa(preset)

	// Play preset
	err := client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "PlayPreset", request, response)
	if err != nil {
		return fmt.Errorf("could not play preset %d: %v", preset, err)
	}

	return nil
}

func getNumberOfPresets(ctx context.Context, client goupnp.ServiceClient) (int, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetNumberOfPresetsValue int
	}{}

	// Send request
	err := client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "GetNumberOfPresets", request, response)
	if err != nil {
		return 0, fmt.Errorf("could not get number of presets: %v", err)
	}

	return response.RetNumberOfPresetsValue, nil
}

func getPreset(ctx context.Context, client goupnp.ServiceClient, num int) (*Preset, error) {
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
	err := client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "GetPreset", request, response)
	if err != nil {
		return nil, fmt.Errorf("could not get preset %d: %v", num, err)
	}

	return &Preset{Number: num, Title: response.RetPresetName, URL: response.RetPresetURL}, nil
}

func getPresets(ctx context.Context, client goupnp.ServiceClient, numPresets int) ([]Preset, error) {
	var pts []Preset

	for i := 1; i <= numPresets; i++ {
		p, err := getPreset(ctx, client, i)
		if err != nil {
			return pts, err
		}
		pts = append(pts, *p)
	}

	return pts, nil
}

func setVolume(ctx context.Context, client goupnp.ServiceClient, volume int) error {
	// Create request
	request := &struct {
		NewVolumeValue string
	}{}
	response := interface{}(nil)

	request.NewVolumeValue = strconv.Itoa(volume)

	// Send request
	err := client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "SetVolume", request, response)
	if err != nil {
		return fmt.Errorf("could not set volume to %d: %v", volume, err)
	}

	return nil
}

func getVolume(ctx context.Context, client goupnp.ServiceClient) (int, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetVolumeValue int
	}{}

	// Send request
	err := client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "GetVolume", request, response)
	if err != nil {
		return 0, fmt.Errorf("could not get volume: %v", err)
	}

	return response.RetVolumeValue, nil
}
