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

		go librato.Librato(metrics.DefaultRegistry, 30e9, config.Email, config.Token, hostname, []float64{95}, time.Millisecond)
	}
}
