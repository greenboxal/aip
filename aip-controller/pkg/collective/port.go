package collective

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type PortID struct {
	forddb.StringResourceID[*Port] `ipld:",inline"`
}

type PortBindingID struct {
	forddb.StringResourceID[*PortBinding] `ipld:",inline"`
}

type Port struct {
	forddb.ResourceBase[PortID, *Port] `json:"metadata"`

	Spec   PortSpec   `json:"spec"`
	Status PortStatus `json:"status"`
}

type PortSpec struct {
	Empty bool `json:"empty"`
}

type PortState string

const (
	PortCreated PortState = "CREATED"
	PortReady   PortState = "READY"
)

type PortStatus struct {
	State     PortState `json:"state"`
	LastError string    `json:"last_error"`
}

type PortBinding struct {
	forddb.ResourceBase[PortBindingID, *PortBinding] `json:"metadata"`

	Spec   PortBindingSpec   `json:"spec"`
	Status PortBindingStatus `json:"status"`
}

type PortBindingSpec struct {
	AgentID AgentID `json:"agent_id"`
	PortID  PortID  `json:"port_id"`
}

type PortBindingState string

const (
	PortBindingCreated PortBindingState = "CREATED"
	PortBindingReady   PortBindingState = "READY"
	PortBindingFailed  PortBindingState = "FAILED"
)

type PortBindingStatus struct {
	State     PortBindingState `json:"state"`
	LastError string           `json:"last_error"`
}
