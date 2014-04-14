package rewrite

import (
	"testing"

	"github.com/huin/mqtt"
	. "launchpad.net/gocheck"
)

func Test3(t *testing.T) { TestingT(t) }

type TopicRewriterSuite struct {
	topicRewriter    TopicRewriter
	topicRewriterDup TopicRewriter

	egressPartRewriter  TopicRewriter
	ingressPartRewriter TopicRewriter
}

var _ = Suite(&TopicRewriterSuite{})

func (s *TopicRewriterSuite) SetUpTest(c *C) {
	s.topicRewriter = &TopicPrefixRewriter{
		Prefix:      "$block",
		Replacement: "$cloud/123",
	}

	s.topicRewriterDup = &TopicPrefixRewriter{
		Prefix:      "block",
		Replacement: "cloud/123",
	}

	s.egressPartRewriter = &TopicPartRewriter{
		Token:     "123",
		Idx:       1,
		Direction: EGRESS,
	}

	s.ingressPartRewriter = &TopicPartRewriter{
		Token:     "123",
		Idx:       1,
		Direction: INGRESS,
	}
}

func (s *TopicRewriterSuite) TestTopicName(c *C) {
	topicName := s.topicRewriter.RewriteTopicName("$block/test/123")
	c.Assert(topicName, Equals, "$cloud/123/test/123")

	// duplicate tokens should only do prefix
	topicNameDup := s.topicRewriterDup.RewriteTopicName("block/something/block")
	c.Assert(topicNameDup, Equals, "cloud/123/something/block")

	// test it ignores and returns the original
	topicNameIngore := s.topicRewriterDup.RewriteTopicName("cloud/something/block")
	c.Assert(topicNameIngore, Equals, "cloud/something/block")

}

func (s *TopicRewriterSuite) TestTopicNames(c *C) {

	testTopicNames := []string{"$block/test/123", "$block/test/456", "$block/test/789"}
	expectedTopicNames := []string{"$cloud/123/test/123", "$cloud/123/test/456", "$cloud/123/test/789"}

	topicNames := s.topicRewriter.RenameTopicNames(testTopicNames)
	c.Assert(topicNames, DeepEquals, expectedTopicNames)
}

func (s *TopicRewriterSuite) TestTopics(c *C) {
	testTopics := []mqtt.TopicQos{
		mqtt.TopicQos{
			Topic: "$block/test/123",
			Qos:   1,
		},
	}
	expectedTopics := []mqtt.TopicQos{
		mqtt.TopicQos{
			Topic: "$cloud/123/test/123",
			Qos:   1,
		},
	}

	topics := s.topicRewriter.RewriteTopics(testTopics)
	c.Assert(topics, DeepEquals, expectedTopics)

}

func (s *TopicRewriterSuite) TestPartTopicName(c *C) {

	topicName := s.ingressPartRewriter.RewriteTopicName("$cloud/test/123")
	c.Assert(topicName, Equals, "$cloud/123/test/123")

	topicName = s.ingressPartRewriter.RewriteTopicName("cloud")
	c.Assert(topicName, Equals, "cloud/123")
}
