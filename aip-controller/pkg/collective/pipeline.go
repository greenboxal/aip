package collective

import (
	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type PipelineID struct {
	forddb2.StringResourceID[*Pipeline]
}

type StageID struct {
	forddb2.StringResourceID[*Stage]
}

type PipelineSpec struct {
	Stages []Stage `json:"stages"`
}

type Pipeline struct {
	forddb2.ResourceMetadata[PipelineID, *Pipeline] `json:"metadata"`

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
	forddb2.ResourceMetadata[StageID, *Stage] `json:"metadata"`

	ID           StageID   `json:"id"`
	AssignedTeam string    `json:"assigned_team"`
	DependsOn    []StageID `json:"depends_on"`
}
