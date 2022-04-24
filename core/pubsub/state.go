package pubsub

import (
	"sync"

	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type StatePub struct {
	subMapMu sync.Mutex
	subsMap  map[*chan state.State]string
}

func NewStatePub() *StatePub {
	return &StatePub{
		subsMap: make(map[*chan state.State]string),
	}
}

func (sp *StatePub) Subscribe(uuid string) (<-chan state.State, func()) {
	sub := make(chan state.State, 1)

	sp.subMapMu.Lock()

	sp.subsMap[&sub] = uuid

	sp.subMapMu.Unlock()

	return sub, sp.unsubscribeFunc(&sub)
}

func (sp *StatePub) unsubscribeFunc(sub *chan state.State) func() {
	return func() {
		sp.subMapMu.Lock()

		delete(sp.subsMap, sub)

		sp.subMapMu.Unlock()
	}
}

func (sp *StatePub) Publish(s state.State) {
	sp.subMapMu.Lock()

	for sub, uuid := range sp.subsMap {
		if uuid == s.UUID {
			// Send latest
			select {
			case *sub <- s:
			default:
				select {
				case <-*sub:
					*sub <- s
				case *sub <- s:
				}
			}
		}
	}

	sp.subMapMu.Unlock()
}
