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
	OperatorTypeNil      OperatorType = "NIL"
	OperatorTypeNotNil   OperatorType = "NOT_NIL"
)

// asdf
const (
	KeyCharset string = "0123456789ABCDEF"
)

// Action asdf
type Action struct {
	Type    ActionType
	Command string
	Args    []string
}

// Answer asdf
type Answer struct {
	Operator OperatorType
	Value    interface{}
}

// HostRegistration asdf
type HostRegistration struct {
	Scenario  string
	HostToken string
	TeamKey   string
}

// Scenario asdf
type Scenario struct {
	ID          uint64
	Name        string
	Description string
	Enabled     bool
}

// ScenarioHostResult asdf
type ScenarioHostResult struct {
	HostToken string
	Timestamp int64
	Findings  []string
}

// Team asdf
type Team struct {
	ID      uint64
	Name    string
	POC     string
	Email   string
	Enabled bool
	Key     string
}
