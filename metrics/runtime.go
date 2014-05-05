package metrics

import (
	"os"
	"runtime"
	"time"

	sigar "github.com/cloudfoundry/gosigar"
	gmetrics "github.com/rcrowley/go-metrics"
)

type MetricsGroup interface {
	Update()
}

type RuntimeMetrics struct {
	Internals, Memory, Cpu, Load MetricsGroup
}

func NewRuntimeMetrics(prefix string) *RuntimeMetrics {
	return &RuntimeMetrics{
		Internals: NewGoInternalMetrics(prefix),
		Memory:    NewProcessMemoryMetrics(prefix),
		Cpu:       NewProcessCpuMetrics(prefix),
		Load:      NewLoadMetrics(prefix),
	}
}

func (rm *RuntimeMetrics) Update() {
	rm.Internals.Update()
	rm.Memory.Update()
	rm.Cpu.Update()
	rm.Load.Update()
}

// system load
type LoadMetrics struct {
	One, Five, Fifteen gmetrics.GaugeFloat64
}

func NewLoadMetrics(prefix string) *LoadMetrics {

	load := &LoadMetrics{
		One:     gmetrics.NewGaugeFloat64(),
		Five:    gmetrics.NewGaugeFloat64(),
		Fifteen: gmetrics.NewGaugeFloat64(),
	}

	gmetrics.Register(prefix+".load.1min", load.One)
	gmetrics.Register(prefix+".load.5min", load.Five)
	gmetrics.Register(prefix+".load.15min", load.Fifteen)

	return load
}

func (lm *LoadMetrics) Update() {

	load := sigar.LoadAverage{}

	err := load.Get()

	if err == nil {
		lm.One.Update(load.One)
		lm.Five.Update(load.Five)
		lm.Fifteen.Update(load.Fifteen)
	}

}

// process memory
type ProcessMemoryMetrics struct {
	Resident, Shared gmetrics.Gauge
	PageFaults       gmetrics.Meter
}

func NewProcessMemoryMetrics(prefix string) *ProcessMemoryMetrics {

	mem := &ProcessMemoryMetrics{
		Resident:   gmetrics.NewGauge(),
		Shared:     gmetrics.NewGauge(),
		PageFaults: gmetrics.NewMeter(),
	}

	gmetrics.Register(prefix+".mem.resident", mem.Resident)
	gmetrics.Register(prefix+".mem.shared", mem.Shared)
	gmetrics.Register(prefix+".mem.pagefaults", mem.PageFaults)

	return mem
}
func (pmm *ProcessMemoryMetrics) Update() {
	pid := os.Getpid()

	mem := sigar.ProcMem{}

	err := mem.Get(pid)

	if err == nil {
		pmm.Resident.Update(int64(mem.Resident))
		pmm.Shared.Update(int64(mem.Share))

		updateMeter(pmm.PageFaults, mem.PageFaults)
	}
}

// process cpu
type ProcessCpuMetrics struct {
	User, Sys, Total gmetrics.Meter
}

func NewProcessCpuMetrics(prefix string) *ProcessCpuMetrics {
	cpu := &ProcessCpuMetrics{
		User:  gmetrics.NewMeter(),
		Sys:   gmetrics.NewMeter(),
		Total: gmetrics.NewMeter(),
	}

	gmetrics.Register(prefix+".cpu.user", cpu.User)
	gmetrics.Register(prefix+".cpu.sys", cpu.Sys)
	gmetrics.Register(prefix+".cpu.total", cpu.Total)

	return cpu
}

func (pcm *ProcessCpuMetrics) Update() {
	pid := os.Getpid()
	cpu := sigar.ProcTime{}

	err := cpu.Get(pid)

	if err == nil {
		updateMeter(pcm.User, cpu.User)
		updateMeter(pcm.Sys, cpu.Sys)
		updateMeter(pcm.Total, cpu.Total)
	}

}

// golang internals
type GoInternalMetrics struct {
	Alloc, TotalAlloc, NumGoroutine gmetrics.Gauge
}

func NewGoInternalMetrics(prefix string) *GoInternalMetrics {
	goint := &GoInternalMetrics{
		Alloc:        gmetrics.NewGauge(),
		TotalAlloc:   gmetrics.NewGauge(),
		NumGoroutine: gmetrics.NewGauge(),
	}

	gmetrics.Register(prefix+".goint.alloc", goint.Alloc)
	gmetrics.Register(prefix+".goint.total_alloc", goint.TotalAlloc)
	gmetrics.Register(prefix+".goint.go_routines", goint.NumGoroutine)

	return goint
}

func (gim *GoInternalMetrics) Update() {

	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)

	gim.Alloc.Update(int64(ms.Alloc))
	gim.TotalAlloc.Update(int64(ms.TotalAlloc))
	gim.NumGoroutine.Update(int64(runtime.NumGoroutine()))

}

func updateMeter(meter gmetrics.Meter, newValue uint64) {
	va := int64(newValue) - meter.Count()
	meter.Mark(int64(va))
}

func StartRuntimeMetricsJob(env string, region string) {

	prefix := buildPrefix(env, region)

	rm := NewRuntimeMetrics(prefix)

	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for _ = range ticker.C {
			rm.Update()
		}
	}()
}
