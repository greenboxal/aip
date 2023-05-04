package collective

import (
	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type ProfileID struct {
	forddb2.StringResourceID[*Profile]
}

type Profile struct {
	forddb2.ResourceMetadata[ProfileID, *Profile] `json:"metadata"`

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
