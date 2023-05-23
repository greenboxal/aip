package models

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-forddb/pkg/jobs"
)

var DomainType = forddb.DefineResourceType[DomainID, *Domain]("domain")
var RouteBindingType = forddb.DefineResourceType[RouteBindingID, *RouteBinding]("route_binding")
var LayoutType = forddb.DefineResourceType[LayoutID, *Layout]("layout")
var PageType = forddb.DefineResourceType[PageID, *Page]("page")
var ImageType = forddb.DefineResourceType[ImageID, *Image]("image")

var GeneratePageJobHandlerID = jobs.DefineHandler[PageSpec, *Page]("wiki.generate-page")
