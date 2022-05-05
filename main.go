package main

import (
	"fmt"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/core/app"
	"github.com/ItsNotGoodName/reciva-web-remote/core/background"
	"github.com/ItsNotGoodName/reciva-web-remote/core/bus"
	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
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
	version    = "dev"
	commit     = ""
	date       = ""
	builtBy    = "unknown"
	releaseURL = ""
	summary    = "dev"
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
		fmt.Printf("Version: %s\nCommit: %s\nDate: %s\nBuilt by: %s\nRelease url: %s\nSummary: %s\n", version, commit, date, builtBy, releaseURL, summary)
		return
	}

	// Right
	presetStore, err := file.NewPresetStore(cfg.ConfigFile)
	if err != nil {
		log.Fatalln("main.main: failed to create preset store:", err)
	}

	// Core
	middlewarePub := pubsub.NewSignalPub()
	middlewareAndPresetStore := middleware.NewPreset(middlewarePub, presetStore)
	radioService := radio.NewRadioService()
	statePub := pubsub.NewStatePub()
	runService := radio.NewRunService(middlewareAndPresetStore, middlewarePub, radioService, statePub)
	createService := radio.NewCreateService(upnpsub.NewControlPoint(upnpsub.WithPort(cfg.CPort)), runService)
	hubService := radio.NewHubService(createService)

	// App
	app := app.New(
		dto.Build{Version: version, Commit: commit, Date: date, BuiltBy: builtBy, ReleaseURL: releaseURL, Summary: summary},
		hubService,
		middlewareAndPresetStore,
		radioService,
	)

	// Bus
	bus := bus.New(app, statePub)

	// Left
	router := router.New(app, bus, cfg.PortStr, presenter.New(json.Render), web.FS())

	// Backgrounds
	backgrounds := []background.Background{
		createService,
		hubService,
		router,
	}

	// Run backgrounds
	background.Run(interrupt.Context(), backgrounds)
}
