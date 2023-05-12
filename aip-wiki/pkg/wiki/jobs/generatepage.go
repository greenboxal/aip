package jobs

import (
	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/memory"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/memoryctx"
	"github.com/greenboxal/aip/aip-sdk/pkg/jobs"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/generators"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki/models"
)

type GeneratePageJobHandler struct {
	pageGenerator *generators.PageGenerator
	db            forddb.Database
	msn           msn.API
}

func NewGeneratePageJobHandler(
	pg *generators.PageGenerator,
	db forddb.Database,
	msn msn.API,
) *GeneratePageJobHandler {
	return &GeneratePageJobHandler{
		db:            db,
		msn:           msn,
		pageGenerator: pg,
	}
}

func (pg *GeneratePageJobHandler) Run(ctx jobs.JobContext) error {
	spec, err := forddb.Convert[models.PageSpec](ctx.Payload().(forddb.RawResource))

	if err != nil {
		return err
	}

	id := models.BuildPageID(spec)

	mem := &memory.ChannelChatMemory{
		ContextKey: chat.ChatHistoryContextKey,
		Messenger:  pg.msn,
		Database:   pg.db,
		Endpoint:   forddb.NewStringID[msn.EndpointID]("bot:" + spec.Voice),
		Channel:    forddb.NewStringID[msn.ChannelID](id.String()),
	}

	mctx := memoryctx.WithMemory(ctx.Context(), mem)

	page, err := pg.pageGenerator.GeneratePage(mctx, spec)

	if err != nil {
		return err
	}

	ctx.SetResult(page)

	return nil
}
