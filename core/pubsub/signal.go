package pubsub

import "sync"

type SignalPub struct {
	subsMapMu sync.Mutex
	subsMap   map[*chan struct{}]struct{}
}

func NewSignalPub() *SignalPub {
	return &SignalPub{
		subsMap: make(map[*chan struct{}]struct{}),
	}
}

func (sp *SignalPub) Publish() {
	sp.subsMapMu.Lock()
	for sub := range sp.subsMap {
		select {
		case *sub <- struct{}{}:
		default:
		}
	}
	sp.subsMapMu.Unlock()
}

func (sp *SignalPub) Subscribe() (<-chan struct{}, func()) {
	sub := make(chan struct{}, 1)
	subRef := &sub

	sp.subsMapMu.Lock()
	sp.subsMap[subRef] = struct{}{}
	sp.subsMapMu.Unlock()

	return sub, sp.unsubscribeFunc(subRef)

}

func (sp *SignalPub) unsubscribeFunc(sub *chan struct{}) func() {
	return func() {
		sp.subsMapMu.Lock()
		delete(sp.subsMap, sub)
		sp.subsMapMu.Unlock()
	}
}
