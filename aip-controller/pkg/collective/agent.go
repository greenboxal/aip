package collective

import (
	forddb "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type AgentID struct {
	forddb.StringResourceID[*Agent]
}

type Agent struct {
	forddb.ResourceBase[AgentID, *Agent] `json:"metadata"`

	Spec   AgentSpec   `json:"spec"`
	Status AgentStatus `json:"status"`
}

type AgentSpec struct {
	GivenName string `json:"given_name"`

	ProfileID ProfileID `json:"profile_id"`
	PortID    string    `json:"port_id"`
	ExtraArgs []string  `json:"extra_args"`
}

type AgentState string

const (
	AgentStateCreated   AgentState = "CREATED"
	AgentStatePending   AgentState = "PENDING"
	AgentStateScheduled AgentState = "SCHEDULED"
	AgentStateReady     AgentState = "READY"
	AgentStateFailed    AgentState = "FAILED"
	AgentStateCompleted AgentState = "COMPLETED"
)

type AgentStatus struct {
	State     AgentState `json:"state"`
	LastError string     `json:"last_error"`
}
