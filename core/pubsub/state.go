package pubsub

import (
	"sync"

	"github.com/ItsNotGoodName/reciva-web-remote/core/state"
)

type StatePub struct {
	subsMu sync.Mutex
	subs   map[*StateSub]string
}

func NewStatePub() *StatePub {
	return &StatePub{
		subs: make(map[*StateSub]string),
	}
}

func (sp *StatePub) Unsubscribe(ss *StateSub) {
	sp.subsMu.Lock()
	delete(sp.subs, ss)
	ss.close()
	sp.subsMu.Unlock()
}

func (sp *StatePub) Subscribe(buffer int, uuid string) *StateSub {
	sub := newFragmentSub(buffer)

	sp.subsMu.Lock()
	sp.subs[sub] = uuid
	sp.subsMu.Unlock()

	return sub
}

func (sp *StatePub) Publish(s state.State) {
	sp.subsMu.Lock()
	for sub, uuid := range sp.subs {
		if uuid == s.UUID || uuid == "" {
			if !sub.send(&s) {
				delete(sp.subs, sub)
				sub.close()
			}
		}
	}
	sp.subsMu.Unlock()
}

type StateSub struct {
	channel chan *state.State
	open    bool
}

func newFragmentSub(buffer int) *StateSub {
	return &StateSub{
		channel: make(chan *state.State, buffer),
		open:    true,
	}
}

func (ss *StateSub) Channel() <-chan *state.State {
	return ss.channel
}

func (ss *StateSub) send(st *state.State) bool {
	select {
	case ss.channel <- st:
		return true
	default:
		return false
	}
}

func (ss *StateSub) close() {
	if ss.open {
		close(ss.channel)
		ss.open = false
	}
}
