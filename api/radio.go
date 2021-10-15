package api

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/gorilla/websocket"
)

type API struct {
	discoverChan    chan chan error
	h               *radio.Hub
	radioMap        map[string]radio.Radio
	radioMapRWMutex sync.RWMutex
}

func (a *API) GetRadio(uuid string) (*radio.Radio, bool) {
	a.radioMapRWMutex.RLock()
	radio, ok := a.radioMap[uuid]
	a.radioMapRWMutex.RUnlock()
	if !ok {
		return nil, ok
	}

	return &radio, ok
}

func (a *API) GetRadioState(ctx context.Context, uuid string) (*radio.State, bool) {
	a.radioMapRWMutex.RLock()
	radio, ok := a.radioMap[uuid]
	a.radioMapRWMutex.RUnlock()
	if !ok {
		return nil, false
	}

	state, err := radio.GetState(ctx)
	if err != nil {
		return nil, false
	}

	return state, ok
}

func (a *API) GetRadioStates(ctx context.Context) []radio.State {
	states := make([]radio.State, 0, len(a.radioMap))

	a.radioMapRWMutex.RLock()
	for _, v := range a.radioMap {
		state, err := v.GetState(ctx)
		if err != nil {
			continue
		}
		states = append(states, *state)
	}
	a.radioMapRWMutex.RUnlock()

	return states
}

func (a *API) IsValidRadio(uuid string) bool {
	a.radioMapRWMutex.RLock()
	_, ok := a.radioMap[uuid]
	a.radioMapRWMutex.RUnlock()
	return ok
}

func (a *API) DiscoverRadios() error {
	d := make(chan error)
	select {
	case a.discoverChan <- d:
		return <-d
	default:
		return errors.New("radios are already being discovered")
	}
}

func (a *API) discoverRadios() error {
	// Discover radios
	radios, err := a.h.NewRadios()
	if err != nil {
		return err
	}

	// Create newRadioMap using radios and radios' uuid
	newRadioMap := make(map[string]radio.Radio)
	for _, v := range radios {
		newRadioMap[v.UUID] = v
	}
	numRadios := len(newRadioMap)

	// Close old radioMap and set new radioMap
	a.radioMapRWMutex.Lock()
	for _, v := range a.radioMap {
		v.Cancel()
	}
	a.radioMap = newRadioMap
	a.radioMapRWMutex.Unlock()

	log.Printf("api.discoverRadios: discovered %d radios", numRadios)
	return nil
}

func (a *API) discoverRadiosLoop() {
	for d := range a.discoverChan {
		d <- a.discoverRadios()
	}
}

func (a *API) HandleWS(conn *websocket.Conn, uuid string) {
	newRadioWS(conn, a, uuid).start()
}

func NewService(h *radio.Hub) *API {
	a := API{
		discoverChan:    make(chan chan error),
		h:               h,
		radioMap:        make(map[string]radio.Radio),
		radioMapRWMutex: sync.RWMutex{},
	}

	go a.discoverRadiosLoop()
	go func() {
		d := make(chan error)
		a.discoverChan <- d
		<-d
	}()

	return &a
}
