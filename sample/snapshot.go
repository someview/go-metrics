package sample

// sampleSnapshot is a read-only copy of another Sample.
type sampleSnapshot struct {
	reqCount int64
	count    int64
	values   []int64
}

func (s *sampleSnapshot) ReqCount() int64 {
	return s.reqCount
}

func NewSampleSnapshot(reqCount, count int64, values []int64) SampleSnapshot {
	return &sampleSnapshot{
		reqCount: reqCount,
		count:    count,
		values:   values,
	}
}

// Count returns the count of inputs at the time the snapshot was taken.
func (s *sampleSnapshot) Count() int64 { return s.count }

// Max returns the maximal value at the time the snapshot was taken.
func (s *sampleSnapshot) Max() int64 { return SampleMax(s.values) }

// Mean returns the mean value at the time the snapshot was taken.
func (s *sampleSnapshot) Mean() float64 { return SampleMean(s.values) }

// Min returns the minimal value at the time the snapshot was taken.
func (s *sampleSnapshot) Min() int64 { return SampleMin(s.values) }

// Percentile returns an arbitrary percentile of values at the time the
// snapshot was taken.
func (s *sampleSnapshot) Percentile(p float64) float64 {
	return SamplePercentile(s.values, p)
}

// Percentiles returns a slice of arbitrary percentiles of values at the time
// the snapshot was taken.
func (s *sampleSnapshot) Percentiles(ps []float64) []float64 {
	return SamplePercentiles(s.values, ps)
}

// Size returns the size of the sample at the time the snapshot was taken.
func (s *sampleSnapshot) Size() int { return len(s.values) }

// StdDev returns the standard deviation of values at the time the snapshot was
// taken.
func (s *sampleSnapshot) StdDev() float64 { return SampleStdDev(s.values) }

// Sum returns the sum of values at the time the snapshot was taken.
func (s *sampleSnapshot) Sum() int64 { return SampleSum(s.values) }

// Variance returns the variance of values at the time the snapshot was taken.
func (s *sampleSnapshot) Variance() float64 { return SampleVariance(s.values) }
