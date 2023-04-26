package ford

import (
	"context"

	"github.com/greenboxal/aip/collective"
)

type ProductionLineStep struct {
	Description *StepDescription
	NextStep    *ProductionLineStep
}

type AgentHandler interface {
	HandleOutput(ctx context.Context, state *TaskState, output *TaskPayload) (*TaskResult, error)
}

type AgentDescription struct {
	GivenName      string
	PrimeDirective string
	Handler        AgentHandler
}

type TeamDescription struct {
	Name        string
	Manager     *AgentDescription
	TeamMembers []*AgentDescription
}

type TaskState struct {
	Description *TaskRequirements
}

type CompletionChecker interface {
	CheckCompletion(ctx context.Context, state *TaskState) (bool, error)
}

type StepDescription struct {
	Name              string
	Requirements      *TaskRequirements
	PointOfContact    *AgentDescription
	AssignedTeam      *TeamDescription
	CompletionChecker CompletionChecker
}

type TaskRequirements struct {
	Description string
}

type TaskResultAction int

const (
	TaskResultActionUnknown TaskResultAction = iota
	TaskResultActionForward
	TaskResultActionBackward
	TaskResultActionAbort
)

type TaskFeedback struct {
	WasSuccessful bool
	Rating        float32
	Feedback      *collective.Message
}

type TaskPayload struct {
	Raw string
}

type TaskResult struct {
	Payload  TaskPayload
	Step     *StepDescription
	Action   TaskResultAction
	Feedback *TaskFeedback
}

type Step interface {
	Forward(ctx context.Context, state *TaskState) (*TaskResult, error)
	Backward(ctx context.Context, state *TaskState, feedback *TaskFeedback) (*TaskResult, error)
}

type ProductionLine interface {
}
