package guage

import (
	"sync/atomic"
)

// Gauges hold an int64 value that can be set arbitrarily.
type Gauge interface {
	Update(int64)
	Snapshot() int64
	SnapShotAndReset() int64
}

// standardGauge is the standard implementation of a Gauge and uses the
// sync/atomic package to manage a single int64 value.
type standardGauge struct {
	value int64
}

// Update updates the gauge's value.
func (g *standardGauge) Update(v int64) {
	atomic.StoreInt64(&g.value, v)
}

func (g *standardGauge) Snapshot() int64 {
	return atomic.LoadInt64(&g.value)
}

func (g *standardGauge) SnapShotAndReset() int64 {
	return atomic.SwapInt64(&g.value, 0)
}

// NewGauge constructs a new standardGauge.
func NewGauge() Gauge {
	return &standardGauge{0}
}
