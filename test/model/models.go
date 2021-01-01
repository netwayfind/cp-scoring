package model

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

// HostTokenRequest asdf
type HostTokenRequest struct {
	Hostname string
}

// HostTokenRegistration asdf
type HostTokenRegistration struct {
	HostToken string
	Scenario  string
	TeamKey   string
}

// Scenario asdf
type Scenario struct {
	ID          uint64
	Name        string
	Description string
	Enabled     bool
}

// ScenarioSummary asdf
type ScenarioSummary struct {
	ID      uint64
	Name    string
	Enabled bool
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

// TeamSummary asdf
type TeamSummary struct {
	ID      uint64
	Name    string
	Enabled bool
}
