package reporter

import (
	"context"
	"github.com/someview/go-metrics/histogram"
	"github.com/someview/go-metrics/sample"
	"testing"
	"time"
)

func TestStdReporter_ReportPeriodically(t *testing.T) {
	metrics := []NamedMetric{
		NewHistogramMetric("disk", histogram.NewHistogram(sample.NewSlidingWindowSample(1))),
	}
	r := NewStdReporter(metrics)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r.UpdateHistogram("disk", 1)
	r.ReportPeriodically(ctx, 1)
}
