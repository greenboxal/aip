package milvus

import (
	"context"
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/require"

	"github.com/greenboxal/aip/pkg/ford/forddb"
	"github.com/greenboxal/aip/pkg/indexing"
	"github.com/greenboxal/aip/pkg/indexing/impl"
	"github.com/greenboxal/aip/pkg/indexing/reducers/summarizers"
	"github.com/greenboxal/aip/pkg/indexing/reducers/tokenizers"
)

func TestMilvusStorage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	oai := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	milvus, err := NewStorage(oai)

	require.Nil(t, err)

	err = milvus.AppendSegment(ctx, &indexing.MemorySegment{
		ResourceMetadata: forddb.ResourceMetadata[indexing.MemorySegmentID, *indexing.MemorySegment]{
			ID: indexing.MemorySegmentID{"head"},
		},

		Memories: []indexing.Memory{
			{
				ResourceMetadata: forddb.ResourceMetadata[indexing.MemoryID, *indexing.Memory]{
					ID: indexing.MemoryID{"head"},
				},

				RootMemoryID:   indexing.MemoryID{"root"},
				BranchMemoryID: indexing.MemoryID{"branch"},
				ParentMemoryID: indexing.MemoryID{"parent"},

				Data: indexing.MemoryData{
					Text: "Hello, world!",
				},
			},
		},
	})

	require.Nil(t, err)
}

func TestMilvusStorageWithIndex(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	oai := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	milvus, err := NewStorage(oai)

	require.Nil(t, err)

	tokenizer, err := tokenizers.TikTokenForModel(openai.GPT3Dot5Turbo)

	require.Nil(t, err)

	index := impl.NewIndex(milvus, indexing.IndexConfiguration{
		Reducer: &summarizers.MipMapSummarizer{
			Tokenizer: tokenizer,

			Summarizer: &summarizers.ChatGptSummarizer{
				Client: oai,
				Model:  openai.GPT3Dot5Turbo,
			},

			MaxTokens: 1024,
			MinTokens: 128,
			MaxLevels: 3,
		},
	})

	sess, err := index.OpenSession(ctx, indexing.SessionOptions{
		RootMemoryID:    indexing.MemoryID{"root"},
		BranchMemoryID:  indexing.MemoryID{"branch"},
		InitialMemoryID: indexing.MemoryID{"initial"},
	})

	require.Nil(t, err)

	for _, data := range testData {
		err = sess.UpdateMemoryData(indexing.NewMemoryDataFromBytes([]byte(data)))

		require.Nil(t, err)
	}

	err = sess.Merge()

	require.Nil(t, err)
}
