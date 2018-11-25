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
		report.Findings = append(report.Findings, auditSoftware(state, template)...)
		report.Findings = append(report.Findings, auditNetworkConnections(state, template)...)
	}

	return report
}

func auditUserAccountState(templateUser model.User, present bool) model.Finding {
	var presentFinding model.Finding
	if templateUser.AccountState == model.ObjectStateAdd && !present {
		presentFinding.Show = false
		presentFinding.Value = 0
		presentFinding.Message = "User not added: " + templateUser.Name
	} else if templateUser.AccountState == model.ObjectStateAdd && present {
		presentFinding.Show = true
		presentFinding.Value = 1
		presentFinding.Message = "User added: " + templateUser.Name
	} else if templateUser.AccountState == model.ObjectStateKeep && !present {
		presentFinding.Show = true
		presentFinding.Value = -1
		presentFinding.Message = "User not present: " + templateUser.Name
	} else if templateUser.AccountState == model.ObjectStateKeep && present {
		presentFinding.Show = false
		presentFinding.Value = 0
		presentFinding.Message = "User present: " + templateUser.Name
	} else if templateUser.AccountState == model.ObjectStateRemove && !present {
		presentFinding.Show = true
		presentFinding.Value = 1
		presentFinding.Message = "User removed: " + templateUser.Name
	} else if templateUser.AccountState == model.ObjectStateRemove && present {
		presentFinding.Show = false
		presentFinding.Value = 0
		presentFinding.Message = "User not removed: " + templateUser.Name
	} else {
		presentFinding.Show = false
		presentFinding.Value = 0
		presentFinding.Message = "Unknown user present state: " + templateUser.Name
	}

	return presentFinding
}

