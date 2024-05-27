//go:build !windows
// +build !windows

package reporter

import (
	"fmt"
	"github.com/someview/go-metrics/counter"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"log/syslog"
	"time"
)

// Output each metric in the given registry to syslog periodically using
// the given syslogger.
func Syslog(r Registry, d time.Duration, w *syslog.Writer) {
	for _ = range time.Tick(d) {
		r.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case counter.Counter:
				w.Info(fmt.Sprintf("counter %s: count: %d", name, metric.Count()))
			case guage.Gauge:
				w.Info(fmt.Sprintf("gauge %s: value: %d", name, metric.Value()))
			case guage.GaugeFloat64:
				w.Info(fmt.Sprintf("gauge %s: value: %f", name, metric.Value()))
			case histogram.Histogram:
				h := metric.Snapshot()
				ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				w.Info(fmt.Sprintf(
					"histogram %s: count: %d min: %d max: %d mean: %.2f stddev: %.2f median: %.2f 75%%: %.2f 95%%: %.2f 99%%: %.2f 99.9%%: %.2f",
					name,
					h.Count(),
					h.Min(),
					h.Max(),
					h.Mean(),
					h.StdDev(),
					ps[0],
					ps[1],
					ps[2],
					ps[3],
					ps[4],
				))
			}
		})
	}
}
