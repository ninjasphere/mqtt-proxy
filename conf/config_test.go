package conf

import (
	. "launchpad.net/gocheck"

	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type LoadConfigurationSuite struct{}

var _ = Suite(&LoadConfigurationSuite{})

func (self *LoadConfigurationSuite) TestConfig(c *C) {
	config := LoadConfiguration("config_test.toml")
	c.Assert(config.BackendServers, DeepEquals, []string{"hosta:8090", "hostb:8090"})
	// check mqtt defaults
	c.Assert(config.Mqtt.ListenAddress, Equals, ":1883")
	c.Assert(config.User, Equals, "guest")
	c.Assert(config.Pass, Equals, "guest")
}
