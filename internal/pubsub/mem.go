package pubsub

import (
	"sync"
	"time"
)

type MemPub struct {
	mu      sync.Mutex
	subsMap map[Topic]*memSub
	timer   *time.Timer
}

type MemSub struct {
	Message  chan Message
	Overflow chan struct{}
}

type memSub struct {
	messageC  chan Message
	overflowC chan struct{}
	next      *memSub
}

func NewMemPub() *MemPub {
	return &MemPub{
		mu:      sync.Mutex{},
		subsMap: make(map[Topic]*memSub),
		timer:   time.NewTimer(0),
	}
}

func (mp *MemPub) Subscribe(topics []Topic) (MemSub, func()) {
	return mp.SubscribeWithBuffer(topics, 100)
}

func (mp *MemPub) SubscribeWithBuffer(topics []Topic, buffer int) (MemSub, func()) {
	sub := MemSub{Message: make(chan Message, buffer), Overflow: make(chan struct{}, 1)}

	mp.mu.Lock()
	subs := mp.subscribe(topics, sub)
	mp.mu.Unlock()

	return sub, mp.unsubscribeFunc(topics, subs)
}

func (mp *MemPub) Resubscribe(topics []Topic, sub MemSub, unsub func()) func() {
	unsub()
Loop:
	for {
		select {
		case <-sub.Message:
		default:
			break Loop
		}
	}
	select {
	case <-sub.Overflow:
	default:
	}

	mp.mu.Lock()
	subs := mp.subscribe(topics, sub)
	mp.mu.Unlock()

	return mp.unsubscribeFunc(topics, subs)
}

func (mp *MemPub) subscribe(topics []Topic, sub MemSub) []*memSub {
	subs := []*memSub{}
	for _, topic := range topics {
		sub := &memSub{messageC: sub.Message, overflowC: sub.Overflow}
		subs = append(subs, sub)
		if next, ok := mp.subsMap[topic]; ok {
			sub.next = next
		}

		mp.subsMap[topic] = sub
	}
	return subs
}

func (mp *MemPub) unsubscribeFunc(topics []Topic, sub []*memSub) func() {
	return func() {
		mp.mu.Lock()
		for i, sub := range sub {
			topic := topics[i]
			// There should only be 1 or 0 sub in next because the unsubscribe function might be called twice
			next := mp.subsMap[topic]
			if next == nil {
				continue
			}
			if next == sub {
				mp.subsMap[topic] = next.next
				continue
			}

			prev := next
			for next = next.next; next != nil; next = next.next {
				if next == sub {
					prev.next = next.next
					break
				}
				prev = next
			}
		}
		mp.mu.Unlock()
	}
}

func (mp *MemPub) publish(topic Topic, data any) {
	msg := Message{Topic: topic, Data: data}

	mp.mu.Lock()
	if sub, ok := mp.subsMap[topic]; ok {
		for sub != nil {
			if !mp.timer.Stop() {
				<-mp.timer.C
			}
			mp.timer.Reset(50 * time.Millisecond)
			select {
			case sub.messageC <- msg:
			case <-mp.timer.C:
				select {
				case sub.overflowC <- struct{}{}:
				default:
				}
			}
			sub = sub.next
		}
	}
	mp.mu.Unlock()
}

type MemPubStat struct {
	Topic    Topic
	SubCount int
}

func (mp *MemPub) Stats() []MemPubStat {
	mp.mu.Lock()
	var stats []MemPubStat
	for topic, next := range mp.subsMap {
		var count int
		for ; next != nil; next = next.next {
			count++
		}

		stats = append(stats, MemPubStat{Topic: topic, SubCount: count})
	}
	mp.mu.Unlock()

	return stats
}
