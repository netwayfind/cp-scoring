package model

type Template struct {
	Name string
	UsersAdd []string
	UsersKeep []string
	UsersRemove []string
	GroupMembersAdd map[string]string
	GroupMembersKeep map[string]string
	GroupMembersRemove map[string]string
}

type Report struct {

}