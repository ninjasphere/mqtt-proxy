package metrics

import (
	"log"
	"time"

	gmetrics "github.com/rcrowley/go-metrics"
)

func logForever(r gmetrics.Registry, d time.Duration) {
	for {
		r.Each(func(name string, i interface{}) {
			switch m := i.(type) {
			case gmetrics.Counter:
				log.Printf("counter %s\n", name)
				log.Printf("  count:       %9d\n", m.Count())
			case gmetrics.Gauge:
				log.Printf("gauge %s\n", name)
				log.Printf("  value:       %9d\n", m.Value())
			case gmetrics.Healthcheck:
				m.Check()
				log.Printf("healthcheck %s\n", name)
				log.Printf("  error:       %v\n", m.Error())
			case gmetrics.Histogram:
				ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				log.Printf("histogram %s\n", name)
				log.Printf("  count:       %9d\n", m.Count())
				log.Printf("  min:         %9d\n", m.Min())
				log.Printf("  max:         %9d\n", m.Max())
				log.Printf("  mean:        %12.2f\n", m.Mean())
				log.Printf("  stddev:      %12.2f\n", m.StdDev())
				log.Printf("  median:      %12.2f\n", ps[0])
				log.Printf("  75%%:         %12.2f\n", ps[1])
				log.Printf("  95%%:         %12.2f\n", ps[2])
				log.Printf("  99%%:         %12.2f\n", ps[3])
				log.Printf("  99.9%%:       %12.2f\n", ps[4])
			case gmetrics.Meter:
				log.Printf("meter %s\n", name)
				log.Printf("  count:       %9d\n", m.Count())
				log.Printf("  1-min rate:  %12.2f\n", m.Rate1())
				log.Printf("  5-min rate:  %12.2f\n", m.Rate5())
				log.Printf("  15-min rate: %12.2f\n", m.Rate15())
				log.Printf("  mean rate:   %12.2f\n", m.RateMean())
			case gmetrics.Timer:
				ps := m.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				log.Printf("timer %s\n", name)
				log.Printf("  count:       %9d\n", m.Count())
				log.Printf("  min:         %9d\n", m.Min())
				log.Printf("  max:         %9d\n", m.Max())
				log.Printf("  mean:        %12.2f\n", m.Mean())
				log.Printf("  stddev:      %12.2f\n", m.StdDev())
				log.Printf("  median:      %12.2f\n", ps[0])
				log.Printf("  75%%:         %12.2f\n", ps[1])
				log.Printf("  95%%:         %12.2f\n", ps[2])
				log.Printf("  99%%:         %12.2f\n", ps[3])
				log.Printf("  99.9%%:       %12.2f\n", ps[4])
				log.Printf("  1-min rate:  %12.2f\n", m.Rate1())
				log.Printf("  5-min rate:  %12.2f\n", m.Rate5())
				log.Printf("  15-min rate: %12.2f\n", m.Rate15())
				log.Printf("  mean rate:   %12.2f\n", m.RateMean())
			}
		})
		time.Sleep(d)
	}
}

func ConsoleOutput() {
	log.Println("starting metrics job")
	go logForever(gmetrics.DefaultRegistry, 10e9)
}
