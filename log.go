package metrics

import (
	"github.com/someview/go-metrics/counter"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"github.com/someview/go-metrics/reporter"
	"time"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

// Log outputs each metric in the given registry periodically using the given logger.
func Log(r reporter.Registry, freq time.Duration, l Logger) {
	LogScaled(r, freq, time.Nanosecond, l)
}

// LogOnCue outputs each metric in the given registry on demand through the channel
// using the given logger
func LogOnCue(r reporter.Registry, ch chan interface{}, l Logger) {
	LogScaledOnCue(r, ch, l)
}

// LogScaled outputs each metric in the given registry periodically using the given
// logger. Print timings in `scale` units (eg time.Millisecond) rather than nanos.
func LogScaled(r reporter.Registry, freq time.Duration, scale time.Duration, l Logger) {
	ch := make(chan interface{})
	go func(channel chan interface{}) {
		for _ = range time.Tick(freq) {
			channel <- struct{}{}
		}
	}(ch)
	LogScaledOnCue(r, ch, l)
}

// LogScaledOnCue outputs each metric in the given registry on demand through the channel
// using the given logger. Print timings in `scale` units (eg time.Millisecond) rather
// than nanos.
func LogScaledOnCue(r reporter.Registry, ch chan interface{}, l Logger) {
	for _ = range ch {
		r.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case counter.Counter:
				l.Printf("counter %s\n", name)
				l.Printf("  count:       %9d\n", metric.Snapshot())
			case guage.Gauge:
				l.Printf("gauge %s\n", name)
				l.Printf("  value:       %9d\n", metric.SnapShotAndReset())
			case guage.GaugeFloat64:
				l.Printf("gauge %s\n", name)
				l.Printf("  value:       %f\n", metric.SnapshotAndReset())
			case histogram.Histogram:
				h := metric.Sample().SnapshotAndReset()
				ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				l.Printf("histogram %s\n", name)
				l.Printf("  count:       %9d\n", h.Count())
				l.Printf("  min:         %9d\n", h.Min())
				l.Printf("  max:         %9d\n", h.Max())
				l.Printf("  mean:        %12.2f\n", h.Mean())
				l.Printf("  stddev:      %12.2f\n", h.StdDev())
				l.Printf("  median:      %12.2f\n", ps[0])
				l.Printf("  75%%:         %12.2f\n", ps[1])
				l.Printf("  95%%:         %12.2f\n", ps[2])
				l.Printf("  99%%:         %12.2f\n", ps[3])
				l.Printf("  99.9%%:       %12.2f\n", ps[4])
			}
		})
	}
}
