package tot

import "github.com/greenboxal/aip/aip-forddb/pkg/forddb"

type NodeID struct {
	forddb.StringResourceID[*Node] `ipld:",inline"`
}

type Node struct {
	forddb.ResourceBase[NodeID, *Node] `json:"metadata"`

	ParentNodeID NodeID `json:"parent_node_id"`
	Height       int    `json:"height"`
	Contents     string `json:"contents"`
	TokenCount   int    `json:"token_count"`
}

func (n *Node) Fork() *Node {
	return &Node{
		ParentNodeID: n.ID,
		Height:       n.Height + 1,
	}
}
