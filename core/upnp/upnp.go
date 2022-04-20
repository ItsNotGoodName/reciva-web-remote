package upnp

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/avast/retry-go/v3"
	"github.com/huin/goupnp"
)

// const controlServiceType = "urn:schemas-upnp-org:service:RenderingControl:1"

const recivaRadioServiceType = "urn:reciva-com:service:RecivaRadio:0.0"

func Discover() ([]goupnp.ServiceClient, []error, error) {
	return goupnp.NewServiceClients(recivaRadioServiceType)
}

var uuidReg = regexp.MustCompile(`(?m)\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`)

func GetUUID(c goupnp.ServiceClient) (string, error) {
	uuid := uuidReg.FindString(c.Location.String())
	if uuid == "" {
		return "", fmt.Errorf("could not find UUID in location: %s", c.Location)
	}

	return uuid, nil
}

func GetName(client goupnp.ServiceClient) string {
	return client.RootDevice.Device.FriendlyName
}

func GetModelName(client goupnp.ServiceClient) string {
	return client.RootDevice.Device.ModelName
}

func GetModelNumber(client goupnp.ServiceClient) string {
	return client.RootDevice.Device.ModelNumber
}

func GetEventURL(client goupnp.ServiceClient) url.URL {
	return client.Service.EventSubURL.URL
}

func SetPowerState(ctx context.Context, client goupnp.ServiceClient, power bool) error {
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
	err := retry.Do(func() error {
		return client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "SetPowerState", request, response)
	}, retry.Context(ctx))
	if err != nil {
		return fmt.Errorf("could not set power to %t: %w", power, err)
	}

	return nil
}

func PlayPreset(ctx context.Context, client goupnp.ServiceClient, preset int) error {
	// Create request
	request := &struct {
		NewPresetNumberValue string
	}{
		NewPresetNumberValue: strconv.Itoa(preset),
	}
	response := interface{}(nil)

	// Send request
	err := retry.Do(func() error {
		return client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "PlayPreset", request, response)
	}, retry.Context(ctx))
	if err != nil {
		return fmt.Errorf("could not play preset %d: %w", preset, err)
	}

	return nil
}

func GetAudioSources(ctx context.Context, client goupnp.ServiceClient) ([]string, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetAudioSourceListValue string
	}{}

	// Send request
	err := retry.Do(func() error {
		return client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "GetAudioSources", request, response)
	}, retry.Context(ctx))
	if err != nil {
		return []string{}, fmt.Errorf("could not get audio sources: %w", err)
	}

	return strings.Split(response.RetAudioSourceListValue, ","), nil
}

func GetNumberOfPresets(ctx context.Context, client goupnp.ServiceClient) (int, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetNumberOfPresetsValue int
	}{}

	// Send request
	err := retry.Do(func() error {
		return client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "GetNumberOfPresets", request, response)
	}, retry.Context(ctx))
	if err != nil {
		return 0, fmt.Errorf("could not get number of presets: %w", err)
	}

	return response.RetNumberOfPresetsValue, nil
}

type GetPresetResponse struct {
	Name string
	URL  string
}

func GetPreset(ctx context.Context, client goupnp.ServiceClient, num int) (*GetPresetResponse, error) {
	// Create request
	request := &struct {
		RetPresetNumberValue string
	}{
		RetPresetNumberValue: strconv.Itoa(num),
	}
	response := &struct {
		RetPresetName string
		RetPresetURL  string
	}{}

	// Send request
	err := retry.Do(func() error {
		return client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "GetPreset", request, response)
	}, retry.Context(ctx))
	if err != nil {
		return nil, fmt.Errorf("could not get preset %d: %w", num, err)
	}

	return &GetPresetResponse{Name: response.RetPresetName, URL: response.RetPresetURL}, nil
}

func SetVolume(ctx context.Context, client goupnp.ServiceClient, volume int) error {
	// Create request
	request := &struct {
		NewVolumeValue string
	}{
		NewVolumeValue: strconv.Itoa(volume),
	}
	response := interface{}(nil)

	// Send request
	err := retry.Do(func() error {
		return client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "SetVolume", request, response)
	}, retry.Context(ctx))
	if err != nil {
		return fmt.Errorf("could not set volume to %d: %w", volume, err)
	}

	return nil
}

func GetVolume(ctx context.Context, client goupnp.ServiceClient) (int, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetVolumeValue int
	}{}

	// Send request
	err := retry.Do(func() error {
		return client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "GetVolume", request, response)
	}, retry.Context(ctx))
	if err != nil {
		return 0, fmt.Errorf("could not get volume: %w", err)
	}

	return response.RetVolumeValue, nil
}
