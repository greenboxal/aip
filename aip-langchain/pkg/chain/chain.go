package chain

type Chain interface {
	Name() string
	Options() Options

	Handler
}

type Handler interface{ Run(ctx ChainContext) error }

type HandlerFunc func(ctx ChainContext) error

func (c HandlerFunc) Run(ctx ChainContext) error {
	return c(ctx)
}

func Func(handler HandlerFunc, options ...Option) Chain {
	options = append(options, WithHandler(handler))

	return New(options...)
}
