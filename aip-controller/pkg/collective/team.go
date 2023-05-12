package collective

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type TeamID struct {
	forddb.StringResourceID[*Team]
}

type Team struct {
	forddb.ResourceBase[TeamID, *Team] `json:"metadata"`

	Spec TeamSpec `json:"spec"`
}
type TeamSpec struct {
	Manager string   `json:"manager"`
	Members []string `json:"members"`
}
