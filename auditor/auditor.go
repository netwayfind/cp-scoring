package auditor

import (
	"strconv"
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
		report.Findings = append(report.Findings, auditScheduledTasks(state, template)...)
		report.Findings = append(report.Findings, auditWindowsFirewallProfiles(state, template)...)
		report.Findings = append(report.Findings, auditWindowsFirewallRules(state, template)...)
	}

	return report
}

func auditUserObjectState(templateUser model.User, present bool) model.Finding {
	var presentFinding model.Finding
	if templateUser.ObjectState == model.ObjectStateAdd {
		if present {
			presentFinding.Show = true
			presentFinding.Value = 1
			presentFinding.Message = "User added: " + templateUser.Name
		} else {
			presentFinding.Show = false
			presentFinding.Value = 0
			presentFinding.Message = "User not added: " + templateUser.Name
		}
	} else if templateUser.ObjectState == model.ObjectStateKeep {
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
	} else if templateUser.ObjectState == model.ObjectStateRemove {
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
			if templateUser.ObjectState == model.ObjectStateKeep {
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
		// all user accounts should have ObjectState as Keep
		foundUsers[user.Name] = user
	}

	for _, templateUser := range template.State.Users {
		user, present := foundUsers[templateUser.Name]

		// check for user account state
		presentFinding := auditUserObjectState(templateUser, present)
		findings = append(findings, presentFinding)

		// no need to check further if user isn't present
		if !present {
			continue
		}

		// no need to check further if user should be removed
		if templateUser.ObjectState == model.ObjectStateRemove {
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

	for templateGroup, templateMembers := range template.State.Groups {
		if groupMembers, present := state.Groups[templateGroup]; present {
			// template group present
			foundMembers := make(map[string]bool)
			for _, groupMember := range groupMembers {
				foundMembers[groupMember.Name] = true
			}
			for _, templateMember := range templateMembers {
				finding := model.Finding{}
				_, memberPresent := foundMembers[templateMember.Name]
				if templateMember.ObjectState == model.ObjectStateAdd {
					if memberPresent {
						finding.Show = true
						finding.Value = 1
						finding.Message = "Group " + templateGroup + ", member added: " + templateMember.Name
					} else {
						finding.Show = false
						finding.Value = 0
						finding.Message = "Group " + templateGroup + ", member not added: " + templateMember.Name
					}
				} else if templateMember.ObjectState == model.ObjectStateKeep {
					if memberPresent {
						finding.Show = false
						finding.Value = 0
						finding.Message = "Group " + templateGroup + ", member found: " + templateMember.Name
					} else {
						finding.Show = true
						finding.Value = -1
						finding.Message = "Group " + templateGroup + ", member not found: " + templateMember.Name
					}
				} else if templateMember.ObjectState == model.ObjectStateRemove {
					if memberPresent {
						finding.Show = false
						finding.Value = 0
						finding.Message = "Group " + templateGroup + ", member not removed: " + templateMember.Name
					} else {
						finding.Show = true
						finding.Value = 1
						finding.Message = "Group " + templateGroup + ", member removed: " + templateMember.Name
					}
				} else {
					finding.Show = false
					finding.Value = 0
					finding.Message = "Group " + templateGroup + ", member unknown state: " + templateMember.Name
				}
				findings = append(findings, finding)
			}
		} else {
			// template group not present
			for _, templateMember := range templateMembers {
				finding := model.Finding{}
				if templateMember.ObjectState == model.ObjectStateAdd {
					finding.Show = false
					finding.Value = 0
					finding.Message = "Group " + templateGroup + ", member not added: " + templateMember.Name
				} else if templateMember.ObjectState == model.ObjectStateKeep {
					finding.Show = true
					finding.Value = -1
					finding.Message = "Group " + templateGroup + ", member not found: " + templateMember.Name
				} else if templateMember.ObjectState == model.ObjectStateRemove {
					finding.Show = true
					finding.Value = 1
					finding.Message = "Group " + templateGroup + ", member removed: " + templateMember.Name
				} else {
					finding.Show = false
					finding.Value = 0
					finding.Message = "Group " + templateGroup + ", member unknown state: " + templateMember.Name
				}
				findings = append(findings, finding)
			}
		}
	}

	return findings
}

func auditProcesses(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	for _, templateProcess := range template.State.Processes {
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

	if templateSoftware.ObjectState == model.ObjectStateAdd {
		if found {
			finding.Show = true
			finding.Value = 1
			finding.Message = "Software added: "
		} else {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Software not added: "
		}
	} else if templateSoftware.ObjectState == model.ObjectStateKeep {
		if found {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Software found: "
		} else {
			finding.Show = true
			finding.Value = -1
			finding.Message = "Software not found: "
		}
	} else if templateSoftware.ObjectState == model.ObjectStateRemove {
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

	for _, templateSoftware := range template.State.Software {
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

	for _, templateConn := range template.State.NetworkConnections {
		finding := auditNetworkConnectionObjectState(templateConn, state)
		findings = append(findings, finding)
	}

	return findings
}

func auditScheduledTaskObjectState(templateTask model.ScheduledTask, state model.State) model.Finding {
	var finding model.Finding

	taskStr := templateTask.Name + " @ " + templateTask.Path + ", " + strconv.FormatBool(templateTask.Enabled)

	found := false
	for _, task := range state.ScheduledTasks {
		if templateTask.Name != task.Name {
			continue
		}
		if templateTask.Enabled != task.Enabled {
			continue
		}
		if len(templateTask.Path) > 0 && templateTask.Path != task.Path {
			continue
		}
		found = true
		break
	}

	if templateTask.ObjectState == model.ObjectStateAdd {
		if found {
			finding.Show = true
			finding.Value = 1
			finding.Message = "Scheduled task added: " + taskStr
		} else {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Scheduled task not added: " + taskStr
		}
	} else if templateTask.ObjectState == model.ObjectStateKeep {
		if found {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Scheduled task found: " + taskStr
		} else {
			finding.Show = true
			finding.Value = -1
			finding.Message = "Scheduled task not found: " + taskStr
		}
	} else if templateTask.ObjectState == model.ObjectStateRemove {
		if found {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Scheduled task not removed: " + taskStr
		} else {
			finding.Show = true
			finding.Value = 1
			finding.Message = "Scheduled task removed: " + taskStr
		}
	} else {
		finding.Show = false
		finding.Value = 0
		finding.Message = "Unknown scheduled task state: " + taskStr
	}

	return finding
}

func auditScheduledTasks(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	for _, templateTask := range template.State.ScheduledTasks {
		finding := auditScheduledTaskObjectState(templateTask, state)
		findings = append(findings, finding)
	}

	return findings
}

func auditWindowsFirewallProfiles(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	foundProfiles := make(map[string]model.WindowsFirewallProfile)
	for _, profile := range state.WindowsFirewallProfiles {
		foundProfiles[profile.Name] = profile
	}

	for _, templateProfile := range template.State.WindowsFirewallProfiles {
		// check profile present
		profile, present := foundProfiles[templateProfile.Name]
		presentFinding := model.Finding{}
		if present {
			presentFinding.Show = false
			presentFinding.Value = 0
			presentFinding.Message = "Windows Firewall Profile found: " + templateProfile.Name
		} else {
			presentFinding.Show = false
			presentFinding.Value = 0
			presentFinding.Message = "Windows Firewall Profile not found: " + templateProfile.Name
		}
		findings = append(findings, presentFinding)
		if !present {
			continue
		}

		// check profile enabled setting
		enabledFinding := model.Finding{}
		if templateProfile.Enabled == profile.Enabled {
			enabledFinding.Show = true
			enabledFinding.Value = 1
			enabledFinding.Message = "Windows Firewall Profile " + profile.Name + " enabled"
		} else {
			enabledFinding.Show = false
			enabledFinding.Value = 0
			enabledFinding.Message = "Windows Firewall Profile " + profile.Name + " not enabled"
		}
		findings = append(findings, enabledFinding)

		// check profile inbound setting
		inboundFinding := model.Finding{}
		if templateProfile.DefaultInboundAction == profile.DefaultInboundAction {
			inboundFinding.Show = true
			inboundFinding.Value = 1
			inboundFinding.Message = "Windows Firewall Profile " + profile.Name + " inbound: " + profile.DefaultInboundAction
		} else {
			inboundFinding.Show = false
			inboundFinding.Value = 0
			inboundFinding.Message = "Windows Firewall Profile " + profile.Name + " inbound not matched: " + profile.DefaultInboundAction
		}
		findings = append(findings, inboundFinding)

		// check profile outbound setting
		outboundFinding := model.Finding{}
		if templateProfile.DefaultOutboundAction == profile.DefaultOutboundAction {
			outboundFinding.Show = true
			outboundFinding.Value = 1
			outboundFinding.Message = "Windows Firewall Profile " + profile.Name + " outbound: " + profile.DefaultOutboundAction
		} else {
			outboundFinding.Show = false
			outboundFinding.Value = 0
			outboundFinding.Message = "Windows Firewall Profile " + profile.Name + " outbound not matched: " + profile.DefaultOutboundAction
		}
		findings = append(findings, outboundFinding)
	}

	return findings
}

func auditWindowsFirewallRuleObjectState(templateRule model.WindowsFirewallRule, state model.State) model.Finding {
	var finding model.Finding

	found := false
	for _, rule := range state.WindowsFirewallRules {
		if len(templateRule.DisplayName) > 0 && templateRule.DisplayName != rule.DisplayName {
			continue
		}
		if templateRule.Enabled != rule.Enabled {
			continue
		}
		if len(templateRule.Protocol) > 0 && templateRule.Protocol != rule.Protocol {
			continue
		}
		if len(templateRule.LocalPort) > 0 && templateRule.LocalPort != rule.LocalPort {
			continue
		}
		if len(templateRule.RemoteAddress) > 0 && templateRule.RemoteAddress != rule.RemoteAddress {
			continue
		}
		if len(templateRule.RemotePort) > 0 && templateRule.RemotePort != rule.RemotePort {
			continue
		}
		if len(templateRule.Direction) > 0 && templateRule.Direction != rule.Direction {
			continue
		}
		if len(templateRule.Action) > 0 && templateRule.Action != rule.Action {
			continue
		}
		found = true
		break
	}

	templateRuleStr := templateRule.DisplayName + ", " + templateRule.Protocol + ", " + templateRule.LocalPort + ", " + templateRule.RemoteAddress + ", " + templateRule.RemotePort + ", " + templateRule.Direction + ", " + templateRule.Action

	if templateRule.ObjectState == model.ObjectStateAdd {
		if found {
			finding.Show = true
			finding.Value = 1
			finding.Message = "Windows Firewall rule added: " + templateRuleStr
		} else {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Windows Firewall rule not added: " + templateRuleStr
		}
	} else if templateRule.ObjectState == model.ObjectStateKeep {
		if found {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Windows Firewall rule found: " + templateRuleStr
		} else {
			finding.Show = true
			finding.Value = -1
			finding.Message = "Windows Firewall rule not found: " + templateRuleStr
		}
	} else if templateRule.ObjectState == model.ObjectStateRemove {
		if found {
			finding.Show = false
			finding.Value = 0
			finding.Message = "Windows Firewall rule not removed: " + templateRuleStr
		} else {
			finding.Show = true
			finding.Value = 1
			finding.Message = "Windows Firewall rule removed: " + templateRuleStr
		}
	} else {
		finding.Show = false
		finding.Value = 0
		finding.Message = "Unknown Windows Firewall rule state: " + templateRuleStr
	}

	return finding
}

func auditWindowsFirewallRules(state model.State, template model.Template) []model.Finding {
	findings := make([]model.Finding, 0)

	for _, templateRule := range template.State.WindowsFirewallRules {
		// check rule present
		presentFinding := auditWindowsFirewallRuleObjectState(templateRule, state)
		findings = append(findings, presentFinding)
	}

	return findings
}
