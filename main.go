package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ninjablocks/mqtt-proxy/conf"
	"github.com/ninjablocks/mqtt-proxy/metrics"
	"github.com/ninjablocks/mqtt-proxy/proxy"
	"github.com/ninjablocks/mqtt-proxy/tcp"
)

var configFile = flag.String("config", "config.toml", "configuration file")
var debug = flag.Bool("debug", false, "enable debugging")
var version = flag.Bool("version", false, "show version")

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	conf := conf.LoadConfiguration(*configFile)

	if *debug {
		log.Printf("[main] conf %+v", conf)
	}

	p := proxy.CreateMQTTProxy(conf)

	// asign the servers
	tcpServer := tcp.CreateTcpServer(p)

	go tcpServer.StartServer(&conf.Mqtt)

	metrics.StartRuntimeMetricsJob(conf.Environment, conf.Region)

	metrics.UploadToInflux(&conf.Influx)

	metrics.UploadToLibrato(&conf.Librato)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	log.Printf("Got signal %s, exiting now", <-c)
}
