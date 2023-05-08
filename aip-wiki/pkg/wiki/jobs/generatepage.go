package jobs

import (
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/jobs"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/generators"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

type GeneratePageJobHandler struct {
	pageGenerator *generators.PageGenerator
}

func NewGeneratePageJobHandler(pg *generators.PageGenerator) *GeneratePageJobHandler {
	return &GeneratePageJobHandler{pageGenerator: pg}
}

func (pg *GeneratePageJobHandler) Run(ctx jobs.JobContext) error {
	spec, err := forddb.Convert[models.PageSpec](ctx.Payload().(forddb.RawResource))

	if err != nil {
		return err
	}

	page, err := pg.pageGenerator.GeneratePage(ctx.Context(), spec)

	if err != nil {
		return err
	}

	ctx.SetResult(page)

	return nil
}
