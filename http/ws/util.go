package ws

import (
	"fmt"
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

func validateTopic(topic pubsub.Topic) (pubsub.Topic, error) {
	switch pubsub.Topic(topic) {
	case pubsub.DiscoverTopic:
		return pubsub.DiscoverTopic, nil
	case pubsub.StateTopic:
		return pubsub.StateTopic, nil
	// case pubsub.ErrorTopic:
	// 	return pubsub.ErrorTopic, nil
	default:
		return "", fmt.Errorf("invalid topic")
	}
}

func filterValidTopics(topics []pubsub.Topic) []pubsub.Topic {
	pubsubTopics := []pubsub.Topic{}
	for _, topic := range topics {
		pubsubTopic, err := validateTopic(topic)
		if err != nil {
			log.Println("http.wsParseTopics: invalid topic:", topic)
			continue
		}

		pubsubTopics = append(pubsubTopics, pubsubTopic)
	}

	return pubsubTopics
}
