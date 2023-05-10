package msn

import "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"

type ChannelID struct {
	forddb.StringResourceID[*Channel]
}

type Channel struct {
	forddb.ResourceBase[ChannelID, *Channel] `json:"metadata"`

	Subscribers []EndpointID `json:"subscribers"`
}
