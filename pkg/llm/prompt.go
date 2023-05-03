package llm

import (
	"bytes"
	"fmt"
	"html/template"
)

type Prompt interface {
	Build(ctx ChainContext) (string, error)
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
	Template *template.Template
}

func NewTemplatePrompt(templateText string) *TemplatePrompt {
	tmpl := template.Must(template.New("TemplatePrompt").Parse(templateText))

	return &TemplatePrompt{
		Template: tmpl,
	}
}

func (t *TemplatePrompt) Build(ctx ChainContext) (string, error) {
	buffer := bytes.NewBuffer(nil)

	if err := t.Template.Execute(buffer, ctx.Inputs()); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

type TemplateContextPrompt struct {
	Key BasicContextKey
}

func TemplateFromContext(key BasicContextKey) *TemplateContextPrompt {
	return &TemplateContextPrompt{
		Key: key,
	}
}

func (t *TemplateContextPrompt) Build(ctx ChainContext) (string, error) {
	value := ctx.Input(t.Key)

	if valueStr, ok := value.(string); ok {
		return valueStr, nil
	}

	if valueStr, ok := value.(fmt.Stringer); ok {
		return valueStr.String(), nil
	}

	return fmt.Sprintf("%v", value), nil
}
