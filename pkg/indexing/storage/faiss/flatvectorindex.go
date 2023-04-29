package faiss

import (
	"runtime"

	"github.com/DataIntelligenceCrew/go-faiss"
)

type FlatVectorIndex struct {
	idx faiss.Index
}

func NewFlatVectorIndex(d int, metric IndexMetric) *FlatVectorIndex {
	var m int

	switch metric {
	case IndexMetricL1:
		m = faiss.MetricL1
	case IndexMetricL2:
		m = faiss.MetricL2
	case IndexMetricLInf:
		m = faiss.MetricLinf
	case IndexMetricLp:
		m = faiss.MetricLp
	case IndexMetricInnerProduct:
		m = faiss.MetricInnerProduct
	case IndexMetricCanberra:
		m = faiss.MetricCanberra
	case IndexMetricBrayCurtis:
		m = faiss.MetricBrayCurtis
	case IndexMetricJensenshannon:
		m = faiss.MetricJensenShannon
	default:
		panic("unsupported metric")
	}

	idx, err := faiss.NewIndexFlat(d, m)

	if err != nil {
		panic(err)
	}

	ti := &FlatVectorIndex{
		idx: idx,
	}

	runtime.SetFinalizer(ti, func(ti *FlatVectorIndex) {
		_ = ti.Close()
	})

	return ti
}

func (ti *FlatVectorIndex) Count() int64 {
	return ti.idx.Ntotal()
}

func (ti *FlatVectorIndex) Dim() int {
	return ti.idx.D()
}

func (ti *FlatVectorIndex) Add(x []float32) error {
	return ti.idx.Add(x)
}

func (ti *FlatVectorIndex) Search(x []float32, k int64) (distances []float32, labels []int64, err error) {
	return ti.idx.Search(x, k)
}

func (ti *FlatVectorIndex) SearchRange(x []float32, radius float32) (r SearchRangeResult, err error) {
	result, err := ti.idx.RangeSearch(x, radius)

	if err != nil {
		return SearchRangeResult{}, err
	}

	defer result.Delete()

	lims := result.Lims()
	ids, distances := result.Labels()

	r.Distances = make([][]float32, len(lims)-1)
	r.Labels = make([]int64, len(lims)-1)

	for i := 0; i < len(r.Distances); i++ {
		r.Distances[i] = distances[lims[i]:lims[i+1]]
		r.Labels[i] = ids[i]
	}

	return SearchRangeResult{}, nil
}

func (ti *FlatVectorIndex) Close() error {
	ti.idx.Delete()

	return nil
}
