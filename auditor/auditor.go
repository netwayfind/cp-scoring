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

	usersPresent := make(map[string]bool)
	for _, user := range state.Users {
		usersPresent[user] = true
	}

	for _, user := range template.UsersAdd {
		var finding model.Finding
		if _, present := usersPresent[user]; present {
			finding.Value = 1
			finding.Hidden = false
			finding.Message = "Added user " + user
		} else {
			finding.Value = 0
			finding.Hidden = true
			finding.Message = "Missing user " + user
		}
		findings = append(findings, finding)
	}

	for _, user := range template.UsersKeep {
		var finding model.Finding
		if _, present := usersPresent[user]; !present {
			finding.Value = -1
			finding.Hidden = false
			finding.Message = "Removed required user " + user
		} else {
			finding.Value = 0
			finding.Hidden = true
			finding.Message = "Found required user " + user
		}
		findings = append(findings, finding)
	}

	for _, user := range template.UsersRemove {
		var finding model.Finding
		if _, present := usersPresent[user]; !present {
			finding.Value = 1
			finding.Hidden = false
			finding.Message = "Removed user " + user
		} else {
			finding.Value = 0
			finding.Hidden = true
			finding.Message = "Did not remove user " + user
		}
		findings = append(findings, finding)
	}

	return findings
}
