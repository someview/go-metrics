package sample

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlidingWindowSample_SnapshotAndReset(t *testing.T) {
	s := NewSlidingWindowSample(1)
	s.Update(1)
	s.Update(2)
	snapshot := s.SnapshotAndReset()
	assert.Equal(t, int64(2), snapshot.ReqCount())
	assert.Equal(t, int64(1), snapshot.Count())
	s.Update(1)
	s.Update(2)
	assert.Equal(t, int64(2), snapshot.ReqCount())
	assert.Equal(t, int64(1), snapshot.Count())
}

func TestSlidingWindowSample_Snapshot(t *testing.T) {
	s := NewSlidingWindowSample(1)
	s.Update(1)
	s.Update(2)
	snapshot := s.Snapshot()
	assert.Equal(t, int64(2), snapshot.ReqCount())
	assert.Equal(t, int64(1), snapshot.Count())
	s.Update(1)
	snapshot = s.Snapshot()
	assert.Equal(t, int64(3), snapshot.ReqCount())
	assert.Equal(t, int64(1), snapshot.Count())
}
