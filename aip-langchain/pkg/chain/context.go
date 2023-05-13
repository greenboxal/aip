package chain

import (
	"context"

	"github.com/greenboxal/aip/aip-langchain/pkg/documents"
)

const DefaultInput ContextKey[string] = "completion_input"
const DefaultOutput ContextKey[string] = "completion_output"

type BasicContextKey interface {
	ContextKey() string
}

type IContextKey[T any] interface {
	BasicContextKey
}

type ContextKey[T any] string

func (c ContextKey[T]) ContextKey() string {
	return string(c)
}

type ChainInput struct {
	Key   BasicContextKey
	Value any
}

type ChainOutput struct {
	Key   BasicContextKey
	Value any
}

type ChainContext interface {
	Context() context.Context
	Documents() documents.Store

	NumIn() int
	NumOut() int

	Inputs() map[BasicContextKey]ChainInput
	Input(name BasicContextKey) (any, bool)

	Outputs() map[BasicContextKey]ChainOutput
	Output(name BasicContextKey) (any, bool)

	SetInput(name BasicContextKey, value any)
	SetOutput(name BasicContextKey, value any)

	ApplyOptions(o ...Option) error
	SubChain(chain Handler, options ...Option) (ChainContext, error)

	Flip()
	Reset()
}

func NewChainContext(ctx context.Context) ChainContext {
	return &chainContext{
		ctx:     ctx,
		inputs:  map[BasicContextKey]ChainInput{},
		outputs: map[BasicContextKey]ChainOutput{},
	}
}

type chainSubContext struct {
	ChainContext

	ctx context.Context
}

func (sc *chainSubContext) Context() context.Context {
	return sc.ctx
}

type chainContext struct {
	ctx context.Context

	documentStore documents.Store

	inputs  map[BasicContextKey]ChainInput
	outputs map[BasicContextKey]ChainOutput

	opts Options
}

func (cctx *chainContext) Context() context.Context                 { return cctx.ctx }
func (cctx *chainContext) Documents() documents.Store               { return cctx.documentStore }
func (cctx *chainContext) NumIn() int                               { return len(cctx.inputs) }
func (cctx *chainContext) NumOut() int                              { return len(cctx.outputs) }
func (cctx *chainContext) Inputs() map[BasicContextKey]ChainInput   { return cctx.inputs }
func (cctx *chainContext) Outputs() map[BasicContextKey]ChainOutput { return cctx.outputs }

func (cctx *chainContext) Input(key BasicContextKey) (any, bool) {
	input, ok := cctx.inputs[key]

	if !ok {
		return nil, false
	}

	return input.Value, true
}

func (cctx *chainContext) Output(key BasicContextKey) (any, bool) {
	input, ok := cctx.outputs[key]

	if !ok {
		return nil, false
	}

	return input.Value, true
}

func (cctx *chainContext) SetInput(name BasicContextKey, value any) {
	cctx.inputs[name] = ChainInput{
		Key:   name,
		Value: value,
	}
}

func (cctx *chainContext) SetOutput(name BasicContextKey, value any) {
	cctx.outputs[name] = ChainOutput{
		Key:   name,
		Value: value,
	}
}

func (cctx *chainContext) Run(chain Handler) error {
	if cctx.opts.IOMap != nil {
		for src, dst := range cctx.opts.IOMap {
			if dst.ToKind == IOKindOutput {
				continue
			}

			v, ok := GetIO(cctx, src.Kind, src.Key)

			if !ok {
				continue
			}

			SetIO(cctx, dst.ToKind, dst.ToKey, dst.Mapper(v))
		}
	}

	if err := chain.Run(cctx); err != nil {
		return err
	}

	if cctx.opts.IOMap != nil {
		for src, dst := range cctx.opts.IOMap {
			if dst.ToKind == IOKindInput {
				continue
			}

			v, ok := GetIO(cctx, src.Kind, src.Key)

			if !ok {
				continue
			}

			SetIO(cctx, dst.ToKind, dst.ToKey, dst.Mapper(v))
		}
	}

	return nil
}

func (cctx *chainContext) ApplyOptions(options ...Option) error {
	for _, opt := range options {
		opt.applyOption(&cctx.opts)
	}

	return nil
}

func (cctx *chainContext) SubChain(chain Handler, options ...Option) (ChainContext, error) {
	sb := cctx.clone()

	if err := sb.ApplyOptions(options...); err != nil {
		return nil, err
	}

	if err := sb.Run(chain); err != nil {
		return nil, err
	}

	return sb, nil
}

func (cctx *chainContext) Flip() {
	for outputKey, output := range cctx.outputs {
		cctx.inputs[outputKey] = ChainInput{
			Key:   outputKey,
			Value: output.Value,
		}

		delete(cctx.outputs, outputKey)
	}
}

func (cctx *chainContext) Reset() {
	cctx.inputs = map[BasicContextKey]ChainInput{}
	cctx.outputs = map[BasicContextKey]ChainOutput{}
}

func (cctx *chainContext) clone() *chainContext {
	c := &chainContext{
		ctx:           cctx.ctx,
		documentStore: cctx.documentStore,

		inputs:  map[BasicContextKey]ChainInput{},
		outputs: map[BasicContextKey]ChainOutput{},
	}

	for k, v := range cctx.inputs {
		c.inputs[k] = v
	}

	for k, v := range cctx.outputs {
		c.outputs[k] = v
	}

	return c
}
