package msn

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type EndpointID struct {
	forddb.StringResourceID[*Endpoint]
}

type Endpoint struct {
	forddb.ResourceBase[EndpointID, *Endpoint] `json:"metadata"`

	Subscriptions []ChannelID `json:"subscriptions"`
}
