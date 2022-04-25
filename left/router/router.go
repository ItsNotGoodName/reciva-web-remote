package router

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/core/app"
	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/left/api"
	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router struct {
	port string
	r    chi.Router
}

func New(port string, h presenter.Presenter, fs fs.FS, hub radio.HubService, radioService radio.RadioService, application *app.App) *Router {
	r := newMux()
	upgrader := newUpgrader()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/ws", api.GetWS(upgrader, api.HandleWS(application)))

		r.Get("/radios", h(api.GetRadios(hub, radioService)))
		r.Post("/radios", h(api.PostRadios(hub)))

		r.Route("/radio/{uuid}", func(r chi.Router) {
			r.Get("/", h(api.RequireRadio(hub, api.GetRadio(radioService))))
			r.Patch("/", h(api.RequireRadio(hub, api.PatchRadio(radioService))))
			r.Post("/", h(api.RequireRadio(hub, api.PostRadio(radioService))))
			r.Post("/volume", h(api.RequireRadio(hub, api.PostRadioVolume(radioService))))
		})
	})

	mountFS(r, fs)

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
