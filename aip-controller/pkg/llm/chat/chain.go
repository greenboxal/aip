package chat

import (
	chain2 "github.com/greenboxal/aip/aip-controller/pkg/llm/chain"
)

type predictChain struct {
	model   LanguageModel
	prompt  Prompt
	outputs []chain2.OutputParser
}

func (p *predictChain) Run(ctx chain2.ChainContext) error {
	prompt, err := p.prompt.Build(ctx)

	if err != nil {
		return err
	}

	result, err := p.model.PredictChat(ctx.Context(), prompt)

	if err != nil {
		return err
	}

	for _, output := range p.outputs {
		err := output.Parse(ctx, result.Entries[0].Content)

		if err != nil {
			return err
		}
	}

	return nil
}

func CompletionMessageParser(key chain2.ContextKey[Message]) chain2.OutputParser {
	return chain2.OutputParserFunc(func(ctx chain2.ChainContext, result string) error {
		msg := Compose(Entry(RoleAI, result))

		ctx.SetOutput(key, msg)

		return nil
	})
}

func Predict(model LanguageModel, prompt Prompt, parsers ...chain2.OutputParser) chain2.Chain {
	if len(parsers) == 0 {
		parsers = []chain2.OutputParser{
			CompletionMessageParser(ChatReplyContextKey),
		}
	}

	return &predictChain{
		model:   model,
		prompt:  prompt,
		outputs: parsers,
	}
}
