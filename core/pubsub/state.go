package pubsub

import (
	"sync"

	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type StatePub struct {
	subsMapMu sync.Mutex
	subsMap   map[*chan state.PubMessage]string
}

func NewStatePub() *StatePub {
	return &StatePub{
		subsMap: make(map[*chan state.PubMessage]string),
	}
}

func (sp *StatePub) Subscribe(uuid string) (<-chan state.PubMessage, func()) {
	sub := make(chan state.PubMessage, 1)
	subRef := &sub

	sp.subsMapMu.Lock()
	sp.subsMap[subRef] = uuid
	sp.subsMapMu.Unlock()

	return sub, sp.unsubscribeFunc(subRef)
}

func (sp *StatePub) unsubscribeFunc(sub *chan state.PubMessage) func() {
	return func() {
		sp.subsMapMu.Lock()
		delete(sp.subsMap, sub)
		sp.subsMapMu.Unlock()

		select {
		case <-*sub:
		default:
		}
	}
}

func (sp *StatePub) Publish(s state.State, changed state.Changed) {
	msg := state.PubMessage{State: s, Changed: changed}

	sp.subsMapMu.Lock()
	for sub, uuid := range sp.subsMap {
		if uuid == s.UUID {
			// Send latest
			select {
			case *sub <- msg:
			default:
				select {
				case old := <-*sub:
					msg.Changed = old.Changed.Merge(msg.Changed)
					*sub <- msg
					msg.Changed = changed
				case *sub <- msg:
				}
			}
		}
	}
	sp.subsMapMu.Unlock()
}
