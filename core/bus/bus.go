package bus

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type BusServiceImpl struct {
	hubService   radio.HubService
	radioService radio.RadioService
	statePub     state.Pub
}

func New(hubService radio.HubService, radioService radio.RadioService, statePub state.Pub) *BusServiceImpl {
	return &BusServiceImpl{
		hubService:   hubService,
		radioService: radioService,
		statePub:     statePub,
	}
}

func (bs *BusServiceImpl) stateGet(ctx context.Context, uuid string) (*state.State, error) {
	radio, err := bs.hubService.Get(uuid)
	if err != nil {
		return nil, err
	}

	return bs.radioService.GetState(ctx, radio)
}

func (bs *BusServiceImpl) Handle(ctx context.Context, readC <-chan Command, writeC chan<- Command) {
	stateUUID := ""
	stateSub, stateUnsub := make(<-chan state.Message), func() {}
	defer stateUnsub()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-stateSub:
			// Write state or state partial
			if state.IsChangedAll(msg.Changed) {
				writeCommand(ctx, writeC, newStateCommand(&msg.State))
			} else {
				writeCommand(ctx, writeC, newStatePartialCommand(state.GetPartial(&msg.State, msg.Changed)))
			}
		case dto := <-readC:
			switch dto.Type {
			case TypeStateSubscribe:
				stateUUID = fmt.Sprint(dto.Slug)
				if stateUUID == "" {
					writeCommand(ctx, writeC, newErrorCommand(fmt.Errorf("invalid uuid")))
					continue
				}

				// Subscribe to state
				stateUnsub()
				stateSub, stateUnsub = bs.statePub.Subscribe(stateUUID)

				// Get state
				state, err := bs.stateGet(ctx, stateUUID)
				if err != nil {
					writeCommand(ctx, writeC, newErrorCommand(err))
					continue
				}

				// Write state
				writeCommand(ctx, writeC, newStateCommand(state))
			case TypeStateUnsubscribe:
				// Unsubscribe from radio
				stateUnsub()
				select {
				case <-stateSub:
				default:
				}
			}
		}
	}
}
