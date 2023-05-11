package milvus

import (
	"context"
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/require"

	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	indexing2 "github.com/greenboxal/aip/aip-controller/pkg/indexing2"
	"github.com/greenboxal/aip/aip-controller/pkg/indexing2/impl"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/summarizers"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/tokenizers"
)

func TestMilvusStorage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	oai := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	milvus, err := NewStorage(oai)

	require.Nil(t, err)

	err = milvus.AppendSegment(ctx, &collective2.MemorySegment{
		ResourceBase: forddb.ResourceBase[collective2.MemorySegmentID, *collective2.MemorySegment]{
			ID: collective2.MemorySegmentID{"head"},
		},

		Memories: []collective2.Memory{
			{
				ResourceBase: forddb.ResourceBase[collective2.MemoryID, *collective2.Memory]{
					ID: collective2.MemoryID{"head"},
				},

				RootMemoryID:   collective2.MemoryID{"root"},
				BranchMemoryID: collective2.MemoryID{"branch"},
				ParentMemoryID: collective2.MemoryID{"parent"},

				Data: collective2.MemoryData{
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

	index := impl.NewIndex(milvus, indexing2.IndexConfiguration{
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

	sess, err := index.OpenSession(ctx, indexing2.SessionOptions{
		RootMemoryID:    collective2.MemoryID{"root"},
		BranchMemoryID:  collective2.MemoryID{"branch"},
		InitialMemoryID: collective2.MemoryID{"initial"},
	})

	require.Nil(t, err)

	for _, data := range testData {
		err = sess.UpdateMemoryData(collective2.NewMemoryDataFromBytes([]byte(data)))

		require.Nil(t, err)
	}

	err = sess.Merge()

	require.Nil(t, err)
}
