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

type HostsTemplates struct {
	HostID     int64
	TemplateID int64
}

type Finding struct {
	Value   int
	Hidden  bool
	Message string
}

type Report struct {
	Findings []Finding
}
