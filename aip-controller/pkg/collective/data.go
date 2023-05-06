package collective

import (
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type MemoryDataID struct {
	forddb.CidResourceID[*Memory]
}

type MemoryData struct {
	forddb.ResourceBase[MemoryDataID, *MemoryData] `json:"metadata"`

	Cid  cid.Cid `json:"cid"`
	Text string  `json:"data"`
}

func NewMemoryData(data string) MemoryData {
	return NewMemoryDataFromBytes([]byte(data))
}

func NewMemoryDataFromBytes(data []byte) MemoryData {
	h, err := multihash.Sum(data, multihash.SHA2_256, -1)

	if err != nil {
		panic(err)
	}

	return MemoryData{
		Cid:  cid.NewCidV1(cid.Raw, h),
		Text: string(data),
	}
}
