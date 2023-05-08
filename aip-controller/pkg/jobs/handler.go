package jobs

import (
	"context"

	"github.com/jbenet/goprocess"
)

type BasicHandlerID interface {
	ID() BasicHandlerID
	String() string
}

type JobHandlerID string

func (id JobHandlerID) ID() BasicHandlerID {
	return BasicHandlerID(id)
}

func (id JobHandlerID) String() string {
	return string(id)
}

type JobHandlerKey[TPayload, TResult any] struct{ JobHandlerID }

type JobHandlerBinding interface {
	ID() BasicHandlerID
	Handler() JobHandler
}

type jobHandlerBinding struct {
	id      BasicHandlerID
	handler JobHandler
}

func (j *jobHandlerBinding) ID() BasicHandlerID {
	return j.id
}

func (j *jobHandlerBinding) Handler() JobHandler {
	return j.handler
}

type JobHandler interface {
	Run(ctx JobContext) error
}

type JobContext interface {
	JobHandle

	Payload() any

	SetError(err error)
	SetResult(result any)
}

type JobHandle interface {
	Job() Job
	Context() context.Context
	Process() goprocess.Process

	Await() (any, error)

	Result() any
	Error() error
}

type TypedJobHandle[TResult any] JobHandle
