package tcp

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/ninjablocks/mqtt-proxy/rewrite"
)

type TcpProxyConn struct {

	// proxy connection
	pConn net.Conn

	// client connection
	cConn net.Conn

	id string

	rewriter *rewrite.MsgRewriter
}

func CreateTcpProxyConn(conn net.Conn, backend string) (*TcpProxyConn, error) {

	addr, err := net.ResolveTCPAddr("tcp", backend)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[serv] Error resolving upstream: %s", err))
	}
	log.Printf("Opening connection to %s", addr)
	tcpconn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("[serv] Error connecting to upstream: %s", err))
	}

	return &TcpProxyConn{cConn: conn, pConn: tcpconn, id: fmt.Sprintf("%s %s", conn.RemoteAddr(), conn.LocalAddr())}, nil

}

func (c *TcpProxyConn) Id() string {
	return c.id
}

func (c *TcpProxyConn) Close() {
	c.pConn.Close()
}
