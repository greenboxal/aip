package cli

import (
	"github.com/ipfs/go-unixfsnode/data/builder"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	"github.com/ipld/go-ipld-prime/traversal"
	builder2 "github.com/ipld/go-ipld-prime/traversal/selector/builder"
	"github.com/spf13/cobra"
)

type IndexAddCommand struct {
}

func NewIndexAddCommand() *IndexAddCommand {
	return &IndexAddCommand{}
}

func (i *IndexAddCommand) Run(cmd *cobra.Command, args []string) error {
	ls := cidlink.DefaultLinkSystem()
	store := memstore.Store{Bag: make(map[string][]byte)}
	ls.SetReadStorage(&store)
	ls.SetWriteStorage(&store)

	rootLink, _, err := builder.BuildUnixFSRecursive(args[0], &ls)

	if err != nil {
		return err
	}

	lctx := linking.LinkContext{Ctx: cmd.Context()}
	rootNode, err := ls.Load(lctx, rootLink, basicnode.Prototype.Any)

	if err != nil {
		return err
	}

	p := &traversal.Progress{
		Cfg: &traversal.Config{
			Ctx: cmd.Context(),

			LinkSystem:        ls,
			LinkVisitOnlyOnce: true,
		},
	}

	sb := builder2.NewSelectorSpecBuilder(basicnode.Prototype.Any)
	sb.ExploreAll(sb.ExploreRecursive())

	err = p.WalkAdv(rootNode, s, func(progress traversal.Progress, node datamodel.Node, reason traversal.VisitReason) error {

	})

	if err != nil {
		return err
	}

	return nil
}
