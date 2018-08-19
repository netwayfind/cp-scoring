package auditor

import (
	"strings"

	"github.com/sumwonyuno/cp-scoring/model"
)

// Audit will audit the given state against the given templates, and then returns a report.
func Audit(state model.State, templates []model.Template) model.Report {
	var report model.Report

	for _, template := range templates {
		report.Findings = append(report.Findings, auditUsers(state, template)...)
		report.Findings = append(report.Findings, auditGroups(state, template)...)
		report.Findings = append(report.Findings, auditProcesses(state, template)...)
	}

	return report
}

func auditUsers(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	foundUsers := make(map[string]model.User)
	for _, user := range state.Users {
		// all users should have AccountPresent as true
		foundUsers[user.Name] = user
	}

	for _, templateUser := range template.Users {
		user, present := foundUsers[templateUser.Name]

		// check for user presence
		var presentFinding model.Finding
		if templateUser.AccountPresent && !present {
			presentFinding.Show = true
			presentFinding.Value = -1
			presentFinding.Message = "Required user missing: " + templateUser.Name
		} else if !templateUser.AccountPresent && !present {
			presentFinding.Show = true
			presentFinding.Value = 1
			presentFinding.Message = "User removed: " + templateUser.Name
		} else if templateUser.AccountPresent {
			presentFinding.Show = true
			presentFinding.Value = 1
			presentFinding.Message = "Required user present: " + templateUser.Name
		} else if !templateUser.AccountActive {
			presentFinding.Show = false
			presentFinding.Value = 0
			presentFinding.Message = "User not removed: " + templateUser.Name
		} else {
			presentFinding.Show = false
			presentFinding.Value = 0
			presentFinding.Message = "Unknown user present state: " + templateUser.Name
		}
		findings = append(findings, presentFinding)

		// check if user is active/unlocked
		var activeFinding model.Finding
		if templateUser.AccountActive && user.AccountActive {
			activeFinding.Show = true
			activeFinding.Value = 1
			activeFinding.Message = "User active: " + templateUser.Name
		} else if templateUser.AccountActive && !user.AccountActive {
			activeFinding.Show = true
			activeFinding.Value = -1
			activeFinding.Message = "User not active: " + templateUser.Name
		} else if !templateUser.AccountActive && user.AccountActive {
			activeFinding.Show = false
			activeFinding.Value = 0
			activeFinding.Message = "User active: " + templateUser.Name
		} else if !templateUser.AccountActive && !user.AccountActive {
			activeFinding.Show = true
			activeFinding.Value = 1
			activeFinding.Message = "User not active: " + templateUser.Name
		} else {
			activeFinding.Show = true
			activeFinding.Value = 0
			activeFinding.Message = "Unknown user active state: " + templateUser.Name
		}
		findings = append(findings, activeFinding)

		// check if user password expires
		var passwordExpiresFinding model.Finding
		if templateUser.PasswordExpires && user.PasswordExpires {
			passwordExpiresFinding.Show = true
			passwordExpiresFinding.Value = 1
			passwordExpiresFinding.Message = "User password expires: " + templateUser.Name
		} else if templateUser.PasswordExpires && !user.PasswordExpires {
			passwordExpiresFinding.Show = false
			passwordExpiresFinding.Value = 0
			passwordExpiresFinding.Message = "User password does not expire: " + templateUser.Name
		} else if !templateUser.PasswordExpires && user.PasswordExpires {
			passwordExpiresFinding.Show = true
			passwordExpiresFinding.Value = -1
			passwordExpiresFinding.Message = "User password expires: " + templateUser.Name
		} else if !templateUser.PasswordExpires && !user.PasswordExpires {
			passwordExpiresFinding.Show = true
			passwordExpiresFinding.Value = 1
			passwordExpiresFinding.Message = "User password does not expire: " + templateUser.Name
		} else {
			passwordExpiresFinding.Show = true
			passwordExpiresFinding.Value = 0
			passwordExpiresFinding.Message = "User password expires unknown: " + templateUser.Name
		}
		findings = append(findings, passwordExpiresFinding)

		// check if user changed password
		var passwordChangedFinding model.Finding
		if templateUser.PasswordLastSet < user.PasswordLastSet {
			passwordChangedFinding.Show = true
			passwordChangedFinding.Value = 1
			passwordChangedFinding.Message = "User password changed: " + templateUser.Name
		} else {
			passwordChangedFinding.Show = false
			passwordChangedFinding.Value = 0
			passwordChangedFinding.Message = "User password not changed: " + templateUser.Name
		}
		findings = append(findings, passwordChangedFinding)
	}

	return findings
}

func auditGroups(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	for group, members := range state.Groups {
		foundMembers := make(map[string]bool)
		for _, member := range members {
			foundMembers[member] = true
		}

		templateMembers, present := template.GroupMembersAdd[group]
		if present {
			for _, templateMember := range templateMembers {
				_, p := foundMembers[templateMember]
				var finding model.Finding
				if p {
					finding.Show = true
					finding.Value = 1
					finding.Message = "User " + templateMember + " added to group " + group
				} else {
					finding.Show = false
					finding.Value = 0
					finding.Message = "User " + templateMember + " not added to group " + group
				}
				findings = append(findings, finding)
			}
		}
		templateMembers, present = template.GroupMembersKeep[group]
		if present {
			for _, templateMember := range templateMembers {
				_, p := foundMembers[templateMember]
				var finding model.Finding
				if p {
					finding.Show = false
					finding.Value = 0
					finding.Message = "User " + templateMember + " in group " + group
				} else {
					finding.Show = true
					finding.Value = -1
					finding.Message = "User " + templateMember + " not in group " + group
				}
				findings = append(findings, finding)
			}
		}
		templateMembers, present = template.GroupMembersRemove[group]
		if present {
			for _, templateMember := range templateMembers {
				_, p := foundMembers[templateMember]
				var finding model.Finding
				if p {
					finding.Show = false
					finding.Value = 0
					finding.Message = "User " + templateMember + " in group " + group
				} else {
					finding.Show = true
					finding.Value = 1
					finding.Message = "User " + templateMember + " removed from group " + group
				}
				findings = append(findings, finding)
			}
		}
	}

	return findings
}

func auditProcesses(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	for _, process := range state.Processes {
		for _, templateProcess := range template.ProcessesAdd {
			var finding model.Finding
			if strings.HasPrefix(process.CommandLine, templateProcess) {
				finding.Show = true
				finding.Value = 1
				finding.Message = "Process found: " + templateProcess
			} else {
				finding.Show = false
				finding.Value = 0
				finding.Message = "Process not found: " + templateProcess
			}
			findings = append(findings, finding)
		}
		for _, templateProcess := range template.ProcessesKeep {
			var finding model.Finding
			if strings.HasPrefix(process.CommandLine, templateProcess) {
				finding.Show = false
				finding.Value = 0
				finding.Message = "Process found: " + templateProcess
			} else {
				finding.Show = true
				finding.Value = -1
				finding.Message = "Process not found: " + templateProcess
			}
			findings = append(findings, finding)
		}
		for _, templateProcess := range template.ProcessesRemove {
			var finding model.Finding
			if strings.HasPrefix(process.CommandLine, templateProcess) {
				finding.Show = false
				finding.Value = 0
				finding.Message = "Process found: " + templateProcess
			} else {
				finding.Show = true
				finding.Value = 1
				finding.Message = "Process removed: " + templateProcess
			}
			findings = append(findings, finding)
		}
	}

	return findings
}
