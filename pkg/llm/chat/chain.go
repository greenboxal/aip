package chat

import (
	"github.com/greenboxal/aip/pkg/llm/chain"
)

type predictChain struct {
	model   LanguageModel
	prompt  Prompt
	outputs []chain.OutputParser
}

func (p *predictChain) Run(ctx chain.ChainContext) error {
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

func CompletionMessageParser(key chain.ContextKey[Message]) chain.OutputParser {
	return chain.OutputParserFunc(func(ctx chain.ChainContext, result string) error {
		msg := Compose(Entry(RoleAI, result))

		ctx.SetOutput(key, msg)

		return nil
	})
}

func Predict(model LanguageModel, prompt Prompt, parsers ...chain.OutputParser) chain.Chain {
	if len(parsers) == 0 {
		parsers = []chain.OutputParser{
			CompletionMessageParser(ChatReplyContextKey),
		}
	}

	return &predictChain{
		model:   model,
		prompt:  prompt,
		outputs: parsers,
	}
}
