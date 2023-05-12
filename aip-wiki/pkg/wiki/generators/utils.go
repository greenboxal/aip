package generators

import (
	"html/template"

	chain2 "github.com/greenboxal/aip/aip-langchain/pkg/llm/chain"
)

func GeneratedHtmlParser(key chain2.ContextKey[string]) chain2.OutputParser {
	return chain2.OutputParserFunc(func(ctx chain2.ChainContext, result string) error {
		ctx.SetOutput(key, result)

		return nil
	})
}

func GoTemplateParser(key chain2.ContextKey[*template.Template]) chain2.OutputParser {
	return chain2.OutputParserFunc(func(ctx chain2.ChainContext, result string) error {
		t, err := template.New("template").Parse(result)

		if err != nil {
			return err
		}

		ctx.SetOutput(key, t)

		return nil
	})
}
