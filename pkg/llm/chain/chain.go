package chain

type Chain interface {
	Run(ctx ChainContext) error
}

type Func func(ctx ChainContext) error

func (c Func) Run(ctx ChainContext) error {
	return c(ctx)
}

func Compose(items ...Chain) Chain {
	return sliceChain(items)
}

type sliceChain []Chain

func (s sliceChain) Run(ctx ChainContext) error {
	for i, item := range s {
		if err := item.Run(ctx); err != nil {
			return err
		}

		if i < len(s)-1 {
			ctx.Flip()
		}
	}

	return nil
}

func WithOptions(options ...SubChainOption) Chain {
	return optionsChain(options)
}

type optionsChain []SubChainOption

func (o optionsChain) Run(ctx ChainContext) error {
	return ctx.ApplyOptions(o...)
}
