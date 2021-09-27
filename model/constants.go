package model

// ActionType asdf
type ActionType string

// asdf
const (
	ActionTypeExec         ActionType = "EXEC"
	ActionTypeFileContains ActionType = "FILE_CONTAINS"
	ActionTypeFileExist    ActionType = "FILE_EXIST"
	ActionTypeFileRegex    ActionType = "FILE_REGEX"
	ActionTypeFileValue    ActionType = "FILE_VALUE"
)

// AuditQueueStatus asdf
type AuditQueueStatus string

// asdf
const (
	AuditQueueStatusReceived AuditQueueStatus = "RECV"
	AuditQueueStatusFailed   AuditQueueStatus = "FAIL"
)

// OperatorType asdf
type OperatorType string

// asdf
const (
	OperatorTypeEqual    OperatorType = "EQUAL"
	OperatorTypeNotEqual OperatorType = "NOT_EQUAL"
)

// Role asdf
type Role string

// asdf
const (
	RoleAdmin Role = "Admin"
	RoleTeam  Role = "Team"
)

// asdf
const (
	AuthCookieName       = "auth"
	JavascriptDateFormat = "Mon, 02 Jan 2006 15:04:05 MST"
	KeyCharset           = "0123456789ABCDEF"
	TeamCookieName       = "team"
)
