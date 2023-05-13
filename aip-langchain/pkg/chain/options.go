package chain

type Options struct {
	Name    string
	Handler Handler
	IOMap   IOMap
}

func (o *Options) Apply(options ...Option) {
	for _, option := range options {
		option.applyOption(o)
	}
}

type Option interface {
	applyOption(options *Options)
}

type OptionFunc func(options *Options)

func (o OptionFunc) applyOption(options *Options) {
	o(options)
}

func NewChainOptions(options ...Option) Options {
	result := Options{}

	result.Apply(options...)

	return result
}

func WithSubChainOptions(options Options) OptionFunc {
	return func(result *Options) {
		*result = options
	}
}

func WithName(name string) OptionFunc {
	return func(options *Options) {
		options.Name = name
	}
}

func WithHandler(handler Handler) OptionFunc {
	return func(options *Options) {
		options.Handler = handler
	}
}

func WithFunc(handler HandlerFunc) OptionFunc {
	return func(options *Options) {
		options.Handler = handler
	}
}
