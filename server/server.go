package server

import (
	"context"
	"log"
	"strconv"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/api"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/store"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	p    *api.PresetAPI
	h    *radio.Hub
	r    *gin.Engine
	u    *websocket.Upgrader
	stop context.CancelFunc
}

func NewServer(cfg *config.Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	s := Server{
		stop: cancel,
	}

	// Create and start controlpoint
	cp := upnpsub.NewControlPointWithPort(cfg.CPort)
	go cp.Start()

	// Create radio hub
	s.h = radio.NewHub(cp)

	// Create store
	store, err := store.NewStore(cfg.ConfigFile)
	if err != nil {
		log.Println("server.NewServer(INFO): store is in readonly mode:", err)
	}

	// Create preset api
	s.p = api.NewPresetAPI(store, s.h)

	// Start hub
	go s.h.Start(ctx)

	// Create websocket upgrader
	s.u = NewUpgrader()

	// Create router
	s.r = NewRouter()

	// Create routes
	s.routes()

	return &s
}

func (s *Server) Start(cfg *config.Config) {
	port := strconv.Itoa(cfg.Port)
	log.Println("Server.Start: starting on port", port)
	PrintAddresses(port)
	log.Fatal("Server.Start", s.r.Run(":"+port))
}

func (s *Server) Stop() error {
	log.Println("Server.Stop: stopping")
	s.stop()
	<-s.h.Done
	return nil
}
