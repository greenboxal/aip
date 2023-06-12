package forddbimpl

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/multiformats/go-multihash"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/objectstore"
)

var basicLinkPrototype = cidlink.LinkPrototype{
	Prefix: cid.Prefix{
		Codec:    cid.DagCBOR,
		MhLength: -1,
		MhType:   multihash.SHA2_256,
		Version:  1,
	},
}

type Database struct {
	logStore    forddb.LogStore
	objectStore objectstore.ObjectStore
	linkSystem  linking.LinkSystem
}

var _ forddb.Database = (*Database)(nil)

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Subscribe(listener forddb.Listener) func() { panic("implement me") }
func (d *Database) LogStore() forddb.LogStore                 { return d.logStore }

func (d *Database) List(ctx context.Context, typ forddb.TypeID, options ...forddb.ListOption) ([]forddb.BasicResource, error) {
	opts := forddb.NewListOptions(typ, options...)
	opts = opts

	panic("implement me")
}

func (d *Database) Get(ctx context.Context, typ forddb.TypeID, id forddb.BasicResourceID, options ...forddb.GetOption) (forddb.BasicResource, error) {
	opts := forddb.NewGetOptions(typ, options...)
	opts = opts
	//TODO implement me
	panic("implement me")
}

func (d *Database) Put(ctx context.Context, resource forddb.BasicResource, options ...forddb.PutOption) (forddb.BasicResource, error) {
	opts := forddb.NewPutOptions(options...)

	opts = opts

	return resource, nil
}

func (d *Database) Delete(ctx context.Context, resource forddb.BasicResource, options ...forddb.DeleteOption) (forddb.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}
