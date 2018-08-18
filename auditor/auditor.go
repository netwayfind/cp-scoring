package auditor

import (
	"github.com/sumwonyuno/cp-scoring/model"
)

// Audit will audit the given state against the given templates, and then returns a report.
func Audit(state model.State, templates []model.Template) model.Report {
	var report model.Report

	for _, template := range templates {
		r := auditUsers(state, template)
		for _, result := range r {
			report.Findings = append(report.Findings, result)
		}
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
