package reporter

import (
	"context"
	"time"
)

type Reporter interface {
	RegisterMetrics([]NamedMetric)
	Metrics() []NamedMetric
	UpdateHistogram(name string, v int64)
	IncGauge(name string, v int64)
	IncCounter(name string, v int64)
	ReportPeriodically(ctx context.Context, interval time.Duration)
}
