package sample

import (
	"runtime"
	"testing"
	"time"
)

// Benchmark{Compute,Copy}{1000,1000000} demonstrate that, even for relatively
// expensive computations like Variance, the cost of copying the Sample, as
// approximated by a make and copy, is much greater than the cost of the
// computation for small samples and only slightly less for large samples.
func BenchmarkCompute1000(b *testing.B) {
	s := make([]int64, 1000)
	for i := 0; i < len(s); i++ {
		s[i] = int64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SampleVariance(s)
	}
}
func BenchmarkCompute1000000(b *testing.B) {
	s := make([]int64, 1000000)
	for i := 0; i < len(s); i++ {
		s[i] = int64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SampleVariance(s)
	}
}
func BenchmarkCopy1000(b *testing.B) {
	s := make([]int64, 1000)
	for i := 0; i < len(s); i++ {
		s[i] = int64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sCopy := make([]int64, len(s))
		copy(sCopy, s)
	}
}
func BenchmarkCopy1000000(b *testing.B) {
	s := make([]int64, 1000000)
	for i := 0; i < len(s); i++ {
		s[i] = int64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sCopy := make([]int64, len(s))
		copy(sCopy, s)
	}
}

func BenchmarkExpDecaySample257(b *testing.B) {
	benchmarkSample(b, NewExpDecaySample(257, 0.015))
}

func BenchmarkExpDecaySample514(b *testing.B) {
	benchmarkSample(b, NewExpDecaySample(514, 0.015))
}

func BenchmarkExpDecaySample1028(b *testing.B) {
	benchmarkSample(b, NewExpDecaySample(1028, 0.015))
}

// This test makes sure that the sample's priority is not amplified by using
// nanosecond duration since start rather than second duration since start.
// The priority becomes +Inf quickly after starting if this is done,
// effectively freezing the set of samples until a rescale step happens.

func TestExpDecaySampleRescale(t *testing.T) {
	s := NewExpDecaySample(2, 0.001).(*ExpDecaySample)
	s.update(time.Now(), 1)
	s.update(time.Now().Add(time.Hour+time.Microsecond), 1)
	for _, v := range s.values.Values() {
		if v.k == 0.0 {
			t.Fatal("v.k == 0.0")
		}
	}
}

func benchmarkSample(b *testing.B, s Sample) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	pauseTotalNs := memStats.PauseTotalNs
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Update(1)
	}
	b.StopTimer()
	runtime.GC()
	runtime.ReadMemStats(&memStats)
	b.Logf("GC cost: %d ns/op", int(memStats.PauseTotalNs-pauseTotalNs)/b.N)
}
