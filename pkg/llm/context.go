package llm

import (
	"context"

	"github.com/greenboxal/aip/pkg/llm/documents"
)

type BasicContextKey interface {
	ContextKey() string
}

type ContextKey[T any] string

func (c ContextKey[T]) ContextKey() string {
	return string(c)
}

type ChainInput struct {
	Key   BasicContextKey
	Value any
}

func Input(key BasicContextKey, value any) ChainInput {
	return ChainInput{
		Key:   key,
		Value: value,
	}
}

func Output(key BasicContextKey, parser OutputParser) *ChainOutput {
	return &ChainOutput{
		Key:    key,
		Parser: parser,
	}
}

type ChainOutput struct {
	Key    BasicContextKey
	Parser OutputParser
	Value  any
}

type ChainContext interface {
	Context() context.Context
	Documents() documents.Store

	NumIn() int
	NumOut() int

	Inputs() map[BasicContextKey]ChainInput
	Input(name BasicContextKey) any

	Outputs() map[BasicContextKey]ChainOutput
	Output(name BasicContextKey) any

	SetInput(name BasicContextKey, value any)
	SetOutput(name BasicContextKey, value any)

	Flip()
	Reset()
}

func GetInput[T any](ctx ChainContext, key ContextKey[T]) T {
	return ctx.Input(key).(T)
}

func GetOutput[T any](ctx ChainContext, key ContextKey[T]) T {
	return ctx.Output(key).(T)
}

func SetInput[T any](ctx ChainContext, key ContextKey[T], value T) {
	ctx.SetInput(key, value)
}

func SetOutput[T any](ctx ChainContext, key ContextKey[T], value T) {
	ctx.SetOutput(key, value)
}

type chainContext struct {
	ctx context.Context

	documentStore documents.Store

	inputs  map[BasicContextKey]ChainInput
	outputs map[BasicContextKey]ChainOutput
}

func NewChainContext(ctx context.Context) ChainContext {
	return &chainContext{
		ctx:     ctx,
		inputs:  map[BasicContextKey]ChainInput{},
		outputs: map[BasicContextKey]ChainOutput{},
	}
}

func (p *chainContext) Context() context.Context                 { return p.ctx }
func (p *chainContext) Documents() documents.Store               { return p.documentStore }
func (p *chainContext) NumIn() int                               { return len(p.inputs) }
func (p *chainContext) NumOut() int                              { return len(p.outputs) }
func (p *chainContext) Inputs() map[BasicContextKey]ChainInput   { return p.inputs }
func (p *chainContext) Outputs() map[BasicContextKey]ChainOutput { return p.outputs }

func (p *chainContext) Input(key BasicContextKey) any {
	input, ok := p.inputs[key]

	if !ok {
		panic("input not found")
	}

	return input.Value
}

func (p *chainContext) Output(key BasicContextKey) any {
	input, ok := p.outputs[key]

	if !ok {
		panic("output not found")
	}

	return input.Value
}

func (p *chainContext) SetInput(name BasicContextKey, value any) {
	p.inputs[name] = ChainInput{
		Key:   name,
		Value: value,
	}
}

func (p *chainContext) SetOutput(name BasicContextKey, value any) {
	p.outputs[name] = ChainOutput{
		Key:   name,
		Value: value,
	}
}

func (p *chainContext) Flip() {
	for outputKey, output := range p.outputs {
		p.inputs[outputKey] = ChainInput{
			Key:   outputKey,
			Value: output.Value,
		}

		delete(p.outputs, outputKey)
	}
}

func (p *chainContext) Reset() {
	p.inputs = map[BasicContextKey]ChainInput{}
	p.outputs = map[BasicContextKey]ChainOutput{}
}
