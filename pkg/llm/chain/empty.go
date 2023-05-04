package chain

import "context"

func newEmptyChain() *emptyChainContext {
	return &emptyChainContext{
		ChainContext: NewChainContext(context.Background()),
	}
}

type emptyChainContext struct {
	ChainContext
}

var emptyChainInstance = newEmptyChain()

func EmptyContext() ChainContext                                       { return emptyChainInstance }
func (e *emptyChainContext) SetInput(name BasicContextKey, value any)  {}
func (e *emptyChainContext) SetOutput(name BasicContextKey, value any) {}
