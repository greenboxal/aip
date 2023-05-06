package documents

import (
	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type DocumentID struct {
	forddb2.StringResourceID[Document]
}

type Document interface {
	forddb2.Resource[DocumentID]

	getNode() *documentNode

	AsText() string
}

type DocumentBase[T Document] struct {
	forddb2.ResourceBase[DocumentID, Document]

	node *documentNode
}

func (doc *DocumentBase[T]) getNode() *documentNode {
	if doc.node == nil {
		doc.node = &documentNode{}
	}

	return doc.node
}
