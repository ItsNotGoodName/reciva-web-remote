package api

import (
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
}

type PresetAPI struct {
	s *store.Store
}

type radioWS struct {
	a        *API
	conn     *websocket.Conn
	hubChan  *chan radio.State
	readChan chan *radio.State
	sendChan chan *radio.State
}
