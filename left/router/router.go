package router

import (
	"context"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/core/bus"
	"github.com/ItsNotGoodName/reciva-web-remote/core/preset"
	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/left/api"
	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func init() {
	mime.AddExtensionType(".js", "application/javascript")
}

type Router struct {
	port string
	r    chi.Router
}

func New(port string, p presenter.Presenter, fs fs.FS, hub radio.HubService, radioService radio.RadioService, busService bus.Service, presetStore preset.PresetStore) *Router {
	r := newRouter()
	upgrader := newUpgrader()

	// A good base middleware stack
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/ws", api.GetWS(upgrader, api.HandleWS(busService)))

		r.Get("/radios", p(api.GetRadios(hub, radioService)))
		r.Post("/radios", p(api.PostRadios(hub)))

		r.Route("/radio/{uuid}", func(r chi.Router) {
			r.Get("/", p(api.RequireRadio(hub, api.GetRadio(radioService))))
			r.Patch("/", p(api.RequireRadio(hub, api.PatchRadio(radioService))))
			r.Post("/", p(api.RequireRadio(hub, api.PostRadio(radioService))))
			r.Post("/volume", p(api.RequireRadio(hub, api.PostRadioVolume(radioService))))
		})

		r.Get("/presets", p(api.GetPresets(presetStore)))

		r.Get("/preset", p(api.GetPreset(presetStore)))
		r.Post("/preset", p(api.PostPreset(presetStore)))
	})

	mountFS(r, fs)

	mountPresets(r, presetStore)

	return &Router{
		port: port,
		r:    r,
	}
}

func (r *Router) Start() {
	printAddresses(r.port)
	if err := http.ListenAndServe(":"+r.port, r.r); err != nil {
		log.Fatalln("router.Router.Start:", err)
	}
}

func (r *Router) Background(ctx context.Context, doneC chan<- struct{}) {
	go r.Start()
	<-ctx.Done()
	doneC <- struct{}{}
}
