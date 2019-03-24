package model

type Template struct {
	ID    uint64
	Name  string
	State State
}

// NewTemplate returns a new Template instance
func NewTemplate() Template {
	return Template{State: State{}}
}

type Finding struct {
	Value   int64
	Show    bool
	Message string
}

type Report struct {
	Timestamp int64
	Findings  []Finding
}

type TeamSummary struct {
	ID   uint64
	Name string
}

type Team struct {
	ID      uint64
	Name    string
	POC     string
	Email   string
	Enabled bool
	Key     string
}

type ScenarioSummary struct {
	ID   uint64
	Name string
}

type Scenario struct {
	ID            uint64
	Name          string
	Description   string
	Enabled       bool
	HostTemplates map[uint64][]uint64
}

type TeamScore struct {
	TeamName  string
	Timestamp int64
	Score     int64
}

type ScenarioHostScore struct {
	ScenarioID uint64
	HostToken  string
	Timestamp  int64
	Score      int64
}

type ScenarioTimeline struct {
	Timestamps []int64
	Scores     []int64
}

type ScenarioHosts struct {
	ScenarioID   uint64
	ScenarioName string
	Hosts        []Host
}
