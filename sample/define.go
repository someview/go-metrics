package sample

// Samples maintain a statistically-significant selection of values from
// a stream.
type Sample interface {
	Clear()
	Snapshot() SampleSnapshot
	SnapshotAndReset() SampleSnapshot
	Update(int64)
}

type SampleSnapshot interface {
	ReqCount() int64 // 请求次数
	Count() int64    // 采样次数
	Max() int64
	Mean() float64
	Min() int64
	Percentile(float64) float64
	Percentiles([]float64) []float64
	Size() int
	Sum() int64
	StdDev() float64
	Variance() float64
}
