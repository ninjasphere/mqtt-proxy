package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/ninjablocks/mqtt-proxy/conf"
	"github.com/ninjablocks/mqtt-proxy/content"
	"github.com/ninjablocks/mqtt-proxy/proxy"
	"github.com/ninjablocks/mqtt-proxy/store"
)

type HttpHanders struct {
	proxy *proxy.MQTTProxy
	store store.Store
}

func CreateHttpHanders(proxy *proxy.MQTTProxy, store store.Store) *HttpHanders {
	return &HttpHanders{
		proxy: proxy,
		store: store,
	}
}

func (h *HttpHanders) StartServer(conf *conf.HttpConfiguration) {

	log.Printf("[http] listening on %s", conf.ListenAddress)

	r := mux.NewRouter()
	fs := http.FileServer(content.ContentBox().HTTPBox())

	// setup the handlers in the router
	r.HandleFunc("/mqtt/{key}", h.mqttHandler)
	r.PathPrefix("/").Handler(fs)

	// configure this router in http
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(conf.ListenAddress, nil))
}

func (h *HttpHanders) mqttHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	user, err := h.store.FindUser(vars["key"])

	if err != nil {
		log.Println("Auth:", err)
		http.Error(w, "Unauthorized", 401)
		return
	}

	log.Printf("[http] User: %v", user)

	conn, err := websocket.Upgrade(w, r, nil, 4096, 4096)
	if err != nil {
		log.Println("Upgrade:", err)
		http.Error(w, "Bad request", 400)
		return
	}

	defer conn.Close()

	// This needs to be distributed across all servers
	backend := h.proxy.Conf.BackendServers[0]

	c, err := CreateWsProxyConn(conn, backend, h.proxy.MqttMsgRewriter(user))

	if err != nil {
		log.Println("Create ProxyConn:", err)
		http.Error(w, "Connect to upstream server failed", 500)
		return
	}

	c.wait.Add(2)

	go c.ReadEgressConn()
	go c.ReadIngressConn()

	c.wait.Wait()

	// TODO Output stats from the proxy connection
	log.Println("[handler] Finished")

}
