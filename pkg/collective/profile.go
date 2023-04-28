package collective

import (
	"encoding/json"

	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type ProfileID struct {
	forddb.StringResourceID[*Profile]
}

type Profile struct {
	forddb.ResourceMetadata[ProfileID, *Profile] `json:"metadata"`

	Spec json.RawMessage `json:"spec"`
}
