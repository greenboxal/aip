package chain

type MappedIO struct {
	FromKind IOKind
	FromKey  BasicContextKey
	ToKind   IOKind
	ToKey    BasicContextKey
	Mapper   func(any) any
}

type SubChainOptions struct {
	IOMap map[IOAddress]MappedIO
}

type SubChainOption func(options *SubChainOptions)

func NewSubChainOptions(options ...SubChainOption) SubChainOptions {
	result := SubChainOptions{}

	for _, option := range options {
		option(&result)
	}

	return result
}

func WithSubChainOptions(options SubChainOptions) SubChainOption {
	return func(result *SubChainOptions) {
		*result = options
	}
}

func MapContext(options ...SubChainOption) Chain {
	return &mapChain{
		options: NewSubChainOptions(options...),
	}
}

type mapChain struct {
	options SubChainOptions
}

func (s *mapChain) Run(ctx ChainContext) error {
	return s.run(ctx, ctx)
}

func (s *mapChain) run(src, dst ChainContext) error {
	for input, output := range s.options.IOMap {
		v, ok := GetIO(src, input.Kind, input.Key)

		if !ok {
			continue
		}

		SetIO(dst, output.ToKind, output.ToKey, output.Mapper(v))
	}

	return nil
}

func SubChain(chain Chain, options ...SubChainOption) Chain {
	return &subChain{
		chain: chain,

		mapChain: mapChain{
			options: NewSubChainOptions(options...),
		},
	}
}

type subChain struct {
	mapChain

	chain Chain
}

func (s *subChain) Run(ctx ChainContext) error {
	sb, err := ctx.SubChain(s.chain, WithSubChainOptions(s.options))

	if err != nil {
		return err
	}

	return s.mapChain.run(sb, ctx)
}

func NestedChain(chain Chain, options ...SubChainOption) Chain {
	return &nestedChain{
		chain: chain,

		mapChain: mapChain{
			options: NewSubChainOptions(options...),
		},
	}
}

type nestedChain struct {
	mapChain

	chain Chain
}

func (s *nestedChain) Run(ctx ChainContext) error {
	if err := s.chain.Run(ctx); err != nil {
		return err
	}

	return s.mapChain.run(ctx, ctx)
}
