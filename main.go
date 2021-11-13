package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/server"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	// Create config
	cfg := config.NewConfig(config.WithFlag)

	// Show version and exit
	if cfg.ShowVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// Show info and exit
	if cfg.ShowInfo {
		fmt.Printf("Version: %s\nCommit: %s\nDate: %s\nBuilt by: %s\n", version, commit, date, builtBy)
		os.Exit(0)
	}

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
		log.Fatal("main(ERROR):", err)
	}
}
