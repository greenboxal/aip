package chain

import (
	"fmt"
	"reflect"

	"github.com/greenboxal/aip/aip-forddb/pkg/tracing"
)

func New(options ...Option) Chain {
	sc := &basicChain{
		opts: NewChainOptions(options...),
	}

	if sc.opts.Name == "" {
		sc.opts.Name = reflect.TypeOf(sc.opts.Handler).String()
	}

	return sc
}

type basicChain struct {
	opts Options
}

func (s *basicChain) Name() string {
	return s.opts.Name
}

func (s *basicChain) Options() Options {
	return s.opts
}

func (s *basicChain) Run(ctx ChainContext) (err error) {
	defer func() {
		if e := recover(); err != nil {
			if er, ok := e.(error); ok {
				err = er
			} else {
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	spanCtx, span := tracing.StartSpan(ctx.Context(), s.Name())
	defer span.End()

	stepCtx := chainSubContext{
		ChainContext: ctx,
		ctx:          spanCtx,
	}

	return s.opts.Handler.Run(&stepCtx)
}
