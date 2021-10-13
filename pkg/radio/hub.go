package radio

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/goupnpsub"
)

func NewHub(cp *goupnpsub.ControlPoint) *Hub {
	h := Hub{
		Register:   make(chan *HubClient),
		Unregister: make(chan *HubClient),
		clients:    make(map[*HubClient]bool),
		cp:         cp,
		stateChan:  make(chan State),
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
		case state := <-h.stateChan:
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

func (h *Hub) NewRadioFromClient(client *Client) (Radio, error) {
	dctx := context.Background()
	dctx, cancel := context.WithCancel(dctx)

	// Create sub
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
		UpdateVolumeChan: make(chan int),
		stateChan:        h.stateChan,
		subscription:     sub,
		state:            &State{UUID: client.UUID},
	}
	go rd.radioLoop(dctx)

	return rd, nil
}

func (h *Hub) NewRadios() ([]Radio, error) {
	// Discover clients
	clients, err := NewClients()
	if err != nil {
		return nil, err
	}

	// Create radios from clients
	radios := make([]Radio, len(clients))
	for i := range radios {
		radio, err := h.NewRadioFromClient(&clients[i])
		if err != nil {
			continue
		}
		radios[i] = radio
	}

	return radios, nil
}
