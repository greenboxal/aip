package chain

type IOKind int

const (
	IOKindContext IOKind = iota
	IOKindInput
	IOKindOutput
)

type IOAddress struct {
	Kind IOKind
	Key  BasicContextKey
}

func MapInput[T any](src, dst ContextKey[T]) SubChainOption {
	return MapIO(IOKindInput, src, IOKindInput, dst)
}

func TransformInput[Src, Dst any](src ContextKey[Src], dst ContextKey[Dst], mapper func(Src) Dst) SubChainOption {
	return TransformIO(IOKindInput, src, IOKindInput, dst, mapper)
}

func MapOutput[T any](src, dst ContextKey[T]) SubChainOption {
	return MapIO(IOKindOutput, src, IOKindOutput, dst)
}

func TransformOutput[Src, Dst any](src ContextKey[Src], dst ContextKey[Dst], mapper func(Src) Dst) SubChainOption {
	return TransformIO(IOKindOutput, src, IOKindOutput, dst, mapper)
}

func MapIO[T any](srcKind IOKind, src ContextKey[T], dstKind IOKind, dst ContextKey[T]) SubChainOption {
	return TransformIO(srcKind, src, dstKind, dst, func(srcValue T) T {
		return srcValue
	})
}

func TransformIO[Src, Dst any, SrcKey IContextKey[Src], DstKey IContextKey[Dst]](
	srcKind IOKind,
	src SrcKey,
	dstKind IOKind,
	dst DstKey,
	fn func(Src) Dst,
) SubChainOption {
	return func(options *SubChainOptions) {
		if options.IOMap == nil {
			options.IOMap = make(map[IOAddress]MappedIO)
		}

		addr := IOAddress{Kind: srcKind, Key: src}

		options.IOMap[addr] = MappedIO{
			FromKind: srcKind,
			FromKey:  src,
			ToKind:   dstKind,
			ToKey:    dst,

			Mapper: func(v any) any {
				return fn(v.(Src))
			},
		}
	}
}

func GetIO(ctx ChainContext, kind IOKind, key BasicContextKey) (any, bool) {
	switch kind {
	case IOKindContext:
		fallthrough
	case IOKindInput:
		return ctx.Input(key)

	case IOKindOutput:
		return ctx.Output(key)
	}

	return nil, false
}

func SetIO(ctx ChainContext, kind IOKind, key BasicContextKey, value any) {
	switch kind {
	case IOKindContext:
		fallthrough
	case IOKindInput:
		ctx.SetInput(key, value)

	case IOKindOutput:
		ctx.SetOutput(key, value)
	}
}
