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
	statePub      state.Pub
}

func NewRunService(middleware state.Middleware, middlewarePub state.MiddlewarePub, radioService RadioService, statePub state.Pub) *RunServiceImpl {
	return &RunServiceImpl{
		middleware:    middleware,
		middlewarePub: middlewarePub,
		radioService:  radioService,
		statePub:      statePub,
	}
}

func (rs *RunServiceImpl) Run(dctx context.Context, radio Radio, s state.State) {
	handle := func(frag state.Fragment, changed int) {
		rs.middleware.Apply(&frag)
		if changed = s.Merge(frag) | changed; changed != 0 {
			rs.statePub.Publish(s, changed)
		}
	}
	handle(s.Fragment(), state.ChangedAll)

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
			handle(s.Fragment(), 0)
		case radio.readC <- s:
		case frag := <-radio.updateC:
			handle(frag, 0)
		case event := <-radio.subscription.Events():
			frag := state.NewFragment(radio.UUID)
			parseEvent(event, &frag)
			handle(frag, 0)
		}
	}
}
