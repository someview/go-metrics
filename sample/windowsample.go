package sample

import (
	"sync"
)

// SlidingWindowSample is a sample that stores values in a ring buffer.
// When get snapshot, it's return snapshot and reset the sampler
type SlidingWindowSample struct {
	mutex    sync.Mutex
	values   []int64
	size     uint64
	index    uint64
	count    int64
	reqCount int64
}

func (s *SlidingWindowSample) SnapshotAndReset() SampleSnapshot {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	res := NewSampleSnapshot(s.reqCount, s.count, s.values[:s.count])
	s.reset()
	return res
}

func NewSlidingWindowSample(size uint64) Sample {
	return &SlidingWindowSample{
		size:   size,
		values: make([]int64, size),
	}
}

// Clear clears all samples.
func (s *SlidingWindowSample) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.reset()
}

// Snapshot returns a read-only copy of the sample.
func (s *SlidingWindowSample) Snapshot() SampleSnapshot {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	res := NewSampleSnapshot(s.reqCount, s.count, s.values[:s.count])
	return res
}

func (s *SlidingWindowSample) reset() {
	s.values = make([]int64, s.size)
	s.index = 0
	s.count = 0
	s.reqCount = 0
}

// Update samples a new value.
func (s *SlidingWindowSample) Update(v int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.reqCount++
	s.values[s.index] = v
	s.index = (s.index + 1) % s.size
	if s.count < int64(s.size) {
		s.count++
	}
}
