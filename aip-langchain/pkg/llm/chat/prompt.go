package chat

import (
	"fmt"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chain"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/tokenizers"
)

type Prompt interface {
	AsPrompt() chain.Prompt
	Build(ctx chain.ChainContext) (Message, error)
}

type PromptFunc func(ctx chain.ChainContext) (Message, error)

func (t PromptFunc) Build(ctx chain.ChainContext) (Message, error) {
	return t(ctx)
}

func ComposeTemplate(entries ...Prompt) Prompt {
	return &HistoryTemplate{
		Entries: entries,
	}
}

func EntryTemplate(role msn.Role, template chain.Prompt) *SingleEntryTemplate {
	return &SingleEntryTemplate{
		Name:     "",
		Role:     role,
		Template: template,
	}
}

func HistoryFromContext(key chain.ContextKey[Message]) *HistoryFromContextTemplate {
	return &HistoryFromContextTemplate{
		Key: key,
	}
}

type HistoryTemplate struct {
	Entries []Prompt

	hasEmptyTokenCount bool
	emptyTokenCount    int
}

func (tp *HistoryTemplate) AsPrompt() chain.Prompt { return (*historyPromptTemplate)(tp) }

func (tp *HistoryTemplate) Build(ctx chain.ChainContext) (Message, error) {
	var entries []MessageEntry

	for _, entry := range tp.Entries {
		message, err := entry.Build(ctx)

		if err != nil {
			return Message{}, err
		}

		entries = append(entries, message.Entries...)
	}

	return Message{Entries: entries}, nil
}

func (tp *HistoryTemplate) GetEmptyTokenCount(tokenizer tokenizers.BasicTokenizer) int {
	if !tp.hasEmptyTokenCount {
		tp.emptyTokenCount = chain.GetEmptyTokenCountNoCache(tokenizer, tp.AsPrompt())
		tp.hasEmptyTokenCount = true
	}

	return tp.emptyTokenCount
}

type historyPromptTemplate HistoryTemplate

func (p *historyPromptTemplate) Build(ctx chain.ChainContext) (string, error) {
	result, err := (*HistoryTemplate)(p).Build(ctx)

	if err != nil {
		return "", err
	}

	return result.String(), nil
}

type SingleEntryTemplate struct {
	Name     string
	Role     msn.Role
	Template chain.Prompt

	hasEmptyTokenCount bool
	emptyTokenCount    int
}

func (tp *SingleEntryTemplate) AsPrompt() chain.Prompt { return (*basicPromptTemplate)(tp) }

func (tp *SingleEntryTemplate) Build(ctx chain.ChainContext) (Message, error) {
	prompt, err := tp.Template.Build(ctx)

	if err != nil {
		return Message{}, err
	}

	return Compose(Entry(tp.Role, prompt)), nil
}

func (tp *SingleEntryTemplate) GetEmptyTokenCount(tokenizer tokenizers.BasicTokenizer) int {
	if !tp.hasEmptyTokenCount {
		tp.emptyTokenCount = chain.GetEmptyTokenCountNoCache(tokenizer, tp.AsPrompt())
		tp.hasEmptyTokenCount = true
	}

	return tp.emptyTokenCount
}

type basicPromptTemplate SingleEntryTemplate

func (tp *basicPromptTemplate) Build(ctx chain.ChainContext) (string, error) {
	result, err := (*SingleEntryTemplate)(tp).Build(ctx)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s: %s", tp.Role, result), nil
}

type HistoryFromContextTemplate struct {
	Key chain.ContextKey[Message]

	hasEmptyTokenCount bool
	emptyTokenCount    int
}

func (tp *HistoryFromContextTemplate) AsPrompt() chain.Prompt { panic("not supported") }

func (tp *HistoryFromContextTemplate) Build(ctx chain.ChainContext) (Message, error) {
	history := chain.Input(ctx, tp.Key)

	return history, nil
}

func (tp *HistoryFromContextTemplate) GetEmptyTokenCount(tokenizer tokenizers.BasicTokenizer) int {
	if !tp.hasEmptyTokenCount {
		tp.emptyTokenCount = chain.GetEmptyTokenCountNoCache(tokenizer, tp.AsPrompt())
		tp.hasEmptyTokenCount = true
	}

	return tp.emptyTokenCount
}
