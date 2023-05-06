package wiki

import (
	openai "github.com/greenboxal/aip/aip-controller/pkg/llm/providers/openai"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/tokenizers"
)

func NewWiki(client *openai.Client) (*Wiki, error) {
	var err error

	w := &Wiki{}
	w.client = client

	w.model = &openai.ChatLanguageModel{
		Client: client,
		Model:  "gpt-3.5-turbo",
	}

	w.tokenizer, err = tokenizers.TikTokenForModel(openai.AdaEmbeddingV2.String())

	if err != nil {
		return nil, err
	}

	return w, nil
}

type Wiki struct {
	client *openai.Client

	model     *openai.ChatLanguageModel
	tokenizer *tokenizers.TikTokenTokenizer
}
