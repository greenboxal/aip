package ford

import (
	"os"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"

	indexing22 "github.com/greenboxal/aip/aip-langchain/pkg/indexing"
	"github.com/greenboxal/aip/aip-langchain/pkg/indexing/impl"
	summarizers2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/summarizers"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/tokenizers"
)

type Manager struct {
	logger *zap.SugaredLogger
	index  *impl.Index
}

func NewManager(
	logger *zap.SugaredLogger,
	storage indexing22.MemoryStorage,
) (*Manager, error) {
	m := &Manager{
		logger: logger.Named("ford"),
	}

	oai := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	tokenizer, err := tokenizers.TikTokenForModel(openai.GPT3Dot5Turbo)

	if err != nil {
		return nil, err
	}

	m.index = impl.NewIndex(storage, indexing22.IndexConfiguration{
		Reducer: &summarizers2.MipMapSummarizer{
			Tokenizer: tokenizer,

			Summarizer: &summarizers2.ChatGptSummarizer{
				Client: oai,
				Model:  openai.GPT3Dot5Turbo,
			},

			MaxTokens: 1024,
			MinTokens: 128,
			MaxLevels: 3,
		},
	})

	return m, nil
}

func (m *Manager) Index() indexing22.Index {
	return m.index
}
