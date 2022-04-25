package pubsub

import (
	"sync"

	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type StatePub struct {
	subMapMu sync.Mutex
	subsMap  map[*chan state.Message]string
}

func NewStatePub() *StatePub {
	return &StatePub{
		subsMap: make(map[*chan state.Message]string),
	}
}

func (sp *StatePub) Subscribe(uuid string) (<-chan state.Message, func()) {
	sub := make(chan state.Message, 1)

	sp.subMapMu.Lock()

	sp.subsMap[&sub] = uuid

	sp.subMapMu.Unlock()

	return sub, sp.unsubscribeFunc(&sub)
}

func (sp *StatePub) unsubscribeFunc(sub *chan state.Message) func() {
	return func() {
		sp.subMapMu.Lock()

		delete(sp.subsMap, sub)

		sp.subMapMu.Unlock()
	}
}

func (sp *StatePub) Publish(s state.State, changed int) {
	sp.subMapMu.Lock()

	msg := state.Message{
		State: s,
	}
	for sub, uuid := range sp.subsMap {
		if uuid == s.UUID {
			msg.Changed = changed

			// Send latest
			select {
			case *sub <- msg:
			default:
				select {
				case old := <-*sub:
					msg.Changed |= old.Changed
					*sub <- msg
				case *sub <- msg:
				}
			}
		}
	}

	sp.subMapMu.Unlock()
}
