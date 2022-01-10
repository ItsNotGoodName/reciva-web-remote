package radio

import (
	"fmt"
	"sync"
)

var ErrSubscriptionClosed = fmt.Errorf("subscription closed")

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

func (p *Pub) Subscribe(sub *Sub) error {
	var err error

	p.subsMu.Lock()
	if sub.open {
		p.subs[sub] = struct{}{}
	} else {
		err = ErrSubscriptionClosed
	}
	p.subsMu.Unlock()

	return err
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
	open    bool
}

func NewSub(handleC chan State) *Sub {
	return &Sub{
		HandleC: handleC,
		open:    true,
	}
}

// handle sends state to subscriber.
func (s *Sub) handle(state *State) bool {
	select {
	case s.HandleC <- *state:
		return true
	default:
		return false
	}
}

// close the sub.
func (s *Sub) close() {
	if s.open {
		close(s.HandleC)
		s.open = false
	}
}
