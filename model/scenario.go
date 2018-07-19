package model

type Template struct {
	Name               string
	UsersAdd           []string
	UsersKeep          []string
	UsersRemove        []string
	GroupMembersAdd    map[string]string
	GroupMembersKeep   map[string]string
	GroupMembersRemove map[string]string
}

type Finding struct {
	Value   int
	Hidden  bool
	Message string
}

type Report struct {
	Findings []Finding
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
}

type Scenario struct {
	ID            int64
	Name          string
	Description   string
	Enabled       bool
	HostTemplates map[int64][]int64
}
