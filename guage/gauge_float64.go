package guage

import (
	"math"
	"sync/atomic"
)

// GaugeFloat64s hold a float64 value that can be set arbitrarily.
type GaugeFloat64 interface {
	Update(float64)
	Snapshot() float64
	SnapshotAndReset() float64
}

// NewGaugeFloat64 constructs a new StandardGaugeFloat64.
func NewGaugeFloat64() GaugeFloat64 {
	return &StandardGaugeFloat64{
		value: 0.0,
	}
}

// StandardGaugeFloat64 is the standard implementation of a GaugeFloat64 and uses
// sync.Mutex to manage a single float64 value.
type StandardGaugeFloat64 struct {
	value uint64
}

func (g *StandardGaugeFloat64) Value() float64 {
	//TODO implement me
	panic("implement me")
}

func (g *StandardGaugeFloat64) Snapshot() float64 {
	return math.Float64frombits(atomic.LoadUint64(&g.value))
}

// Update updates the gauge's value.
func (g *StandardGaugeFloat64) Update(v float64) {
	atomic.StoreUint64(&g.value, math.Float64bits(v))
}
