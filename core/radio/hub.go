package radio

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ItsNotGoodName/reciva-web-remote/core/upnp"
)

var (
	ErrHubDiscovering   = fmt.Errorf("hub is discovering")
	ErrHubServiceClosed = fmt.Errorf("hub service closed")
	ErrRadioNotFound    = fmt.Errorf("radio not found")
)

type HubServiceImpl struct {
	discoverC    chan chan discoverResponse
	doneC        chan struct{}
	radioService CreateService

	radioMapMu sync.RWMutex
	radiosMap  map[string]Radio
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

func (hs *HubServiceImpl) Discover() (int, error) {
	resC := make(chan discoverResponse)
	select {
	case hs.discoverC <- resC:
		res := <-resC
		return res.count, res.err
	case <-hs.doneC:
		return 0, ErrHubServiceClosed
	default:
		return 0, ErrHubDiscovering
	}
}

func (hs *HubServiceImpl) Background(ctx context.Context, doneC chan<- struct{}) {
	var oldCancel context.CancelFunc = func() {}
	discover := func() (int, error) {
		log.Println("radio.HubService.Background: Discovering radios...")

		// Discover clients
		clients, _, err := upnp.Discover()
		if err != nil {
			return 0, err
		}

		newCtx, newCancel := context.WithCancel(ctx)

		// Create radios
		var radios []Radio
		for _, client := range clients {
			radio, err := hs.radioService.Create(newCtx, client)
			if err != nil {
				fmt.Println("radio.HubService.Background:", err)
				continue
			}

			radios = append(radios, radio)
		}

		// Create radios map
		radiosMap := make(map[string]Radio)
		for _, radio := range radios {
			radiosMap[radio.UUID] = radio
		}

		// Set radios map
		hs.radioMapMu.Lock()
		oldRadiosMap := hs.radiosMap
		hs.radiosMap = radiosMap
		hs.radioMapMu.Unlock()

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
			hs.radioMapMu.RLock()
			for _, r := range hs.radiosMap {
				<-r.Done()
			}
			hs.radioMapMu.RUnlock()

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
	hs.radioMapMu.RLock()
	r, ok := hs.radiosMap[uuid]
	hs.radioMapMu.RUnlock()
	if !ok {
		return Radio{}, ErrRadioNotFound
	}

	return r, nil
}

func (hs *HubServiceImpl) List() []Radio {
	hs.radioMapMu.RLock()
	var radios []Radio
	for _, r := range hs.radiosMap {
		radios = append(radios, r)
	}
	hs.radioMapMu.RUnlock()

	return radios
}
