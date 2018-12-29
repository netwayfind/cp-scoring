package model

type Template struct {
	ID    int64
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
	ID   int64
	Name string
}

type Team struct {
	ID      int64
	Name    string
	POC     string
	Email   string
	Enabled bool
	Key     string
}

type ScenarioSummary struct {
	ID   int64
	Name string
}

type Scenario struct {
	ID            int64
	Name          string
	Description   string
	Enabled       bool
	HostTemplates map[int64][]int64
}

type ScenarioLatestScore struct {
	TeamName  string
	Timestamp int64
	Score     int64
}

type ScenarioScore struct {
	ScenarioID int64
	TeamID     int64
	HostID     int64
	Timestamp  int64
	Score      int64
}

type ScenarioTimeline struct {
	Timestamps []int64
	Scores     []int64
}

type ScenarioHosts struct {
	ScenarioID   int64
	ScenarioName string
	Hosts        []Host
}
