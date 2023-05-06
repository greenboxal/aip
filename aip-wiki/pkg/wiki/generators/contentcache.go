package generators

import (
	"context"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

type ContentCache struct {
	db forddb.Database
}

func NewContentCache(
	db forddb.Database,
) *ContentCache {
	return &ContentCache{
		db: db,
	}
}

func (pm *ContentCache) GetPageByID(ctx context.Context, id models.PageID) (*models.Page, error) {
	return forddb.Get[*models.Page](ctx, pm.db, id)
}

func (pm *ContentCache) GetPage(ctx context.Context, spec models.PageSpec) (*models.Page, error) {
	id := models.BuildPageID(spec)

	return forddb.Get[*models.Page](ctx, pm.db, id)
}

func (pm *ContentCache) PutPage(ctx context.Context, page *models.Page) (*models.Page, error) {
	id := models.BuildPageID(page.Spec)
	page.ResourceBase.ID = id

	return forddb.Put(pm.db, page)
}

func (pm *ContentCache) GetImage(ctx context.Context, spec models.ImageSpec) (*models.Image, error) {
	id := models.BuildImageID(spec)

	return forddb.Get[*models.Image](ctx, pm.db, id)
}

func (pm *ContentCache) PutImage(ctx context.Context, image *models.Image) (*models.Image, error) {
	id := models.BuildImageID(image.Spec)
	image.ID = id

	return forddb.Put(pm.db, image)
}
