package radio

import (
	"context"
	"log"
	"sync"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/huin/goupnp"
)

type Hub struct {
	cp           *upnpsub.ControlPoint // cp is used to create subscriptions.
	discoverChan chan chan error       // discoverChan is used to discover radios.

	radiosMu sync.RWMutex      // radiosMu is used to protect radios map.
	radios   map[string]*Radio // radios is used to store all the Radios.
}

func NewHub(cp *upnpsub.ControlPoint) *Hub {
	return &Hub{
		cp:           cp,
		discoverChan: make(chan chan error),
		radios:       make(map[string]*Radio),
		radiosMu:     sync.RWMutex{},
	}
}

func (h *Hub) GetRadios() []*Radio {
	h.radiosMu.RLock()
	var radios []*Radio
	for _, r := range h.radios {
		radios = append(radios, r)
	}
	h.radiosMu.RUnlock()

	return radios
}

func (h *Hub) GetRadio(uuid string) (*Radio, error) {
	h.radiosMu.RLock()
	r, ok := h.radios[uuid]
	h.radiosMu.RUnlock()
	if !ok {
		return nil, ErrRadioNotFound
	}

	return r, nil
}

func (h *Hub) IsValidRadio(uuid string) bool {
	h.radiosMu.RLock()
	_, ok := h.radios[uuid]
	h.radiosMu.RUnlock()
	return ok
}

func (h *Hub) Start(ctx context.Context) {
	discover := func(cancel context.CancelFunc) (context.CancelFunc, error) {
		newCtx := context.Background()
		newCtx, newCancel := context.WithCancel(newCtx)
		newRadios, err := h.discover(newCtx)
		if err != nil {
			newCancel()
			return nil, err
		}

		radios := make(map[string]*Radio)
		for _, r := range newRadios {
			radios[r.UUID] = r
		}

		if cancel != nil {
			cancel()
		}

		h.radiosMu.Lock()
		h.radios = radios
		h.radiosMu.Unlock()

		return newCancel, nil
	}

	cancel, err := discover(nil)
	if err != nil {
		log.Println("Hub.run(ERROR):", err)
	}

	for {
		select {
		case <-ctx.Done():
			if cancel != nil {
				cancel()
			}
		case errC := <-h.discoverChan:
			cancel, err = discover(cancel)
			errC <- err
		}
	}
}

func (h *Hub) Discover() error {
	errChan := make(chan error)
	select {
	case h.discoverChan <- errChan:
		return <-errChan
	default:
		return ErrDiscovering
	}
}

func (h *Hub) discover(ctx context.Context) ([]*Radio, error) {
	clients, _, err := goupnp.NewServiceClients(radioServiceType)
	if err != nil {
		return nil, err
	}

	rds := make(chan *Radio)
	var wg sync.WaitGroup

	for i := range clients {
		wg.Add(1)

		go (func(idx int) {
			radio, err := h.newRadio(ctx, clients[idx])
			if err != nil {
				log.Println("Hub.discover(ERROR):", err)
			} else {
				rds <- radio
			}
			wg.Done()
		})(i)
	}

	go (func() {
		wg.Wait()
		close(rds)
	})()

	var radios []*Radio
	for r := range rds {
		radios = append(radios, r)
	}

	return radios, nil
}

func (h *Hub) newRadio(ctx context.Context, client goupnp.ServiceClient) (*Radio, error) {
	// Create sub
	sub, err := h.cp.NewSubscription(ctx, &client.Service.EventSubURL.URL)
	if err != nil {
		return nil, err
	}

	// Create state
	state, err := newStateFromClient(ctx, client)
	if err != nil {
		return nil, err
	}

	// Create and start radio
	rd := newRadio(state.UUID, client, sub)
	go rd.run(ctx, *state)

	return rd, nil
}
