package chain

import (
	"fmt"
	"reflect"

	"github.com/greenboxal/aip/aip-langchain/pkg/tracing"
)

type Hook interface {
	BeforeAll(ctx ChainContext) error
	BeforeOne(ctx ChainContext) error
	AfterOne(ctx ChainContext) error
	AfterAll(ctx ChainContext) error
}

type Chain interface{ Handler }
type Handler interface{ Run(ctx ChainContext) error }

type HandlerFunc func(ctx ChainContext) error

func (c HandlerFunc) Run(ctx ChainContext) error {
	return c(ctx)
}

func Func(handler HandlerFunc) Chain {
	return handler
}

func Sequential(items ...Handler) Chain {
	return sequentialChain(items)
}

type sequentialChain []Handler

func (s sequentialChain) Run(ctx ChainContext) error {
	for i, item := range s {
		run := func() (err error) {
			defer func() {
				if e := recover(); err != nil {
					if er, ok := e.(error); ok {
						err = er
					} else {
						err = fmt.Errorf("%v", e)
					}
				}
			}()

			spanCtx, span := tracing.Start(ctx.Context(), reflect.TypeOf(item).String())
			defer span.End()

			stepCtx := chainSubContext{
				ChainContext: ctx,
				ctx:          spanCtx,
			}

			return item.Run(&stepCtx)
		}

		if err := run(); err != nil {
			return err
		}

		if i < len(s)-1 {
			ctx.Flip()
		}
	}

	return nil
}

func WithOptions(options ...SubChainOption) Handler {
	return optionsChain(options)
}

type optionsChain []SubChainOption

func (o optionsChain) Run(ctx ChainContext) error {
	return ctx.ApplyOptions(o...)
}

type DefinedChain struct {
	Name     string
	Elements []Handler
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
