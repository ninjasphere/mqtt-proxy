package tcp

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/huin/mqtt"
	"github.com/ninjablocks/mqtt-proxy/conf"
	"github.com/ninjablocks/mqtt-proxy/proxy"
	"github.com/ninjablocks/mqtt-proxy/store"
)

type TcpServer struct {
	proxy *proxy.MQTTProxy
	store store.Store
}

func CreateTcpServer(proxy *proxy.MQTTProxy, store store.Store) *TcpServer {
	return &TcpServer{
		proxy: proxy,
		store: store,
	}
}

func (t *TcpServer) StartServer(conf *conf.MqttConfiguration) {

	log.Printf("[tcp] listening on %s", conf.ListenAddress)

	cert, err := tls.LoadX509KeyPair(conf.Cert, conf.Key)

	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", conf.ListenAddress, &config)

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

func (t *TcpServer) clientHandler(conn net.Conn) {

	var authUser *store.User

	log.Printf("[serv] read connect msg")
	msg, err := mqtt.DecodeOneMessage(conn, nil)

	if err != nil {
		conn.Close()
		return
	}
	switch msg := msg.(type) {
	case *mqtt.Connect:

		authUser, err = t.store.FindUser(msg.Username)

		if err != nil {
			log.Printf("[serv] Error authenticating connection - %s", err)
			sendBadUsernameOrPassword(conn)
			return
		}

	default:
		// anything else is bad
		conn.Close()
		return
	}

	// This needs to be distributed across all servers
	backend := t.proxy.Conf.BackendServers[0]

	proxyConn, err := CreateTcpProxyConn(conn, backend, t.proxy.MqttMsgRewriter(authUser))

	if err != nil {
		log.Println("[serv] Create ProxyConn:", err)
		sendServerUnavailable(conn)
		return
	}

	proxyConn.wait.Add(2)

	log.Printf("[serv] start readers")
	go proxyConn.ReadEgressConn()
	go proxyConn.ReadIngressConn()

	proxyConn.wait.Wait()

	// TODO Output stats from the proxy connection
	log.Println("[serv] Finished")

}

func sendBadUsernameOrPassword(conn net.Conn) {
	log.Printf("[serv] bad username / password %s %s", conn.LocalAddr(), conn.RemoteAddr())
	connAck := &mqtt.ConnAck{
		ReturnCode: mqtt.RetCodeBadUsernameOrPassword,
	}
	connAck.Encode(conn)
	conn.Close()
}

func sendServerUnavailable(conn net.Conn) {
	log.Printf("[serv] server unavailable %s %s", conn.LocalAddr(), conn.RemoteAddr())
	connAck := &mqtt.ConnAck{
		ReturnCode: mqtt.RetCodeServerUnavailable,
	}
	connAck.Encode(conn)
	conn.Close()
}
