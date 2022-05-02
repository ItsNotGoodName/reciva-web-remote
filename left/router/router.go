package router

import (
	"context"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
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

func New(app dto.App, bus dto.Bus, port string, p presenter.Presenter, fs fs.FS) *Router {
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
		// Build
		r.Get("/build", p(api.GetBuild(app)))

		// Radios
		r.Get("/radios", p(api.GetRadios(app)))
		r.Post("/radios", p(api.PostRadios(app)))

		getState := p(api.RequireUUID(api.GetState(app)))
		patchState := p(api.RequireUUID(api.PatchState(app)))

		// Radio
		r.Route("/radio/{uuid}", func(r chi.Router) {
			r.Get("/", p(api.RequireUUID(api.GetRadio(app))))
			r.Get("/state", getState)
			r.Patch("/state", patchState)
			r.Post("/subscription", p(api.RequireUUID(api.PostRadioSubscription(app))))
			r.Post("/volume", p(api.RequireUUID(api.PostRadioVolume(app))))
		})

		// States
		r.Get("/states", p(api.GetStates(app)))

		// State
		r.Route("/state/{uuid}", func(r chi.Router) {
			r.Get("/", getState)
			r.Patch("/", patchState)
		})

		// Presets
		r.Get("/presets", p(api.GetPresets(app)))

		// Preset
		r.Get("/preset", p(api.GetPreset(app)))
		r.Post("/preset", p(api.PostPreset(app)))

		// WS
		r.Get("/ws", api.GetWS(upgrader, api.HandleWS(bus)))
	})

	mountFS(r, fs)

	mountPresets(r, app)

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
	doneC <- struct{}{}
}
