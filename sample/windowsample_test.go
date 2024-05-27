package sample

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlidingWindowSample_Snapshot(t *testing.T) {
	s := NewSlidingWindowSample(1)
	s.Update(1)
	s.Update(2)
	snapshot := s.SnapshotAndReset()
	assert.Equal(t, int64(2), snapshot.ReqCount())
	assert.Equal(t, int64(1), snapshot.Count())
}
