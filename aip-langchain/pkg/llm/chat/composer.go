package chat

import (
	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/tokenizers"
)

type MessageComposer struct {
	Tokenizer tokenizers.BasicTokenizer

	llm.Constraints
	llm.Statistics

	entries []MessageEntry
}

func (mc *MessageComposer) Append(entry MessageEntry) error {
	count, err := entry.CalculateTokenCount(mc.Tokenizer)

	if err != nil {
		return err
	}

	mc.entries = append(mc.entries, entry)
	mc.TokenCount += count

	return nil
}

func (mc *MessageComposer) RemainingTokens() int {
	return mc.MaxTokens - mc.TokenCount
}

func (mc *MessageComposer) Validate() error {
	return nil
}

func (mc *MessageComposer) Build() (Message, error) {
	if err := mc.Validate(); err != nil {
		return Message{}, err
	}

	m := Message{
		Entries: mc.entries,
	}

	return m, nil
}
