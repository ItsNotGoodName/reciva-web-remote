package radio

import (
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/sig"
)

type RunServiceImpl struct {
	statePub      state.StatePub
	middleware    state.Middleware
	middlewarePub *sig.Pub
}

func NewRunService(statePub state.StatePub, middleware state.Middleware, middlewarePub *sig.Pub) *RunServiceImpl {
	return &RunServiceImpl{
		statePub:      statePub,
		middleware:    middleware,
		middlewarePub: middlewarePub,
	}
}

func (rs *RunServiceImpl) Run(radio Radio, s state.State) {
	handle := func(frag state.Fragment) {
		rs.middleware.Apply(&frag)
		if s.Merge(frag) {
			rs.statePub.Publish(s)
		}
	}
	handle(s.Fragment())

	// Middleware signal
	middlewareSub, middlewareUnsub := rs.middlewarePub.Subscribe()
	defer middlewareUnsub()

	for {
		select {
		case <-radio.Done():
			return
		case <-middlewareSub:
			handle(s.Fragment())
		case radio.readC <- s:
		case fragment := <-radio.updateC:
			handle(fragment)
		case event := <-radio.subscription.Events():
			fragment := state.NewFragment(radio.UUID)
			parseEvent(event, &fragment)
			handle(fragment)
		}
	}
}
