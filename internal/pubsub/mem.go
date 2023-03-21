package pubsub

import "sync"

type MemPub struct {
	subsMapMu sync.Mutex
	subsMap   map[Topic]*memSub
}

type memSub struct {
	messageC chan Message
	next     *memSub
}

func NewMemPub() *MemPub {
	return &MemPub{
		subsMapMu: sync.Mutex{},
		subsMap:   make(map[Topic]*memSub),
	}
}

func (mp *MemPub) Subscribe(topics []Topic) (chan Message, func()) {
	return mp.SubscribeWithBuffer(topics, 100)
}

func (mp *MemPub) SubscribeWithBuffer(topics []Topic, buffer int) (chan Message, func()) {
	messageC := make(chan Message, buffer)

	mp.subsMapMu.Lock()
	subs := mp.subscribe(topics, messageC)
	mp.subsMapMu.Unlock()

	return messageC, mp.unsubscribeFunc(topics, subs)
}

func (mp *MemPub) Resubscribe(topics []Topic, messageC chan Message, unsub func()) func() {
	unsub()
Loop:
	for {
		select {
		case <-messageC:
		default:
			break Loop
		}
	}

	mp.subsMapMu.Lock()
	subs := mp.subscribe(topics, messageC)
	mp.subsMapMu.Unlock()

	return mp.unsubscribeFunc(topics, subs)
}

func (mp *MemPub) subscribe(topics []Topic, messageC chan Message) []*memSub {
	subs := []*memSub{}
	for _, topic := range topics {
		sub := &memSub{messageC: messageC}
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
		mp.subsMapMu.Lock()
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
		mp.subsMapMu.Unlock()
	}
}

func (mp *MemPub) publish(topic Topic, data any) {
	msg := Message{Topic: topic, Data: data}

	mp.subsMapMu.Lock()
	if sub, ok := mp.subsMap[topic]; ok {
		for sub != nil {
			select {
			case sub.messageC <- msg:
			default:
			}
			sub = sub.next
		}
	}
	mp.subsMapMu.Unlock()
}

type MemPubStat struct {
	Topic    Topic
	SubCount int
}

func (mp *MemPub) Stats() []MemPubStat {
	mp.subsMapMu.Lock()
	var stats []MemPubStat
	for topic, next := range mp.subsMap {
		var count int
		for ; next != nil; next = next.next {
			count++
		}

		stats = append(stats, MemPubStat{Topic: topic, SubCount: count})
	}
	mp.subsMapMu.Unlock()

	return stats
}
