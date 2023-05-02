package collective

import "github.com/greenboxal/aip/pkg/ford/forddb"

type TeamID = forddb.StringResourceID[*Team]

type Team struct {
	forddb.ResourceMetadata[TeamID, *Team] `json:"metadata"`

	Spec TeamSpec `json:"spec"`
}
type TeamSpec struct {
	Manager string   `json:"manager"`
	Members []string `json:"members"`
}
