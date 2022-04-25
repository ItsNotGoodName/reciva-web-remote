package app

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type App struct {
	hubService   radio.HubService
	radioService radio.RadioService
	statePub     state.StatePub
}

func New(hubService radio.HubService, radioService radio.RadioService, statePub state.StatePub) *App {
	return &App{
		hubService:   hubService,
		radioService: radioService,
		statePub:     statePub,
	}
}

func (a *App) StateGet(ctx context.Context, uuid string) (*state.State, error) {
	radio, err := a.hubService.Get(uuid)
	if err != nil {
		return nil, err
	}

	return a.radioService.GetState(ctx, radio)
}

func (a *App) Bus(ctx context.Context, readC <-chan Command, writeC chan<- Command) {
	stateUUID := ""
	stateSub, stateUnsub := make(<-chan state.State), func() {}
	defer stateUnsub()
	var stateLast *state.State

	for {
		select {
		case <-ctx.Done():
			return
		case state := <-stateSub:
			// Write state
			writeCommand(ctx, writeC, NewStateCommand(&state))
		case dto := <-readC:
			switch dto.Type {
			case TypeStateSubscribe:
				stateUUID = fmt.Sprint(dto.Slug)
				if stateUUID == "" {
					writeCommand(ctx, writeC, NewErrorCommand(fmt.Errorf("invalid uuid")))
					continue
				}

				// Subscribe to state
				stateUnsub()
				stateSub, stateUnsub = a.statePub.Subscribe(stateUUID)

				// Get state
				var err error
				stateLast, err = a.StateGet(ctx, stateUUID)
				if err != nil {
					writeCommand(ctx, writeC, NewErrorCommand(err))
					continue
				}

				// Write state
				writeCommand(ctx, writeC, NewStateCommand(stateLast))
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
