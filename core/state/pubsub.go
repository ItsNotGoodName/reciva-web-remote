package state

import (
	"sync"
)

type PubImpl struct {
	subsMu sync.Mutex
	subs   map[*SubImpl]string
}

func NewPub() *PubImpl {
	return &PubImpl{
		subsMu: sync.Mutex{},
		subs:   make(map[*SubImpl]string),
	}
}

func (p *PubImpl) Unsubscribe(sub *SubImpl) {
	p.subsMu.Lock()
	delete(p.subs, sub)
	sub.close()
	p.subsMu.Unlock()
}

func (p *PubImpl) Subscribe(buffer int, uuid string) *SubImpl {
	sub := newSub(buffer)

	p.subsMu.Lock()
	p.subs[sub] = uuid
	p.subsMu.Unlock()

	return sub
}

func (p *PubImpl) Publish(f Fragment) {
	p.subsMu.Lock()
	for sub, uuid := range p.subs {
		if uuid == f.UUID || uuid == "" {
			if !sub.send(&f) {
				delete(p.subs, sub)
				sub.close()
			}
		}
	}
	p.subsMu.Unlock()
}

type SubImpl struct {
	channel chan *Fragment
	open    bool
}

func newSub(buffer int) *SubImpl {
	return &SubImpl{
		channel: make(chan *Fragment, buffer),
		open:    true,
	}
}

func (s *SubImpl) Channel() <-chan *Fragment {
	return s.channel
}

func (s *SubImpl) send(f *Fragment) bool {
	select {
	case s.channel <- f:
		return true
	default:
		return false
	}
}

func (s *SubImpl) close() {
	if s.open {
		close(s.channel)
		s.open = false
	}
}
