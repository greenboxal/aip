package documents

import "github.com/greenboxal/aip/pkg/ford/forddb"

type DocumentID struct {
	forddb.StringResourceID[Document]
}

type Document interface {
	forddb.Resource[DocumentID]

	getNode() *documentNode

	AsText() string
}

type DocumentBase[T Document] struct {
	forddb.ResourceMetadata[DocumentID, Document]

	node *documentNode
}

func (doc *DocumentBase[T]) getNode() *documentNode {
	if doc.node == nil {
		doc.node = &documentNode{}
	}

	return doc.node
}
