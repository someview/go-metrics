package state

import (
	"github.com/someview/go-metrics"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"github.com/someview/go-metrics/sample"
	"github.com/someview/go-metrics/timer"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

var (
	memStats       runtime.MemStats
	runtimeMetrics struct {
		MemStats struct {
			Alloc         guage.Gauge
			BuckHashSys   guage.Gauge
			DebugGC       guage.Gauge
			EnableGC      guage.Gauge
			Frees         guage.Gauge
			HeapAlloc     guage.Gauge
			HeapIdle      guage.Gauge
			HeapInuse     guage.Gauge
			HeapObjects   guage.Gauge
			HeapReleased  guage.Gauge
			HeapSys       guage.Gauge
			LastGC        guage.Gauge
			Lookups       guage.Gauge
			Mallocs       guage.Gauge
			MCacheInuse   guage.Gauge
			MCacheSys     guage.Gauge
			MSpanInuse    guage.Gauge
			MSpanSys      guage.Gauge
			NextGC        guage.Gauge
			NumGC         guage.Gauge
			GCCPUFraction guage.GaugeFloat64
			PauseNs       histogram.Histogram
			PauseTotalNs  guage.Gauge
			StackInuse    guage.Gauge
			StackSys      guage.Gauge
			Sys           guage.Gauge
			TotalAlloc    guage.Gauge
		}
		NumCgoCall   guage.Gauge
		NumGoroutine guage.Gauge
		NumThread    guage.Gauge
		ReadMemStats timer.Timer
	}
	frees       uint64
	lookups     uint64
	mallocs     uint64
	numGC       uint32
	numCgoCalls int64

	threadCreateProfile        = pprof.Lookup("threadcreate")
	registerRuntimeMetricsOnce = sync.Once{}
)

// Capture new values for the Go state statistics exported in
// state.MemStats.  This is designed to be called as a goroutine.
func CaptureRuntimeMemStats(r metrics.Registry, d time.Duration) {
	for _ = range time.Tick(d) {
		CaptureRuntimeMemStatsOnce(r)
	}
}

// Capture new values for the Go state statistics exported in
// state.MemStats.  This is designed to be called in a background
// goroutine.  Giving a registry which has not been given to
// RegisterRuntimeMemStats will panic.
//
// Be very careful with this because state.ReadMemStats calls the C
// functions state·semacquire(&state·worldsema) and state·stoptheworld()
// and that last one does what it says on the tin.
func CaptureRuntimeMemStatsOnce(r metrics.Registry) {
	t := time.Now()
	runtime.ReadMemStats(&memStats) // This takes 50-200us.
	runtimeMetrics.ReadMemStats.UpdateSince(t)

	runtimeMetrics.MemStats.Alloc.Update(int64(memStats.Alloc))
	runtimeMetrics.MemStats.BuckHashSys.Update(int64(memStats.BuckHashSys))
	if memStats.DebugGC {
		runtimeMetrics.MemStats.DebugGC.Update(1)
	} else {
		runtimeMetrics.MemStats.DebugGC.Update(0)
	}
	if memStats.EnableGC {
		runtimeMetrics.MemStats.EnableGC.Update(1)
	} else {
		runtimeMetrics.MemStats.EnableGC.Update(0)
	}

	runtimeMetrics.MemStats.Frees.Update(int64(memStats.Frees - frees))
	runtimeMetrics.MemStats.HeapAlloc.Update(int64(memStats.HeapAlloc))
	runtimeMetrics.MemStats.HeapIdle.Update(int64(memStats.HeapIdle))
	runtimeMetrics.MemStats.HeapInuse.Update(int64(memStats.HeapInuse))
	runtimeMetrics.MemStats.HeapObjects.Update(int64(memStats.HeapObjects))
	runtimeMetrics.MemStats.HeapReleased.Update(int64(memStats.HeapReleased))
	runtimeMetrics.MemStats.HeapSys.Update(int64(memStats.HeapSys))
	runtimeMetrics.MemStats.LastGC.Update(int64(memStats.LastGC))
	runtimeMetrics.MemStats.Lookups.Update(int64(memStats.Lookups - lookups))
	runtimeMetrics.MemStats.Mallocs.Update(int64(memStats.Mallocs - mallocs))
	runtimeMetrics.MemStats.MCacheInuse.Update(int64(memStats.MCacheInuse))
	runtimeMetrics.MemStats.MCacheSys.Update(int64(memStats.MCacheSys))
	runtimeMetrics.MemStats.MSpanInuse.Update(int64(memStats.MSpanInuse))
	runtimeMetrics.MemStats.MSpanSys.Update(int64(memStats.MSpanSys))
	runtimeMetrics.MemStats.NextGC.Update(int64(memStats.NextGC))
	runtimeMetrics.MemStats.NumGC.Update(int64(memStats.NumGC - numGC))
	runtimeMetrics.MemStats.GCCPUFraction.Update(gcCPUFraction(&memStats))

	// <https://code.google.com/p/go/source/browse/src/pkg/runtime/mgc0.c>
	i := numGC % uint32(len(memStats.PauseNs))
	ii := memStats.NumGC % uint32(len(memStats.PauseNs))
	if memStats.NumGC-numGC >= uint32(len(memStats.PauseNs)) {
		for i = 0; i < uint32(len(memStats.PauseNs)); i++ {
			runtimeMetrics.MemStats.PauseNs.Update(int64(memStats.PauseNs[i]))
		}
	} else {
		if i > ii {
			for ; i < uint32(len(memStats.PauseNs)); i++ {
				runtimeMetrics.MemStats.PauseNs.Update(int64(memStats.PauseNs[i]))
			}
			i = 0
		}
		for ; i < ii; i++ {
			runtimeMetrics.MemStats.PauseNs.Update(int64(memStats.PauseNs[i]))
		}
	}
	frees = memStats.Frees
	lookups = memStats.Lookups
	mallocs = memStats.Mallocs
	numGC = memStats.NumGC

	runtimeMetrics.MemStats.PauseTotalNs.Update(int64(memStats.PauseTotalNs))
	runtimeMetrics.MemStats.StackInuse.Update(int64(memStats.StackInuse))
	runtimeMetrics.MemStats.StackSys.Update(int64(memStats.StackSys))
	runtimeMetrics.MemStats.Sys.Update(int64(memStats.Sys))
	runtimeMetrics.MemStats.TotalAlloc.Update(int64(memStats.TotalAlloc))

	currentNumCgoCalls := numCgoCall()
	runtimeMetrics.NumCgoCall.Update(currentNumCgoCalls - numCgoCalls)
	numCgoCalls = currentNumCgoCalls

	runtimeMetrics.NumGoroutine.Update(int64(runtime.NumGoroutine()))

	runtimeMetrics.NumThread.Update(int64(threadCreateProfile.Count()))
}

// Register runtimeMetrics for the Go state statistics exported in state and
// specifically state.MemStats.  The runtimeMetrics are named by their
// fully-qualified Go symbols, i.e. state.MemStats.Alloc.
func RegisterRuntimeMemStats(r metrics.Registry) {
	registerRuntimeMetricsOnce.Do(func() {
		runtimeMetrics.MemStats.Alloc = guage.NewGauge()
		runtimeMetrics.MemStats.BuckHashSys = guage.NewGauge()
		runtimeMetrics.MemStats.DebugGC = guage.NewGauge()
		runtimeMetrics.MemStats.EnableGC = guage.NewGauge()
		runtimeMetrics.MemStats.Frees = guage.NewGauge()
		runtimeMetrics.MemStats.HeapAlloc = guage.NewGauge()
		runtimeMetrics.MemStats.HeapIdle = guage.NewGauge()
		runtimeMetrics.MemStats.HeapInuse = guage.NewGauge()
		runtimeMetrics.MemStats.HeapObjects = guage.NewGauge()
		runtimeMetrics.MemStats.HeapReleased = guage.NewGauge()
		runtimeMetrics.MemStats.HeapSys = guage.NewGauge()
		runtimeMetrics.MemStats.LastGC = guage.NewGauge()
		runtimeMetrics.MemStats.Lookups = guage.NewGauge()
		runtimeMetrics.MemStats.Mallocs = guage.NewGauge()
		runtimeMetrics.MemStats.MCacheInuse = guage.NewGauge()
		runtimeMetrics.MemStats.MCacheSys = guage.NewGauge()
		runtimeMetrics.MemStats.MSpanInuse = guage.NewGauge()
		runtimeMetrics.MemStats.MSpanSys = guage.NewGauge()
		runtimeMetrics.MemStats.NextGC = guage.NewGauge()
		runtimeMetrics.MemStats.NumGC = guage.NewGauge()
		runtimeMetrics.MemStats.GCCPUFraction = guage.NewGaugeFloat64()
		runtimeMetrics.MemStats.PauseNs = histogram.NewHistogram(sample.NewExpDecaySample(1028, 0.015))
		runtimeMetrics.MemStats.PauseTotalNs = guage.NewGauge()
		runtimeMetrics.MemStats.StackInuse = guage.NewGauge()
		runtimeMetrics.MemStats.StackSys = guage.NewGauge()
		runtimeMetrics.MemStats.Sys = guage.NewGauge()
		runtimeMetrics.MemStats.TotalAlloc = guage.NewGauge()
		runtimeMetrics.NumCgoCall = guage.NewGauge()
		runtimeMetrics.NumGoroutine = guage.NewGauge()
		runtimeMetrics.NumThread = guage.NewGauge()
		runtimeMetrics.ReadMemStats = timer.NewTimer()

		r.Register("state.MemStats.Alloc", runtimeMetrics.MemStats.Alloc)
		r.Register("state.MemStats.BuckHashSys", runtimeMetrics.MemStats.BuckHashSys)
		r.Register("state.MemStats.DebugGC", runtimeMetrics.MemStats.DebugGC)
		r.Register("state.MemStats.EnableGC", runtimeMetrics.MemStats.EnableGC)
		r.Register("state.MemStats.Frees", runtimeMetrics.MemStats.Frees)
		r.Register("state.MemStats.HeapAlloc", runtimeMetrics.MemStats.HeapAlloc)
		r.Register("state.MemStats.HeapIdle", runtimeMetrics.MemStats.HeapIdle)
		r.Register("state.MemStats.HeapInuse", runtimeMetrics.MemStats.HeapInuse)
		r.Register("state.MemStats.HeapObjects", runtimeMetrics.MemStats.HeapObjects)
		r.Register("state.MemStats.HeapReleased", runtimeMetrics.MemStats.HeapReleased)
		r.Register("state.MemStats.HeapSys", runtimeMetrics.MemStats.HeapSys)
		r.Register("state.MemStats.LastGC", runtimeMetrics.MemStats.LastGC)
		r.Register("state.MemStats.Lookups", runtimeMetrics.MemStats.Lookups)
		r.Register("state.MemStats.Mallocs", runtimeMetrics.MemStats.Mallocs)
		r.Register("state.MemStats.MCacheInuse", runtimeMetrics.MemStats.MCacheInuse)
		r.Register("state.MemStats.MCacheSys", runtimeMetrics.MemStats.MCacheSys)
		r.Register("state.MemStats.MSpanInuse", runtimeMetrics.MemStats.MSpanInuse)
		r.Register("state.MemStats.MSpanSys", runtimeMetrics.MemStats.MSpanSys)
		r.Register("state.MemStats.NextGC", runtimeMetrics.MemStats.NextGC)
		r.Register("state.MemStats.NumGC", runtimeMetrics.MemStats.NumGC)
		r.Register("state.MemStats.GCCPUFraction", runtimeMetrics.MemStats.GCCPUFraction)
		r.Register("state.MemStats.PauseNs", runtimeMetrics.MemStats.PauseNs)
		r.Register("state.MemStats.PauseTotalNs", runtimeMetrics.MemStats.PauseTotalNs)
		r.Register("state.MemStats.StackInuse", runtimeMetrics.MemStats.StackInuse)
		r.Register("state.MemStats.StackSys", runtimeMetrics.MemStats.StackSys)
		r.Register("state.MemStats.Sys", runtimeMetrics.MemStats.Sys)
		r.Register("state.MemStats.TotalAlloc", runtimeMetrics.MemStats.TotalAlloc)
		r.Register("state.NumCgoCall", runtimeMetrics.NumCgoCall)
		r.Register("state.NumGoroutine", runtimeMetrics.NumGoroutine)
		r.Register("state.NumThread", runtimeMetrics.NumThread)
		r.Register("state.ReadMemStats", runtimeMetrics.ReadMemStats)
	})
}
