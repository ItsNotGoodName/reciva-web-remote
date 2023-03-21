package ws

import (
	"log"

	"github.com/ItsNotGoodName/reciva-web-remote/internal/pubsub"
)

type CommandSubscribe struct {
	Topics []pubsub.Topic `json:"topics" validate:"required"`
}

type CommandState struct {
	UUID    string `json:"uuid" validate:"required"`
	Partial bool   `json:"partial" validate:"required"`
}

type Command struct {
	Subscribe *CommandSubscribe `json:"subscribe"`
	State     *CommandState     `json:"state"`
}

type Event struct {
	Topic pubsub.Topic `json:"topic" validate:"required"`
	Data  any          `json:"data" validate:"required"`
}

func uniqueTopics(topics []pubsub.Topic) []pubsub.Topic {
	pubsubTopics := []pubsub.Topic{}
	for _, top := range topics {
		if top.In(pubsubTopics) {
			log.Println("ws.uniqueTopics: received duplicate topic:", top)
		} else {
			pubsubTopics = append(pubsubTopics, top)
		}
	}

	return pubsubTopics
}

func filterTopics(topics []pubsub.Topic) []pubsub.Topic {
	if length := len(topics); length > 16 {
		log.Println("ws.filterTopics: received invalid topics length:", length)
		return []pubsub.Topic{}
	}

	pubsubTopics := []pubsub.Topic{}
	for _, topic := range topics {

		if topic == pubsub.TopicStale {
			pubsubTopics = append(pubsubTopics, pubsub.TopicStaleRadios, pubsub.TopicStaleStateHook)
		} else {
			switch pubsub.Topic(topic) {
			case pubsub.TopicDiscover:
				break
			case pubsub.TopicState:
				break
			case pubsub.TopicError:
				break
			default:
				log.Println("ws.filterTopics: invalid topic:", topic)
				continue
			}

			pubsubTopics = append(pubsubTopics, topic)
		}
	}

	return pubsubTopics
}
