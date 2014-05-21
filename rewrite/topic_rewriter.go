package rewrite

import (
	"log"
	"strings"

	"github.com/wolfeidau/mqtt"
)

var INGRESS = 1
var EGRESS = 2

type TopicRewriter interface {
	RewriteTopicName(topic string) string
	RenameTopicNames(topicNames []string) []string
	RewriteTopics(topics []mqtt.TopicQos) []mqtt.TopicQos
}

// rewriter which inserts a partition after the first token in the topic.
type TopicPartRewriter struct {
	Token     string
	Idx       int
	Direction int
}

func NewTopicPartRewriter(token string, idx int, direction int) *TopicPartRewriter {
	return &TopicPartRewriter{
		Token:     token,
		Idx:       idx,
		Direction: direction,
	}
}

func (tppw *TopicPartRewriter) RewriteTopicName(topic string) string {
	switch tppw.Direction {
	case INGRESS:
		return insertToken(topic, tppw.Token, tppw.Idx)
	case EGRESS:
		return removeToken(topic, tppw.Token, tppw.Idx)
	}
	return topic
}

func (tppw *TopicPartRewriter) RenameTopicNames(topicNames []string) []string {
	for i, _ := range topicNames {
		topicNames[i] = tppw.RewriteTopicName(topicNames[i])
	}
	return topicNames
}

func (tppw *TopicPartRewriter) RewriteTopics(topics []mqtt.TopicQos) []mqtt.TopicQos {
	for i, _ := range topics {
		topics[i].Topic = tppw.RewriteTopicName(topics[i].Topic)
	}
	return topics
}

func insertToken(topic string, token string, idx int) string {

	tokens := strings.Split(topic, "/")

	// append an empty token to the end
	tokens = append(tokens, "")

	// check the index will be valid
	if len(tokens) > idx {
		// copy the array over all members after the index supplied
		copy(tokens[idx+1:], tokens[idx:])
		tokens[idx] = token

		return strings.Join(tokens, "/")
	} else {
		log.Printf("[topic] Index less than tokens %d %s ", idx, topic)
		return topic
	}

}

func removeToken(topic string, token string, idx int) string {
	tokens := strings.Split(topic, "/")

	if tokens[idx] == token {
		return strings.Join(append(tokens[:idx], tokens[idx+1:]...), "/")
	} else {
		log.Printf("[topic] token not found %d %s %s", idx, topic, token)
		return topic
	}
}
