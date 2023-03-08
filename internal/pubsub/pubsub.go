package pubsub

import (
	"sync"
)

var DefaultPub *Pub = NewPub()

type Topic string

type Message struct {
	Topic Topic
	Data  interface{}
}

type Sub struct {
	messageC chan Message
	next     *Sub
}

type Pub struct {
	subsMapMu sync.Mutex
	subsMap   map[Topic]*Sub
}

func NewPub() *Pub {
	return &Pub{
		subsMapMu: sync.Mutex{},
		subsMap:   make(map[Topic]*Sub),
	}
}

func (sp *Pub) Subscribe(topics []Topic) (<-chan Message, func()) {
	sub := &Sub{messageC: make(chan Message, 100)}

	sp.subsMapMu.Lock()
	for _, topic := range topics {
		if next, ok := sp.subsMap[topic]; ok {
			sub.next = next
		}

		sp.subsMap[topic] = sub
	}
	sp.subsMapMu.Unlock()

	return sub.messageC, sp.unsubscribeFunc(topics, sub)
}

func (sp *Pub) unsubscribeFunc(topics []Topic, sub *Sub) func() {
	return func() {
		sp.subsMapMu.Lock()
		for _, topic := range topics {
			next := sp.subsMap[topic]
			if next == sub {
				sp.subsMap[topic] = nil
				break
			}

			prev := next
			for next = next.next; next != nil; next = next.next {
				if next == sub {
					prev.next = next.next
					break
				}
				prev = next
			}

			sub.next = next
		}
		sp.subsMapMu.Unlock()
	}
}

func (sp *Pub) Publish(topic Topic, data interface{}) {
	msg := Message{Topic: topic, Data: data}

	sp.subsMapMu.Lock()
	if sub, ok := sp.subsMap[topic]; ok {
		for sub != nil {
			select {
			case sub.messageC <- msg:
			default:
			}
			sub = sub.next
		}
	}
	sp.subsMapMu.Unlock()
}
