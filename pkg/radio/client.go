package radio

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/huin/goupnp"
)

const radioServiceType = "urn:reciva-com:service:RecivaRadio:0.0"

// const controlServiceType = "urn:schemas-upnp-org:service:RenderingControl:1"

func NewClientFromUrl(radioUrl *url.URL) (*Client, error) {
	clients, err := goupnp.NewServiceClientsByURL(radioUrl, radioServiceType)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(clients) != 1 {
		return nil, errors.New("service client not found")
	}

	uuid, ok := getServiceClientUUID(&clients[0])

	if !ok {
		return nil, errors.New("uuid not found from service client")
	}

	return &Client{clients[0], uuid}, nil
}

func NewClients() ([]Client, error) {
	// Discover
	clients, _, err := goupnp.NewServiceClients(radioServiceType)
	if err != nil {
		return nil, err
	}

	// Create radios array
	radios := make([]Client, len(clients))
	for i := 0; i < len(clients); i++ {
		uuid, ok := getServiceClientUUID(&clients[i])

		if !ok {
			log.Println("NewClients(WARNING): could not find uuid, ignoring Client")
			continue
		}

		radios[i] = Client{clients[i], uuid}
	}

	return radios, nil
}

func (client *Client) SetPowerState(ctx context.Context, power bool) error {
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
	return client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "SetPowerState", request, response)
}

func (client *Client) PlayPreset(ctx context.Context, preset int) error {
	// Create request
	request := &struct {
		NewPresetNumberValue string
	}{}
	response := interface{}(nil)

	request.NewPresetNumberValue = fmt.Sprint(preset)

	// Send request
	return client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "PlayPreset", request, response)
}

func (client *Client) GetNumberOfPresets(ctx context.Context) (int, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetNumberOfPresetsValue int
	}{}

	// Send request
	err := client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "GetNumberOfPresets", request, response)

	// Return number of presets
	if err != nil {
		return 0, err
	}
	return response.RetNumberOfPresetsValue, nil
}

func (client *Client) SetVolume(ctx context.Context, volume int) error {
	// Create request
	request := &struct {
		NewVolumeValue string
	}{}
	response := interface{}(nil)

	request.NewVolumeValue = fmt.Sprint(volume)

	// Send request
	return client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "SetVolume", request, response)
}

func (client *Client) GetVolume(ctx context.Context) (int, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetVolumeValue int
	}{}

	// Send request
	err := client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "GetVolume", request, response)

	// Return volume
	if err != nil {
		return 0, err
	}
	return response.RetVolumeValue, nil
}

func (client *Client) GetPowerState(ctx context.Context) (bool, error) {
	// Create request
	request := interface{}(nil)
	response := &struct {
		RetPowerStateValue string
	}{}

	// Send request
	err := client.SOAPClient.PerformActionCtx(ctx, client.Service.ServiceType, "GetPowerState", request, response)

	// Parse response
	if err != nil {
		return false, err
	}
	if response.RetPowerStateValue == "On" {
		return true, nil
	} else if response.RetPowerStateValue == "Off" {
		return false, nil
	}
	return false, fmt.Errorf("invalid power state value, %s", response.RetPowerStateValue)
}
