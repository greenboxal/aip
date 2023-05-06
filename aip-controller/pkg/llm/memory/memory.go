package memory

type AttentionPartial struct {
	Text string
}

type AttentionContext interface {
	Consume(text string) error
	BuildPartial(targetTokenCount int) (AttentionPartial, error)
}

type AttentionContextImpl struct {
}

func (a *AttentionContextImpl) Consume(text string) error {
	//TODO implement me
	panic("implement me")
}

func (a *AttentionContextImpl) BuildPartial(targetTokenCount int) (AttentionPartial, error) {
	//TODO implement me
	panic("implement me")
}
