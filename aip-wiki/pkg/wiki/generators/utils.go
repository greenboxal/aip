package generators

import (
	"html/template"

	"github.com/greenboxal/aip/aip-langchain/pkg/chain"
)

func GeneratedHtmlParser(key chain.ContextKey[string]) chain.OutputParser {
	return chain.OutputParserFunc(func(ctx chain.ChainContext, result string) error {
		ctx.SetOutput(key, result)

		return nil
	})
}

func GoTemplateParser(key chain.ContextKey[*template.Template]) chain.OutputParser {
	return chain.OutputParserFunc(func(ctx chain.ChainContext, result string) error {
		t, err := template.New("template").Parse(result)

		if err != nil {
			return err
		}

		ctx.SetOutput(key, t)

		return nil
	})
}
