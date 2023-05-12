package chain

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/greenboxal/aip/aip-langchain/pkg/llm/tokenizers"
)

type PromptOption func(*PromptOptions)

func NewPromptOptions(options ...PromptOption) *PromptOptions {
	opts := &PromptOptions{}

	for _, option := range options {
		option(opts)
	}

	return opts
}

func WithRequirement(kind IOKind, keys ...BasicContextKey) PromptOption {
	return func(opts *PromptOptions) {
		for _, k := range keys {
			opts.Requirements = append(opts.Requirements, IOAddress{Kind: kind, Key: k})
		}
	}
}

func WithRequiredInput(keys ...BasicContextKey) PromptOption {
	return WithRequirement(IOKindInput, keys...)
}

func WithRequiredContext(keys ...BasicContextKey) PromptOption {
	return WithRequirement(IOKindContext, keys...)
}

func WithRequiredOutput(keys ...BasicContextKey) PromptOption {
	return WithRequirement(IOKindOutput, keys...)
}

type PromptOptions struct {
	Requirements []IOAddress
}

type Prompt interface {
	Build(ctx ChainContext) (string, error)
}

type HasGetEmptyTokenCount interface {
	GetEmptyTokenCount(tokenizer tokenizers.BasicTokenizer) int
}

type PromptFunc func(ctx ChainContext) (string, error)

func (t PromptFunc) Build(ctx ChainContext) (string, error) {
	return t(ctx)
}

func Static(text string) Prompt {
	return PromptFunc(func(ctx ChainContext) (string, error) {
		return text, nil
	})
}

type TemplatePrompt struct {
	Template     *template.Template
	Requirements []IOAddress

	emptyTokenCount int
}

func NewTemplatePrompt(templateText string, options ...PromptOption) *TemplatePrompt {
	opts := NewPromptOptions(options...)
	tmpl := template.Must(template.New("TemplatePrompt").Parse(templateText))

	return &TemplatePrompt{
		Template:     tmpl,
		Requirements: opts.Requirements,

		emptyTokenCount: -1,
	}
}

func (t *TemplatePrompt) GetEmptyTokenCount(tokenizer tokenizers.BasicTokenizer) int {
	if t.emptyTokenCount == -1 {
		t.emptyTokenCount = GetEmptyTokenCountNoCache(tokenizer, t)
	}

	return t.emptyTokenCount
}

func (t *TemplatePrompt) Build(ctx ChainContext) (string, error) {
	buffer := bytes.NewBuffer(nil)

	inputs := map[string]interface{}{}

	for _, req := range t.Requirements {
		if req.Kind == IOKindOutput {
			continue
		}

		value, ok := GetIO(ctx, req.Kind, req.Key)

		if !ok {
			panic("missing requirement")
		}

		inputs[req.Key.ContextKey()] = value
	}

	if err := t.Template.Execute(buffer, inputs); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

type TemplateContextPrompt struct {
	Key BasicContextKey

	emptyTokenCount int
}

func TemplateFromContext(key BasicContextKey) *TemplateContextPrompt {
	return &TemplateContextPrompt{
		Key:             key,
		emptyTokenCount: -1,
	}
}

func (t *TemplateContextPrompt) Build(ctx ChainContext) (string, error) {
	value, ok := ctx.Input(t.Key)

	if !ok {
		return "", fmt.Errorf("missing input %s", t.Key)
	}

	if valueStr, ok := value.(string); ok {
		return valueStr, nil
	}

	if valueStr, ok := value.(fmt.Stringer); ok {
		return valueStr.String(), nil
	}

	return fmt.Sprintf("%v", value), nil
}

func (t *TemplateContextPrompt) GetEmptyTokenCount(tokenizer tokenizers.BasicTokenizer) int {
	if t.emptyTokenCount == -1 {
		t.emptyTokenCount = GetEmptyTokenCountNoCache(tokenizer, t)
	}

	return t.emptyTokenCount
}

func GetEmptyTokenCount(tokenizer tokenizers.BasicTokenizer, t Prompt) int {
	if hasGetEmptyTokenCount, ok := t.(HasGetEmptyTokenCount); ok {
		return hasGetEmptyTokenCount.GetEmptyTokenCount(tokenizer)
	}

	return GetEmptyTokenCountNoCache(tokenizer, t)
}

func GetEmptyTokenCountNoCache(tokenizer tokenizers.BasicTokenizer, t Prompt) int {
	str, err := t.Build(EmptyContext())

	if err != nil {
		panic(err)
	}

	count, err := tokenizer.Count(str)

	if err != nil {
		panic(err)
	}

	return count
}
