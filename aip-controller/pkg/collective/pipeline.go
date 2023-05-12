package collective

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type PipelineID struct {
	forddb.StringResourceID[*Pipeline]
}

type StageID struct {
	forddb.StringResourceID[*Stage]
}

type PipelineSpec struct {
	Stages []Stage `json:"stages"`
}

type Pipeline struct {
	forddb.ResourceBase[PipelineID, *Pipeline] `json:"metadata"`

	Spec PipelineSpec `json:"spec"`
}

func (p *Pipeline) GetStage(id StageID) *Stage {
	for _, stage := range p.Spec.Stages {
		if stage.ID == id {
			return &stage
		}
	}

	return nil
}

type Stage struct {
	forddb.ResourceBase[StageID, *Stage] `json:"metadata"`

	ID           StageID   `json:"id"`
	AssignedTeam string    `json:"assigned_team"`
	DependsOn    []StageID `json:"depends_on"`
}
