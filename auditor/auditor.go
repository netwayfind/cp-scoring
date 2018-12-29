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
	if templateUser.AccountState == model.ObjectStateAdd {
		if present {
			presentFinding.Show = true
			presentFinding.Value = 1
			presentFinding.Message = "User added: " + templateUser.Name
		} else {
			presentFinding.Show = false
			presentFinding.Value = 0
			presentFinding.Message = "User not added: " + templateUser.Name
		}
	} else if templateUser.AccountState == model.ObjectStateKeep {
		if present {
			// don't show indication that user should be kept
			presentFinding.Show = false
			presentFinding.Value = 0
			presentFinding.Message = "User present: " + templateUser.Name
		} else {
			presentFinding.Show = true
			presentFinding.Value = -1
			presentFinding.Message = "User not present: " + templateUser.Name
		}
	} else if templateUser.AccountState == model.ObjectStateRemove {
		if present {
			presentFinding.Show = false
			presentFinding.Value = 0
			presentFinding.Message = "User not removed: " + templateUser.Name
		} else {
			presentFinding.Show = true
			presentFinding.Value = 1
			presentFinding.Message = "User removed: " + templateUser.Name
		}
	} else {
		presentFinding.Show = false
		presentFinding.Value = 0
		presentFinding.Message = "Unknown user account state: " + templateUser.Name
	}

	return presentFinding
}

func auditUserAccountActive(templateUser model.User, user model.User) model.Finding {
	var activeFinding model.Finding
	if templateUser.AccountActive {
		if user.AccountActive {
			activeFinding.Message = "User active: " + templateUser.Name
			// if user is supposed to be kept, don't show this and don't add points
			if templateUser.AccountState == model.ObjectStateKeep {
				activeFinding.Show = false
				activeFinding.Value = 0
			} else {
				activeFinding.Show = true
				activeFinding.Value = 1
			}
		} else {
			activeFinding.Show = true
			activeFinding.Value = -1
			activeFinding.Message = "User not active: " + templateUser.Name
		}
	} else {
		if user.AccountActive {
			activeFinding.Show = false
			activeFinding.Value = 0
			activeFinding.Message = "User active: " + templateUser.Name
		} else {
			activeFinding.Show = true
			activeFinding.Value = 1
			activeFinding.Message = "User not active: " + templateUser.Name
		}
	}

	return activeFinding
}

func auditUserAccountExpire(templateUser model.User, user model.User) model.Finding {
	var accountExpiresFinding model.Finding
	if templateUser.AccountExpires {
		if user.AccountExpires {
			accountExpiresFinding.Show = true
			accountExpiresFinding.Value = 1
			accountExpiresFinding.Message = "User account expires: " + templateUser.Name
		} else {
			accountExpiresFinding.Show = false
			accountExpiresFinding.Value = 0
			accountExpiresFinding.Message = "User account does not expire: " + templateUser.Name
		}
	} else {
		if user.AccountExpires {
			accountExpiresFinding.Show = true
			accountExpiresFinding.Value = -1
			accountExpiresFinding.Message = "User account expires: " + templateUser.Name
		} else {
			accountExpiresFinding.Show = false
			accountExpiresFinding.Value = 0
			accountExpiresFinding.Message = "User account does not expire: " + templateUser.Name
		}
	}

	return accountExpiresFinding
}

func auditUserPasswordExpire(templateUser model.User, user model.User) model.Finding {
	var passwordExpiresFinding model.Finding
	if templateUser.PasswordExpires {
		if user.PasswordExpires {
			passwordExpiresFinding.Show = true
			passwordExpiresFinding.Value = 1
			passwordExpiresFinding.Message = "User password expires: " + templateUser.Name
		} else {
			passwordExpiresFinding.Show = false
			passwordExpiresFinding.Value = 0
			passwordExpiresFinding.Message = "User password does not expire: " + templateUser.Name
		}
	} else {
		if user.PasswordExpires {
			passwordExpiresFinding.Show = true
			passwordExpiresFinding.Value = -1
			passwordExpiresFinding.Message = "User password expires: " + templateUser.Name
		} else {
			passwordExpiresFinding.Show = false
			passwordExpiresFinding.Value = 0
			passwordExpiresFinding.Message = "User password does not expire: " + templateUser.Name
		}
	}

	return passwordExpiresFinding
}

func auditUserPasswordChange(templateUser model.User, user model.User) model.Finding {
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

	return passwordChangedFinding
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

		// no need to check further if user isn't present
		if !present {
			continue
		}

		// no need to check further if user should be removed
		if templateUser.AccountState == model.ObjectStateRemove {
			continue
		}

		// check if user account active
		activeFinding := auditUserAccountActive(templateUser, user)
		findings = append(findings, activeFinding)

		// check if user account expires
		accountExpiresFinding := auditUserAccountExpire(templateUser, user)
		findings = append(findings, accountExpiresFinding)

		// check if user password expires
		passwordExpiresFinding := auditUserPasswordExpire(templateUser, user)
		findings = append(findings, passwordExpiresFinding)

		// check if user changed password
		passwordChangedFinding := auditUserPasswordChange(templateUser, user)
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

	for _, templateProcess := range template.Processes {
		finding := auditProcessState(templateProcess, state)
		findings = append(findings, finding)
	}

	return findings
}

