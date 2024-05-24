package counter

import (
	"sync/atomic"
)

// Counter hold an int64 value that can be incremented and decremented.
type Counter interface {
	Dec(int64)
	Inc(int64)
	Snapshot() int64
	SnapshotAndReset() int64
}

// StandardCounter is the standard implementation of a Counter and uses the
// sync/atomic package to manage a single int64 value.
type StandardCounter struct {
	count int64
}

// NewCounter constructs a new StandardCounter.
func NewCounter() Counter {
	return &StandardCounter{0}
}

// SnapshotAndReset 为 StandardCounter 类型的实例创建一个快照，并重置计数器。
// 参数: c *StandardCounter - 指向当前计数器实例的指针。
// 返回值: CounterSnapshot - 计数器快照，包含重置前的计数器值
func (c *StandardCounter) SnapshotAndReset() int64 {
	return atomic.SwapInt64(&c.count, 0)
}

// Dec decrements the counter by the given amount.
func (c *StandardCounter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

// Inc increments the counter by the given amount.
func (c *StandardCounter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}

func (c *StandardCounter) Swap(i int64) int64 {
	return atomic.SwapInt64(&c.count, i)
}

// Snapshot returns a read-only copy of the counter.
func (c *StandardCounter) Snapshot() int64 {
	return atomic.LoadInt64(&c.count)
}
