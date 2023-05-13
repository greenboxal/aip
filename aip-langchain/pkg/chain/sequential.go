package chain

func Sequential(steps ...Handler) SequentialHandler {
	return steps
}

type SequentialHandler []Handler

func (s SequentialHandler) applyOption(options *Options) {
	options.Handler = s
}

func (s SequentialHandler) Run(ctx ChainContext) error {
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
