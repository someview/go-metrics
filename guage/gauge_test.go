package guage

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func BenchmarkGuage(b *testing.B) {
	g := NewGauge()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Update(int64(i))
	}
}

// exercise race detector
func TestGaugeConcurrency(t *testing.T) {
	rand.Seed(time.Now().Unix())
	g := NewGauge()
	wg := &sync.WaitGroup{}
	reps := 100
	for i := 0; i < reps; i++ {
		wg.Add(1)
		go func(g Gauge, wg *sync.WaitGroup) {
			g.Update(rand.Int63())
			wg.Done()
		}(g, wg)
	}
	wg.Wait()
}
