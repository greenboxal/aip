package milvus

import (
	"context"
	"os"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Config) SetDefaults() {
	if c.Endpoint == "" {
		c.Endpoint = os.Getenv("MILVUS_ENDPOINT")
	}

	if c.Username == "" {
		c.Username = os.Getenv("MILVUS_USERNAME")
	}

	if c.Password == "" {
		c.Password = os.Getenv("MILVUS_PASSWORD")
	}
}

type Milvus struct {
	client client.Client
}

func NewMilvus(
	lc fx.Lifecycle,
	config *Config,
) (*Milvus, error) {
	m := &Milvus{}

	c, err := client.NewDefaultGrpcClientWithURI(
		context.Background(),
		config.Endpoint,
		config.Username,
		config.Password,
	)

	if err != nil {
		return nil, err
	}

	m.client = c

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return m.Close()
		},
	})

	return m, nil
}

func (m *Milvus) Client() client.Client {
	return m.client
}

func (m *Milvus) Close() error {
	return m.client.Close()
}

func (m *Milvus) IndexDocumentChunks(ctx context.Context, chunks []*vectorstore.DocumentChunk) error {
	documentIds := make([]string, len(chunks))
	documentTypes := make([]string, len(chunks))
	chunkIds := make([]int32, len(chunks))
	chunkContents := make([]string, len(chunks))
	chunkEmbeddings := make([][]float32, len(chunks))

	for i, v := range chunks {
		documentIds[i] = v.ID
		documentTypes[i] = v.Type
		chunkIds[i] = int32(v.Chunk)
		chunkContents[i] = v.Content
		chunkEmbeddings[i] = v.Embeddings.Embeddings
	}

	columns := []entity.Column{
		entity.NewColumnVarChar("document_id", documentIds),
		entity.NewColumnVarChar("document_type", documentTypes),
		entity.NewColumnInt32("chunk_id", chunkIds),
		entity.NewColumnVarChar("chunk_content", chunkContents),
		entity.NewColumnFloatVector("chunk_embeddings", 1536, chunkEmbeddings),
	}

	_, err := m.client.Insert(ctx, "global_index", "_default", columns...)

	if err != nil {
		return err
	}

	return nil
}
