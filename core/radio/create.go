package radio

import (
	"context"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
	"github.com/ItsNotGoodName/reciva-web-remote/core/upnp"
	"github.com/huin/goupnp"
)

type CreateServiceImpl struct {
	controlPoint upnpsub.ControlPoint
	radioService RunService
}

func NewCreateService(controlPoint upnpsub.ControlPoint, radioService RunService) *CreateServiceImpl {
	return &CreateServiceImpl{
		controlPoint: controlPoint,
		radioService: radioService,
	}
}

func (cs *CreateServiceImpl) Create(ctx context.Context, client goupnp.ServiceClient) (Radio, error) {
	// Get UUID
	uuid, err := upnp.GetUUID(client)
	if err != nil {
		return Radio{}, err
	}

	// Create state
	s := state.New(uuid, upnp.GetName(client), upnp.GetModelName(client), upnp.GetModelNumber(client))

	// Get and set volume
	volume, err := upnp.GetVolume(ctx, client)
	if err != nil {
		return Radio{}, err
	}
	s.SetVolume(volume)

	// Get and parse presets count
	presetsCount, err := upnp.GetNumberOfPresets(ctx, client)
	if err != nil {
		return Radio{}, err
	}
	if presetsCount, err = state.ParsePresetsCount(presetsCount); err != nil {
		return Radio{}, err
	}

	// Get and set presets
	var presets []state.Preset
	for i := 1; i <= presetsCount; i++ {
		p, err := upnp.GetPreset(ctx, client, i)
		if err != nil {
			return Radio{}, err
		}

		presets = append(presets, state.NewPreset(i, p.Name, p.URL))
	}
	s.SetPresets(presets)

	// Get audio sources
	audioSources, err := upnp.GetAudioSources(ctx, client)
	if err != nil {
		return Radio{}, err
	}
	s.SetAudioSources(audioSources)

	// Create subscription
	eventURL := upnp.GetEventURL(client)
	sub, err := cs.controlPoint.Subscribe(ctx, &eventURL)
	if err != nil {
		return Radio{}, err
	}

	// Create and run radio
	radio := new(sub, uuid, client)
	go cs.radioService.Run(radio, s)

	return radio, nil
}
