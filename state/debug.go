package state

import (
	"github.com/someview/go-metrics"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"github.com/someview/go-metrics/sample"
	"github.com/someview/go-metrics/timer"
	"runtime/debug"
	"sync"
	"time"
)

var (
	debugMetrics struct {
		GCStats struct {
			LastGC guage.Gauge
			NumGC  guage.Gauge
			Pause  histogram.Histogram
			//PauseQuantiles Histogram
			PauseTotal guage.Gauge
		}
		ReadGCStats timer.Timer
	}
	gcStats                  debug.GCStats
	registerDebugMetricsOnce = sync.Once{}
)

// Capture new values for the Go garbage collector statistics exported in
// debug.GCStats.  This is designed to be called as a goroutine.
func CaptureDebugGCStats(r metrics.Registry, d time.Duration) {
	for _ = range time.Tick(d) {
		CaptureDebugGCStatsOnce(r)
	}
}

// Capture new values for the Go garbage collector statistics exported in
// debug.GCStats.  This is designed to be called in a background goroutine.
// Giving a registry which has not been given to RegisterDebugGCStats will
// panic.
//
// Be careful (but much less so) with this because debug.ReadGCStats calls
// the C function state·lock(state·mheap) which, while not a stop-the-world
// operation, isn't something you want to be doing all the time.
func CaptureDebugGCStatsOnce(r metrics.Registry) {
	lastGC := gcStats.LastGC
	t := time.Now()
	debug.ReadGCStats(&gcStats)
	debugMetrics.ReadGCStats.UpdateSince(t)

	debugMetrics.GCStats.LastGC.Update(int64(gcStats.LastGC.UnixNano()))
	debugMetrics.GCStats.NumGC.Update(int64(gcStats.NumGC))
	if lastGC != gcStats.LastGC && 0 < len(gcStats.Pause) {
		debugMetrics.GCStats.Pause.Update(int64(gcStats.Pause[0]))
	}
	//debugMetrics.GCStats.PauseQuantiles.Update(gcStats.PauseQuantiles)
	debugMetrics.GCStats.PauseTotal.Update(int64(gcStats.PauseTotal))
}

// Register metrics for the Go garbage collector statistics exported in
// debug.GCStats.  The metrics are named by their fully-qualified Go symbols,
// i.e. debug.GCStats.PauseTotal.
func RegisterDebugGCStats(r metrics.Registry) {
	registerDebugMetricsOnce.Do(func() {
		debugMetrics.GCStats.LastGC = guage.NewGauge()
		debugMetrics.GCStats.NumGC = guage.NewGauge()
		debugMetrics.GCStats.Pause = histogram.NewHistogram(sample.NewExpDecaySample(1028, 0.015))
		//debugMetrics.GCStats.PauseQuantiles = NewHistogram(NewExpDecaySample(1028, 0.015))
		debugMetrics.GCStats.PauseTotal = guage.NewGauge()
		debugMetrics.ReadGCStats = timer.NewTimer()

		r.Register("debug.GCStats.LastGC", debugMetrics.GCStats.LastGC)
		r.Register("debug.GCStats.NumGC", debugMetrics.GCStats.NumGC)
		r.Register("debug.GCStats.Pause", debugMetrics.GCStats.Pause)
		//r.Register("debug.GCStats.PauseQuantiles", debugMetrics.GCStats.PauseQuantiles)
		r.Register("debug.GCStats.PauseTotal", debugMetrics.GCStats.PauseTotal)
		r.Register("debug.ReadGCStats", debugMetrics.ReadGCStats)
	})
}

// Allocate an initial slice for gcStats.Pause to avoid allocations during
// normal operation.
func init() {
	gcStats.Pause = make([]time.Duration, 11)
}
