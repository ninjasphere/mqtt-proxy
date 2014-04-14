package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/ninjablocks/mqtt-proxy/conf"
	"github.com/ninjablocks/mqtt-proxy/proxy"
	"github.com/ninjablocks/mqtt-proxy/store"
	"github.com/ninjablocks/mqtt-proxy/tcp"
	"github.com/ninjablocks/mqtt-proxy/ws"
)

var configFile = flag.String("config", "config.toml", "configuration file")
var debug = flag.Bool("debug", false, "enable debugging")

func main() {
	flag.Parse()

	conf := conf.LoadConfiguration(*configFile)

	if *debug {
		log.Printf("[main] conf %+v", conf)
	}

	proxy := proxy.CreateMQTTProxy(conf)

	store := store.NewMysqlStore(&conf.Mysql)

	handlers := ws.CreateHttpHanders(proxy, store)
	tcpServer := tcp.CreateTcpServer(proxy, store)

	go handlers.StartServer(&conf.Http)

	go tcpServer.StartServer(&conf.Mqtt)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	log.Printf("Got signal %s, exiting now", <-c)
}
