package radio

import (
	"context"
	"errors"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/huin/goupnp"
)

func NewHub(cp *upnpsub.ControlPoint) *Hub {
	h := Hub{
		cp:  cp,
		ops: make(chan func(map[*chan State]bool)),
	}
	go h.hubLoop()
	return &h
}

func (h *Hub) NewRadios() ([]Radio, error) {
	// Discover clients
	clients, _, err := goupnp.NewServiceClients(radioServiceType)
	if err != nil {
		return nil, err
	}

	// Create radios array
	var radios []Radio
	for i := range clients {
		radio, err := h.NewRadio(clients[i])
		if err != nil {
			log.Println("Hub.NewRadios:", err)
			continue
		}

		radios = append(radios, *radio)
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
		Cancel:            cancel,
		Client:            client,
		Subscription:      sub,
		UUID:              uuid,
		emitState:         h.EmitState,
		ctx:               ctx,
		getStateChan:      make(chan State),
		state:             NewState(uuid),
		updatePresetsChan: make(chan []Preset),
		updateVolumeChan:  make(chan int),
	}
	go rd.radioLoop()

	return &rd, nil
}

func (h *Hub) EmitState(state *State) {
	h.ops <- func(m map[*chan State]bool) {
		for client := range m {
			select {
			case *client <- *state:
			default:
				delete(m, client)
				close(*client)
				log.Println("Hub.hubLoop: client deleted")
			}
		}
	}
}

func (h *Hub) AddClient(client *chan State) {
	h.ops <- func(m map[*chan State]bool) {
		m[client] = true
		log.Println("Hub.hubLoop: client registered")
	}
}

func (h *Hub) RemoveClient(client *chan State) {
	h.ops <- func(m map[*chan State]bool) {
		if _, ok := m[client]; ok {
			delete(m, client)
			close(*client)
			log.Println("Hub.hubLoop: client unregistered")
		}
	}
}

func (h *Hub) hubLoop() {
	log.Println("Hub.hubLoop: started")
	m := make(map[*chan State]bool)
	for op := range h.ops {
		op(m)
	}
}
