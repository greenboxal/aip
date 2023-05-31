package cli

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var IndexGroup = &cobra.Command{
	Use: "index",
}

var Module = fx.Module(
	"cli",

	fx.Provide(New),

	RegisterGroup(IndexGroup),

	SubCommand[*IndexAddCommand](IndexGroup, &cobra.Command{
		Use: "add [path]",
	}, NewIndexAddCommand),
)

type Handler interface {
	Run(cmd *cobra.Command, args []string) error
}

type Option[T Handler] func(cmd *cobra.Command, h Handler)

func WithParent(parent *cobra.Command) Option[Handler] {
	return func(cmd *cobra.Command, h Handler) {
		parent.AddCommand(cmd)
	}
}

func Command[THandler Handler](
	cmd *cobra.Command,
	constructor any,
	options ...Option[THandler],
) fx.Option {
	return fx.Options(
		fx.Provide(constructor),

		fx.Invoke(func(m *Manager, handler THandler) {
			cmd.RunE = handler.Run

			for _, option := range options {
				option(cmd, handler)
			}

			if cmd.Parent() == nil {
				m.Root.AddCommand(cmd)
			}
		}),
	)
}

func SubCommand[THandler Handler](
	parent *cobra.Command,
	cmd *cobra.Command,
	constructor any,
	options ...Option[THandler],
) fx.Option {
	return fx.Options(
		fx.Provide(constructor),

		fx.Invoke(func(m *Manager, handler THandler) {
			cmd.RunE = handler.Run

			for _, option := range options {
				option(cmd, handler)
			}

			if cmd.Parent() == nil {
				parent.AddCommand(cmd)
			}
		}),
	)
}

func RegisterGroup(
	cmd *cobra.Command,
	options ...Option[Handler],
) fx.Option {
	return fx.Invoke(func(m *Manager) {
		if cmd.Parent() == nil {
			for _, option := range options {
				option(cmd, nil)
			}

			m.Root.AddCommand(cmd)
		}
	})
}
