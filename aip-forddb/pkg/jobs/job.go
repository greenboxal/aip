package jobs

import (
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type JobID struct {
	forddb.StringResourceID[*Job] `ipld:",inline"`
}

type Job struct {
	forddb.ResourceBase[JobID, *Job] `json:"metadata"`

	Spec   JobSpec   `json:"spec"`
	Status JobStatus `json:"status"`
}

type JobState string

const (
	JobStateInvalid   JobState = ""
	JobStateCreated   JobState = "CREATED"
	JobStateScheduled JobState = "SCHEDULED"
	JobStatePending   JobState = "PENDING"
	JobStateRunning   JobState = "RUNNING"
	JobStateFailed    JobState = "FAILED"
	JobStateCompleted JobState = "COMPLETED"
)

type JobSpec struct {
	Payload any    `json:"payload"`
	Handler string `json:"handler"`
}

type JobStatus struct {
	State     JobState `json:"state"`
	Progress  any      `json:"progress"`
	Result    any      `json:"result"`
	LastError string   `json:"last_error"`
}

func init() {
	forddb.DefineResourceType[JobID, *Job]("job")
}
