package collective

import "github.com/greenboxal/aip/pkg/ford/forddb"

type ProfileID struct {
	forddb.StringResourceID[*Profile]
}

type PipelineID struct {
	forddb.StringResourceID[*Pipeline]
}

type StageID struct {
	forddb.StringResourceID[*Stage]
}

type TeamID struct {
	forddb.StringResourceID[*Team]
}

type AgentID struct {
	forddb.StringResourceID[*Agent]
}

type TaskID struct {
	forddb.StringResourceID[*Task]
}

type PortID struct {
	forddb.StringResourceID[*Port]
}

type PipelineSpec struct {
	Agents []Agent `json:"agents"`
	Teams  []Team  `json:"teams"`
	Stages []Stage `json:"stages"`
}

type Pipeline struct {
	forddb.ResourceMetadata[PipelineID, *Pipeline] `json:"metadata"`

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

type Team struct {
	forddb.ResourceMetadata[TeamID, *Team] `json:"metadata"`

	Manager string   `json:"manager"`
	Members []string `json:"members"`
}

type Stage struct {
	forddb.ResourceMetadata[StageID, *Stage] `json:"metadata"`

	AssignedTeam string    `json:"assigned_team"`
	DependsOn    []StageID `json:"depends_on"`
}

type Task struct {
	forddb.ResourceMetadata[TaskID, *Task] `json:"metadata"`

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

type Port struct {
	forddb.ResourceMetadata[PortID, *Port] `json:"metadata"`

	Spec   PortSpec   `json:"spec"`
	Status PortStatus `json:"status"`
}

type PortSpec struct {
}

type PortState string

const (
	PortCreated PortState = "CREATED"
	PortReady   PortState = "READY"
)

type PortStatus struct {
	State PortState `json:"state"`
}

type Agent struct {
	forddb.ResourceMetadata[AgentID, *Agent] `json:"metadata"`

	Spec   AgentSpec   `json:"spec"`
	Status AgentStatus `json:"status"`
}

type AgentSpec struct {
	GivenName string `json:"given_name"`

	ProfileID ProfileID `json:"profile_id"`
	PortID    string    `json:"port_id"`
}

type AgentState string

const (
	AgentStateCreated   AgentState = "CREATED"
	AgentStatePending   AgentState = "PENDING"
	AgentStateScheduled AgentState = "SCHEDULED"
	AgentStateReady     AgentState = "READY"
	AgentStateFailed    AgentState = "FAILED"
	AgentStateCompleted AgentState = "COMPLETED"
)

type AgentStatus struct {
	State     AgentState `json:"state"`
	LastError string     `json:"last_error"`
}

type Profile struct {
	forddb.ResourceMetadata[ProfileID, *Profile] `json:"metadata"`
}

func init() {
	forddb.DefineResourceType[PipelineID, *Pipeline]("pipeline")
	forddb.DefineResourceType[TaskID, *Task]("task")
	forddb.DefineResourceType[AgentID, *Agent]("agent")
	forddb.DefineResourceType[TeamID, *Team]("team")
	forddb.DefineResourceType[StageID, *Stage]("stage")
	forddb.DefineResourceType[PortID, *Port]("port")
	forddb.DefineResourceType[ProfileID, *Profile]("profile")
}
