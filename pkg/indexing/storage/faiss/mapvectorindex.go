package faiss

type MapVectorIndex[K comparable] struct {
	idx    VectorIndex
	labels []K
}

func (idx *MapVectorIndex[K]) Add(labels []K, x []float32) error {
	if len(labels) != len(x) {
		return ErrInvalidLength
	}

	if err := idx.idx.Add(x); err != nil {
		return err
	}

	idx.labels = append(idx.labels, labels...)

	return nil
}

func (idx *MapVectorIndex[K]) Search(x []float32, k int64) (distances []float32, labels []K, err error) {
	distances, idLabels, err := idx.idx.Search(x, k)

	if err != nil {
		return nil, nil, err
	}

	labels = make([]K, len(idLabels))

	for i, l := range idLabels {
		labels[i] = idx.labels[l]
	}

	return distances, labels, nil
}

func (idx *MapVectorIndex[K]) SearchRange(x []float32, radius float32) (r MapSearchRangeResult[K], err error) {
	result, err := idx.idx.SearchRange(x, radius)

	if err != nil {
		return MapSearchRangeResult[K]{}, nil
	}

	r.Distances = result.Distances

	for i, l := range result.Labels {
		r.Labels[i] = idx.labels[l]
	}

	return r, nil
}

func (idx *MapVectorIndex[K]) Dim() int {
	return idx.idx.Dim()
}

func (idx *MapVectorIndex[K]) Count() int64 {
	return idx.idx.Count()
}

func NewMapVectorIndex[K comparable](idx VectorIndex) *MapVectorIndex[K] {
	return &MapVectorIndex[K]{idx: idx}
}
