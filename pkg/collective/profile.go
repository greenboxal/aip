package collective

import (
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type ProfileID struct {
	forddb.StringResourceID[*Profile]
}

type Profile struct {
	forddb.ResourceMetadata[ProfileID, *Profile] `json:"metadata"`

	Spec ProfileSpec `json:"spec"`

	// FIXME: Move to spec
	Aptitudes []ProfileAptitude `json:"aptitudes"`
}

type ProfileSpec struct {
	Directive string `json:"directive"`
}

type ProfileAptitude struct {
	Description string `json:"description"`
}
