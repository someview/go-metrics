package storage

import (
	"encoding/json"
	"github.com/someview/go-metrics"
	"io"
	"time"
)

// MarshalJSON returns a byte slice containing a JSON representation of all
// the metrics in the Registry.
func (r *metrics.StandardRegistry) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.GetAll())
}

// WriteJSON writes metrics from the given registry  periodically to the
// specified io.Writer as JSON.
func WriteJSON(r metrics.Registry, d time.Duration, w io.Writer) {
	for _ = range time.Tick(d) {
		WriteJSONOnce(r, w)
	}
}

// WriteJSONOnce writes metrics from the given registry to the specified
// io.Writer as JSON.
func WriteJSONOnce(r metrics.Registry, w io.Writer) {
	json.NewEncoder(w).Encode(r)
}

func (p *metrics.PrefixedRegistry) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.GetAll())
}
