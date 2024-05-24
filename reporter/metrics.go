package reporter

import (
	"github.com/someview/go-metrics/counter"
	"github.com/someview/go-metrics/guage"
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
