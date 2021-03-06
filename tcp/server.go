package tcp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/ninjablocks/mqtt-proxy/conf"
	"github.com/ninjablocks/mqtt-proxy/proxy"
	"github.com/ninjablocks/mqtt-proxy/store"
	"github.com/ninjablocks/mqtt-proxy/util"
	"github.com/wolfeidau/mqtt"
)

type TcpServer struct {
	proxy *proxy.MQTTProxy
	store store.Store
}

func CreateTcpServer(proxy *proxy.MQTTProxy) *TcpServer {

	store := store.NewMysqlStore(&proxy.Conf.MqttStoreMysql)

	return &TcpServer{
		proxy: proxy,
		store: store,
	}
}

func (t *TcpServer) StartServer(conf *conf.MqttConfiguration) {

	log.Printf("[tcp] listening on %s", conf.ListenAddress)

	listener, err := t.startListener(conf)

	if err != nil {
		log.Fatalln("error listening:", err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Client error: %s", err)
		} else {
			go t.clientHandler(conn)
		}
	}

}

func (t *TcpServer) startListener(conf *conf.MqttConfiguration) (net.Listener, error) {
	if conf.Cert != "" {
		cert, err := tls.LoadX509KeyPair(conf.Cert, conf.Key)

		if err != nil {
			log.Fatalf("server: loadkeys: %s", err)
		}
		log.Println("[serv] Starting tls listener")

		config := tls.Config{Certificates: []tls.Certificate{cert}}

		return tls.Listen("tcp", conf.ListenAddress, &config)
	} else {
		log.Println("[serv] Starting tcp listener")
		return net.Listen("tcp", conf.ListenAddress)
	}

}

func (t *TcpServer) clientHandler(conn net.Conn) {

	log.Printf("[serv] client connection opened - %s", conn.RemoteAddr())

	defer conn.Close()

	t.proxy.RegisterSession(conn)
	defer t.proxy.UnRegisterSession(conn)

	// create channels for the return messages from the client
	cmr := util.CreateMqttTcpMessageReader(conn, t.proxy.Conf.GetReadTimeout())

	go cmr.ReadMqttMessages()

	// This needs to be distributed across all servers
	backend := t.proxy.Conf.BackendServers[0]

	p, err := CreateTcpProxyConn(conn, backend)

	if err != nil {
		log.Printf("[serv] Error creating proxy connection - %s", err)
		sendServerUnavailable(conn)
		return
	}

	defer p.Close()

	t.proxy.Metrics.Connects.Mark(1)

	// do the authentication up front before going into normal operation
	if err = t.handleAuth(cmr, p); err != nil {
		log.Printf("[serv] Error authenticating connection - %s", err)
		// be very careful and clear on the error type as we are saying
		// for sure these credentials are not valid.
		if err == store.ErrUserNotFound {
			sendBadUsernameOrPassword(p.cConn)
		} else {
			sendServerUnavailable(conn)
		}
		return
	}

	// create channels for the return messages from the backend
	pmr := util.CreateMqttTcpMessageReader(p.pConn, t.proxy.Conf.GetReadTimeout())

	go pmr.ReadMqttMessages()

Loop:
	for {

		select {

		case msg := <-cmr.InMsgs:

			//util.DebugMQTT("client", conn, msg)
			msg = p.rewriter.RewriteIngress(msg)

			t.updateMsgCount(msg)

			// write to the proxy connection
			len, err := msg.Encode(p.pConn)

			if err != nil {
				log.Printf("[serv] proxy connection error - %s", err)
				break Loop
			}
			t.updateMsgBodySize(len)

		case err := <-cmr.InErrors:
			if err == io.EOF {
				log.Printf("[serv] client closed connection")
			} else {
				log.Printf("[serv] client connection read error - %s", err)
			}
			break Loop

		case msg := <-pmr.InMsgs:

			//util.DebugMQTT("proxy", conn, msg)
			msg = p.rewriter.RewriteEgress(msg)

			switch msg := msg.(type) {
			case *mqtt.ConnAck:
				log.Printf("[serv] got connack for %s", conn.RemoteAddr())
				log.Printf("[serv] connack %+v", msg)
			case *mqtt.Disconnect:
				log.Printf("[serv] got disconnect for %s", conn.RemoteAddr())
				log.Printf("[serv] disconnect %+v", msg)
			}
			t.proxy.Metrics.MsgReply.Mark(1)

			// write to the client connection
			len, err := msg.Encode(p.cConn)
			if err != nil {
				log.Printf("[serv] proxy connection error - %s", err)
				break Loop
			}
			t.updateMsgBodySize(len)

		case err := <-pmr.InErrors:
			if err == io.EOF {
				log.Printf("[serv] proxy connection closed by backend server")
			} else {
				log.Printf("[serv] proxy connection read error - %s", err)

			}
			break Loop
		}

	}

}

func (t *TcpServer) handleAuth(cmr *util.MqttTcpMessageReader, proxyConn *TcpProxyConn) error {

	select {
	case msg := <-cmr.InMsgs:

		//util.DebugMQTT("auth", proxyConn.cConn, msg)
		t.updateMsgCount(msg)

		switch cmsg := msg.(type) {
		case *mqtt.Connect:

			authUser, err := t.store.FindUser(cmsg.Username)

			if err != nil {
				log.Printf("[serv] authentication failed for %s - %s", authUser, err)
				return err
			}

			proxyConn.rewriter = t.proxy.MqttMsgRewriter(authUser)

			msg = proxyConn.rewriter.RewriteIngress(msg)

			len, err := msg.Encode(proxyConn.pConn)

			if err != nil {
				log.Printf("[serv] proxy connection error - %s", err)
				log.Println(spew.Sprintf("msg %v", msg))
				return err
			}

			t.updateMsgBodySize(len)

			return nil

		}
		// anything else is bad
		return errors.New(fmt.Sprintf("expected connect got - %v", reflect.TypeOf(msg)))

	case err := <-cmr.InErrors:
		log.Printf("connection error ocurred during authentication - %s", err)
		return err
	}

}

func (t *TcpServer) updateMsgCount(msg mqtt.Message) {
	t.proxy.Metrics.Msgs.Mark(1)
}

func (t *TcpServer) updateMsgBodySize(len int) {
	t.proxy.Metrics.MsgBodySize.Update(int64(len))
}

func sendBadUsernameOrPassword(conn net.Conn) {
	log.Printf("[serv] bad username / password %s %s", conn.LocalAddr(), conn.RemoteAddr())
	connAck := &mqtt.ConnAck{
		ReturnCode: mqtt.RetCodeBadUsernameOrPassword,
	}
	connAck.Encode(conn)
}

func sendServerUnavailable(conn net.Conn) {
	log.Printf("[serv] server unavailable %s %s", conn.LocalAddr(), conn.RemoteAddr())
	connAck := &mqtt.ConnAck{
		ReturnCode: mqtt.RetCodeServerUnavailable,
	}
	connAck.Encode(conn)
}
