package agents

type Agent interface {
}

type Message struct {
}

type Session interface {
	Send(msg Message) error
}
