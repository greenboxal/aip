package faiss

import "errors"

type InverseVectorIndex struct {
	idx    VectorIndex
	labels []int64
}

var ErrInvalidLength = errors.New("invalid length")

func (idx *InverseVectorIndex) Add(labels []int64, x []float32) error {
	if len(labels) != len(x) {
		return ErrInvalidLength
	}
	if err := idx.idx.Add(x); err != nil {
		return err
	}

	idx.labels = append(idx.labels, labels...)

	return nil
}

func (idx *InverseVectorIndex) Search(x []float32, k int64) (distances []float32, labels []int64, err error) {
	distances, labels, err = idx.idx.Search(x, k)

	if err != nil {
		return nil, nil, err
	}

	for i, l := range labels {
		labels[i] = idx.labels[l]
	}

	return distances, labels, nil
}

func (idx *InverseVectorIndex) SearchRange(x []float32, radius float32) (r SearchRangeResult, err error) {
	result, err := idx.idx.SearchRange(x, radius)

	if err != nil {
		return SearchRangeResult{}, nil
	}

	for i, l := range result.Labels {
		result.Labels[i] = idx.labels[l]
	}

	return result, nil
}

func (idx *InverseVectorIndex) Dim() int {
	return idx.idx.Dim()
}

func (idx *InverseVectorIndex) Count() int64 {
	return idx.idx.Count()
}

func NewInverseVectorIndex(idx VectorIndex) *InverseVectorIndex {
	return &InverseVectorIndex{idx: idx}
}
