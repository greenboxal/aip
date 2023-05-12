package chain

func Input[T any](ctx ChainContext, key ContextKey[T]) (def T) {
	val, ok := ctx.Input(key)

	if !ok {
		return def
	}

	return val.(T)
}

func Output[T any](ctx ChainContext, key ContextKey[T]) (def T) {
	val, ok := ctx.Output(key)

	if !ok {
		return def
	}

	return val.(T)
}

func SetInput[T any](ctx ChainContext, key ContextKey[T], value T) {
	ctx.SetInput(key, value)
}

func SetOutput[T any](ctx ChainContext, key ContextKey[T], value T) {
	ctx.SetOutput(key, value)
}
