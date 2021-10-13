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
		Register:   make(chan *HubClient),
		Unregister: make(chan *HubClient),
		clients:    make(map[*HubClient]bool),
		cp:         cp,
		allStateChan:  make(chan State),
	}
	go h.hubLoop()
	return &h
}

func (h *Hub) hubLoop() {
	log.Println("hubLoop: started")
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			log.Println("hubLoop: registered")
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Println("hubLoop: unregistered")
			}
		case state := <-h.allStateChan:
			for client := range h.clients {
				select {
				case client.Send <- state:
				default:
					delete(h.clients, client)
					close(client.Send)
					log.Println("hubLoop: client deleted")
				}
			}
		}
	}
}

func (h *Hub) NewRadioFromClient(client goupnp.ServiceClient) (Radio, error) {
	// Get UUID from client
	uuid, ok := getServiceClientUUID(&client)
	if !ok {
		return Radio{}, errors.New("could not find uuid from client")
	}

	// Create sub
	dctx := context.Background()
	dctx, cancel := context.WithCancel(dctx)
	sub, err := h.cp.NewSubscription(dctx, &client.Service.EventSubURL.URL)
	if err != nil {
		cancel()
		return Radio{}, err
	}

	// Create rd and start radioLoop
	rd := Radio{
		Cancel:           cancel,
		Client:           client,
		GetStateChan:     make(chan State),
		UUID:             uuid,
		UpdateVolumeChan: make(chan int),
		allStateChan:        h.allStateChan,
		subscription:     sub,
		state:            &State{UUID: uuid},
	}
	go rd.radioLoop(dctx)

	return rd, nil
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
			log.Println(err)
			continue
		}

		radios[i] = radio
	}

	return radios, nil
}
