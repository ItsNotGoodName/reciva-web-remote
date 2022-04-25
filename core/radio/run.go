package radio

import (
	"context"
	"log"
	"time"

	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type RunServiceImpl struct {
	middleware    state.Middleware
	middlewarePub state.MiddlewarePub
	radioService  RadioService
	statePub      state.StatePub
}

func NewRunService(middleware state.Middleware, middlewarePub state.MiddlewarePub, radioService RadioService, statePub state.StatePub) *RunServiceImpl {
	return &RunServiceImpl{
		middleware:    middleware,
		middlewarePub: middlewarePub,
		radioService:  radioService,
		statePub:      statePub,
	}
}

func (rs *RunServiceImpl) Run(dctx context.Context, radio Radio, s state.State) {
	handle := func(frag state.Fragment) {
		rs.middleware.Apply(&frag)
		if s.Merge(frag) {
			rs.statePub.Publish(s)
		}
	}
	handle(s.Fragment())

	middlewareSub, middlewareUnsub := rs.middlewarePub.Subscribe()
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		middlewareUnsub()
		ticker.Stop()
	}()

	for {
		select {
		case <-radio.Done():
			return
		case <-ticker.C:
			go func() {
				if err := rs.radioService.RefreshVolume(dctx, radio); err != nil {
					log.Println("radio.RunService.Run: failed to refresh volume:", err)
				}
			}()
		case <-middlewareSub:
			handle(s.Fragment())
		case radio.readC <- s:
		case frag := <-radio.updateC:
			handle(frag)
		case event := <-radio.subscription.Events():
			fragment := state.NewFragment(radio.UUID)
			parseEvent(event, &fragment)
			handle(fragment)
		}
	}
}
