package chat

import (
	"fmt"

	"github.com/greenboxal/aip/pkg/llm"
)

type Prompt interface {
	AsPrompt() llm.Prompt
	Build(ctx llm.ChainContext) (Message, error)
}

type PromptFunc func(ctx llm.ChainContext) (Message, error)

func (t PromptFunc) Build(ctx llm.ChainContext) (Message, error) {
	return t(ctx)
}

func ComposeTemplate(entries ...Prompt) Prompt {
	return &HistoryTemplate{
		Entries: entries,
	}
}

func EntryTemplate(name string, role Role, template llm.Prompt) *SingleEntryTemplate {
	return &SingleEntryTemplate{
		Name:     name,
		Role:     role,
		Template: template,
	}
}

func HistoryFromContext(key llm.ContextKey[Message]) *HistoryFromContextTemplate {
	return &HistoryFromContextTemplate{
		Key: key,
	}
}

type HistoryTemplate struct {
	Entries []Prompt
}

func (h *HistoryTemplate) AsPrompt() llm.Prompt { return (*historyPromptTemplate)(h) }

func (h *HistoryTemplate) Build(ctx llm.ChainContext) (Message, error) {
	var entries []MessageEntry

	for _, entry := range h.Entries {
		message, err := entry.Build(ctx)

		if err != nil {
			return Message{}, err
		}

		entries = append(entries, message.Entries...)
	}

	return Message{Entries: entries}, nil
}

type historyPromptTemplate HistoryTemplate

func (b *historyPromptTemplate) Build(ctx llm.ChainContext) (string, error) {
	result, err := (*HistoryTemplate)(b).Build(ctx)

	if err != nil {
		return "", err
	}

	return result.String(), nil
}

type SingleEntryTemplate struct {
	Name     string
	Role     Role
	Template llm.Prompt
}

func (p *SingleEntryTemplate) AsPrompt() llm.Prompt { return (*basicPromptTemplate)(p) }

func (p *SingleEntryTemplate) Build(ctx llm.ChainContext) (Message, error) {
	prompt, err := p.Template.Build(ctx)

	if err != nil {
		return Message{}, err
	}

	return Compose(Entry(p.Role, prompt)), nil
}

type basicPromptTemplate SingleEntryTemplate

func (b *basicPromptTemplate) Build(ctx llm.ChainContext) (string, error) {
	result, err := (*SingleEntryTemplate)(b).Build(ctx)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s: %s", b.Role, result), nil
}

type HistoryFromContextTemplate struct {
	Key llm.ContextKey[Message]
}

func (b *HistoryFromContextTemplate) AsPrompt() llm.Prompt { panic("not supported") }

func (b *HistoryFromContextTemplate) Build(ctx llm.ChainContext) (Message, error) {
	history := llm.GetInput(ctx, b.Key)

	return history, nil
}
