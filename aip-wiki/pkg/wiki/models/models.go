package models

import (
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/jobs"
)

var GeneratePageJobHandlerID = jobs.DefineHandler[PageSpec, *Page]("wiki.generate-page")

func init() {
	forddb.DefineResourceType[PageID, *Page]("page")
	forddb.DefineResourceType[ImageID, *Image]("image")
}
