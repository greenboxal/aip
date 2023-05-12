package msn

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type ChannelID struct {
	forddb.StringResourceID[*Channel]
}

type Channel struct {
	forddb.ResourceBase[ChannelID, *Channel] `json:"metadata"`

	Subscribers []EndpointID `json:"subscribers"`
}
