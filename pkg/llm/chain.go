package llm

type Chainable interface {
	Run(ctx ChainContext) error
}

type ChainFunc func(ctx ChainContext) error

func (c ChainFunc) Run(ctx ChainContext) error {
	return c(ctx)
}

func Chain(items ...Chainable) Chainable {
	return sliceChain(items)
}

type sliceChain []Chainable

func (s sliceChain) Run(ctx ChainContext) error {
	for i, item := range s {
		if err := item.Run(ctx); err != nil {
			return err
		}

		if i < len(s)-1 {
			ctx.Flip()
		}
	}

	return nil
}

type predictChain struct {
	model   LanguageModel
	prompt  Prompt
	parsers []OutputParser
}

func (p *predictChain) Run(ctx ChainContext) error {
	prompt, err := p.prompt.Build(ctx)

	if err != nil {
		return err
	}

	result, err := p.model.Predict(ctx.Context(), prompt)

	if err != nil {
		return err
	}

	for _, output := range p.parsers {
		err := output.Parse(ctx, result)

		if err != nil {
			return err
		}
	}

	return nil
}

const CompletionInputContextKey ContextKey[string] = "completion_input"
const CompletionOutputContextKey ContextKey[string] = "completion_output"

func Predict(model LanguageModel, prompt Prompt, parsers ...OutputParser) Chainable {
	if len(parsers) == 0 {
		parsers = []OutputParser{StringCompletionParser(CompletionOutputContextKey)}
	}

	return &predictChain{
		model:   model,
		prompt:  prompt,
		parsers: parsers,
	}
}
