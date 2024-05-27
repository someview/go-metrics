package histogram

import (
	"github.com/someview/go-metrics/sample"
	"testing"
)

func BenchmarkHistogram(b *testing.B) {
	h := NewHistogram(sample.NewSlidingWindowSample(100))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Update(int64(i))
	}
}
