package guage

import (
	"testing"
)

func BenchmarkGuageFloat64(b *testing.B) {
	g := NewGaugeFloat64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Update(float64(i))
	}
}

func BenchmarkGuageFloat64Parallel(b *testing.B) {
	g := NewGaugeFloat64()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			g.Update(float64(1))
		}
	})
}
