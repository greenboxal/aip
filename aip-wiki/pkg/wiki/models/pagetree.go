package models

import "github.com/greenboxal/aip/aip-forddb/pkg/forddb"

type PageTreeID struct {
	forddb.StringResourceID[*PageTree] `ipld:",inline"`
}

type PageTree struct {
	forddb.ContentAddressedResourceBase[PageTreeID, *PageTree] `json:"metadata"`

	Spec   PageTreeSpec   `json:"spec"`
	Status PageTreeStatus `json:"status"`
}

type PageTreeSpec struct {
	BasePageID PageID `json:"base_page_id"`
}

type PageTreeStatus struct {
	Root *PageTreeNode `json:"root"`
}

type PageTreeNode struct {
}

func (p *PageTree) GetContentAddressableRoot() any {
	return p.Spec
}
