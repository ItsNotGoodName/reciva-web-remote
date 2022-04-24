package main

import (
	"fmt"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/core/background"
	"github.com/ItsNotGoodName/reciva-web-remote/core/middleware"
	"github.com/ItsNotGoodName/reciva-web-remote/core/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/left/json"
	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
	"github.com/ItsNotGoodName/reciva-web-remote/left/router"
	"github.com/ItsNotGoodName/reciva-web-remote/left/web"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/interrupt"
	"github.com/ItsNotGoodName/reciva-web-remote/right/file"
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

	// Backgrounds
	var backgrounds []background.Background

	// Right
	presetStore, err := file.NewPresetStore(cfg.ConfigFile)
	if err != nil {
		log.Fatal("Failed to create preset store:", err)
	}

	// Core
	middlewarePub := pubsub.NewSignalPub()
	middlewareAndPresetStore := middleware.NewPreset(middlewarePub, presetStore)
	statePub := pubsub.NewStatePub()
	runService := radio.NewRunService(statePub, middlewareAndPresetStore, middlewarePub)
	radioService := radio.NewRadioService()
	createService := radio.NewCreateService(upnpsub.NewControlPoint(upnpsub.WithPort(cfg.CPort)), runService)
	backgrounds = append(backgrounds, createService)
	hubService := radio.NewHubService(createService)
	backgrounds = append(backgrounds, hubService)

	// Left
	router := router.New(cfg.PortStr, presenter.New(json.Render), web.FS(), hubService, radioService)
	backgrounds = append(backgrounds, router)

	// Run backgrounds
	background.Run(interrupt.Context(), backgrounds)
}
