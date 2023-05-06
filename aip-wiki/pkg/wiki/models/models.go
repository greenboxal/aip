package models

import "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"

func init() {
	forddb.DefineResourceType[PageID, *Page]("page")
	forddb.DefineResourceType[ImageID, *Image]("image")
}
