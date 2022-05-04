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
	handle := func(c state.Changed) {
		c = c.Merge(rs.middleware.Apply(dctx, &s, c))
		if c != 0 {
			rs.statePub.Publish(s, c)
		}
	}
	handle(state.ChangedAll)

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
			handle(state.ChangedAll)
		case radio.stateC <- s:
		case fn := <-radio.updateFnC:
			handle(fn(&s))
		case event := <-radio.subscription.Events():
			handle(parseEvent(event, &s))
		}
	}
}
