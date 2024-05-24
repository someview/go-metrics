package timer

import (
	"github.com/someview/go-metrics/histogram"
	"github.com/someview/go-metrics/meter"
	"github.com/someview/go-metrics/sample"
	"sync"
	"time"
)

// Timers capture the duration and rate of events.
type Timer interface {
	Count() int64
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	Percentiles([]float64) []float64
	Rate1() float64
	Rate5() float64
	Rate15() float64
	RateMean() float64
	Snapshot() Timer
	StdDev() float64
	Stop()
	Sum() int64
	Time(func())
	Update(time.Duration)
	UpdateSince(time.Time)
	Variance() float64
}

// NewCustomTimer constructs a new StandardTimer from a Histogram and a Meter.
// Be sure to call Stop() once the timer is of no use to allow for garbage collection.
func NewCustomTimer(h histogram.Histogram, m meter.Meter) Timer {
	return &StandardTimer{
		histogram: h,
		meter:     m,
	}
}

// NewTimer constructs a new StandardTimer using an exponentially-decaying
// sample with the same reservoir size and alpha as UNIX load averages.
// Be sure to call Stop() once the timer is of no use to allow for garbage collection.
func NewTimer() Timer {
	return &StandardTimer{
		histogram: histogram.NewHistogram(sample.NewExpDecaySample(1028, 0.015)),
		meter:     meter.NewMeter(),
	}
}

// NilTimer is a no-op Timer.
type NilTimer struct {
	h histogram.Histogram
	m meter.Meter
}

// StandardTimer is the standard implementation of a Timer and uses a Histogram
// and Meter.
type StandardTimer struct {
	histogram histogram.Histogram
	meter     meter.Meter
	mutex     sync.Mutex
}

// Count returns the number of events recorded.
func (t *StandardTimer) Count() int64 {
	return t.histogram.Count()
}

// Max returns the maximum value in the sample.
func (t *StandardTimer) Max() int64 {
	return t.histogram.Max()
}

// Mean returns the mean of the values in the sample.
func (t *StandardTimer) Mean() float64 {
	return t.histogram.Mean()
}

// Min returns the minimum value in the sample.
func (t *StandardTimer) Min() int64 {
	return t.histogram.Min()
}

// Percentile returns an arbitrary percentile of the values in the sample.
func (t *StandardTimer) Percentile(p float64) float64 {
	return t.histogram.Percentile(p)
}

// Percentiles returns a slice of arbitrary percentiles of the values in the
// sample.
func (t *StandardTimer) Percentiles(ps []float64) []float64 {
	return t.histogram.Percentiles(ps)
}

// Rate1 returns the one-minute moving average rate of events per second.
func (t *StandardTimer) Rate1() float64 {
	return t.meter.Rate1()
}

// Rate5 returns the five-minute moving average rate of events per second.
func (t *StandardTimer) Rate5() float64 {
	return t.meter.Rate5()
}

// Rate15 returns the fifteen-minute moving average rate of events per second.
func (t *StandardTimer) Rate15() float64 {
	return t.meter.Rate15()
}

// RateMean returns the meter's mean rate of events per second.
func (t *StandardTimer) RateMean() float64 {
	return t.meter.RateMean()
}

// Snapshot returns a read-only copy of the timer.
func (t *StandardTimer) Snapshot() Timer {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return &TimerSnapshot{
		histogram: t.histogram.Snapshot().(*histogram.HistogramSnapshot),
		meter:     t.meter.Snapshot().(*meter.MeterSnapshot),
	}
}

// StdDev returns the standard deviation of the values in the sample.
func (t *StandardTimer) StdDev() float64 {
	return t.histogram.StdDev()
}

// Stop stops the meter.
func (t *StandardTimer) Stop() {
	t.meter.Stop()
}

// Sum returns the sum in the sample.
func (t *StandardTimer) Sum() int64 {
	return t.histogram.Sum()
}

// Record the duration of the execution of the given function.
func (t *StandardTimer) Time(f func()) {
	ts := time.Now()
	f()
	t.Update(time.Since(ts))
}

// Record the duration of an event.
func (t *StandardTimer) Update(d time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.histogram.Update(int64(d))
	t.meter.Mark(1)
}

// Record the duration of an event that started at a time and ends now.
func (t *StandardTimer) UpdateSince(ts time.Time) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.histogram.Update(int64(time.Since(ts)))
	t.meter.Mark(1)
}

// Variance returns the variance of the values in the sample.
func (t *StandardTimer) Variance() float64 {
	return t.histogram.Variance()
}

// TimerSnapshot is a read-only copy of another Timer.
type TimerSnapshot struct {
	histogram *histogram.HistogramSnapshot
	meter     *meter.MeterSnapshot
}

// Count returns the number of events recorded at the time the snapshot was
// taken.
func (t *TimerSnapshot) Count() int64 { return t.histogram.Count() }

// Max returns the maximum value at the time the snapshot was taken.
func (t *TimerSnapshot) Max() int64 { return t.histogram.Max() }

// Mean returns the mean value at the time the snapshot was taken.
func (t *TimerSnapshot) Mean() float64 { return t.histogram.Mean() }

// Min returns the minimum value at the time the snapshot was taken.
func (t *TimerSnapshot) Min() int64 { return t.histogram.Min() }

// Percentile returns an arbitrary percentile of sampled values at the time the
// snapshot was taken.
func (t *TimerSnapshot) Percentile(p float64) float64 {
	return t.histogram.Percentile(p)
}

// Percentiles returns a slice of arbitrary percentiles of sampled values at
// the time the snapshot was taken.
func (t *TimerSnapshot) Percentiles(ps []float64) []float64 {
	return t.histogram.Percentiles(ps)
}

// Rate1 returns the one-minute moving average rate of events per second at the
// time the snapshot was taken.
func (t *TimerSnapshot) Rate1() float64 { return t.meter.Rate1() }

// Rate5 returns the five-minute moving average rate of events per second at
// the time the snapshot was taken.
func (t *TimerSnapshot) Rate5() float64 { return t.meter.Rate5() }

// Rate15 returns the fifteen-minute moving average rate of events per second
// at the time the snapshot was taken.
func (t *TimerSnapshot) Rate15() float64 { return t.meter.Rate15() }

// RateMean returns the meter's mean rate of events per second at the time the
// snapshot was taken.
func (t *TimerSnapshot) RateMean() float64 { return t.meter.RateMean() }

// Snapshot returns the snapshot.
func (t *TimerSnapshot) Snapshot() Timer { return t }

// StdDev returns the standard deviation of the values at the time the snapshot
// was taken.
func (t *TimerSnapshot) StdDev() float64 { return t.histogram.StdDev() }

// Stop is a no-op.
func (t *TimerSnapshot) Stop() {}

// Sum returns the sum at the time the snapshot was taken.
func (t *TimerSnapshot) Sum() int64 { return t.histogram.Sum() }

// Time panics.
func (*TimerSnapshot) Time(func()) {
	panic("Time called on a TimerSnapshot")
}

// Update panics.
func (*TimerSnapshot) Update(time.Duration) {
	panic("Update called on a TimerSnapshot")
}

// UpdateSince panics.
func (*TimerSnapshot) UpdateSince(time.Time) {
	panic("UpdateSince called on a TimerSnapshot")
}

// Variance returns the variance of the values at the time the snapshot was
// taken.
func (t *TimerSnapshot) Variance() float64 { return t.histogram.Variance() }