package collective

import "github.com/greenboxal/aip/pkg/ford/forddb"

type PortID struct {
	forddb.StringResourceID[*Port]
}

type Port struct {
	forddb.ResourceMetadata[PortID, *Port] `json:"metadata"`

	Spec   PortSpec   `json:"spec"`
	Status PortStatus `json:"status"`
}

type PortSpec struct {
}

type PortState string

const (
	PortCreated PortState = "CREATED"
	PortReady   PortState = "READY"
)

type PortStatus struct {
	State PortState `json:"state"`
}
