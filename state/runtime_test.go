package state

import (
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/reporter"
	"runtime"
	"testing"
	"time"
)

func TestRuntimeMemStatsDoubleRegister(t *testing.T) {
	r := reporter.NewRegistry()
	RegisterRuntimeMemStats(r)
	storedGauge := r.Get("state.MemStats.LastGC").(guage.Gauge)

	runtime.GC()
	CaptureRuntimeMemStatsOnce(r)

	firstGC := storedGauge.Value()
	if 0 == firstGC {
		t.Errorf("firstGC got %d, expected timestamp > 0", firstGC)
	}

	time.Sleep(time.Millisecond)

	RegisterRuntimeMemStats(r)
	runtime.GC()
	CaptureRuntimeMemStatsOnce(r)
	if lastGC := storedGauge.Value(); firstGC == lastGC {
		t.Errorf("lastGC got %d, expected a higher timestamp value", lastGC)
	}
}

func BenchmarkRuntimeMemStats(b *testing.B) {
	r := reporter.NewRegistry()
	RegisterRuntimeMemStats(r)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CaptureRuntimeMemStatsOnce(r)
	}
}

func TestRuntimeMemStats(t *testing.T) {
	r := reporter.NewRegistry()
	RegisterRuntimeMemStats(r)
	CaptureRuntimeMemStatsOnce(r)
	zero := runtimeMetrics.MemStats.PauseNs.Count() // Get a "zero" since GC may have state before these tests.
	runtime.GC()
	CaptureRuntimeMemStatsOnce(r)
	if count := runtimeMetrics.MemStats.PauseNs.Count(); 1 != count-zero {
		t.Fatal(count - zero)
	}
	runtime.GC()
	runtime.GC()
	CaptureRuntimeMemStatsOnce(r)
	if count := runtimeMetrics.MemStats.PauseNs.Count(); 3 != count-zero {
		t.Fatal(count - zero)
	}
	for i := 0; i < 256; i++ {
		runtime.GC()
	}
	CaptureRuntimeMemStatsOnce(r)
	if count := runtimeMetrics.MemStats.PauseNs.Count(); 259 != count-zero {
		t.Fatal(count - zero)
	}
	for i := 0; i < 257; i++ {
		runtime.GC()
	}
	CaptureRuntimeMemStatsOnce(r)
	if count := runtimeMetrics.MemStats.PauseNs.Count(); 515 != count-zero { // We lost one because there were too many GCs between captures.
		t.Fatal(count - zero)
	}
}

func TestRuntimeMemStatsNumThread(t *testing.T) {
	r := reporter.NewRegistry()
	RegisterRuntimeMemStats(r)
	CaptureRuntimeMemStatsOnce(r)

	if value := runtimeMetrics.NumThread.Value(); value < 1 {
		t.Fatalf("got NumThread: %d, wanted at least 1", value)
	}
}

func TestRuntimeMemStatsBlocking(t *testing.T) {
	if g := runtime.GOMAXPROCS(0); g < 2 {
		t.Skipf("skipping TestRuntimeMemStatsBlocking with GOMAXPROCS=%d\n", g)
	}
	ch := make(chan int)
	go testRuntimeMemStatsBlocking(ch)
	var memStats runtime.MemStats
	t0 := time.Now()
	runtime.ReadMemStats(&memStats)
	t1 := time.Now()
	t.Log("i++ during state.ReadMemStats:", <-ch)
	go testRuntimeMemStatsBlocking(ch)
	d := t1.Sub(t0)
	t.Log(d)
	time.Sleep(d)
	t.Log("i++ during time.Sleep:", <-ch)
}

func testRuntimeMemStatsBlocking(ch chan int) {
	i := 0
	for {
		select {
		case ch <- i:
			return
		default:
			i++
		}
	}
}