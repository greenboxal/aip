package storage

type IndexMetric int

const (
	IndexMetricNone IndexMetric = iota
	IndexMetricL1
	IndexMetricL2
	IndexMetricLInf
	IndexMetricLp
	IndexMetricCanberra
	IndexMetricBrayCurtis
	IndexMetricJensenshannon
	IndexMetricInnerProduct
)

type SearchRangeResult struct {
	Distances [][]float32
	Labels    []int64
}

type MapSearchRangeResult[K comparable] struct {
	Distances [][]float32
	Labels    []K
}

type VectorIndex interface {
	Add(x []float32) error
	Search(x []float32, k int64) (distances []float32, labels []int64, err error)
	SearchRange(x []float32, radius float32) (r SearchRangeResult, err error)

	Dim() int
	Count() int64
}
