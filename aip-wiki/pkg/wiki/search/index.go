package search

type Request struct {
}

type Hit struct {
}

type Result struct {
	Hits []Hit
}

type Index interface {
	Search(req Request) (*Result, error)
}
