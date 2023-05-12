package cms

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type Document interface {
	DocumentID() forddb.BasicResourceID

	AsText() string
}

type DocumentIndexer struct {
}

func (di *DocumentIndexer) Index() {
}
