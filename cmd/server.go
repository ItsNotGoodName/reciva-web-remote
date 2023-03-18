package cmd

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/config"
	"github.com/ItsNotGoodName/reciva-web-remote/http"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/middleware"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/store"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/upnp"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/background"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/interrupt"
	"github.com/ItsNotGoodName/reciva-web-remote/web"
)

func Server(cfg *config.Config) {
	ctx, cancel := context.WithCancel(interrupt.Context())
	defer cancel()

	store := middleware.NewStore(store.Must(store.NewFile(cfg.ConfigFile)))

	hub := hub.New()
	controlPoint := upnpsub.NewControlPoint(upnpsub.WithPort(cfg.CPort))
	stateHook := middleware.NewStateHook(store)
	discoverer := radio.NewDiscoverer(hub, controlPoint, stateHook)
	api := http.NewAPI(hub, discoverer, store)
	router := http.Router(api, web.FS())

	<-background.Run(ctx, []background.Background{
		hub,
		upnp.NewBackgroundControlPoint(controlPoint),
		discoverer,
		background.NewFunction(func(ctx context.Context) {
			if err := discoverer.Discover(ctx, true); err != nil {
				log.Println("cmd.Server:", err)
			}
		}),
		background.NewFunction(func(ctx context.Context) {
			log.Println("cmd.Server:", http.Start(router, cfg.Port))
			cancel()
		}),
	})
}
