// Package sig is a signal pubsub that notifies subscribers in other goroutines.
package sig

import "sync"

type Pub struct {
	subsMapMu sync.Mutex
	subsMap   map[*chan struct{}]struct{}
}

func NewPub() *Pub {
	return &Pub{
		subsMap: make(map[*chan struct{}]struct{}),
	}
}

func (p *Pub) Publish() {
	p.subsMapMu.Lock()

	for sub := range p.subsMap {
		select {
		case *sub <- struct{}{}:
		default:
		}
	}

	p.subsMapMu.Unlock()
}

func (p *Pub) Subscribe() (<-chan struct{}, func()) {
	p.subsMapMu.Lock()

	sub := make(chan struct{}, 1)
	p.subsMap[&sub] = struct{}{}

	p.subsMapMu.Unlock()

	return sub, p.unsubscribeFunc(&sub)

}

func (p *Pub) unsubscribeFunc(ch *chan struct{}) func() {
	return func() {
		p.subsMapMu.Lock()

		delete(p.subsMap, ch)

		p.subsMapMu.Unlock()
	}
}
