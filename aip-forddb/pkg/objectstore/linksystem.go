package objectstore

import (
	"github.com/ipld/go-ipld-prime/linking"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-sdk/pkg/network/ipfs"
)

func NewLinkSystem(
	db forddb.Database,
	storage ObjectStore,
	ipfs *ipfs.Manager,
) linking.LinkSystem {
	base := ipfs.LinkSystem()
	lsys := base

	return lsys
}
