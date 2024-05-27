package reporter

import (
	"context"
	"github.com/someview/go-metrics/counter"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"github.com/someview/go-metrics/sample"
	"testing"
	"time"
)

func TestStdReporter_ReportPeriodically(t *testing.T) {
	r := NewStdReporter(DefaultRegistry)
	DefaultRegistry.Register("cpu", guage.NewGauge())
	DefaultRegistry.Register("mem", counter.NewCounter())
	DefaultRegistry.Register("disk", histogram.NewHistogram(sample.NewSlidingWindowSample(1)))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r.IncGauge("cpu", 1)
	r.IncCounter("mem", 1)
	r.UpdateHistogram("disk", 1)
	r.ReportPeriodically(ctx, 1)
}
