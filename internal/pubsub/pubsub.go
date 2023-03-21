package pubsub

var DefaultPub *MemPub = NewMemPub()

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

type Pub interface {
	publish(topic Topic, data any)
}
