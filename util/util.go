package util

import (
	"log"
	"net"
	"reflect"

	"github.com/huin/mqtt"
	"github.com/ninjablocks/mqtt-proxy/proxy"
)

type MqttTcpMessageReader struct {
	Tcpconn  net.Conn
	InMsgs   chan mqtt.Message
	InErrors chan error
}

func CreateMqttTcpMessageReader(tcpconn net.Conn) *MqttTcpMessageReader {
	return &MqttTcpMessageReader{
		Tcpconn:  tcpconn,
		InMsgs:   make(chan mqtt.Message),
		InErrors: make(chan error),
	}
}

func (m *MqttTcpMessageReader) ReadMqttMessages() {

	for {
		msg, err := mqtt.DecodeOneMessage(m.Tcpconn, nil)
		if err != nil {
			m.InErrors <- err
			break
		}
		m.InMsgs <- msg
	}
}

func (m *MqttTcpMessageReader) SendMqttMessage(msg mqtt.Message) {
}
func IsMqttDisconnect(msg mqtt.Message) bool {
	return reflect.TypeOf(msg) == reflect.TypeOf(mqtt.MsgDisconnect)
}

func DebugMQTTMsg(tag string, c proxy.ProxyConn, msg mqtt.Message) {
	log.Printf("[%s] (%s) %s", tag, c.Id(), reflect.TypeOf(msg))
}

func DebugMQTT(tag string, c net.Conn, msg mqtt.Message) {
	log.Printf("[%s] (%s) %s", tag, c.RemoteAddr(), reflect.TypeOf(msg))
}
