package collective

import "github.com/greenboxal/aip/pkg/ford/forddb"

func init() {
	forddb.DefineResourceType[PortID, *Port]("port")
	forddb.DefineResourceType[AgentID, *Agent]("agent")
	forddb.DefineResourceType[ProfileID, *Profile]("profile")
	forddb.DefineResourceType[TeamID, *Team]("team")
	forddb.DefineResourceType[PipelineID, *Pipeline]("pipeline")
	forddb.DefineResourceType[TaskID, *Task]("task")
	forddb.DefineResourceType[StageID, *Stage]("stage")
}
