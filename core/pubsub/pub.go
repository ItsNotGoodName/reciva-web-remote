package pubsub

import (
	"sync"

	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type Pub struct {
	subsMu sync.Mutex
	subs   map[*Sub]string
}

func NewPub() *Pub {
	return &Pub{
		subsMu: sync.Mutex{},
		subs:   make(map[*Sub]string),
	}
}

func (p *Pub) Unsubscribe(sub *Sub) {
	p.subsMu.Lock()
	delete(p.subs, sub)
	sub.close()
	p.subsMu.Unlock()
}

func (p *Pub) Subscribe(buffer int, uuid string) *Sub {
	sub := newSub(buffer)

	p.subsMu.Lock()
	p.subs[sub] = uuid
	p.subsMu.Unlock()

	return sub
}

func (p *Pub) Publish(f state.Fragment) {
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
