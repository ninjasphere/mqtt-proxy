package rewrite

import (
	"fmt"
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
	Direction int
}

func NewTopicPartRewriter(token string, direction int) *TopicPartRewriter {
	return &TopicPartRewriter{
		Token:     token,
		Direction: direction,
	}
}

func (tppw *TopicPartRewriter) RewriteTopicName(topic string) string {
	switch tppw.Direction {
	case INGRESS:
		return insertToken(topic, tppw.Token)
	case EGRESS:
		return removeToken(topic, tppw.Token)
	}
	return topic
}

func (tppw *TopicPartRewriter) RenameTopicNames(topicNames []string) []string {
	for i := range topicNames {
		topicNames[i] = tppw.RewriteTopicName(topicNames[i])
	}
	return topicNames
}

func (tppw *TopicPartRewriter) RewriteTopics(topics []mqtt.TopicQos) []mqtt.TopicQos {
	for i := range topics {
		topics[i].Topic = tppw.RewriteTopicName(topics[i].Topic)
	}
	return topics
}

func insertToken(topic string, token string) string {
	return fmt.Sprintf("%s/%s", token, topic)
}

func removeToken(topic string, token string) string {
	tokens := strings.Split(topic, "/")

	if tokens[0] == token {
		return strings.Join(tokens[1:], "/")
	}

	log.Printf("[topic] token not found %d %s", topic, token)
	return topic
}
