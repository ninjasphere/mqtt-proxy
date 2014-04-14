package proxy

import (
	"github.com/ninjablocks/mqtt-proxy/conf"
	"github.com/ninjablocks/mqtt-proxy/rewrite"
	"github.com/ninjablocks/mqtt-proxy/store"
)

type ProxyConn interface {
	Id() string
	Close()
}

type MQTTProxy struct {
	Conf        *conf.Configuration
	connections map[string]ProxyConn
}

func CreateMQTTProxy(conf *conf.Configuration) *MQTTProxy {
	return &MQTTProxy{
		Conf: conf,
	}
}

func (p *MQTTProxy) mqttCredentialsRewriter(user *store.User) rewrite.CredentialsRewriter {
	return rewrite.NewCredentialsReplaceRewriter(p.Conf.User, p.Conf.Pass, user.UserId)
}

func (p *MQTTProxy) mqttTopicRewriter(mqttId string, direction int) rewrite.TopicRewriter {
	return rewrite.NewTopicPartRewriter(mqttId, 1, direction)
}

func (p *MQTTProxy) MqttMsgRewriter(user *store.User) *rewrite.MsgRewriter {
	return rewrite.CreatMsgRewriter(p.mqttCredentialsRewriter(user), p.mqttTopicRewriter(user.MqttId, rewrite.INGRESS), p.mqttTopicRewriter(user.MqttId, rewrite.EGRESS))
}
