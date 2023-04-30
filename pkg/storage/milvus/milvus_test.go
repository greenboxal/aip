package milvus

import (
	"context"
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/require"

	"github.com/greenboxal/aip/pkg/collective"
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

	err = milvus.AppendSegment(ctx, &collective.MemorySegment{
		ResourceMetadata: forddb.ResourceMetadata[collective.MemorySegmentID, *collective.MemorySegment]{
			ID: collective.MemorySegmentID{"head"},
		},

		Memories: []collective.Memory{
			{
				ResourceMetadata: forddb.ResourceMetadata[collective.MemoryID, *collective.Memory]{
					ID: collective.MemoryID{"head"},
				},

				RootMemoryID:   collective.MemoryID{"root"},
				BranchMemoryID: collective.MemoryID{"branch"},
				ParentMemoryID: collective.MemoryID{"parent"},

				Data: collective.MemoryData{
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
		RootMemoryID:    collective.MemoryID{"root"},
		BranchMemoryID:  collective.MemoryID{"branch"},
		InitialMemoryID: collective.MemoryID{"initial"},
	})

	require.Nil(t, err)

	for _, data := range testData {
		err = sess.UpdateMemoryData(collective.NewMemoryDataFromBytes([]byte(data)))

		require.Nil(t, err)
	}

	err = sess.Merge()

	require.Nil(t, err)
}
