package ws

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/huin/mqtt"
	"github.com/ninjablocks/mqtt-proxy/rewrite"
	"github.com/ninjablocks/mqtt-proxy/util"
)

type WsProxyConn struct {
	tcpconn net.Conn
	wsconn  *websocket.Conn

	id       string
	wait     sync.WaitGroup
	rewriter *rewrite.MsgRewriter

	// mutex for the closed flag
	mutex  sync.Mutex
	closed bool
}

func CreateWsProxyConn(conn *websocket.Conn, backend string, rewriter *rewrite.MsgRewriter) (*WsProxyConn, error) {

	addr, err := net.ResolveTCPAddr("tcp", backend)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error resolving upstream: %s", err))
	}

	tcpconn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error connecting to upstream: %s", err))
	}

	return &WsProxyConn{wsconn: conn, tcpconn: tcpconn, id: fmt.Sprintf("%s %s", conn.RemoteAddr(), conn.LocalAddr()), rewriter: rewriter, closed: false}, nil

}

func (c *WsProxyConn) Id() string {
	return c.id
}

func (c *WsProxyConn) ReadEgressConn() {

	defer c.wait.Done()

	// cleanup once we are done
	defer c.Close()

	for {
		msg, err := mqtt.DecodeOneMessage(c.tcpconn, nil)
		if err != nil {
			if err != io.EOF {
				log.Println("[mqtt] Bad Message:", err)
			}
			break

		}
		util.DebugMQTTMsg("mqtt", c, msg)

		msg = c.rewriter.RewriteEgress(msg)

		encodedBuf := new(bytes.Buffer)

		if err := msg.Encode(encodedBuf); err != nil {
			log.Printf("[mqtt] Send failed: %s", err)
			break
		} else {
			c.wsconn.WriteMessage(websocket.BinaryMessage, encodedBuf.Bytes())
		}
	}

}

func (c *WsProxyConn) ReadIngressConn() {

	defer c.wait.Done()

	// cleanup once we are done
	defer c.Close()

	for {

		// need to rate limit these incoming messages
		mt, b, err := c.wsconn.ReadMessage()

		if err != nil {
			if err != io.EOF {
				log.Println("NextReader:", err)
			}
			break
		}

		if mt == websocket.BinaryMessage {

			msg, err := mqtt.DecodeOneMessage(bytes.NewReader(b), nil)

			util.DebugMQTTMsg("ws", c, msg)

			msg = c.rewriter.RewriteIngress(msg)

			if err != nil {
				log.Println("[ws] Unable to decode msg:", err)
				break
			}
			if err := msg.Encode(c.tcpconn); err != nil {
				log.Println("[ws] Send to upstream failed:", err)
				break
			}
			if util.IsMqttDisconnect(msg) {
				break
			}

		}
	}

}

func (c *WsProxyConn) Close() {

	c.mutex.Lock()

	if !c.closed {
		log.Printf("[ws] (%s) Close", c.id)
		c.tcpconn.Close()
		c.closed = true
	}

	c.mutex.Unlock()

}
