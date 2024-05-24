package reporter

import (
	"github.com/someview/go-metrics/counter"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"github.com/someview/go-metrics/meter"
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

// GetOrRegisterHistogram returns an existing Histogram or constructs and
// registers a new StandardHistogram.
func GetOrRegisterHistogram(name string, r Registry, s sample.Sample) histogram.Histogram {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, func() histogram.Histogram { return histogram.NewHistogram(s) }).(histogram.Histogram)
}

// GetOrRegisterMeter returns an existing Meter or constructs and registers a
// new StandardMeter.
// Be sure to unregister the meter from the registry once it is of no use to
// allow for garbage collection.
func GetOrRegisterMeter(name string, r Registry) meter.Meter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, meter.NewMeter).(meter.Meter)
}
