package histogram

import (
	. "github.com/someview/go-metrics/sample"
)

// StandardHistogram is the standard implementation of a Histogram and uses a
// Sample to bound its memory use.
type StandardHistogram struct {
	sample Sample
}

func NewHistogram(s Sample) Histogram {
	return &StandardHistogram{
		sample: s,
	}
}

func (h *StandardHistogram) Update(i int64) {
	h.sample.Update(i)
}

// Clear clears the histogram and its sample.
func (h *StandardHistogram) Clear() { h.sample.Clear() }

// Sample returns the Sample underlying the histogram.
func (h *StandardHistogram) Sample() Sample { return h.sample }
