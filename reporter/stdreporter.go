package reporter

import (
	"context"
	"github.com/someview/go-metrics/counter"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"log/slog"
	"time"
)

type NamedMetric struct {
	name string
	m    interface{}
}

func (n *NamedMetric) Value() any {
	return n.m
}

func (n *NamedMetric) Name() string {
	return n.name
}

type stdReporter struct {
	r Registry
}

func (s *stdReporter) Metrics() []NamedMetric {
	return s.r.GetAll()
}

func (s *stdReporter) UpdateHistogram(name string, v int64) {
	s.r.Get(name).(histogram.Histogram).Update(v)
}

func (s *stdReporter) IncGauge(name string, v int64) {
	s.r.Get(name).(guage.Gauge).Inc(v)
}

func (s *stdReporter) IncCounter(name string, v int64) {
	s.r.Get(name).(counter.Counter).Inc(v)
}

func (s *stdReporter) ReportPeriodically(ctx context.Context, interval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			for _, metricVal := range s.Metrics() {
				name := metricVal.Name()
				metric := metricVal.Value()
				switch instance := metric.(type) {
				case counter.Counter:
					slog.Info("counter", slog.String("name", name), slog.Int64("val", instance.Snapshot()))
				case guage.Gauge:
					slog.Info("gauge", slog.String("name", name), slog.Int64("val", instance.SnapShotAndReset()))
				case guage.GaugeFloat64:
					slog.Info("gauge", slog.String("name", name), slog.Float64("val", instance.SnapshotAndReset()))
				case histogram.Histogram:
					h := instance.Sample().SnapshotAndReset()
					ps := h.Percentiles([]float64{0.5, 0.95, 0.99, 0.999})
					slog.Info(
						"histogram",
						slog.String("name", name),
						slog.Int64("count", h.ReqCount()),
						slog.Int64("sample", h.Count()),
						slog.Int64("min", h.Min()),
						slog.Int64("max", h.Max()),
						slog.Float64("mean", h.Mean()),
						slog.Float64("stddev", h.StdDev()),
						slog.Float64("50%", ps[0]),
						slog.Float64("95%", ps[1]),
						slog.Float64("99%", ps[2]),
						slog.Float64("99.9%", ps[3]),
					)
				}
			}
		}
	}
}

func NewStdReporter(r Registry) Reporter {
	return &stdReporter{
		r: r,
	}
}
