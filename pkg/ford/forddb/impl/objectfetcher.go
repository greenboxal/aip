package forddbimpl

import (
	"context"
	"runtime"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"

	"github.com/greenboxal/aip/pkg/ford/forddb"
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
	storage forddb.Storage
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
			if res, err := of.fetchResource(ctx, req.storage, req.slot.id); err != nil {
				req.slot.setError(err)
			} else {
				_, err = req.slot.Update(ctx, res, forddb.PutOptions{
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
	storage forddb.Storage,
	id forddb.BasicResourceID,
) (forddb.BasicResource, error) {
	return storage.Get(ctx, id.BasicResourceType().GetID(), id)
}
