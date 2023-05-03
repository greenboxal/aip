package documents

import (
	"context"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type DocumentPointer = forddb.ResourcePointer[DocumentID, Document]

type Store interface {
	Put(ctx context.Context, document Document) (DocumentPointer, error)
}
