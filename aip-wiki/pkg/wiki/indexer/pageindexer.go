package indexer

import (
	"context"
	"reflect"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

const PageIndexerStreamID = "aip-wiki/indexer/page"

var PageResourceTypeID = forddb.TypeSystem().LookupByIDType(reflect.TypeOf((*models.PageID)(nil)).Elem())

type PageIndexer struct {
	forddb.LogConsumer

	db forddb.Database
}

func NewPageIndexer(db forddb.Database) *PageIndexer {
	pi := &PageIndexer{}
	pi.db = db
	pi.LogStore = db.LogStore()
	pi.StreamID = PageIndexerStreamID
	pi.Handler = pi.handleStream
	return pi
}

func (i *PageIndexer) handleStream(ctx context.Context, record *forddb.LogEntryRecord) error {
	if record.Type != PageResourceTypeID.GetResourceTypeID() {
		return nil
	}

	switch record.Kind {
	case forddb.LogEntryKindSet:

	}

	return nil
}
