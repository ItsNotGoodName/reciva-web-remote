package router

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/core/radio"
	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type Router struct {
	port string
	r    chi.Router
}

func New(port string, h presenter.Presenter, hub radio.HubService, radioService radio.RadioService) *Router {
	r := chi.NewRouter()
	upgrader := websocket.Upgrader{}

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/ws", GetWS(hub, upgrader))

		r.Get("/radios", h(GetRadios(hub, radioService)))
		r.Post("/radios", h(PostRadios(hub)))

		r.Route("/radio/{uuid}", func(r chi.Router) {
			r.Get("/", h(RequireRadio(hub, GetRadio(radioService))))
			r.Patch("/", h(RequireRadio(hub, PatchRadio(radioService))))
			r.Post("/", h(RequireRadio(hub, PostRadio(radioService))))
			r.Post("/volume", h(RequireRadio(hub, PostRadioVolume(radioService))))
		})
	})

	return &Router{
		port: port,
		r:    r,
	}
}

func (r *Router) MountFS(fs fs.FS) {
	mountFS(r.r, fs)
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
