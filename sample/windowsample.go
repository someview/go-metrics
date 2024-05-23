package sample

import (
	"sync"
)

type SlidingWindowSample struct {
	mutex  sync.Mutex
	values []int64
	size   int
	index  int
	count  int64
}

func NewSlidingWindowSample(size int) Sample {
	return &SlidingWindowSample{
		size:   size,
		values: make([]int64, size),
	}
}

// Clear clears all samples.
func (s *SlidingWindowSample) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.count = 0
	s.values = make([]int64, 0, s.size)
}

// Count returns the number of samples recorded, which may exceed the
// reservoir size.
func (s *SlidingWindowSample) Count() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.count
}

// Max returns the maximum value in the sample, which may not be the maximum
// value ever to be part of the sample.
func (s *SlidingWindowSample) Max() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleMax(s.values)
}

// Mean returns the mean of the values in the sample.
func (s *SlidingWindowSample) Mean() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleMean(s.values)
}

// Min returns the minimum value in the sample, which may not be the minimum
// value ever to be part of the sample.
func (s *SlidingWindowSample) Min() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleMin(s.values)
}

// Percentile returns an arbitrary percentile of values in the sample.
func (s *SlidingWindowSample) Percentile(p float64) float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SamplePercentile(s.values, p)
}

// Percentiles returns a slice of arbitrary percentiles of values in the
// sample.
func (s *SlidingWindowSample) Percentiles(ps []float64) []float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SamplePercentiles(s.values, ps)
}

// Size returns the size of the sample, which is at most the reservoir size.
func (s *SlidingWindowSample) Size() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return len(s.values)
}

// Snapshot returns a read-only copy of the sample.
func (s *SlidingWindowSample) Snapshot() Sample {
	s.mutex.Lock()
	values := make([]int64, len(s.values))
	copy(values, s.values)
	res := NewSampleSnapshot(s.count, values)
	s.mutex.Unlock()
	return res
}

// StdDev returns the standard deviation of the values in the sample.
func (s *SlidingWindowSample) StdDev() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleStdDev(s.values)
}

// Sum returns the sum of the values in the sample.
func (s *SlidingWindowSample) Sum() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleSum(s.values)
}

// Update samples a new value.
func (s *SlidingWindowSample) Update(v int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.values[s.index] = v
	s.index = (s.index + 1) % s.size
	if s.count < int64(s.size) {
		s.count++
	}
}

// Values returns a copy of the values in the sample.
func (s *SlidingWindowSample) Values() []int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	values := make([]int64, len(s.values))
	copy(values, s.values)
	return values
}

// Variance returns the variance of the values in the sample.
func (s *SlidingWindowSample) Variance() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return SampleVariance(s.values)
}
