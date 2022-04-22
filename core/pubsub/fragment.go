package pubsub

import (
	"sync"

	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type FragmentPub struct {
	subsMu sync.Mutex
	subs   map[*FragmentSub]string
}

func NewFragmentPub() *FragmentPub {
	return &FragmentPub{
		subsMu: sync.Mutex{},
		subs:   make(map[*FragmentSub]string),
	}
}

func (p *FragmentPub) Unsubscribe(sub *FragmentSub) {
	p.subsMu.Lock()
	delete(p.subs, sub)
	sub.close()
	p.subsMu.Unlock()
}

func (p *FragmentPub) Subscribe(buffer int, uuid string) *FragmentSub {
	sub := newFragmentSub(buffer)

	p.subsMu.Lock()
	p.subs[sub] = uuid
	p.subsMu.Unlock()

	return sub
}

func (p *FragmentPub) Publish(frag state.Fragment) {
	p.subsMu.Lock()
	for sub, uuid := range p.subs {
		if uuid == frag.UUID || uuid == "" {
			if !sub.send(&frag) {
				delete(p.subs, sub)
				sub.close()
			}
		}
	}
	p.subsMu.Unlock()
}

type FragmentSub struct {
	channel chan *state.Fragment
	open    bool
}

func newFragmentSub(buffer int) *FragmentSub {
	return &FragmentSub{
		channel: make(chan *state.Fragment, buffer),
		open:    true,
	}
}

func (s *FragmentSub) Channel() <-chan *state.Fragment {
	return s.channel
}

func (s *FragmentSub) send(f *state.Fragment) bool {
	select {
	case s.channel <- f:
		return true
	default:
		return false
	}
}

func (s *FragmentSub) close() {
	if s.open {
		close(s.channel)
		s.open = false
	}
}
