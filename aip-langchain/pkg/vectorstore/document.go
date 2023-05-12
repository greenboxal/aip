package vectorstore

import "github.com/greenboxal/aip/aip-langchain/pkg/llm"

type Document struct {
	DocumentReference

	Content string
}

type DocumentReference struct {
	ID   string
	Type string
}

type DocumentChunkReference struct {
	ID    string
	Type  string
	Chunk int
}

type DocumentChunk struct {
	DocumentChunkReference

	Content    string
	Embeddings llm.Embeddings
}

type IndexedDocument struct {
	DocumentReference

	Chunks []DocumentChunkReference
}
