package rewrite

import (
	"testing"

	. "launchpad.net/gocheck"
)

func Test3(t *testing.T) { TestingT(t) }

type TopicRewriterSuite struct {
	egressPartRewriter  TopicRewriter
	ingressPartRewriter TopicRewriter
}

var _ = Suite(&TopicRewriterSuite{})

func (s *TopicRewriterSuite) SetUpTest(c *C) {

	s.egressPartRewriter = &TopicPartRewriter{
		Token:     "123",
		Direction: EGRESS,
	}

	s.ingressPartRewriter = &TopicPartRewriter{
		Token:     "123",
		Direction: INGRESS,
	}
}

func (s *TopicRewriterSuite) TestPartTopicName(c *C) {

	topicName := s.ingressPartRewriter.RewriteTopicName("$cloud/test/123")
	c.Assert(topicName, Equals, "123/$cloud/test/123")

	topicName = s.ingressPartRewriter.RewriteTopicName("cloud")
	c.Assert(topicName, Equals, "123/cloud")
}
