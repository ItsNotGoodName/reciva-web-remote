package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/server"
)

func main() {
	// Create config
	cfg := config.NewConfig(config.WithFlag)

	// Create server
	s := server.NewServer(cfg)

	// Start server
	go s.Start(cfg)

	// Listen for interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Shutdown server
	if err := s.Stop(); err != nil {
		log.Fatal("main:", err)
	}
}
