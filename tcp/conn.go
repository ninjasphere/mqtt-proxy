package tcp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/huin/mqtt"
	"github.com/ninjablocks/mqtt-proxy/rewrite"
	"github.com/ninjablocks/mqtt-proxy/util"
)

type TcpProxyConn struct {

	// proxy connection
	pConn net.Conn

	// client connection
	cConn net.Conn

	id       string
	wait     sync.WaitGroup
	rewriter *rewrite.MsgRewriter

	// mutex for the closed flag
	mutex  sync.Mutex
	closed bool
}

func CreateTcpProxyConn(conn net.Conn, backend string, rewriter *rewrite.MsgRewriter) (*TcpProxyConn, error) {

	addr, err := net.ResolveTCPAddr("tcp", backend)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error resolving upstream: %s", err))
	}

	tcpconn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error connecting to upstream: %s", err))
	}

	return &TcpProxyConn{cConn: conn, pConn: tcpconn, id: fmt.Sprintf("%s %s", conn.RemoteAddr(), conn.LocalAddr()), rewriter: rewriter, closed: false}, nil

}

func (c *TcpProxyConn) Id() string {
	return c.id
}

func (c *TcpProxyConn) ReadEgressConn() {

	defer c.wait.Done()

	// cleanup once we are done
	defer c.Close()

Loop:
	for {
		msg, err := mqtt.DecodeOneMessage(c.pConn, nil)
		if err != nil {
			if err != io.EOF {
				log.Println("[mqtt] Bad Egress Message:", err)
				break Loop
			}
			break

		}
		util.DebugMQTTMsg("tcp out", c, msg)

		msg = c.rewriter.RewriteEgress(msg)

		if err := msg.Encode(c.cConn); err != nil {
			log.Printf("[mqtt] Send failed: %s", err)
			break
		}
	}
}

func (c *TcpProxyConn) ReadIngressConn() {

	defer c.wait.Done()

	// cleanup once we are done
	defer c.Close()

	for {
		msg, err := mqtt.DecodeOneMessage(c.cConn, nil)
		if err != nil {
			if err != io.EOF {
				log.Println("[mqtt] Bad Ingress Message:", err)
			}
			break
		}
		util.DebugMQTTMsg("tcp in", c, msg)

		msg = c.rewriter.RewriteIngress(msg)

		if err := msg.Encode(c.pConn); err != nil {
			log.Printf("[tcp] Send failed: %s", err)
			break
		}

	}
}

func (c *TcpProxyConn) Close() {
	log.Printf("[tcp] (%s) Close", c.id)

	c.mutex.Lock()

	if !c.closed {
		c.cConn.Close()
		c.pConn.Close()
		c.closed = true
	}

	c.mutex.Unlock()

}
