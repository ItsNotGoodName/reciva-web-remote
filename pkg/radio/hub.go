package radio

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/huin/goupnp"
)

type Hub struct {
	DoneC     chan struct{}         // DoneC is closed when all radios are stopped.
	Pub       *Pub                  // Pub is the state change event publisher
	mutator   MutatorPort           // mutator is used to mutate the state of the radio.
	cp        *upnpsub.ControlPoint // cp is used to create subscriptions.
	discoverC chan chan error       // discoverC is used to discover radios.

	radiosMu sync.RWMutex      // radiosMu is used to protect radios map.
	radios   map[string]*Radio // radios is used to store all the Radios.
}

func NewHub(cp *upnpsub.ControlPoint) *Hub {
	return NewHubWithMutator(cp, NewMutator())
}

func NewHubWithMutator(cp *upnpsub.ControlPoint, mutator MutatorPort) *Hub {
	return &Hub{
		DoneC:     make(chan struct{}),
		Pub:       newPub(),
		cp:        cp,
		discoverC: make(chan chan error),
		mutator:   mutator,
		radios:    make(map[string]*Radio),
		radiosMu:  sync.RWMutex{},
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

func (h *Hub) GetRadioState(ctx context.Context, uuid string) (*State, error) {
	h.radiosMu.RLock()
	radio, ok := h.radios[uuid]
	h.radiosMu.RUnlock()
	if !ok {
		return nil, ErrRadioNotFound
	}

	return radio.GetState(ctx)
}

func (h *Hub) GetRadioStates(ctx context.Context) []State {
	states := make([]State, 0, len(h.radios))

	h.radiosMu.RLock()
	for _, v := range h.radios {
		state, err := v.GetState(ctx)
		if err != nil {
			continue
		}
		states = append(states, *state)
	}
	h.radiosMu.RUnlock()

	return states
}

func (h *Hub) IsValidRadio(uuid string) bool {
	h.radiosMu.RLock()
	_, ok := h.radios[uuid]
	h.radiosMu.RUnlock()
	return ok
}

func (h *Hub) mutatorStart(ctx context.Context) {
	c := h.mutator.GetTrigger()
	for {
		select {
		case <-ctx.Done():
			return
		case <-c:
			h.radiosMu.RLock()
			for _, v := range h.radios {
				v.Mutate()
			}
			h.radiosMu.RUnlock()
		}
	}
}

func (h *Hub) Start(ctx context.Context) {
	go h.mutatorStart(ctx)

	firstRun := true

	discover := func(cancel context.CancelFunc) (context.CancelFunc, error) {
		radioCTX, newCancel := context.WithCancel(ctx)
		newRadios, err := h.discover(radioCTX)
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

		if !firstRun {
			// Sleep for a while to allow radios to get their initial state and to prevent spamming the discovery channel.
			select {
			case <-ctx.Done():
				newCancel()
				return nil, ctx.Err()
			case <-time.After(3 * time.Second):
			}
			firstRun = false
		}

		h.radiosMu.Lock()
		h.radios = radios
		h.radiosMu.Unlock()

		return newCancel, nil
	}

	cancel, err := discover(nil)
	if err != nil {
		log.Println("Hub.Start(ERROR):", err)
	}

	for {
		select {
		case <-ctx.Done():
			h.radiosMu.RLock()
			for _, radio := range h.radios {
				<-radio.DoneC
			}
			close(h.DoneC)
			h.radiosMu.RUnlock()
			log.Println("Hub.Start: stopped")
			return
		case errC := <-h.discoverC:
			cancel, err = discover(cancel)
			errC <- err
		}
	}
}

func (h *Hub) Discover() error {
	errChan := make(chan error)
	select {
	case h.discoverC <- errChan:
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
				if err != context.Canceled {
					log.Println("Hub.discover(ERROR):", err)
				}
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
	rd := newRadio(state.UUID, client, sub, h.Pub)
	go rd.start(ctx, *state, h.mutator)

	return rd, nil
}
