package forddb

type ResourceScope string

const (
	ResourceScopeInvalid ResourceScope = ""
	// ResourceScopeGlobal resources are persisted to the global database/log and are persistent.
	ResourceScopeGlobal ResourceScope = "GLOBAL"
	// ResourceScopeLocal resources are persisted to the local database/log and are persistent.
	ResourceScopeLocal ResourceScope = "LOCAL"
	// ResourceScopeEphemeral are in-memory ephemeral resources that are not persisted to the database.
	ResourceScopeEphemeral ResourceScope = "EPHEMERAL"
)

type SortOrder string

const (
	Asc  SortOrder = "ASC"
	Desc SortOrder = "DESC"
)

//go:generate go run github.com/dmarkham/enumer -type=
type OnConflict int

const (
	OnConflictError OnConflict = iota
	OnConflictOptimistic
	OnConflictLatestWins
	OnConflictReplace
)

type ReadConsistencyLevel int

const (
	ReadConsistencyNoCache ReadConsistencyLevel = iota
	ReadConsistencyCachedStrong
	ReadConsistencyCachedDirty
)

type Kind int

const (
	KindInvalid Kind = iota
	KindId
	KindResource
	KindValue
	KindPointer
)
