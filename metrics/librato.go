package metrics

import (
	"log"
	"os"
	"time"

	"github.com/ninjablocks/mqtt-proxy/conf"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/librato"
)

func UploadToLibrato(config *conf.LibratoConfiguration) {
	if config.Email != "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("Unable to retrieve a hostname %s", err)
		}

		go librato.Librato(metrics.DefaultRegistry,
			30e9,          // interval
			config.Email,  // account email addres
			config.Token,  // auth token
			hostname,      // source
			[]float64{95}, // precentiles to send
			time.Millisecond)
	}
}
