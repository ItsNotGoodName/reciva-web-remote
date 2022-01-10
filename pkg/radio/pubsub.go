package radio

import (
	"sync"
)

type Pub struct {
	subsMu sync.Mutex
	subs   map[*Sub]struct{}
}

func newPub() *Pub {
	return &Pub{
		subsMu: sync.Mutex{},
		subs:   make(map[*Sub]struct{}),
	}
}

func (p *Pub) Unsubscribe(sub *Sub) {
	p.subsMu.Lock()
	delete(p.subs, sub)
	sub.close()
	p.subsMu.Unlock()
}

func (p *Pub) Subscribe(sub *Sub) {
	p.subsMu.Lock()
	p.subs[sub] = struct{}{}
	p.subsMu.Unlock()
}

func (p *Pub) publish(state *State) {
	p.subsMu.Lock()

	for sub := range p.subs {
		if !sub.handle(state) {
			delete(p.subs, sub)
			sub.close()
		}
	}

	p.subsMu.Unlock()
}

type Sub struct {
	HandleC chan State

	mu     sync.RWMutex
	closed bool
}

func NewSub(handleC chan State) *Sub {
	return &Sub{
		mu:      sync.RWMutex{},
		closed:  false,
		HandleC: handleC,
	}
}

func (s *Sub) handle(state *State) bool {
	s.mu.RLock()
	select {
	case s.HandleC <- *state:
		s.mu.RUnlock()
		return true
	default:
		s.mu.RUnlock()
		return false
	}
}

func (s *Sub) close() {
	s.mu.RLock()
	if !s.closed {
		close(s.HandleC)
		s.closed = true
	}
	s.mu.RUnlock()
}
