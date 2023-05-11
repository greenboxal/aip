package ford

import (
	"os"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"

	indexing2 "github.com/greenboxal/aip/aip-controller/pkg/indexing2"
	"github.com/greenboxal/aip/aip-controller/pkg/indexing2/impl"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/summarizers"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/tokenizers"
)

type Manager struct {
	logger *zap.SugaredLogger
	index  *impl.Index
}

func NewManager(
	logger *zap.SugaredLogger,
	storage indexing2.MemoryStorage,
) (*Manager, error) {
	m := &Manager{
		logger: logger.Named("ford"),
	}

	oai := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	tokenizer, err := tokenizers.TikTokenForModel(openai.GPT3Dot5Turbo)

	if err != nil {
		return nil, err
	}

	m.index = impl.NewIndex(storage, indexing2.IndexConfiguration{
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

	return m, nil
}

func (m *Manager) Index() indexing2.Index {
	return m.index
}
