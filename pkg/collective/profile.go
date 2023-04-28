package collective

import "github.com/greenboxal/aip/pkg/ford/forddb"

type ProfileID struct {
	forddb.StringResourceID[*Profile]
}

type Profile struct {
	forddb.ResourceMetadata[ProfileID, *Profile] `json:"metadata"`
}
