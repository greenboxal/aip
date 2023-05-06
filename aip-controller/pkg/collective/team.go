package collective

import (
	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type TeamID struct {
	forddb2.StringResourceID[*Team]
}

type Team struct {
	forddb2.ResourceBase[TeamID, *Team] `json:"metadata"`

	Spec TeamSpec `json:"spec"`
}
type TeamSpec struct {
	Manager string   `json:"manager"`
	Members []string `json:"members"`
}
