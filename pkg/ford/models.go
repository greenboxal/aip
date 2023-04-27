package ford

type ProfileID string
type PipelineID string
type StageID string
type TeamID string
type AgentID string
type TaskID string

type PipelineSpec struct {
	PipelineMetadata `json:"pipeline"`

	Agents []Agent `json:"agents"`
	Teams  []Team  `json:"teams"`
	Stages []Stage `json:"stages"`
}

type PipelineMetadata struct {
	Name string `json:"name"`
}

type Pipeline struct {
	ResourceBase[PipelineID, Pipeline]

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

type Agent struct {
	ResourceBase[AgentID, Agent]

	Name      string    `json:"name"`
	ProfileID ProfileID `json:"profile_id"`
}

type Team struct {
	ResourceBase[TeamID, Team]

	Name    string   `json:"name"`
	Manager string   `json:"manager"`
	Members []string `json:"members"`
}

type Stage struct {
	ResourceBase[StageID, Stage]

	Name         string   `json:"name"`
	AssignedTeam string   `json:"assigned_team"`
	DependsOn    []string `json:"depends_on"`
}

type Task struct {
	ResourceBase[TeamID, Team]

	Spec   TaskSpec   `json:"spec"`
	Status TaskStatus `json:"status"`
}

type TaskMetadata struct {
	Name string `json:"name"`
}

type TaskSpec struct {
	TaskMetadata `json:"task"`

	PipelineID    PipelineID `json:"pipeline_id"`
	OutputStageID StageID    `json:"output_stage_id"`

	Description string `json:"description"`
}

type TaskStatus struct {
	Phases []TaskPhaseStatus `json:"phases"`
	Phase  string            `json:"phase"`
}

type TaskPhaseState string

const (
	TaskPhaseStatePending    TaskPhaseState = "pending"
	TaskPhaseStateScheduled  TaskPhaseState = "scheduled"
	TaskPhaseStateInProgress TaskPhaseState = "in_progress"
	TaskPhaseStateCompleted  TaskPhaseState = "completed"
)

type TaskPhaseStatus struct {
	Name    string         `json:"name"`
	State   TaskPhaseState `json:"state"`
	StageID StageID        `json:"stage_id"`
}
