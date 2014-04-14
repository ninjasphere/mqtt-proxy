package rewrite

import (
	"testing"

	"github.com/huin/mqtt"

	. "launchpad.net/gocheck"
)

func Test4(t *testing.T) { TestingT(t) }

type CredentialsRewriterSuite struct {
	credentialsRewriter CredentialsRewriter
}

var _ = Suite(&CredentialsRewriterSuite{})

func (s *CredentialsRewriterSuite) SetUpTest(c *C) {
	s.credentialsRewriter = NewCredentialsReplaceRewriter("user", "pass", 1)
}

func (s *CredentialsRewriterSuite) TestCredsRewrite(c *C) {

	// connection request message
	connect := createConnectMessage("bob", "11223344", true, true, "abc")
	expectedConnect := createConnectMessage("user", "pass", true, true, "1-abc")

	modConnect := s.credentialsRewriter.RewriteCredentials(connect)
	c.Assert(modConnect, DeepEquals, expectedConnect)

}

func (s *CredentialsRewriterSuite) TestCredsRewriteJustUser(c *C) {

	// connection request message
	connect := createConnectMessage("bob", "", true, false, "abc")
	expectedConnect := createConnectMessage("user", "pass", true, true, "1-abc")

	modConnect := s.credentialsRewriter.RewriteCredentials(connect)
	c.Assert(modConnect, DeepEquals, expectedConnect)

}

func createConnectMessage(user string, pass string, userFlag bool, passFlag bool, clientId string) *mqtt.Connect {
	return &mqtt.Connect{
		ProtocolName:    "MQIsdp",
		ProtocolVersion: 3,
		UsernameFlag:    userFlag,
		PasswordFlag:    passFlag,
		WillRetain:      false,
		WillQos:         1,
		WillFlag:        true,
		CleanSession:    true,
		KeepAliveTimer:  10,
		ClientId:        clientId,
		WillTopic:       "topic",
		WillMessage:     "message",
		Username:        user,
		Password:        pass,
	}
}
