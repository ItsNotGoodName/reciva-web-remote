// Packages signal is a pubsub that notifies subscribers in other goroutines.
package signal

import "sync"

type Pub struct {
	subsMapMu sync.RWMutex
	subsMap   map[*Sub]struct{}
}

func NewPub() *Pub {
	return &Pub{
		subsMap: make(map[*Sub]struct{}),
	}
}

func (p *Pub) Publish() {
	p.subsMapMu.RLock()

	for sub := range p.subsMap {
		select {
		case sub.ch <- struct{}{}:
		default:
		}
	}

	p.subsMapMu.RUnlock()
}

func (p *Pub) Subscribe() *Sub {
	p.subsMapMu.Lock()

	sub := newSub()
	p.subsMap[sub] = struct{}{}

	p.subsMapMu.Unlock()

	return sub
}

func (p *Pub) Unsubscribe(sub *Sub) {
	p.subsMapMu.Lock()

	delete(p.subsMap, sub)

	p.subsMapMu.Unlock()
}
