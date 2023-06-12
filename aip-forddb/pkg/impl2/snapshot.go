package forddbimpl

import (
	"context"
	"fmt"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/linking"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type Snapshot struct {
	db *Database

	linkSystem linking.LinkSystem

	parentSnapshots []ObjectID
}

func (p *Snapshot) List(ctx context.Context, typ forddb.TypeID, options *forddb.ListOptions) ([]forddb.BasicResource, error) {
	panic("implement me")
}

func (p *Snapshot) loadLink(ctx context.Context, typ forddb.BasicResourceType, link ipld.Link) (forddb.BasicResource, error) {
	lctx := linking.LinkContext{Ctx: ctx}
	node, err := p.linkSystem.Load(lctx, link, typ.ActualType().IpldPrototype())

	if err != nil {
		return nil, err
	}

	resource, ok := typesystem.TryUnwrap[forddb.BasicResource](node)

	if !ok {
		return nil, fmt.Errorf("expected node to be a BasicResource, got %T", node)
	}

	return resource, nil
}

func (p *Snapshot) computeLink(resource forddb.BasicResource) (ipld.Link, error) {
	node := typesystem.Wrap(resource)
	link, err := p.linkSystem.ComputeLink(basicLinkPrototype, node)

	if err != nil {
		return nil, err
	}

	return link, nil
}

func (p *Snapshot) storeLink(ctx context.Context, resource forddb.BasicResource) (ipld.Link, error) {
	lctx := linking.LinkContext{Ctx: ctx}
	node := typesystem.Wrap(resource)

	link, err := p.linkSystem.Store(lctx, basicLinkPrototype, node)

	if err != nil {
		return nil, err
	}

	return link, nil
}
