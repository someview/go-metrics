package histogram

import "github.com/someview/go-metrics/sample"

// Histograms calculate distribution statistics from a series of int64 values.
type Histogram interface {
	Clear()
	Sample() sample.Sample
	Update(int64)
}
