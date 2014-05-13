package proxy

import (
	"fmt"
	"net"

	"github.com/ninjablocks/mqtt-proxy/conf"
	"github.com/ninjablocks/mqtt-proxy/metrics"
	"github.com/ninjablocks/mqtt-proxy/rewrite"
	"github.com/ninjablocks/mqtt-proxy/store"
)

type ProxyConn interface {
	Id() string
	Close()
}

type MQTTProxy struct {
	Conf        *conf.Configuration
	connections map[string]net.Conn
	Metrics     metrics.ProxyMetrics
}

func CreateMQTTProxy(conf *conf.Configuration) *MQTTProxy {
	p := &MQTTProxy{
		Conf:        conf,
		Metrics:     metrics.NewProxyMetrics(conf.Environment, conf.Region),
		connections: make(map[string]net.Conn),
	}

	return p
}

func (p *MQTTProxy) RegisterSession(conn net.Conn) {
	id := fmt.Sprintf("%s %s", conn.RemoteAddr(), conn.LocalAddr())
	p.connections[id] = conn
	p.Metrics.Connections.Update(int64(len(p.connections)))
}

func (p *MQTTProxy) UnRegisterSession(conn net.Conn) {
	id := fmt.Sprintf("%s %s", conn.RemoteAddr(), conn.LocalAddr())
	delete(p.connections, id)
	p.Metrics.Connections.Update(int64(len(p.connections)))
}

func (p *MQTTProxy) mqttCredentialsRewriter(user *store.User) rewrite.CredentialsRewriter {
	return rewrite.NewCredentialsReplaceRewriter(p.Conf.User, p.Conf.Pass, user.UserId)
}

func (p *MQTTProxy) mqttTopicRewriter(mqttId string, direction int) rewrite.TopicRewriter {
	return rewrite.NewTopicPartRewriter(mqttId, 1, direction)
}

func (p *MQTTProxy) MqttMsgRewriter(user *store.User) *rewrite.MsgRewriter {
	return rewrite.CreatMsgRewriter(p.mqttCredentialsRewriter(user), p.mqttTopicRewriter(user.UserId, rewrite.INGRESS), p.mqttTopicRewriter(user.UserId, rewrite.EGRESS))
}
