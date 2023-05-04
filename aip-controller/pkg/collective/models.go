package collective

import (
	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

func init() {
	forddb.DefineResourceType[PortID, *Port]("port")
	forddb.DefineResourceType[ProfileID, *Profile]("profile")
	forddb.DefineResourceType[AgentID, *Agent]("agent")
	forddb.DefineResourceType[PortBindingID, *PortBinding]("port_binding")

	forddb.DefineResourceType[PipelineID, *Pipeline]("pipeline")
	forddb.DefineResourceType[StageID, *Stage]("stage")
	forddb.DefineResourceType[TeamID, *Team]("team")
	forddb.DefineResourceType[TaskID, *Task]("task")

	forddb.DefineResourceType[MemoryDataID, *MemoryData]("memory_data")
	forddb.DefineResourceType[MemoryID, *Memory]("memory")
	forddb.DefineResourceType[MemorySegmentID, *MemorySegment]("memory_segment")
}
