//go:build !dev

package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

func newMux() *chi.Mux {
	return chi.NewMux()
}

func newUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{}
}
