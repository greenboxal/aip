package cli

import (
	"github.com/ipfs/go-unixfsnode/data/builder"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/storage/memstore"
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

	_, _, err := builder.BuildUnixFSRecursive(args[0], &ls)

	if err != nil {
		return err
	}

	return nil
}
