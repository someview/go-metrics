package main

import (
	"fmt"
	"github.com/someview/go-metrics/counter"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"github.com/someview/go-metrics/reporter"
	"github.com/someview/go-metrics/sample"
	"time"
)

func main() {
	r := reporter.NewRegistry()
	for i := 0; i < 10000; i++ {
		r.Register(fmt.Sprintf("counter-%d", i), counter.NewCounter())
		r.Register(fmt.Sprintf("gauge-%d", i), guage.NewGauge())
		r.Register(fmt.Sprintf("gaugefloat64-%d", i), guage.NewGaugeFloat64())
		r.Register(fmt.Sprintf("histogram-uniform-%d", i), histogram.NewHistogram(sample.NewSlidingWindowSample(1028)))
		r.Register(fmt.Sprintf("histogram-exp-%d", i), histogram.NewHistogram(sample.NewExpDecaySample(1028, 0.015)))
	}
	time.Sleep(600e9)
}
