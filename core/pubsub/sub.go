package pubsub

import "github.com/ItsNotGoodName/reciva-web-remote/core/state"

type Sub struct {
	channel chan *state.Fragment
	open    bool
}

func newSub(buffer int) *Sub {
	return &Sub{
		channel: make(chan *state.Fragment, buffer),
		open:    true,
	}
}

func (s *Sub) Channel() <-chan *state.Fragment {
	return s.channel
}

func (s *Sub) send(f *state.Fragment) bool {
	select {
	case s.channel <- f:
		return true
	default:
		return false
	}
}

func (s *Sub) close() {
	if s.open {
		close(s.channel)
		s.open = false
	}
}
