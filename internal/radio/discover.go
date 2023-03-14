package radio

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/model"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/pubsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/upnp"
	"github.com/sethvargo/go-retry"
)

type DefaultStateHook struct{}

func (DefaultStateHook) OnChanged(ctx context.Context, s *state.State, c state.Changed) state.Changed {
	return c
}
func (DefaultStateHook) OnStart(ctx context.Context, s *state.State, c state.Changed) state.Changed {
	return c
}

type Discoverer struct {
	mu           sync.Mutex
	ctxC         chan context.Context
	hub          *hub.Hub
	controlPoint upnpsub.ControlPoint
	stateHook    StateHook
}

func NewDiscoverer(hub *hub.Hub, controlPoint upnpsub.ControlPoint, stateHook StateHook) *Discoverer {
	return &Discoverer{
		mu:           sync.Mutex{},
		ctxC:         make(chan context.Context),
		hub:          hub,
		controlPoint: controlPoint,
		stateHook:    stateHook,
	}
}

func (d *Discoverer) Background(ctx context.Context, doneC chan<- struct{}) {
	ctxDoneC := make(chan struct{})
	go func() {
		for {
			select {
			case <-ctxDoneC:
				return
			case d.ctxC <- ctx:
			}
		}
	}()

	// Wait for context
	<-ctx.Done()
	// Stop further discoveries
	d.mu.Lock()
	// Stop ctx provider
	close(ctxDoneC)
	// Done
	doneC <- struct{}{}
}

func (d *Discoverer) Discovering() bool {
	if d.mu.TryLock() {
		d.mu.Unlock()
		return false
	}

	return true
}

func (d *Discoverer) Discover(ctx context.Context) error {
	if !d.mu.TryLock() {
		return internal.ErrHubDiscovering
	}
	defer d.mu.Unlock()

	pubsub.DefaultPub.Publish(pubsub.DiscoverTopic, pubsub.DiscoverMessage{Discovering: true})
	defer pubsub.DefaultPub.Publish(pubsub.DiscoverTopic, pubsub.DiscoverMessage{Discovering: false})
	defer pubsub.DefaultPub.Publish(pubsub.StaleTopic, model.StaleRadios)

	radioContext := <-d.ctxC
	return discover(ctx, d.hub, func(ctx context.Context, reciva upnp.Reciva) error {
		return create(ctx, radioContext, reciva, d.hub, d.controlPoint, d.stateHook)
	})
}

func discover(ctx context.Context, h *hub.Hub, create func(createContext context.Context, reciva upnp.Reciva) error) error {
	recivas, err := upnp.Discover()
	if err != nil {
		return err
	}

	// Create radios concurrently
	var wg sync.WaitGroup
	createContext, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()
	for i := range recivas {
		wg.Add(1)
		go (func(idx int) {
			// Create radio
			if err := create(createContext, recivas[idx]); err != nil {
				log.Println("radio.Discover:", err)
			}

			wg.Done()
		})(i)
	}
	wg.Wait()

	return nil
}

func create(ctx context.Context, radioContext context.Context, reciva upnp.Reciva, h *hub.Hub, cp upnpsub.ControlPoint, stateHook StateHook) error {
	// Get UUID
	uuid, err := reciva.GetUUID()
	if err != nil {
		return err
	}

	// Get name
	name := reciva.GetName()

	// Get audio sources
	audioSources, err := reciva.GetAudioSources(ctx)
	if err != nil {
		return err
	}

	// Create state
	s := state.New(uuid, name, reciva.GetModelName(), reciva.GetModelNumber(), audioSources)

	// Get and set volume
	volume, err := reciva.GetVolume(ctx)
	if err != nil {
		return err
	}
	s.SetVolume(volume)

	// Get and parse presets count
	presetsCount, err := reciva.GetNumberOfPresets(ctx)
	if err != nil {
		return err
	}
	if presetsCount, err = state.ParsePresetsCount(presetsCount); err != nil {
		return err
	}

	// Get and set presets
	var presets []state.Preset
	for i := 1; i <= presetsCount; i++ {
		p, err := reciva.GetPreset(ctx, i)
		if err != nil {
			return err
		}

		presets = append(presets, state.NewPreset(i, p.Name, p.URL))
	}
	s.SetPresets(presets)

	radioContext, close := context.WithCancel(radioContext)

	// Create subscription
	eventURL := reciva.GetEventURL()
	var sub upnpsub.Subscription
	err = retry.Do(ctx, retry.WithMaxRetries(3, retry.NewFibonacci(time.Second)), func(ctx context.Context) error {
		var err2 error
		sub, err2 = cp.Subscribe(radioContext, &eventURL)
		return retry.RetryableError(err2)
	})
	if err != nil {
		close()
		return err
	}

	// Create and run radio
	stateC := make(hub.RadioStateC)
	updateFnC := make(hub.RadioUpdateFnC)

	radio, err := h.Create(uuid, name, reciva, sub, stateC, updateFnC, close)
	if err != nil {
		close()
		return err
	}

	go run(radioContext, radio, s, stateC, updateFnC, stateHook)

	return nil
}
