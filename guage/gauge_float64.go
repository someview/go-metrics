package guage

import (
	"math"
	"sync/atomic"
)

// GaugeFloat64s hold a float64 value that can be set arbitrarily.
type GaugeFloat64 interface {
	Snapshot() GaugeFloat64
	Update(float64)
	Value() float64
}

// NewGaugeFloat64 constructs a new StandardGaugeFloat64.
func NewGaugeFloat64() GaugeFloat64 {
	return &StandardGaugeFloat64{
		value: 0.0,
	}
}

// NewFunctionalGauge constructs a new FunctionalGauge.
func NewFunctionalGaugeFloat64(f func() float64) GaugeFloat64 {
	return &FunctionalGaugeFloat64{value: f}
}

// GaugeFloat64Snapshot is a read-only copy of another GaugeFloat64.
type GaugeFloat64Snapshot float64

// Snapshot returns the snapshot.
func (g GaugeFloat64Snapshot) Snapshot() GaugeFloat64 { return g }

// Update panics.
func (GaugeFloat64Snapshot) Update(float64) {
	panic("Update called on a GaugeFloat64Snapshot")
}

// Value returns the value at the time the snapshot was taken.
func (g GaugeFloat64Snapshot) Value() float64 { return float64(g) }

// StandardGaugeFloat64 is the standard implementation of a GaugeFloat64 and uses
// sync.Mutex to manage a single float64 value.
type StandardGaugeFloat64 struct {
	value uint64
}

// Snapshot returns a read-only copy of the gauge.
func (g *StandardGaugeFloat64) Snapshot() GaugeFloat64 {
	return GaugeFloat64Snapshot(g.Value())
}

// Update updates the gauge's value.
func (g *StandardGaugeFloat64) Update(v float64) {
	atomic.StoreUint64(&g.value, math.Float64bits(v))
}

// Value returns the gauge's current value.
func (g *StandardGaugeFloat64) Value() float64 {
	return math.Float64frombits(atomic.LoadUint64(&g.value))
}

// FunctionalGaugeFloat64 returns value from given function
type FunctionalGaugeFloat64 struct {
	value func() float64
}

// Value returns the gauge's current value.
func (g FunctionalGaugeFloat64) Value() float64 {
	return g.value()
}

// Snapshot returns the snapshot.
func (g FunctionalGaugeFloat64) Snapshot() GaugeFloat64 { return GaugeFloat64Snapshot(g.Value()) }

// Update panics.
func (FunctionalGaugeFloat64) Update(float64) {
	panic("Update called on a FunctionalGaugeFloat64")
}
