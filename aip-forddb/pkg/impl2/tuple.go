package forddbimpl

import (
	"context"
	"sync"

	"github.com/ipld/go-ipld-prime"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type Tuple struct {
	m sync.RWMutex

	partition *Snapshot

	typ forddb.BasicResourceType
	id  forddb.BasicResourceID

	currentLink  ipld.Link
	currentValue ipld.Node
}

func (t *Tuple) Refresh(ctx context.Context) (forddb.BasicResource, error) {
	panic("implement me")
}

func (t *Tuple) Get(ctx context.Context, options *forddb.GetOptions) (forddb.BasicResource, error) {
	panic("implement me")
}

func (t *Tuple) Update(ctx context.Context, resource forddb.BasicResource, options *forddb.PutOptions) (forddb.BasicResource, error) {
	/*link, err := t.partition.storeLink(ctx, resource)

	if err != nil {
		return nil, err
	}

	meta := resource.GetResourceMetadata()
	meta.*/
	panic("implement me")
}

func (t *Tuple) Delete(ctx context.Context, resource forddb.BasicResource, options *forddb.DeleteOptions) (forddb.BasicResource, error) {
	panic("implement me")
}
