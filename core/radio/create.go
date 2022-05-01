package radio

import (
	"context"
	"log"
	"time"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
	"github.com/ItsNotGoodName/reciva-web-remote/core/upnp"
	"github.com/huin/goupnp"
)

type CreateServiceImpl struct {
	controlPoint upnpsub.ControlPoint
	radioService RunService
}

func NewCreateService(controlPoint upnpsub.ControlPoint, runService RunService) *CreateServiceImpl {
	return &CreateServiceImpl{
		controlPoint: controlPoint,
		radioService: runService,
	}
}

func (cs *CreateServiceImpl) Background(ctx context.Context, doneC chan<- struct{}) {
	go func() {
		if err := upnpsub.ListenAndServe("", cs.controlPoint); err != nil {
			log.Fatalln("Failed to start control point:", err)
		}
	}()

	<-ctx.Done()
	doneC <- struct{}{}
}

func (cs *CreateServiceImpl) Create(ctx context.Context, dctx context.Context, client goupnp.ServiceClient) (Radio, error) {
	// Get UUID
	uuid, err := upnp.GetUUID(client)
	if err != nil {
		return Radio{}, err
	}

	// Get name
	name := upnp.GetName(client)

	// Create state
	s := state.New(uuid, name, upnp.GetModelName(client), upnp.GetModelNumber(client))

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

	// Back off for a second to prevent subscription from failing
	// TODO: Find a better way to do this
	time.Sleep(time.Second)

	// Create subscription
	eventURL := upnp.GetEventURL(client)
	sub, err := cs.controlPoint.Subscribe(dctx, &eventURL)
	if err != nil {
		return Radio{}, err
	}

	// Create and run radio
	radio := new(uuid, name, client, sub)
	go cs.radioService.Run(dctx, radio, s)

	return radio, nil
}
