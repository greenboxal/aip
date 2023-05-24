package indexer

import (
	"context"
	"reflect"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
	"github.com/greenboxal/aip/aip-langchain/pkg/providers/openai"
	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

const PageIndexerStreamID = "aip-wiki/indexer/page"

var PageResourceTypeID = forddb.TypeSystem().LookupByIDType(reflect.TypeOf((*models.PageID)(nil)).Elem())

type PageIndexer struct {
	forddb.LogConsumer

	db       forddb.Database
	index    vectorstore.Indexer
	embedder llm.Embedder
}

func NewPageIndexer(
	db forddb.Database,
	oai *openai.Client,
	indexProvider vectorstore.Provider,
) *PageIndexer {
	index := indexProvider.Collection("global_index", 1536)

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
		page := record.Current.(*models.Page)

		if page.Status.Markdown == "" {
			return nil
		}

		doc := &vectorstore.Document{}
		doc.ID = page.GetResourceID().String()
		doc.Type = page.GetResourceType().GetResourceID().Name()
		doc.Content = page.Status.Markdown

		_, err := i.index.IndexDocument(
			ctx,
			doc,
			vectorstore.WithIndexEmbedder(i.embedder),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
