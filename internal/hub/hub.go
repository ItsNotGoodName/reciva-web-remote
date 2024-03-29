package hub

import (
	"context"
	"sort"
	"sync"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/upnp"
)

type RadioStateC = chan state.State
type RadioUpdateFnC = chan func(*state.State) state.Changed

type Radio struct {
	UUID         string               // UUID of the radio.
	Name         string               // Name of the radio.
	Reciva       upnp.Reciva          // Reciva is the UPnP client.
	Subscription upnpsub.Subscription // Subscription to the UPnP event publisher.
	stateC       RadioStateC          // stateC is used to read the state.
	updateFnC    RadioUpdateFnC       // updateFnC is used to update state.
	close        context.CancelFunc   // close is used shutdown the radio connection.
}

func (r *Radio) Done() <-chan struct{} {
	return r.Subscription.Done()
}

func (r *Radio) State(ctx context.Context) (*state.State, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-r.Done():
		return nil, internal.ErrRadioClosed
	case state := <-r.stateC:
		return &state, nil
	}
}

func (r *Radio) Update(ctx context.Context, updateFn func(*state.State) state.Changed) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.Done():
		return internal.ErrRadioClosed
	case r.updateFnC <- updateFn:
		return nil
	}
}

type Hub struct {
	doneC chan struct{}

	radiosMapMu sync.RWMutex
	radiosMap   map[string]Radio
}

func New() *Hub {
	return &Hub{
		doneC:       make(chan struct{}),
		radiosMapMu: sync.RWMutex{},
		radiosMap:   make(map[string]Radio),
	}
}

func (h *Hub) Background(ctx context.Context, doneC chan<- struct{}) {
	// Wait for context
	<-ctx.Done()

	// Close radios
	h.radiosMapMu.RLock()
	for _, r := range h.radiosMap {
		<-r.Done()
	}
	h.radiosMap = make(map[string]Radio)
	// Prevent creating new radios
	close(h.doneC)
	h.radiosMapMu.RUnlock()

	// Done
	doneC <- struct{}{}
}

func (h *Hub) Create(uuid, name string, reciva upnp.Reciva, subscription upnpsub.Subscription, stateC RadioStateC, updateFnC RadioUpdateFnC, close context.CancelFunc) (Radio, error) {
	h.radiosMapMu.Lock()
	select {
	case <-h.doneC:
		h.radiosMapMu.Unlock()
		return Radio{}, internal.ErrHubClosed
	default:
	}
	h.delete(uuid)

	r := Radio{
		UUID:         uuid,
		Name:         name,
		Reciva:       reciva,
		Subscription: subscription,
		stateC:       stateC,
		updateFnC:    updateFnC,
		close:        close,
	}

	h.radiosMap[uuid] = r
	h.radiosMapMu.Unlock()

	return r, nil
}

func (h *Hub) Delete(uuid string) error {
	h.radiosMapMu.Lock()
	err := h.delete(uuid)
	h.radiosMapMu.Unlock()

	return err
}

func (h *Hub) delete(uuid string) error {
	r, ok := h.radiosMap[uuid]
	if !ok {
		return internal.ErrRadioNotFound
	}

	r.close()
	delete(h.radiosMap, uuid)

	return nil
}

func (h *Hub) Get(uuid string) (Radio, error) {
	h.radiosMapMu.RLock()
	r, ok := h.radiosMap[uuid]
	h.radiosMapMu.RUnlock()
	if !ok {
		return Radio{}, internal.ErrRadioNotFound
	}

	return r, nil
}

func (h *Hub) Exists(uuid string) bool {
	h.radiosMapMu.RLock()
	_, ok := h.radiosMap[uuid]
	h.radiosMapMu.RUnlock()
	return ok
}

func (h *Hub) List() []Radio {
	h.radiosMapMu.RLock()
	var radios []Radio
	for _, r := range h.radiosMap {
		radios = append(radios, r)
	}
	h.radiosMapMu.RUnlock()

	// Sort radios
	sort.Slice(radios, func(i, j int) bool {
		return radios[i].UUID < radios[j].UUID
	})

	return radios
}
