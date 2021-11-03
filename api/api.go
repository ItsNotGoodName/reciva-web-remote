package api

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/store"
	"github.com/gorilla/websocket"
)

type API struct {
	discoverChan    chan chan error
	h               *radio.Hub
	radioMap        map[string]radio.Radio
	radioMapRWMutex sync.RWMutex
	s               *store.Store
}

func NewAPI(h *radio.Hub, s *store.Store) *API {
	a := API{
		discoverChan:    make(chan chan error),
		h:               h,
		radioMap:        make(map[string]radio.Radio),
		radioMapRWMutex: sync.RWMutex{},
		s:               s,
	}

	log.Println("API.NewAPI: discovering radios...")
	a.discoverRadios()

	go a.apiLoop()

	return &a
}

func (a *API) HandleWS(conn *websocket.Conn, uuid string) {
	newRadioWS(conn, a).start(uuid)
}

func (a *API) apiLoop() {
	log.Println("API.apiLoop: started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		select {
		case d := <-a.discoverChan:
			d <- a.discoverRadios()
		case <-c:
			a.radioMapRWMutex.Lock()

			for _, v := range a.radioMap {
				v.Cancel()
			}
			for _, v := range a.radioMap {
				<-v.Subscription.Done
			}

			os.Exit(0)
		}
	}
}
