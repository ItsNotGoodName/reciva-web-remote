package cmd

import (
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/http"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/background"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/middleware"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/store"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/upnp"
	"github.com/ItsNotGoodName/reciva-web-remote/left/web"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/interrupt"
)

func Server(cfg *config.Config) {
	ctx := interrupt.Context()

	store := middleware.NewStore(store.Must(store.NewFile(cfg.ConfigFile)))

	hub := hub.New()
	controlPoint := upnpsub.NewControlPoint(upnpsub.WithPort(cfg.CPort))
	stateHook := middleware.NewStateHook(store)
	discoverer := radio.NewDiscoverer(hub, controlPoint, stateHook)
	api := http.NewAPI(hub, discoverer, store)

	go func() {
		if err := discoverer.Discover(ctx); err != nil {
			log.Println("cmd.Server:", err)
		}
	}()
	go http.Start(api, cfg.Port, web.FS())

	<-background.Run(ctx, []background.Background{hub, upnp.NewBackgroundControlPoint(controlPoint), discoverer})
}
