package documents

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type DocumentID struct {
	forddb.StringResourceID[Document] `ipld:",inline"`
}

type Document interface {
	forddb.Resource[DocumentID]

	getNode() *documentNode

	AsText() string
}

type DocumentBase[T Document] struct {
	forddb.ResourceBase[DocumentID, Document]

	node *documentNode
}

func (doc *DocumentBase[T]) getNode() *documentNode {
	if doc.node == nil {
		doc.node = &documentNode{}
	}

	return doc.node
}
