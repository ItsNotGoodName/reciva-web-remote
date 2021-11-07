package api

import (
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/store"
	"github.com/gorilla/websocket"
)

type PresetAPI struct {
	s *store.Store
	h *radio.Hub
}

type radioWS struct {
	h        *radio.Hub
	conn     *websocket.Conn
	hubChan  *chan radio.State
	readChan chan *radio.State
	sendChan chan *radio.State
}
