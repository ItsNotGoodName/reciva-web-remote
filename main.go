package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/router"
	"github.com/ItsNotGoodName/reciva-web-remote/server"
	"github.com/ItsNotGoodName/reciva-web-remote/store"
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
		return
	}

	// Show info and exit
	if cfg.ShowInfo {
		fmt.Printf("Version: %s\nCommit: %s\nDate: %s\nBuilt by: %s\n", version, commit, date, builtBy)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Create and start preset store
	presetStore := store.New(cfg.ConfigFile)
	go presetStore.Start(ctx)

	// Create Mutator
	mutator := store.NewMutator(presetStore)

	// Create and start control point
	cp := upnpsub.NewControlPointWithPort(cfg.CPort)
	go cp.Start()

	// Create and start hub
	hub := radio.NewHubWithMutator(cp, mutator)
	go hub.Start(ctx)

	// Create engine
	engine := router.NewEngine()

	// Create WS upgrader
	upgrader := router.NewUpgrader()

	// Create routes
	server.Route(
		engine,
		upgrader,
		hub,
		presetStore,
	)

	// Start engine
	go router.Start(engine, cfg.PortStr)

	// Listen for interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Stop hub and store
	cancel()
	<-hub.DoneC
	<-presetStore.DoneC
}
