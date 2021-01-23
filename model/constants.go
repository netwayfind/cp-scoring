package model

// ActionType asdf
type ActionType string

// asdf
const (
	ActionTypeExec      ActionType = "EXEC"
	ActionTypeFileExist ActionType = "FILE_EXIST"
	ActionTypeFileRegex ActionType = "FILE_REGEX"
	ActionTypeFileValue ActionType = "FILE_VALUE"
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
	AuthCookieName        = "auth"
	KeyCharset     string = "0123456789ABCDEF"
	TeamCookieName        = "team"
)
