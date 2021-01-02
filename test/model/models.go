package model

// Action asdf
type Action struct {
	Type    ActionType
	Command string
	Args    []string
}

// Answer asdf
type Answer struct {
	Operator    OperatorType
	Value       interface{}
	Description string
	Points      int
}

// AuditAnswerResults asdf
type AuditAnswerResults struct {
	ScenarioID     uint64
	TeamID         uint64
	HostToken      string
	Timestamp      int64
	CheckResultsID uint64
	Score          int
	AnswerResults  []bool
}

// AuditCheckResults asdf
type AuditCheckResults struct {
	ScenarioID   uint64
	HostToken    string
	Timestamp    int64
	CheckResults []string
}

// HostTokenRequest asdf
type HostTokenRequest struct {
	Hostname string
}

// HostTokenRegistration asdf
type HostTokenRegistration struct {
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

// ScenarioScore asdf
type ScenarioScore struct {
	TeamName  string
	Hostname  string
	Score     int
	Timestamp int64
}

// ScenarioSummary asdf
type ScenarioSummary struct {
	ID      uint64
	Name    string
	Enabled bool
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
