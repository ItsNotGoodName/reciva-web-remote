package pubsub

import (
	"sync"
)

var DefaultPub *Pub = NewPub()

type Topic string

func (t Topic) In(topics []Topic) bool {
	for _, topic := range topics {
		if topic == t {
			return true
		}
	}
	return false
}

type Message struct {
	Topic Topic
	Data  interface{}
}

type subscriber struct {
	messageC chan Message
	next     *subscriber
}

type Pub struct {
	subsMapMu sync.Mutex
	subsMap   map[Topic]*subscriber
}

func NewPub() *Pub {
	return &Pub{
		subsMapMu: sync.Mutex{},
		subsMap:   make(map[Topic]*subscriber),
	}
}

func (p *Pub) Subscribe(topics []Topic) (chan Message, func()) {
	return p.SubscribeWithBuffer(topics, 100)
}

func (p *Pub) SubscribeWithBuffer(topics []Topic, buffer int) (chan Message, func()) {
	messageC := make(chan Message, buffer)

	p.subsMapMu.Lock()
	subs := p.subscribe(topics, messageC)
	p.subsMapMu.Unlock()

	return messageC, p.unsubscribeFunc(topics, subs)
}

func (p *Pub) Resubscribe(topics []Topic, messageC chan Message, unsub func()) func() {
	unsub()
Loop:
	for {
		select {
		case <-messageC:
		default:
			break Loop
		}
	}

	p.subsMapMu.Lock()
	subs := p.subscribe(topics, messageC)
	p.subsMapMu.Unlock()

	return p.unsubscribeFunc(topics, subs)
}

func (p *Pub) subscribe(topics []Topic, messageC chan Message) []*subscriber {
	subs := []*subscriber{}
	for _, topic := range topics {
		sub := &subscriber{messageC: messageC}
		subs = append(subs, sub)
		if next, ok := p.subsMap[topic]; ok {
			sub.next = next
		}

		p.subsMap[topic] = sub
	}
	return subs
}

func (p *Pub) unsubscribeFunc(topics []Topic, sub []*subscriber) func() {
	return func() {
		p.subsMapMu.Lock()
		for i, sub := range sub {
			topic := topics[i]
			// There should only be 1 or 0 sub in next because the unsubscribe function might be called twice
			next := p.subsMap[topic]
			if next == nil {
				continue
			}
			if next == sub {
				p.subsMap[topic] = next.next
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
		p.subsMapMu.Unlock()
	}
}

func (p *Pub) Publish(topic Topic, data any) {
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

// func stats(subs map[Topic]*Sub) {
// 	var topics []Topic
// 	for topic := range subs {
// 		topics = append(topics, topic)
// 	}
// 	fmt.Printf("	TOPICS: %v\n", topics)

// 	for _, topic := range topics {
// 		next := subs[topic]
// 		if next == nil {
// 			fmt.Printf("	TOPIC_SUBS: %s: 0 subs\n", topic)
// 			continue
// 		}

// 		count := 1
// 		for next = next.next; next != nil; next = next.next {
// 			count++
// 		}
// 		fmt.Printf("	TOPIC_SUBS: %s: %d subs\n", topic, count)
// 	}
// }
