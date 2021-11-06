package radio

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/huin/goupnp"
)

func NewHub(cp *upnpsub.ControlPoint) *Hub {
	h := Hub{
		PresetMutator: func(ctx context.Context, p *Preset) {
			p.Name = p.Title
		},
		cp:           cp,
		discoverChan: make(chan chan error),
		radios:       make(map[string]*Radio),
		radiosMu:     sync.RWMutex{},
		stateOPS:     make(chan func(map[*chan State]bool)),
		stopChan:     make(chan chan error),
	}
	return &h
}

func (h *Hub) Start() error {
	go h.discoverLoop()
	go h.stateLoop()

	// Discover radios
	errChan := make(chan error)
	h.discoverChan <- errChan
	return <-errChan
}

func (h *Hub) NewRadios() ([]*Radio, error) {
	// Discover clients
	clients, _, err := goupnp.NewServiceClients(radioServiceType)
	if err != nil {
		return nil, err
	}

	// Create radios array
	var radios []*Radio
	for i := range clients {
		radio, err := h.NewRadio(clients[i])
		if err != nil {
			log.Println("Hub.NewRadios:", err)
			continue
		}

		radios = append(radios, radio)
	}

	return radios, nil
}

func (h *Hub) NewRadio(client goupnp.ServiceClient) (*Radio, error) {
	// Get UUID from client
	uuid, ok := getServiceClientUUID(&client)
	if !ok {
		return nil, errors.New("could not find uuid from client")
	}

	// Create sub
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sub, err := h.cp.NewSubscription(ctx, &client.Service.EventSubURL.URL)
	if err != nil {
		cancel()
		return nil, err
	}

	// Create rd and start radioLoop
	rd := Radio{
		Cancel:           cancel,
		Client:           client,
		Subscription:     sub,
		UUID:             uuid,
		ctx:              ctx,
		getStateChan:     make(chan State),
		h:                h,
		state:            NewState(uuid),
		updateVolumeChan: make(chan int),
	}
	go rd.radioLoop()

	return &rd, nil
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

func (h *Hub) GetRadio(uuid string) (*Radio, bool) {
	h.radiosMu.RLock()
	r, ok := h.radios[uuid]
	h.radiosMu.RUnlock()
	return r, ok
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

func (h *Hub) Discover() error {
	errChan := make(chan error)
	select {
	case h.discoverChan <- errChan:
		return <-errChan
	default:
		return ErrDiscovering
	}
}

// Stop closes radios and makes further discovery not possible.
func (h *Hub) Stop() error {
	errChan := make(chan error)
	h.stopChan <- errChan
	return <-errChan
}

func (h *Hub) discoverLoop() {
	stopped := false
	firstDiscover := true
	for {
		select {
		case d := <-h.discoverChan:
			if stopped {
				d <- errors.New("discoverLoop is stopped")
				continue
			}

			// Discover radios
			radios, err := h.NewRadios()
			if err != nil {
				d <- err
			}

			// Create new radios map
			newRadioMap := make(map[string]*Radio)
			for _, v := range radios {
				newRadioMap[v.UUID] = v
			}
			numRadios := len(newRadioMap)

			// Wait for radios to update their state and also rate limit discovery
			if !firstDiscover {
				<-time.After(time.Second * 5)
			} else {
				firstDiscover = false
			}

			// Swap old and new radios map
			h.radiosMu.Lock()
			oldRadioMap := h.radios
			h.radios = newRadioMap
			h.radiosMu.Unlock()

			// Close old radios
			for _, v := range oldRadioMap {
				v.Cancel()
			}

			log.Printf("Hub.discoverLoop: discovered %d radios", numRadios)
			d <- nil
		case s := <-h.stopChan:
			newRadioMap := make(map[string]*Radio)

			h.radiosMu.Lock()
			oldRadioMap := h.radios
			h.radios = newRadioMap
			h.radiosMu.Unlock()

			for _, v := range oldRadioMap {
				v.Cancel()
			}
			for _, v := range oldRadioMap {
				<-v.Subscription.Done
			}
			stopped = true

			s <- nil
		}
	}
}
