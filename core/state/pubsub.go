package state

import (
	"fmt"
	"sync"
)

var ErrSubClosed = fmt.Errorf("sub closed")

type PubImpl struct {
	subsMu sync.Mutex
	subs   map[*SubImpl]struct{}
}

func NewPub() *PubImpl {
	return &PubImpl{
		subsMu: sync.Mutex{},
		subs:   make(map[*SubImpl]struct{}),
	}
}

func (p *PubImpl) Unsubscribe(sub *SubImpl) {
	p.subsMu.Lock()
	delete(p.subs, sub)
	sub.close()
	p.subsMu.Unlock()
}

func (p *PubImpl) Subscribe(buffer int) (*SubImpl, error) {
	sub := newSub(buffer)

	var err error
	p.subsMu.Lock()
	if sub.open {
		p.subs[sub] = struct{}{}
	} else {
		err = ErrSubClosed
	}
	p.subsMu.Unlock()

	return sub, err
}

func (p *PubImpl) Publish(f Fragment) {
	p.subsMu.Lock()
	for sub := range p.subs {
		if !sub.send(&f) {
			delete(p.subs, sub)
			sub.close()
		}
	}
	p.subsMu.Unlock()
}

type SubImpl struct {
	ch   chan *Fragment
	open bool
}

func newSub(buffer int) *SubImpl {
	return &SubImpl{
		ch:   make(chan *Fragment, buffer),
		open: true,
	}
}

// Channel returns the sub's channel.
func (s *SubImpl) Channel() <-chan *Fragment {
	return s.ch
}

// send sends fragment to subscriber.
func (s *SubImpl) send(f *Fragment) bool {
	select {
	case s.ch <- f:
		return true
	default:
		return false
	}
}

// close the sub.
func (s *SubImpl) close() {
	if s.open {
		close(s.ch)
		s.open = false
	}
}
