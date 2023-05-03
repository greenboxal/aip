package llm

type OutputParser interface {
	Parse(ctx ChainContext, result string) error
}

type OutputParserFunc func(ctx ChainContext, result string) error

func (o OutputParserFunc) Parse(ctx ChainContext, result string) error {
	return o(ctx, result)
}

func NoopParser() OutputParser {
	return OutputParserFunc(func(ctx ChainContext, result string) error {
		return nil
	})
}

func StringCompletionParser(key BasicContextKey) OutputParser {
	return OutputParserFunc(func(ctx ChainContext, result string) error {
		ctx.SetOutput(key, result)

		return nil
	})
}
