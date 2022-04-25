package radio

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
	"github.com/huin/goupnp"
)

var (
	ErrRadioClosed = fmt.Errorf("radio closed")
)

type (
	HubService interface {
		Discover() (int, error)
		Get(uuid string) (Radio, error)
		List() []Radio
	}

	CreateService interface {
		Create(dctx context.Context, client goupnp.ServiceClient) (Radio, error)
	}

	RunService interface {
		Run(radio Radio, state state.State)
	}

	RadioService interface {
		GetState(ctx context.Context, radio Radio) (*state.State, error)
		PlayPreset(ctx context.Context, radio Radio, preset int) error
		Refresh(ctx context.Context, radio Radio) error
		RefreshVolume(ctx context.Context, radio Radio) error
		SetAudioSource(ctx context.Context, radio Radio, audioSource string) error
		SetPower(ctx context.Context, radio Radio, power bool) error
		SetVolume(ctx context.Context, radio Radio, volume int) error
	}

	Radio struct {
		UUID         string               // UUID of the radio.
		client       goupnp.ServiceClient // client is the SOAP client.
		readC        chan state.State     // readC is used to read the state.
		subscription upnpsub.Subscription // subscription to the UPnP event publisher.
		updateC      chan state.Fragment  // updateC is used to update state.
	}
)

func new(subscription upnpsub.Subscription, uuid string, client goupnp.ServiceClient) Radio {
	return Radio{
		UUID:         uuid,
		client:       client,
		readC:        make(chan state.State),
		subscription: subscription,
		updateC:      make(chan state.Fragment),
	}
}

func (r *Radio) Done() <-chan struct{} {
	return r.subscription.Done()
}

func (r *Radio) read(ctx context.Context) (*state.State, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-r.Done():
		return nil, ErrRadioClosed
	case state := <-r.readC:
		return &state, nil
	}
}

func (r *Radio) update(ctx context.Context, frag state.Fragment) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.Done():
		return ErrRadioClosed
	case r.updateC <- frag:
		return nil
	}
}
