package chain

type Hook interface {
	BeforeAll(ctx ChainContext) error
	BeforeOne(ctx ChainContext) error
	AfterOne(ctx ChainContext) error
	AfterAll(ctx ChainContext) error
}

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

type DefinedChain struct {
	Name     string
	Elements []Chain
	Hooks    []Hook
}

func (dc *DefinedChain) Run(ctx ChainContext) error {
	for _, hook := range dc.Hooks {
		if err := hook.BeforeAll(ctx); err != nil {
			return err
		}
	}

	for i, item := range dc.Elements {
		for _, hook := range dc.Hooks {
			if err := hook.BeforeOne(ctx); err != nil {
				return err
			}
		}

		if err := item.Run(ctx); err != nil {
			return err
		}

		if i < len(dc.Elements)-1 {
			ctx.Flip()
		}

		for _, hook := range dc.Hooks {
			if err := hook.AfterOne(ctx); err != nil {
				return err
			}
		}
	}

	for _, hook := range dc.Hooks {
		if err := hook.AfterAll(ctx); err != nil {
			return err
		}
	}

	return nil
}
