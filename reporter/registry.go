package reporter

import (
	"fmt"
	"github.com/someview/go-metrics/counter"
	"github.com/someview/go-metrics/guage"
	"github.com/someview/go-metrics/histogram"
	"reflect"
	"strings"
	"sync"
)

// DuplicateMetric is the error returned by Registry.Register when a metric
// already exists.  If you mean to Register that metric you must first
// Unregister the existing metric.
type DuplicateMetric string

func (err DuplicateMetric) Error() string {
	return fmt.Sprintf("duplicate metric: %s", string(err))
}

// A Registry holds references to a set of metrics by name and can iterate
// over them, calling callback functions provided by the user.
//
// This is an interface so as to encourage other structs to implement
// the Registry API as appropriate.
type Registry interface {

	// Call the given function for each registered metric.
	Each(func(string, interface{}))

	// Get the metric by the given name or nil if none is registered.
	Get(string) interface{}

	// GetAll metrics in the Registry.
	GetAll() map[string]map[string]interface{}

	// Gets an existing metric or registers the given one.
	// The interface can be the metric to register if not found in registry,
	// or a function returning the metric for lazy instantiation.
	GetOrRegister(string, interface{}) interface{}

	// Register the given metric under the given name.
	Register(string, interface{}) error

	// Unregister the metric with the given name.
	Unregister(string)

	// Unregister all metrics.  (Mostly for testing.)
	UnregisterAll()
}

// The standard implementation of a Registry is a mutex-protected map
// of names to metrics.
type StandardRegistry struct {
	metrics map[string]interface{}
	mutex   sync.RWMutex
}

// Create a new registry.
func NewRegistry() Registry {
	return &StandardRegistry{metrics: make(map[string]interface{})}
}

// Call the given function for each registered metric.
func (r *StandardRegistry) Each(f func(string, interface{})) {
	metrics := r.registered()
	for i := range metrics {
		kv := &metrics[i]
		f(kv.name, kv.value)
	}
}

// Get the metric by the given name or nil if none is registered.
func (r *StandardRegistry) Get(name string) interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.metrics[name]
}

// Gets an existing metric or creates and registers a new one. Threadsafe
// alternative to calling Get and Register on failure.
// The interface can be the metric to register if not found in registry,
// or a function returning the metric for lazy instantiation.
func (r *StandardRegistry) GetOrRegister(name string, i interface{}) interface{} {
	// access the read lock first which should be re-entrant
	r.mutex.RLock()
	metric, ok := r.metrics[name]
	r.mutex.RUnlock()
	if ok {
		return metric
	}

	// only take the write lock if we'll be modifying the metrics map
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if metric, ok := r.metrics[name]; ok {
		return metric
	}
	if v := reflect.ValueOf(i); v.Kind() == reflect.Func {
		i = v.Call(nil)[0].Interface()
	}
	r.register(name, i)
	return i
}

// Register the given metric under the given name.  Returns a DuplicateMetric
// if a metric by the given name is already registered.
func (r *StandardRegistry) Register(name string, i interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.register(name, i)
}

// GetAll metrics in the Registry
func (r *StandardRegistry) GetAll() map[string]map[string]interface{} {
	data := make(map[string]map[string]interface{})
	r.Each(func(name string, i interface{}) {
		values := make(map[string]interface{})
		switch metric := i.(type) {
		case counter.Counter:
			values["count"] = metric.Snapshot()
		case guage.Gauge:
			values["value"] = metric.SnapShotAndReset()
		case guage.GaugeFloat64:
			values["value"] = metric.SnapshotAndReset()
		case histogram.Histogram:
			h := metric.Sample().SnapshotAndReset()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] = h.Count()
			values["min"] = h.Min()
			values["max"] = h.Max()
			values["mean"] = h.Mean()
			values["stddev"] = h.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
		}
		data[name] = values
	})
	return data
}

// Unregister the metric with the given name.
func (r *StandardRegistry) Unregister(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.metrics, name)
}

// Unregister all metrics.  (Mostly for testing.)
func (r *StandardRegistry) UnregisterAll() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for name, _ := range r.metrics {
		delete(r.metrics, name)
	}
}

