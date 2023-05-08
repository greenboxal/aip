package forddb

import (
	"context"

	"github.com/ipfs/go-cid"
)

type ContentAddressableStorage interface {
	GetBytes(ctx context.Context, id cid.Cid) ([]byte, error)
	PutBytes(ctx context.Context, data []byte) (cid.Cid, error)
}
