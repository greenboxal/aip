package exp

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-sdk/pkg/cli"
	"github.com/greenboxal/aip/aip-wiki/exp/mdcodegen"
)

var Module = fx.Module(
	"wiki/exp",

	cli.Command[*EmbeddingTestCommand](&cobra.Command{
		Use: "embedding-test",
	}, NewEmbeddingTestCommand),

	cli.Command[*mdcodegen.MarkdownCodeGenCommand](&cobra.Command{
		Use: "mdcodegen <file>",
	}, mdcodegen.NewMarkdownCodeGenCommand),
)
