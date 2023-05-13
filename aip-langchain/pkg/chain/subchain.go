package chain

func MapContext(options ...Option) Handler {
	return &mapChain{
		options: NewChainOptions(options...),
	}
}

type mapChain struct {
	options Options
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

func SubChain(chain Handler, options ...Option) Handler {
	return &subChain{
		chain: chain,

		mapChain: mapChain{
			options: NewChainOptions(options...),
		},
	}
}

type subChain struct {
	mapChain

	chain Handler
}

func (s *subChain) Run(ctx ChainContext) error {
	sb, err := ctx.SubChain(s.chain, WithSubChainOptions(s.options))

	if err != nil {
		return err
	}

	return s.mapChain.run(sb, ctx)
}

func NestedChain(chain Handler, options ...Option) Handler {
	return &nestedChain{
		chain: chain,

		mapChain: mapChain{
			options: NewChainOptions(options...),
		},
	}
}

type nestedChain struct {
	mapChain

	chain Handler
}

func (s *nestedChain) Run(ctx ChainContext) error {
	if err := s.chain.Run(ctx); err != nil {
		return err
	}

	return s.mapChain.run(ctx, ctx)
}
