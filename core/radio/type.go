package radio

import (
	"context"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
	"github.com/ItsNotGoodName/reciva-web-remote/core/upnp"
)

type (
	HubService interface {
		Discover(force bool) (int, error)
		Get(uuid string) (Radio, error)
		List() []Radio
	}

	CreateService interface {
		Create(ctx context.Context, dctx context.Context, reciva upnp.Reciva) (Radio, error)
	}

	RunService interface {
		Run(dctx context.Context, radio Radio, state state.State)
	}

	RadioService interface {
		GetState(ctx context.Context, radio Radio) (*state.State, error)
		PlayPreset(ctx context.Context, radio Radio, preset int) error
		RefreshSubscription(ctx context.Context, radio Radio) error
		RefreshVolume(ctx context.Context, radio Radio) error
		SetAudioSource(ctx context.Context, radio Radio, audioSource string) error
		SetPower(ctx context.Context, radio Radio, power bool) error
		SetVolume(ctx context.Context, radio Radio, volume int) error
	}

	Radio struct {
		UUID         string                                // UUID of the radio.
		Name         string                                // Name of the radio.
		reciva       upnp.Reciva                           // reciva is the UPnP client.
		stateC       chan state.State                      // stateC is used to read the state.
		subscription upnpsub.Subscription                  // subscription to the UPnP event publisher.
		updateFnC    chan func(*state.State) state.Changed // updateFnC is used to update state.
	}
)

func new(uuid, name string, reciva upnp.Reciva, subscription upnpsub.Subscription) Radio {
	return Radio{
		UUID:         uuid,
		Name:         name,
		reciva:       reciva,
		stateC:       make(chan state.State),
		subscription: subscription,
		updateFnC:    make(chan func(*state.State) state.Changed),
	}
}

func (r *Radio) Done() <-chan struct{} {
	return r.subscription.Done()
}

func (r *Radio) state(ctx context.Context) (*state.State, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-r.Done():
		return nil, core.ErrRadioClosed
	case state := <-r.stateC:
		return &state, nil
	}
}

func (r *Radio) update(ctx context.Context, updateFn func(*state.State) state.Changed) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.Done():
		return core.ErrRadioClosed
	case r.updateFnC <- updateFn:
		return nil
	}
}
