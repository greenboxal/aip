package chat

import "github.com/greenboxal/aip/pkg/llm"

type predictChain struct {
	model   LanguageModel
	prompt  Prompt
	outputs []llm.OutputParser
}

func (p *predictChain) Run(ctx llm.ChainContext) error {
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

func CompletionMessageParser(key llm.ContextKey[Message]) llm.OutputParser {
	return llm.OutputParserFunc(func(ctx llm.ChainContext, result string) error {
		msg := Compose(Entry(RoleAI, result))

		ctx.SetOutput(key, msg)

		return nil
	})
}

func Predict(model LanguageModel, prompt Prompt, parsers ...llm.OutputParser) llm.Chainable {
	if len(parsers) == 0 {
		parsers = []llm.OutputParser{
			CompletionMessageParser(ChatReplyContextKey),
		}
	}

	return &predictChain{
		model:   model,
		prompt:  prompt,
		outputs: parsers,
	}
}
