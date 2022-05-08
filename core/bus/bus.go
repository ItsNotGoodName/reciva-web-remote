package bus

import (
	"context"

	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type Bus struct {
	app      dto.App
	statePub state.Pub
}

func New(app dto.App, statePub state.Pub) *Bus {
	return &Bus{
		app:      app,
		statePub: statePub,
	}
}

func (b *Bus) Handle(ctx context.Context, readC <-chan dto.Command, writeC chan<- dto.Command) {
	stateSub, stateUnsub := make(<-chan state.PubMessage), func() {}
	defer func() {
		stateUnsub()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-stateSub:
			// Write state or state partial
			if msg.Changed.Is(state.ChangedAll) {
				writeCommand(ctx, writeC, newStateCommand(&msg.State))
			} else {
				writeCommand(ctx, writeC, newStatePartialCommand(state.GetPartial(&msg.State, msg.Changed)))
			}
		case c := <-readC:
			switch c.Type {
			case dto.CommandTypeStateSubscribe:
				stateUUID, err := parseStateSubscribe(c.Slug)
				if err != nil {
					writeCommand(ctx, writeC, newErrorCommand(err))
					continue
				}

				// Subscribe to state
				stateUnsub()
				stateSub, stateUnsub = b.statePub.Subscribe(stateUUID)

				// Get state
				res, err := b.app.StateGet(ctx, &dto.StateRequest{UUID: stateUUID})
				if err != nil {
					writeCommand(ctx, writeC, newErrorCommand(err))
					continue
				}

				// Write state
				writeCommand(ctx, writeC, newStateCommand(&res.State))
			case dto.CommandTypeStateUnsubscribe:
				// Unsubscribe from radio
				stateUnsub()
			}
		}
	}
}
