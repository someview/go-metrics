package metrics

import (
	"github.com/someview/go-metrics/counter"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"github.com/someview/go-metrics/meter"
	"github.com/someview/go-metrics/reporter"
	"github.com/someview/go-metrics/timer"
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
	LogScaledOnCue(r, ch, time.Nanosecond, l)
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
	LogScaledOnCue(r, ch, scale, l)
}

// LogScaledOnCue outputs each metric in the given registry on demand through the channel
// using the given logger. Print timings in `scale` units (eg time.Millisecond) rather
// than nanos.
func LogScaledOnCue(r reporter.Registry, ch chan interface{}, scale time.Duration, l Logger) {
	du := float64(scale)
	duSuffix := scale.String()[1:]

	for _ = range ch {
		r.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case counter.Counter:
				l.Printf("counter %s\n", name)
				l.Printf("  count:       %9d\n", metric.Count())
			case guage.Gauge:
				l.Printf("gauge %s\n", name)
				l.Printf("  value:       %9d\n", metric.Value())
			case guage.GaugeFloat64:
				l.Printf("gauge %s\n", name)
				l.Printf("  value:       %f\n", metric.Value())
			case histogram.Histogram:
				h := metric.Snapshot()
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
			case meter.Meter:
				m := metric.Snapshot()
				l.Printf("meter %s\n", name)
				l.Printf("  count:       %9d\n", m.Count())
				l.Printf("  1-min rate:  %12.2f\n", m.Rate1())
				l.Printf("  5-min rate:  %12.2f\n", m.Rate5())
				l.Printf("  15-min rate: %12.2f\n", m.Rate15())
				l.Printf("  mean rate:   %12.2f\n", m.RateMean())
			case timer.Timer:
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				l.Printf("timer %s\n", name)
				l.Printf("  count:       %9d\n", t.Count())
				l.Printf("  min:         %12.2f%s\n", float64(t.Min())/du, duSuffix)
				l.Printf("  max:         %12.2f%s\n", float64(t.Max())/du, duSuffix)
				l.Printf("  mean:        %12.2f%s\n", t.Mean()/du, duSuffix)
				l.Printf("  stddev:      %12.2f%s\n", t.StdDev()/du, duSuffix)
				l.Printf("  median:      %12.2f%s\n", ps[0]/du, duSuffix)
				l.Printf("  75%%:         %12.2f%s\n", ps[1]/du, duSuffix)
				l.Printf("  95%%:         %12.2f%s\n", ps[2]/du, duSuffix)
				l.Printf("  99%%:         %12.2f%s\n", ps[3]/du, duSuffix)
				l.Printf("  99.9%%:       %12.2f%s\n", ps[4]/du, duSuffix)
				l.Printf("  1-min rate:  %12.2f\n", t.Rate1())
				l.Printf("  5-min rate:  %12.2f\n", t.Rate5())
				l.Printf("  15-min rate: %12.2f\n", t.Rate15())
				l.Printf("  mean rate:   %12.2f\n", t.RateMean())
			}
		})
	}
}
