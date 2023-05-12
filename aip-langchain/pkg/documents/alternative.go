package documents

type BasicAlternative interface {
	CreateForDocument(document Document) (any, error)
}

type AlternativeFunc func(document Document) (any, error)

func (f AlternativeFunc) CreateForDocument(document Document) (any, error) {
	return f(document)
}

type AlternativeKey[T any] interface {
	BasicAlternative
}

func Alternative[T any](doc Document, key AlternativeKey[T]) (def T, _ error) {
	node := doc.getNode()

	if node.cachedAlternatives == nil {
		node.cachedAlternatives = make(map[BasicAlternative]any)
	}

	if alt, ok := node.cachedAlternatives[key]; ok {
		return alt.(T), nil
	}

	alt, err := key.CreateForDocument(doc)

	if err != nil {
		return def, err
	}

	node.cachedAlternatives[key] = alt

	return alt.(T), nil
}
