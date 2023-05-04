package chain

import (
	"github.com/greenboxal/aip/aip-controller/pkg/llm"
)

type predictChain struct {
	model   llm.LanguageModel
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

func Predict(model llm.LanguageModel, prompt Prompt, parsers ...OutputParser) Chain {
	if len(parsers) == 0 {
		parsers = []OutputParser{StringCompletionParser(DefaultOutput)}
	}

	return &predictChain{
		model:   model,
		prompt:  prompt,
		parsers: parsers,
	}
}