func auditProcessState(templateProcess model.Process, state model.State) model.Finding {
	var finding model.Finding

	found := false
	for _, process := range state.Processes {
		// only need prefix match
		if strings.HasPrefix(process.CommandLine, templateProcess.CommandLine) {
			found = true
			break
		}
	}

	if templateProcess.ObjectState == model.ObjectStateAdd {
		if found {
			finding.Show = true
			finding.Value = 1
			finding.Message = "Process added: "
		} else {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Process not added: "
		}
	} else if templateProcess.ObjectState == model.ObjectStateKeep {
		if found {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Process found: "
		} else {
			finding.Show = true
			finding.Value = -1
			finding.Message = "Process not found: "
		}
	} else if templateProcess.ObjectState == model.ObjectStateRemove {
		if found {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Process not removed: "
		} else {
			finding.Show = true
			finding.Value = 1
			finding.Message = "Process removed: "
		}
	} else {
		finding.Show = false
		finding.Value = 0
		finding.Message = "Unknown process state: "
	}

	finding.Message += templateProcess.CommandLine

	return finding
}

func auditSoftwareState(templateSoftware model.Software, state model.State) model.Finding {
	var finding model.Finding

	found := false
	for _, software := range state.Software {
		if len(templateSoftware.Version) == 0 {
			if software.Name == templateSoftware.Name {
				found = true
				break
			}
		} else {
			if software.Name == templateSoftware.Name && software.Version == templateSoftware.Version {
				found = true
				break
			}
		}
	}

	if templateSoftware.SoftwareState == model.ObjectStateAdd {
		if found {
			finding.Show = true
			finding.Value = 1
			finding.Message = "Software added: "
		} else {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Software not added: "
		}
	} else if templateSoftware.SoftwareState == model.ObjectStateKeep {
		if found {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Software found: "
		} else {
			finding.Show = true
			finding.Value = -1
			finding.Message = "Software not found: "
		}
	} else if templateSoftware.SoftwareState == model.ObjectStateRemove {
		if found {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Software not removed: "
		} else {
			finding.Show = true
			finding.Value = 1
			finding.Message = "Software removed: "
		}
	} else {
		finding.Show = false
		finding.Value = 0
		finding.Message = "Unknown software state: "
	}

	softwareText := templateSoftware.Name
	if len(templateSoftware.Version) > 0 {
		softwareText += ", " + templateSoftware.Version
	}
	finding.Message += softwareText

	return finding
}

func auditSoftware(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	for _, templateSoftware := range template.Software {
		finding := auditSoftwareState(templateSoftware, state)
		findings = append(findings, finding)
	}

	return findings
}

func compareNetworkConnection(templateConn model.NetworkConnection, conn model.NetworkConnection) bool {
	if len(templateConn.Protocol) > 0 && templateConn.Protocol != conn.Protocol {
		return false
	}
	if len(templateConn.LocalAddress) > 0 && templateConn.LocalAddress != conn.LocalAddress {
		return false
	}
	if len(templateConn.LocalPort) > 0 && templateConn.LocalPort != conn.LocalPort {
		return false
	}
	if len(templateConn.RemoteAddress) > 0 && templateConn.RemoteAddress != conn.RemoteAddress {
		return false
	}
	if len(templateConn.RemotePort) > 0 && templateConn.RemotePort != conn.RemotePort {
		return false
	}
	// if here, then matched all
	return true
}

func auditNetworkConnectionObjectState(templateConn model.NetworkConnection, state model.State) model.Finding {
	var finding model.Finding

	connStr := templateConn.Protocol + " " + templateConn.LocalAddress + ":" + templateConn.LocalPort + " " + templateConn.RemoteAddress + ":" + templateConn.RemotePort

	found := false
	for _, conn := range state.NetworkConnections {
		if compareNetworkConnection(templateConn, conn) {
			found = true
			break
		}
	}

	if templateConn.ObjectState == model.ObjectStateAdd {
		if found {
			finding.Show = true
			finding.Value = 1
			finding.Message = "Network connection added: " + connStr
		} else {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Network connection not added: " + connStr
		}
	} else if templateConn.ObjectState == model.ObjectStateKeep {
		if found {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Network connection found: " + connStr
		} else {
			finding.Show = true
			finding.Value = -1
			finding.Message = "Network connection not found: " + connStr
		}
	} else if templateConn.ObjectState == model.ObjectStateRemove {
		if found {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Network connection not removed: " + connStr
		} else {
			finding.Show = true
			finding.Value = 1
			finding.Message = "Network connection removed: " + connStr
		}
	} else {
		finding.Show = false
		finding.Value = 0
		finding.Message = "Unknown network connection state: " + connStr
	}

	return finding
}

func auditNetworkConnections(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	for _, templateConn := range template.NetworkConns {
		finding := auditNetworkConnectionObjectState(templateConn, state)
		findings = append(findings, finding)
	}

	return findings
}
