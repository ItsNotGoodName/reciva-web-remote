//go:build !dev

package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

func newRouter() chi.Router {
	return chi.NewRouter()
}

func newUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{}
}
