package forddbimpl

import (
	"context"
	"runtime"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/objectstore"
)

type objectFetcher struct {
	db *database

	requestCh chan fetchResourceRequest
}

func newObjectFetcher(db *database) *objectFetcher {
	return &objectFetcher{
		db: db,

		requestCh: make(chan fetchResourceRequest, 1024),
	}
}

type fetchResourceRequest struct {
	storage objectstore.ObjectStore
	slot    *resourceSlot
}

func (of *objectFetcher) Run(proc goprocess.Process) {
	for i := 0; i < runtime.NumCPU(); i++ {
		proc.Go(of.runWorker)
	}

	if err := proc.CloseAfterChildren(); err != nil {
		panic(err)
	}
}

func (of *objectFetcher) runWorker(proc goprocess.Process) {
	ctx := goprocessctx.OnClosingContext(proc)

	for {
		select {
		case <-proc.Closing():
			return

		case req := <-of.requestCh:
			if raw, err := of.fetchResource(ctx, req.storage, req.slot.id); err != nil {
				req.slot.setError(err)
			} else {
				_, err = req.slot.Update(ctx, raw, forddb.PutOptions{
					OnConflict: forddb.OnConflictLatestWins,
				})

				if err != nil {
					req.slot.setError(err)
				}
			}
		}
	}
}

func (of *objectFetcher) fetchResource(
	ctx context.Context,
	storage objectstore.ObjectStore,
	id forddb.BasicResourceID,
) (forddb.RawResource, error) {
	return storage.Get(ctx, id.BasicResourceType().GetResourceID(), id, forddb.GetOptions{})
}
