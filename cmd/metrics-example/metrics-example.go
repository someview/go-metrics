package main

import (
	"github.com/someview/go-metrics"
	"github.com/someview/go-metrics/counter"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"github.com/someview/go-metrics/meter"
	"github.com/someview/go-metrics/reporter"
	"github.com/someview/go-metrics/sample"
	"github.com/someview/go-metrics/state"
	"github.com/someview/go-metrics/timer"
	"log"
	"os"
	// "syslog"
	"time"
)

const fanout = 10

func main() {

	r := reporter.NewRegistry()

	c := counter.NewCounter()
	r.Register("foo", c)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				c.Dec(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				c.Inc(47)
				time.Sleep(400e6)
			}
		}()
	}

	g := guage.NewGauge()
	r.Register("bar", g)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				g.Update(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				g.Update(47)
				time.Sleep(400e6)
			}
		}()
	}

	gf := guage.NewGaugeFloat64()
	r.Register("barfloat64", gf)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				g.Update(19.0)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				g.Update(47.0)
				time.Sleep(400e6)
			}
		}()
	}

	s := sample.NewExpDecaySample(1028, 0.015)
	//s := metrics.NewUniformSample(1028)
	h := histogram.NewHistogram(s)
	r.Register("bang", h)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				h.Update(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				h.Update(47)
				time.Sleep(400e6)
			}
		}()
	}

	m := meter.NewMeter()
	r.Register("quux", m)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				m.Mark(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				m.Mark(47)
				time.Sleep(400e6)
			}
		}()
	}

	t := timer.NewTimer()
	r.Register("hooah", t)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				t.Time(func() { time.Sleep(300e6) })
			}
		}()
		go func() {
			for {
				t.Time(func() { time.Sleep(400e6) })
			}
		}()
	}

	state.RegisterDebugGCStats(r)
	go state.CaptureDebugGCStats(r, 5e9)

	state.RegisterRuntimeMemStats(r)
	go state.CaptureRuntimeMemStats(r, 5e9)

	metrics.Log(r, 60e9, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	/*
		w, err := syslog.Dial("unixgram", "/dev/log", syslog.LOG_INFO, "metrics")
		if nil != err { log.Fatalln(err) }
		metrics.Syslog(r, 60e9, w)
	*/

	/*
		addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
		metrics.Graphite(r, 10e9, "metrics", addr)
	*/

	/*
		stathat.Stathat(r, 10e9, "example@example.com")
	*/

}
