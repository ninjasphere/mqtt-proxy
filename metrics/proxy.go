package metrics

import (
	"strings"

	"github.com/ninjablocks/mqtt-proxy/conf"
	gmetrics "github.com/rcrowley/go-metrics"
)

// Record proxy related meters to enable monitoring
// of throughput and volume.
type ProxyMetrics struct {
	Msgs        gmetrics.Meter
	MsgReply    gmetrics.Meter
	MsgForward  gmetrics.Meter
	MsgBodySize gmetrics.Histogram
	Connects    gmetrics.Meter
	Connections gmetrics.Gauge
}

// conf.Environment, conf.Region
func NewProxyMetrics(env string, region string) ProxyMetrics {

	prefix := buildPrefix(env, region)

	pm := ProxyMetrics{
		Msgs:        gmetrics.NewMeter(),
		MsgReply:    gmetrics.NewMeter(),
		MsgForward:  gmetrics.NewMeter(),
		MsgBodySize: gmetrics.NewHistogram(gmetrics.NewExpDecaySample(1028, 0.015)),
		Connects:    gmetrics.NewMeter(),
		Connections: gmetrics.NewGauge(),
	}

	gmetrics.Register(prefix+".proxy.msgs", pm.Msgs)
	gmetrics.Register(prefix+".proxy.msg_reply", pm.Msgs)
	gmetrics.Register(prefix+".proxy.msg_forward", pm.Msgs)
	gmetrics.Register(prefix+".proxy.msg_body_size", pm.MsgBodySize)
	gmetrics.Register(prefix+".proxy.connects", pm.Connects)
	gmetrics.Register(prefix+".proxy.connections", pm.Connections)

	return pm
}

func buildPrefix(env string, region string) string {
	return strings.Join([]string{region, env, "mqtt-proxy"}, ".")
}

func StartMetricsJobs(config *conf.Configuration) {
	StartRuntimeMetricsJob(config.Environment, config.Region)
	UploadToInflux(&config.Influx)
	UploadToLibrato(&config.Librato)
}
