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

type RadioWS struct {
	h        *radio.Hub
	conn     *websocket.Conn
	sub      *radio.Sub
	readChan chan *radio.State
	sendChan chan *radio.State
}
