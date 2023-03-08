package upnp

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/huin/goupnp"
	"github.com/sethvargo/go-retry"
)

// const controlServiceType = "urn:schemas-upnp-org:service:RenderingControl:1"

const recivaRadioServiceType = "urn:reciva-com:service:RecivaRadio:0.0"

func Discover() ([]Reciva, error) {
	recivaRadios, _, err := goupnp.NewServiceClients(recivaRadioServiceType)
	if err != nil {
		return nil, fmt.Errorf("could not discover Reciva radios: %w", err)
	}

	var radios []Reciva
	for _, recivaRadioClient := range recivaRadios {
		radios = append(radios, *newReciva(recivaRadioClient))
	}

	return radios, nil
}

type Reciva struct {
	recivaRadio goupnp.ServiceClient
}

func newReciva(recivaRadio goupnp.ServiceClient) *Reciva {
	return &Reciva{
		recivaRadio: recivaRadio,
	}
}

func (r *Reciva) performAction(ctx context.Context, action string, request interface{}, response interface{}) error {
	return retry.Do(ctx, retry.WithMaxRetries(3, retry.NewFibonacci(time.Second)), func(ctx context.Context) error {
		err := r.recivaRadio.SOAPClient.PerformActionCtx(ctx, r.recivaRadio.Service.ServiceType, action, request, response)
		if err != nil {
			if strings.Contains(err.Error(), "goupnp: error performing SOAP HTTP request:") { // Retry if it is a http error
				return retry.RetryableError(err)
			}
		}

		return err
	})
}

var uuidReg = regexp.MustCompile(`(?m)\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`)

func (r *Reciva) GetUUID() (string, error) {
	uuid := uuidReg.FindString(r.recivaRadio.Location.String())
	if uuid == "" {
		return "", fmt.Errorf("could not find UUID in location: %s", r.recivaRadio.Location)
	}

	return uuid, nil
}

func (r *Reciva) GetName() string {
	return r.recivaRadio.RootDevice.Device.FriendlyName
}

func (r *Reciva) GetModelName() string {
	return r.recivaRadio.RootDevice.Device.ModelName
}

func (r *Reciva) GetModelNumber() string {
	return r.recivaRadio.RootDevice.Device.ModelNumber
}

func (r *Reciva) GetEventURL() url.URL {
	return r.recivaRadio.Service.EventSubURL.URL
}

func (r *Reciva) SetPowerState(ctx context.Context, power bool) error {
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
	if err := r.performAction(ctx, "SetPowerState", request, response); err != nil {
		return fmt.Errorf("could not set power to %t: %w", power, err)
	}

	return nil
}

func (r *Reciva) PlayPreset(ctx context.Context, preset int) error {
	// Create request
	request := &struct {
		NewPresetNumberValue string
	}{
		NewPresetNumberValue: strconv.Itoa(preset),
	}
	response := interface{}(nil)

	// Send request
	if err := r.performAction(ctx, "PlayPreset", request, response); err != nil {
		return fmt.Errorf("could not play preset %d: %w", preset, err)
	}

	return nil
}

func (r *Reciva) SetAudioSource(ctx context.Context, audioSource string) error {
	// Create request
	request := &struct {
		NewAudioSourceValue string
	}{
		NewAudioSourceValue: audioSource,
	}
	response := interface{}(nil)

	// Send request
	if err := r.performAction(ctx, "SetAudioSource", request, response); err != nil {
		return fmt.Errorf("could not set audio source: %w", err)
	}

	return nil
}

func (r *Reciva) GetAudioSources(ctx context.Context) ([]string, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetAudioSourceListValue string
	}{}

	// Send request
	if err := r.performAction(ctx, "GetAudioSources", request, response); err != nil {
		return []string{}, fmt.Errorf("could not get audio sources: %w", err)
	}

	return strings.Split(response.RetAudioSourceListValue, ","), nil
}

func (r *Reciva) GetNumberOfPresets(ctx context.Context) (int, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetNumberOfPresetsValue int
	}{}

	// Send request
	if err := r.performAction(ctx, "GetNumberOfPresets", request, response); err != nil {
		return 0, fmt.Errorf("could not get number of presets: %w", err)
	}

	return response.RetNumberOfPresetsValue, nil
}

type GetPresetResponse struct {
	Name string
	URL  string
}

func (r *Reciva) GetPreset(ctx context.Context, num int) (*GetPresetResponse, error) {
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
	if err := r.performAction(ctx, "GetPreset", request, response); err != nil {
		return nil, fmt.Errorf("could not get preset %d: %w", num, err)
	}

	return &GetPresetResponse{Name: response.RetPresetName, URL: response.RetPresetURL}, nil
}

func (r *Reciva) SetVolume(ctx context.Context, volume int) error {
	// Create request
	request := &struct {
		NewVolumeValue string
	}{
		NewVolumeValue: strconv.Itoa(volume),
	}
	response := interface{}(nil)

	// Send request
	if err := r.performAction(ctx, "SetVolume", request, response); err != nil {
		return fmt.Errorf("could not set volume to %d: %w", volume, err)
	}

	return nil
}

func (r *Reciva) GetVolume(ctx context.Context) (int, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetVolumeValue int
	}{}

	// Send request
	if err := r.performAction(ctx, "GetVolume", request, response); err != nil {
		return 0, fmt.Errorf("could not get volume: %w", err)
	}

	return response.RetVolumeValue, nil
}