func auditUsers(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	foundUsers := make(map[string]model.User)
	for _, user := range state.Users {
		// all user accounts should have AccountState as Keep
		foundUsers[user.Name] = user
	}

	for _, templateUser := range template.Users {
		user, present := foundUsers[templateUser.Name]

		// check for user account state
		presentFinding := auditUserAccountState(templateUser, present)
		findings = append(findings, presentFinding)

		// check if user account active
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

		// check if user account expires
		var accountExpiresFinding model.Finding
		if templateUser.AccountExpires && user.AccountExpires {
			accountExpiresFinding.Show = true
			accountExpiresFinding.Value = 1
			accountExpiresFinding.Message = "User account expires: " + templateUser.Name
		} else if templateUser.AccountExpires && !user.AccountExpires {
			accountExpiresFinding.Show = false
			accountExpiresFinding.Value = 0
			accountExpiresFinding.Message = "User account does not expire: " + templateUser.Name
		} else if !templateUser.AccountExpires && user.AccountExpires {
			accountExpiresFinding.Show = true
			accountExpiresFinding.Value = -1
			accountExpiresFinding.Message = "User account expires: " + templateUser.Name
		} else if !templateUser.AccountExpires && !user.AccountExpires {
			accountExpiresFinding.Show = true
			accountExpiresFinding.Value = 1
			accountExpiresFinding.Message = "User account does not expire: " + templateUser.Name
		} else {
			accountExpiresFinding.Show = true
			accountExpiresFinding.Value = 0
			accountExpiresFinding.Message = "User account expires unknown: " + templateUser.Name
		}
		findings = append(findings, accountExpiresFinding)

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

	for _, templateProcess := range template.ProcessesAdd {
		for _, process := range state.Processes {
			if !strings.HasPrefix(process.CommandLine, templateProcess) {
				continue
			}
			var finding model.Finding
			finding.Show = true
			finding.Value = 1
			finding.Message = "Process found: " + templateProcess
			findings = append(findings, finding)
			break
		}
	}

	for _, templateProcess := range template.ProcessesKeep {
		match := false
		for _, process := range state.Processes {
			if !strings.HasPrefix(process.CommandLine, templateProcess) {
				continue
			}
			match = true
			break
		}
		if !match {
			var finding model.Finding
			finding.Show = true
			finding.Value = -1
			finding.Message = "Process not found: " + templateProcess
			findings = append(findings, finding)
		}
	}

	for _, templateProcess := range template.ProcessesRemove {
		match := false
		for _, process := range state.Processes {
			if !strings.HasPrefix(process.CommandLine, templateProcess) {
				continue
			}
			match = true
			break
		}
		if !match {
			var finding model.Finding
			finding.Show = true
			finding.Value = 1
			finding.Message = "Process removed: " + templateProcess
			findings = append(findings, finding)
		}
	}

	return findings
}

func auditSoftware(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	// keep track of software and version from state
	software := make(map[string]string)
	for _, sw := range state.Software {
		software[sw.Name] = sw.Version
	}

	for _, templateSoftware := range template.SoftwareAdd {
		version, present := software[templateSoftware.Name]
		if present {
			if templateSoftware.Version == version {
				var finding model.Finding
				finding.Show = true
				finding.Value = 1
				finding.Message = "Software added: " + templateSoftware.Name + ", " + templateSoftware.Version
				findings = append(findings, finding)
			}
		}
	}

	for _, templateSoftware := range template.SoftwareKeep {
		version, present := software[templateSoftware.Name]
		if present {
			if templateSoftware.Version != version {
				var finding model.Finding
				finding.Show = true
				finding.Value = -1
				finding.Message = "Software removed: " + templateSoftware.Name + ", " + templateSoftware.Version
				findings = append(findings, finding)
			}
		} else {
			var finding model.Finding
			finding.Show = true
			finding.Value = -1
			finding.Message = "Software removed: " + templateSoftware.Name + ", " + templateSoftware.Version
			findings = append(findings, finding)
		}
	}

	for _, templateSoftware := range template.SoftwareRemove {
		version, present := software[templateSoftware.Name]
		if present {
			if templateSoftware.Version != version {
				var finding model.Finding
				finding.Show = true
				finding.Value = 1
				finding.Message = "Software removed: " + templateSoftware.Name + ", " + templateSoftware.Version
				findings = append(findings, finding)
			}
		} else {
			var finding model.Finding
			finding.Show = true
			finding.Value = 1
			finding.Message = "Software removed: " + templateSoftware.Name + ", " + templateSoftware.Version
			findings = append(findings, finding)
		}
	}

	return findings
}

func auditNetworkConnections(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	for _, templateConn := range template.NetworkConnsAdd {
		connStr := templateConn.Protocol + " " + templateConn.LocalAddress + ":" + templateConn.LocalPort + " " + templateConn.RemoteAddress + ":" + templateConn.RemotePort
		for _, conn := range state.NetworkConnections {
			if len(templateConn.Protocol) > 0 && templateConn.Protocol != conn.Protocol {
				continue
			}
			if len(templateConn.LocalAddress) > 0 && templateConn.LocalAddress != conn.LocalAddress {
				continue
			}
			if len(templateConn.LocalPort) > 0 && templateConn.LocalPort != conn.LocalPort {
				continue
			}
			if len(templateConn.RemoteAddress) > 0 && templateConn.RemoteAddress != conn.RemoteAddress {
				continue
			}
			if len(templateConn.RemotePort) > 0 && templateConn.RemotePort != conn.RemotePort {
				continue
			}
			// if here, then matched all
			var finding model.Finding
			finding.Show = true
			finding.Value = 1
			finding.Message = "Network connection added: " + connStr
			findings = append(findings, finding)
			break
		}
	}

	for _, templateConn := range template.NetworkConnsKeep {
		connStr := templateConn.Protocol + " " + templateConn.LocalAddress + ":" + templateConn.LocalPort + " " + templateConn.RemoteAddress + ":" + templateConn.RemotePort
		match := false
		for _, conn := range state.NetworkConnections {
			if len(templateConn.Protocol) > 0 && templateConn.Protocol != conn.Protocol {
				continue
			}
			if len(templateConn.LocalAddress) > 0 && templateConn.LocalAddress != conn.LocalAddress {
				continue
			}
			if len(templateConn.LocalPort) > 0 && templateConn.LocalPort != conn.LocalPort {
				continue
			}
			if len(templateConn.RemoteAddress) > 0 && templateConn.RemoteAddress != conn.RemoteAddress {
				continue
			}
			if len(templateConn.RemotePort) > 0 && templateConn.RemotePort != conn.RemotePort {
				continue
			}
			// if here, then above matched
			match = true
			break
		}
		if !match {
			var finding model.Finding
			finding.Show = true
			finding.Value = -1
			finding.Message = "Network connection missing: " + connStr
			findings = append(findings, finding)
		}
	}

	for _, templateConn := range template.NetworkConnsRemove {
		connStr := templateConn.Protocol + " " + templateConn.LocalAddress + ":" + templateConn.LocalPort + " " + templateConn.RemoteAddress + ":" + templateConn.RemotePort
		match := false
		for _, conn := range state.NetworkConnections {
			if len(templateConn.Protocol) > 0 && templateConn.Protocol != conn.Protocol {
				continue
			}
			if len(templateConn.LocalAddress) > 0 && templateConn.LocalAddress != conn.LocalAddress {
				continue
			}
			if len(templateConn.LocalPort) > 0 && templateConn.LocalPort != conn.LocalPort {
				continue
			}
			if len(templateConn.RemoteAddress) > 0 && templateConn.RemoteAddress != conn.RemoteAddress {
				continue
			}
			if len(templateConn.RemotePort) > 0 && templateConn.RemotePort != conn.RemotePort {
				continue
			}
			// if here, then above matched
			match = true
			break
		}
		if !match {
			var finding model.Finding
			finding.Show = true
			finding.Value = 1
			finding.Message = "Network connection removed: " + connStr
			findings = append(findings, finding)
		}
	}

	return findings
}
