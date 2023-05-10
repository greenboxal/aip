package forddb

import (
	"context"
	"encoding/json"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multihash"
)

type ContentAddressableStorage interface {
	GetBytes(ctx context.Context, id cid.Cid) ([]byte, error)
	PutBytes(ctx context.Context, data []byte) (cid.Cid, error)
}
type BasicContentAddressedResource interface {
	GetContentAddressableRoot() any
}

type ContentAddressedResource[ID BasicResourceID] interface {
	Resource[ID]

	BasicContentAddressedResource
}

type ContentAddressedResourceBase[ID ResourceID[T], T ContentAddressedResource[ID]] struct {
	ResourceBase[ID, T]
}

func (r *ContentAddressedResourceBase[ID, T]) OnBeforeSave(self BasicResource) {
	car := self.(BasicContentAddressedResource)

	r.ID = CreateContentAddressableID[ID](car.GetContentAddressableRoot())

	r.ResourceBase.OnBeforeSave(self)
}

func CreateContentAddressableID[ID BasicResourceID](spec any) ID {
	data, err := json.Marshal(spec)

	if err != nil {
		panic(err)
	}

	h, err := multihash.Sum(data, multihash.SHA1, -1)

	if err != nil {
		panic(err)
	}

	b, err := multibase.Encode(multibase.Base36, h)

	if err != nil {
		panic(err)
	}

	return NewStringID[ID](b)
}
