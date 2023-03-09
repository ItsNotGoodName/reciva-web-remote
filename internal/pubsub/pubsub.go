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

func (p *Pub) Subscribe(topics []Topic) (<-chan Message, func()) {
	messageC := make(chan Message, 100)
	subs := []*Sub{}

	p.subsMapMu.Lock()
	for _, topic := range topics {
		sub := &Sub{messageC: messageC}
		subs = append(subs, sub)
		if next, ok := p.subsMap[topic]; ok {
			sub.next = next
		}

		p.subsMap[topic] = sub
	}
	p.subsMapMu.Unlock()

	return messageC, p.unsubscribeFunc(topics, subs)
}

func (p *Pub) unsubscribeFunc(topics []Topic, sub []*Sub) func() {
	return func() {
		p.subsMapMu.Lock()
		for i, sub := range sub {
			topic := topics[i]
			next := p.subsMap[topic]
			// Sub found in first place
			if next == sub {
				p.subsMap[topic] = next.next
				continue
			}

			// Will never be nil
			prev := next

			for next = next.next; next != nil; next = next.next {
				// Sub found in second place or more
				if next == sub {
					prev.next = next.next
					break
				}
			}
		}
		p.subsMapMu.Unlock()
	}
}

func (p *Pub) Publish(topic Topic, data interface{}) {
	msg := Message{Topic: topic, Data: data}

	p.subsMapMu.Lock()
	if sub, ok := p.subsMap[topic]; ok {
		for sub != nil {
			select {
			case sub.messageC <- msg:
			default:
			}
			sub = sub.next
		}
	}
	p.subsMapMu.Unlock()
}
