package radio

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/ItsNotGoodName/reciva-web-remote/core/upnp"
)

type HubServiceImpl struct {
	discoverC    chan chan discoverResponse
	doneC        chan struct{}
	radioService CreateService

	radiosMapMu sync.RWMutex
	radiosMap   map[string]Radio
}

type discoverResponse struct {
	count int
	err   error
}

func NewHubService(createService CreateService) *HubServiceImpl {
	return &HubServiceImpl{
		discoverC:    make(chan chan discoverResponse),
		doneC:        make(chan struct{}),
		radioService: createService,
		radiosMap:    make(map[string]Radio),
	}
}

func (hs *HubServiceImpl) Discover(force bool) (int, error) {
	resC := make(chan discoverResponse)
	select {
	case hs.discoverC <- resC:
		res := <-resC
		return res.count, res.err
	case <-hs.doneC:
		return 0, core.ErrHubServiceClosed
	default:
		return 0, core.ErrHubDiscovering
	}
}

func (hs *HubServiceImpl) Background(ctx context.Context, doneC chan<- struct{}) {
	var oldCancel context.CancelFunc = func() {}
	discover := func() (int, error) {
		log.Println("radio.HubService.Background: Discovering radios...")

		// Discover radios
		recivas, err := upnp.Discover()
		if err != nil {
			return 0, err
		}

		// Radio run context
		newCtx, newCancel := context.WithCancel(ctx)

		// Create radios concurrently
		radioC := make(chan Radio)
		var wg sync.WaitGroup
		for i := range recivas {
			wg.Add(1)
			go (func(idx int) {
				// Timeout for creating radio
				sctx, cancel := context.WithTimeout(ctx, 25*time.Second)
				defer cancel()

				// Create radio
				radio, err := hs.radioService.Create(sctx, newCtx, recivas[idx])
				if err != nil {
					fmt.Println("radio.HubService.Background:", err)
				} else {
					radioC <- radio
				}

				wg.Done()
			})(i)
		}
		go (func() {
			wg.Wait()
			close(radioC)
		})()

		// Collect radios
		var radios []Radio
		for r := range radioC {
			radios = append(radios, r)
		}

		// Create radios map
		radiosMap := make(map[string]Radio)
		for _, radio := range radios {
			radiosMap[radio.UUID] = radio
		}

		// Set radios map
		hs.radiosMapMu.Lock()
		oldRadiosMap := hs.radiosMap
		hs.radiosMap = radiosMap
		hs.radiosMapMu.Unlock()

		// Close old radios
		oldCancel()
		for _, r := range oldRadiosMap {
			<-r.Done()
		}

		// Set old cancel
		oldCancel = newCancel

		count := len(radios)

		log.Println("radio.HubService.Background:", count, "radio(s) discovered")

		return count, nil
	}
	if _, err := discover(); err != nil {
		log.Println("radio.HubService.Background:", err)
	}

	for {
		select {
		case <-ctx.Done():
			close(hs.doneC)

			// Close radios
			oldCancel()
			hs.radiosMapMu.RLock()
			for _, r := range hs.radiosMap {
				<-r.Done()
			}
			hs.radiosMapMu.RUnlock()

			doneC <- struct{}{}
			return
		case resC := <-hs.discoverC:
			count, err := discover()
			if err != nil {
				log.Println("radio.HubService.Background:", err)
			}
			resC <- discoverResponse{count: count, err: err}
		}
	}
}

func (hs *HubServiceImpl) Get(uuid string) (Radio, error) {
	hs.radiosMapMu.RLock()
	r, ok := hs.radiosMap[uuid]
	hs.radiosMapMu.RUnlock()
	if !ok {
		return Radio{}, core.ErrRadioNotFound
	}

	return r, nil
}

func (hs *HubServiceImpl) List() []Radio {
	hs.radiosMapMu.RLock()
	var radios []Radio
	for _, r := range hs.radiosMap {
		radios = append(radios, r)
	}
	hs.radiosMapMu.RUnlock()

	// Sort radios
	sort.Slice(radios, func(i, j int) bool {
		return radios[i].UUID < radios[j].UUID
	})

	return radios
}
