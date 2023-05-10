package models

import (
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/jobs"
)

func init() {
	forddb.DefineResourceType[DomainID, *Domain]("domain")
	forddb.DefineResourceType[LayoutID, *Layout]("layout")
	forddb.DefineResourceType[PageID, *Page]("page")
	forddb.DefineResourceType[ImageID, *Image]("image")
}

var GeneratePageJobHandlerID = jobs.DefineHandler[PageSpec, *Page]("wiki.generate-page")
