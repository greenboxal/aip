package utils

import "go.uber.org/fx"

type BindingRegistry[TBinding any] interface {
	Register(item TBinding)
	Bindings() []TBinding
}

type bindingRegistry[TBinding any] struct {
	bindings []TBinding
}

func (br *bindingRegistry[TBinding]) Register(item TBinding) {
	br.bindings = append(br.bindings, item)
}

func (br *bindingRegistry[TBinding]) Bindings() []TBinding {
	return br.bindings
}

func WithBindingRegistry[TBinding any](name string) fx.Option {
	return fx.Provide(fx.Annotate(func(bindings []TBinding) BindingRegistry[TBinding] {
		return &bindingRegistry[TBinding]{bindings: bindings}
	}, fx.ParamTags(`group:"`+name+`"`)))
}

func WithBinding[TBinding any](name string, constructor any) fx.Option {
	return fx.Provide(fx.Annotate(constructor, fx.As((*TBinding)(nil)), fx.ResultTags(`group:"`+name+`"`)))
}
