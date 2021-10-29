package radio

import (
	"context"
	"errors"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
	"github.com/huin/goupnp"
)

func NewHub(cp *goupnpsub.ControlPoint) *Hub {
	h := Hub{
		Register:         make(chan *chan State),
		Unregister:       make(chan *chan State),
		clients:          make(map[*chan State]bool),
		cp:               cp,
		receiveStateChan: make(chan State),
	}
	go h.hubLoop()
	return &h
}

func (h *Hub) hubLoop() {
	log.Println("Hub.hubLoop: started")
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			log.Println("Hub.hubLoop: registered")
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(*client)
				log.Println("Hub.hubLoop: unregistered")
			}
		case state := <-h.receiveStateChan:
			for client := range h.clients {
				select {
				case *client <- state:
				default:
					delete(h.clients, client)
					close(*client)
					log.Println("Hub.hubLoop: client deleted")
				}
			}
		}
	}
}

func (h *Hub) NewRadioFromClient(client goupnp.ServiceClient) (*Radio, error) {
	// Get UUID from client
	uuid, ok := getServiceClientUUID(&client)
	if !ok {
		return nil, errors.New("could not find uuid from client")
	}

	// Create sub
	dctx := context.Background()
	dctx, cancel := context.WithCancel(dctx)
	sub, err := h.cp.NewSubscription(dctx, &client.Service.EventSubURL.URL)
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
		sendStateChan:     h.receiveStateChan,
		dctx:              dctx,
		getStateChan:      make(chan State),
		state:             NewState(uuid),
		updatePresetsChan: make(chan []Preset),
		updateVolumeChan:  make(chan int),
	}
	go rd.radioLoop()

	return &rd, nil
}

func (h *Hub) NewRadios() ([]Radio, error) {
	// Discover clients
	clients, _, err := goupnp.NewServiceClients(radioServiceType)
	if err != nil {
		return nil, err
	}

	// Create radios array
	radios := make([]Radio, len(clients))
	for i := range clients {
		radio, err := h.NewRadioFromClient(clients[i])
		if err != nil {
			log.Println("Hub.NewRadios:", err)
			continue
		}

		radios[i] = *radio
	}

	return radios, nil
}
