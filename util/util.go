package util

import (
	"log"
	"reflect"

	"github.com/huin/mqtt"
	"github.com/ninjablocks/mqtt-proxy/proxy"
)

func IsMqttDisconnect(msg mqtt.Message) bool {
	return reflect.TypeOf(msg) == reflect.TypeOf(mqtt.MsgDisconnect)
}

func DebugMQTTMsg(tag string, c proxy.ProxyConn, msg mqtt.Message) {
	log.Printf("[%s] (%s) %s", tag, c.Id(), reflect.TypeOf(msg))
}
