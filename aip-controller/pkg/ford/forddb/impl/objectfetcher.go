package forddbimpl

import (
	"context"
	"runtime"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"

	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
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
	storage forddb2.Storage
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
				_, err = req.slot.Update(ctx, res, forddb2.PutOptions{
					OnConflict: forddb2.OnConflictLatestWins,
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
	storage forddb2.Storage,
	id forddb2.BasicResourceID,
) (forddb2.BasicResource, error) {
	return storage.Get(ctx, id.BasicResourceType().GetID(), id)
}