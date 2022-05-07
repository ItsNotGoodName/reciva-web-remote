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
		r.Get("/build", p(api.Build(app)))

		// Radio
		r.Route("/radios", func(r chi.Router) {
			r.Get("/", p(api.RadioList(app)))
			r.Post("/", p(api.RadioDiscover(app)))

			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", p(api.UUIDRequire(api.RadioGet(app))))
				r.Post("/subscription", p(api.UUIDRequire(api.RadioRefreshSubscription(app))))
				r.Post("/volume", p(api.UUIDRequire(api.RadioRefreshVolume(app))))
			})
		})

		// States
		r.Route("/states", func(r chi.Router) {
			r.Get("/", p(api.StateList(app)))
			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", p(api.UUIDRequire(api.StateGet(app))))
				r.Patch("/", p(api.UUIDRequire(api.StatePatch(app))))
			})
		})

		// Presets
		r.Route("/presets", func(r chi.Router) {
			r.Get("/", p(api.PresetList(app)))
			r.Post("/", p(api.PresetUpdate(app)))
			r.Get("/*", p(api.PresetGet(app)))
		})

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
