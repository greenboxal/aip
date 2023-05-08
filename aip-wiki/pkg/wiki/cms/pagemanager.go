package cms

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/multiformats/go-multihash"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/jobs"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/generators"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

type PageManager struct {
	db forddb.Database

	fm    *FileManager
	jm    *jobs.Manager
	cache *generators.ContentCache

	pageGenerator  *generators.PageGenerator
	imageGenerator *generators.ImageGenerator
}

func NewPageManager(
	db forddb.Database,
	fm *FileManager,
	jm *jobs.Manager,
	cache *generators.ContentCache,
	pageGenerator *generators.PageGenerator,
	imageGenerator *generators.ImageGenerator,
) *PageManager {
	return &PageManager{
		db:             db,
		fm:             fm,
		jm:             jm,
		cache:          cache,
		pageGenerator:  pageGenerator,
		imageGenerator: imageGenerator,
	}
}

func (pm *PageManager) GetPageByID(ctx context.Context, id models.PageID) (*models.Page, error) {
	return forddb.Get[*models.Page](ctx, pm.db, id)
}

func (pm *PageManager) GetPage(ctx context.Context, spec models.PageSpec) (*models.Page, error) {
	page, err := pm.cache.GetPage(ctx, spec)

	if forddb.IsNotFound(err) {
		job, err := jobs.DispatchJob(
			ctx,
			pm.jm,
			models.GeneratePageJobHandlerID,
			spec,
		)

		if err != nil {
			return nil, err
		}

		return jobs.Await(job)
	} else if err != nil {
		return nil, err
	}

	return page, nil
}

func (pm *PageManager) GetImage(ctx context.Context, spec models.ImageSpec) (*models.Image, error) {
	page, err := pm.cache.GetImage(ctx, spec)

	if forddb.IsNotFound(err) {
		return pm.GenerateImage(ctx, spec)
	} else if err != nil {
		return nil, err
	}

	return page, nil
}

func (pm *PageManager) GenerateImage(ctx context.Context, spec models.ImageSpec) (*models.Image, error) {
	id := models.BuildImageID(spec)

	status, err := pm.imageGenerator.GetImage(ctx, spec)

	if err != nil {
		return nil, err
	}

	response, err := http.Get(status.URL)

	if err != nil {
		return nil, err
	}

	tempFileName := fmt.Sprintf("temp-%s-%d", id.String(), time.Now().UnixNano())

	writer := pm.fm.OpenWriter(ctx, tempFileName)
	reader := io.TeeReader(response.Body, writer)

	h, err := multihash.SumStream(reader, multihash.SHA2_256, -1)

	if err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	fileName := "images/" + h.B58String() + path.Ext(status.URL)

	if err := pm.fm.Rename(ctx, tempFileName, fileName); err != nil {
		return nil, err
	}

	status.URL = "https://cdn.desciclo.ai/" + fileName

	result := &models.Image{
		Spec:   spec,
		Status: status,
	}

	result.ID = id

	return pm.cache.PutImage(ctx, result)
}
