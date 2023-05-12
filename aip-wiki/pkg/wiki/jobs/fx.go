package jobs

import (
	"go.uber.org/fx"

	"github.com/greenboxal/aip/aip-sdk/pkg/jobs"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

var Module = fx.Module(
	"wiki/jobs",

	jobs.ProvideJobHandler[*GeneratePageJobHandler](models.GeneratePageJobHandlerID, NewGeneratePageJobHandler),
)
