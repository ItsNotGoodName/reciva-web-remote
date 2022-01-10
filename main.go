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

	// Create and start control point
	cp := upnpsub.NewControlPointWithPort(cfg.CPort)
	go cp.Start()

	// Create and start hub
	ctx, cancel := context.WithCancel(context.Background())
	hub := radio.NewHub(cp)
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
	)

	// Start engine
	go router.Start(engine, cfg.PortStr)

	// Listen for interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Stop hub
	cancel()
	<-hub.Done
}
