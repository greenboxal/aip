package vectorstore

type Collection interface {
	Indexer
	Retriever
}
