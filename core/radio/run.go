package radio

import (
	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
	"github.com/ItsNotGoodName/reciva-web-remote/pkg/sig"
)

type RunServiceImpl struct {
	fragmentPub   state.FragmentPub
	middleware    state.Middleware
	middlewarePub *sig.Pub
}

func NewRunService(fragmentPub state.FragmentPub, middleware state.Middleware, middlewarePub *sig.Pub) *RunServiceImpl {
	return &RunServiceImpl{
		fragmentPub:   fragmentPub,
		middleware:    middleware,
		middlewarePub: middlewarePub,
	}
}

func (rs *RunServiceImpl) Run(radio Radio, s state.State) {
	handleWithoutMiddleware := func(frag state.Fragment) {
		if f, changed := s.Merge(frag); changed {
			rs.fragmentPub.Publish(f)
		}
	}
	handle := func(frag state.Fragment) {
		rs.middleware.Fragment(&frag)
		handleWithoutMiddleware(frag)
	}

	s.Merge(rs.middleware.FragmentFromState(s))
	middlewareSub := rs.middlewarePub.Subscribe()

	for {
		select {
		case <-radio.Done():
			rs.middlewarePub.Unsubscribe(middlewareSub)
			return
		case <-middlewareSub.Channel():
			handleWithoutMiddleware(rs.middleware.FragmentFromState(s))
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
