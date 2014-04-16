package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ninjablocks/mqtt-proxy/conf"
	"github.com/ninjablocks/mqtt-proxy/proxy"
	"github.com/ninjablocks/mqtt-proxy/tcp"
	"github.com/ninjablocks/mqtt-proxy/ws"
)

var configFile = flag.String("config", "config.toml", "configuration file")
var debug = flag.Bool("debug", false, "enable debugging")
var version = flag.Bool("version", false, "show version")

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Git: %s\n", GitCommit)
		os.Exit(0)
	}

	conf := conf.LoadConfiguration(*configFile)

	if *debug {
		log.Printf("[main] conf %+v", conf)
	}

	proxy := proxy.CreateMQTTProxy(conf)

	handlers := ws.CreateHttpHanders(proxy)
	tcpServer := tcp.CreateTcpServer(proxy)

	go handlers.StartServer(&conf.Http)

	go tcpServer.StartServer(&conf.Mqtt)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	log.Printf("Got signal %s, exiting now", <-c)
}
