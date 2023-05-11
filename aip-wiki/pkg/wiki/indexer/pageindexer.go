package indexer

import (
	"context"
	"reflect"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/indexing"
	"github.com/greenboxal/aip/aip-controller/pkg/llm"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/providers/openai"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

const PageIndexerStreamID = "aip-wiki/indexer/page"

var PageResourceTypeID = forddb.TypeSystem().LookupByIDType(reflect.TypeOf((*models.PageID)(nil)).Elem())

type PageIndexer struct {
	forddb.LogConsumer

	db       forddb.Database
	index    indexing.Provider
	embedder llm.Embedder
}

func NewPageIndexer(
	db forddb.Database,
	oai *openai.Client,
	index indexing.Provider,
) *PageIndexer {
	pi := &PageIndexer{}
	pi.db = db
	pi.index = index
	pi.embedder = &openai.Embedder{Client: oai, Model: openai.AdaEmbeddingV2}
	pi.LogStore = db.LogStore()
	pi.StreamID = PageIndexerStreamID
	pi.Handler = pi.handleStream
	return pi
}

func (i *PageIndexer) handleStream(ctx context.Context, record *forddb.LogEntryRecord) error {
	if record.Type != PageResourceTypeID.GetResourceID() {
		return nil
	}

	switch record.Kind {
	case forddb.LogEntryKindSet:
		page, err := forddb.Convert[*models.Page](record.Current)

		if err != nil {
			return err
		}

		if page.Status.Markdown == "" {
			return nil
		}

		doc := &indexing.Document{}
		doc.ID = page.GetResourceID().String()
		doc.Type = page.GetResourceType().Name()
		doc.Content = page.Status.Markdown

		_, err = i.index.IndexDocument(
			ctx,
			doc,
			indexing.WithIndexEmbedder(i.embedder),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
