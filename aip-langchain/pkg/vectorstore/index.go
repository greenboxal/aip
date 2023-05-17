package vectorstore

type Index interface {
	Indexer
	Retriever

	Close() error
}