func (r *StandardRegistry) register(name string, i interface{}) error {
	if _, ok := r.metrics[name]; ok {
		return DuplicateMetric(name)
	}
	switch i.(type) {
	case counter.Counter, guage.Gauge, guage.GaugeFloat64, histogram.Histogram:
		r.metrics[name] = i
	}
	return nil
}

type metricKV struct {
	name  string
	value interface{}
}

func (r *StandardRegistry) registered() []metricKV {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	metrics := make([]metricKV, 0, len(r.metrics))
	for name, i := range r.metrics {
		metrics = append(metrics, metricKV{
			name:  name,
			value: i,
		})
	}
	return metrics
}

// Stoppable defines the metrics which has to be stopped.
type Stoppable interface {
	Stop()
}

type PrefixedRegistry struct {
	underlying Registry
	prefix     string
}

func NewPrefixedRegistry(prefix string) Registry {
	return &PrefixedRegistry{
		underlying: NewRegistry(),
		prefix:     prefix,
	}
}

func NewPrefixedChildRegistry(parent Registry, prefix string) Registry {
	return &PrefixedRegistry{
		underlying: parent,
		prefix:     prefix,
	}
}

// Call the given function for each registered metric.
func (r *PrefixedRegistry) Each(fn func(string, interface{})) {
	wrappedFn := func(prefix string) func(string, interface{}) {
		return func(name string, iface interface{}) {
			if strings.HasPrefix(name, prefix) {
				fn(name, iface)
			} else {
				return
			}
		}
	}

	baseRegistry, prefix := findPrefix(r, "")
	baseRegistry.Each(wrappedFn(prefix))
}

func findPrefix(registry Registry, prefix string) (Registry, string) {
	switch r := registry.(type) {
	case *PrefixedRegistry:
		return findPrefix(r.underlying, r.prefix+prefix)
	case *StandardRegistry:
		return r, prefix
	}
	return nil, ""
}

// Get the metric by the given name or nil if none is registered.
func (r *PrefixedRegistry) Get(name string) interface{} {
	realName := r.prefix + name
	return r.underlying.Get(realName)
}

// Gets an existing metric or registers the given one.
// The interface can be the metric to register if not found in registry,
// or a function returning the metric for lazy instantiation.
func (r *PrefixedRegistry) GetOrRegister(name string, metric interface{}) interface{} {
	realName := r.prefix + name
	return r.underlying.GetOrRegister(realName, metric)
}

// Register the given metric under the given name. The name will be prefixed.
func (r *PrefixedRegistry) Register(name string, metric interface{}) error {
	realName := r.prefix + name
	return r.underlying.Register(realName, metric)
}

// GetAll metrics in the Registry
func (r *PrefixedRegistry) GetAll() map[string]map[string]interface{} {
	return r.underlying.GetAll()
}

// Unregister the metric with the given name. The name will be prefixed.
func (r *PrefixedRegistry) Unregister(name string) {
	realName := r.prefix + name
	r.underlying.Unregister(realName)
}

// Unregister all metrics.  (Mostly for testing.)
func (r *PrefixedRegistry) UnregisterAll() {
	r.underlying.UnregisterAll()
}

var DefaultRegistry Registry = NewRegistry()

// Call the given function for each registered metric.
func Each(f func(string, interface{})) {
	DefaultRegistry.Each(f)
}

// Get the metric by the given name or nil if none is registered.
func Get(name string) interface{} {
	return DefaultRegistry.Get(name)
}

// Gets an existing metric or creates and registers a new one. Threadsafe
// alternative to calling Get and Register on failure.
func GetOrRegister(name string, i interface{}) interface{} {
	return DefaultRegistry.GetOrRegister(name, i)
}

// Register the given metric under the given name.  Returns a DuplicateMetric
// if a metric by the given name is already registered.
func Register(name string, i interface{}) error {
	return DefaultRegistry.Register(name, i)
}

// Register the given metric under the given name.  Panics if a metric by the
// given name is already registered.
func MustRegister(name string, i interface{}) {
	if err := Register(name, i); err != nil {
		panic(err)
	}
}

// Unregister the metric with the given name.
func Unregister(name string) {
	DefaultRegistry.Unregister(name)
}
