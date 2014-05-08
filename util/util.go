package util

import (
	"log"
	"net"
	"reflect"
	"time"

	"github.com/ninjablocks/mqtt-proxy/proxy"
	"github.com/wolfeidau/mqtt"
)

type MqttTcpMessageReader struct {
	Tcpconn  net.Conn
	InMsgs   chan mqtt.Message
	InErrors chan error
}

func CreateMqttTcpMessageReader(tcpconn net.Conn) *MqttTcpMessageReader {
	return &MqttTcpMessageReader{
		Tcpconn:  tcpconn,
		InMsgs:   make(chan mqtt.Message, 1),
		InErrors: make(chan error, 1),
	}
}

func (m *MqttTcpMessageReader) ReadMqttMessages() {

	defer log.Println("[serv] Reader done - ", m.Tcpconn.RemoteAddr())

	m.Tcpconn.SetReadDeadline(time.Now().Add(proxy.DefaultReadTimeout))

	for {
		msg, err := mqtt.DecodeOneMessage(m.Tcpconn, nil)
		if err != nil {
			m.InErrors <- err
			break
		} else {
			m.InMsgs <- msg
		}
	}

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
