package reporter

import (
	"github.com/someview/go-metrics/counter"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"github.com/someview/go-metrics/sample"
)

// GetOrRegisterCounter returns an existing Counter or constructs and registers
// a new StandardCounter.
func GetOrRegisterCounter(name string, r Registry) counter.Counter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, counter.NewCounter).(counter.Counter)
}

// GetOrRegisterGauge returns an existing Gauge or constructs and registers a
// new StandardGauge.
func GetOrRegisterGauge(name string, r Registry) guage.Gauge {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, guage.NewGauge).(guage.Gauge)
}

// GetOrRegisterGaugeFloat64 returns an existing GaugeFloat64 or constructs and registers a
// new StandardGaugeFloat64.
func GetOrRegisterGaugeFloat64(name string, r Registry) guage.GaugeFloat64 {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, guage.NewGaugeFloat64()).(guage.GaugeFloat64)
}

// GetOrRegisterHistogram returns an existing Histogram or constructs and
// registers a new StandardHistogram.
func GetOrRegisterHistogram(name string, r Registry, s sample.Sample) histogram.Histogram {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, func() histogram.Histogram { return histogram.NewHistogram(s) }).(histogram.Histogram)
}
