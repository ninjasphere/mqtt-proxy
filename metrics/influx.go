package metrics

import (
	"github.com/ninjablocks/mqtt-proxy/conf"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/influxdb"
)

func UploadToInflux(config *conf.InfluxConfiguration) {
	if config.Host != "" {
		go influxdb.Influxdb(metrics.DefaultRegistry, 10e9, &influxdb.Config{
			Host:     config.Host,
			Database: config.Database,
			Username: config.User,
			Password: config.Pass,
		})
	}
}
