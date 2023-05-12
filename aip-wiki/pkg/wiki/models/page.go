package models

import (
	"encoding/json"

	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multihash"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type PageID struct {
	forddb.StringResourceID[*Page]
}

type Page struct {
	forddb.ContentAddressedResourceBase[PageID, *Page] `json:"metadata"`

	Spec   PageSpec   `json:"spec"`
	Status PageStatus `json:"status"`
}

func (p *Page) GetContentAddressableRoot() any {
	return p.Spec
}

type PageLink struct {
	Title string `json:"title"`
	To    string `json:"to"`
}

type PageImage struct {
	Title  string `json:"title"`
	Source string `json:"source"`
}

type PageStatus struct {
	Markdown string `json:"markdown"`
	HTML     string `json:"html"`

	Links  []PageLink  `json:"links"`
	Images []PageImage `json:"images"`
}

type PageSpec struct {
	Format      string `json:"format"`
	Layout      string `json:"layout"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Voice       string `json:"voice"`
	Language    string `json:"language"`
	BasePage    PageID `json:"base_page_id,omitempty"`
}

func BuildPageID(spec PageSpec) PageID {
	data, err := json.Marshal(spec)

	if err != nil {
		panic(err)
	}

	h, err := multihash.Sum(data, multihash.SHA1, -1)

	if err != nil {
		panic(err)
	}

	b, err := multibase.Encode(multibase.Base36, h)

	if err != nil {
		panic(err)
	}

	return forddb.NewStringID[PageID](b)
}
