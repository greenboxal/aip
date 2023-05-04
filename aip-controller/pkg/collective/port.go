package collective

import (
	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type PortID struct {
	forddb2.StringResourceID[*Port]
}

type PortBindingID struct {
	forddb2.StringResourceID[*PortBinding]
}

type Port struct {
	forddb2.ResourceMetadata[PortID, *Port] `json:"metadata"`

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
	forddb2.ResourceMetadata[PortBindingID, *PortBinding] `json:"metadata"`

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
