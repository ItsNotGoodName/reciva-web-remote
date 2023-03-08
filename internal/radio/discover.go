package radio

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ItsNotGoodName/go-upnpsub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/hub"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/state"
	"github.com/ItsNotGoodName/reciva-web-remote/internal/upnp"
	"github.com/sethvargo/go-retry"
)

func Discover(ctx context.Context, h *hub.Hub, cp upnpsub.ControlPoint) error {
	if !h.DiscoverMu.TryLock() {
		return internal.ErrHubDiscovering
	}
	defer h.DiscoverMu.Unlock()

	recivas, err := upnp.Discover()
	if err != nil {
		return err
	}

	// Create radios concurrently
	var wg sync.WaitGroup
	hubContext := h.Context()
	createContext, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()
	for i := range recivas {
		wg.Add(1)
		go (func(idx int) {
			// Create radio
			if err := createAndRun(createContext, hubContext, h, cp, recivas[idx]); err != nil {
				log.Println("radio.Discover:", err)
			}

			wg.Done()
		})(i)
	}
	wg.Wait()

	return nil
}

func createAndRun(ctx context.Context, radioContext context.Context, h *hub.Hub, cp upnpsub.ControlPoint, reciva upnp.Reciva) error {
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
	radio := h.Create(uuid, name, reciva, sub, stateC, updateFnC, close)
	go run(radioContext, radio, s, stateC, updateFnC)

	return nil
}
