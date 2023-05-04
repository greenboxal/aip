package collective

import (
	forddb2 "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type TaskID struct {
	forddb2.StringResourceID[*Task]
}

type Task struct {
	forddb2.ResourceMetadata[TaskID, *Task] `json:"metadata"`

	Spec   TaskSpec   `json:"spec"`
	Status TaskStatus `json:"status"`
}

type TaskSpec struct {
	PipelineID    PipelineID `json:"pipeline_id"`
	OutputStageID StageID    `json:"output_stage_id"`

	Description string `json:"description"`
}

type TaskState string

const (
	TaskStateCreated    TaskState = "created"
	TaskStatePending    TaskState = "pending"
	TaskStateScheduled  TaskState = "scheduled"
	TaskStateInProgress TaskState = "in_progress"
	TaskStateCompleted  TaskState = "completed"
)

type TaskStatus struct {
	Phases []TaskPhaseStatus `json:"phases"`
	Phase  string            `json:"phase"`
	State  TaskState         `json:"state"`
}

func (s *TaskStatus) GetTaskStatus(id StageID) *TaskPhaseStatus {
	for i, phase := range s.Phases {
		if phase.StageID == id {
			return &s.Phases[i]
		}
	}

	return nil
}

func (s *TaskStatus) GetOrCreateTaskStatus(id StageID) *TaskPhaseStatus {
	for i, phase := range s.Phases {
		if phase.StageID == id {
			return &s.Phases[i]
		}
	}

	status := TaskPhaseStatus{
		StageID: id,
		State:   TaskPhaseStateCreated,
	}

	index := len(s.Phases)
	s.Phases = append(s.Phases, status)

	return &s.Phases[index]
}

type TaskPhaseState string

const (
	TaskPhaseStateCreated    TaskPhaseState = "CREATED"
	TaskPhaseStatePending    TaskPhaseState = "PENDING"
	TaskPhaseStateScheduled  TaskPhaseState = "SCHEDULED"
	TaskPhaseStateInProgress TaskPhaseState = "IN_PROGRESS"
	TaskPhaseStateCompleted  TaskPhaseState = "COMPLETED"
)

type TaskPhaseStatus struct {
	State   TaskPhaseState `json:"state"`
	StageID StageID        `json:"stage_id"`
}
