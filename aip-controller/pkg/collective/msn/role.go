package msn

type Role string

const (
	RoleInvalid  Role = ""
	RoleSystem   Role = "System"
	RoleUser     Role = "User"
	RoleAI       Role = "AI"
	RoleFunction Role = "Function"
)
