package auditor

import (
	"testing"

	"github.com/sumwonyuno/cp-scoring/model"
)

func checkFinding(t *testing.T, finding model.Finding, show bool, value int64, message string) {
	t.Log("---------------")
	t.Log(show)
	t.Log(finding.Show)
	t.Log(finding.Value)
	t.Log(value)
	t.Log(message)
	t.Log(finding.Message)
	if finding.Show != show {
		t.Fatal("Finding show does not match")
	}
	if finding.Value != value {
		t.Fatal("Finding value does not match")
	}
	if finding.Message != message {
		t.Fatal("Finding message does not match")
	}
}

func TestAuditUserPresent(t *testing.T) {
	// template needs account present, account not present
	templateUser := model.User{Name: "user1", AccountPresent: true}
	finding := auditUserPresent(templateUser, false)
	checkFinding(t, finding, true, -1, "User removed: user1")

	// template needs account present, account present
	templateUser = model.User{Name: "user1", AccountPresent: true}
	finding = auditUserPresent(templateUser, true)
	checkFinding(t, finding, true, 1, "User present: user1")

	// TODO: case where account present has 0 value

	// template needs account not present, account not present
	templateUser = model.User{Name: "user1", AccountPresent: false}
	finding = auditUserPresent(templateUser, false)
	checkFinding(t, finding, true, 1, "User removed: user1")

	// template needs account not present, account present
	templateUser = model.User{Name: "user1", AccountPresent: false}
	finding = auditUserPresent(templateUser, true)
	checkFinding(t, finding, false, 0, "User present: user1")
}

func TestAuditUsers(t *testing.T) {
	state := model.State{}

	// no users in state, no users in template
	state.Users = make([]model.User, 0)
	template := model.Template{}
	findings := auditUsers(state, template)
	if len(findings) != 0 {
		t.Fatal("Expected 0 findings")
	}

	// user in state, no users in template
	state.Users = make([]model.User, 0)
	user1 := model.User{Name: "user1"}
	state.Users = append(state.Users, user1)
	template = model.Template{}
	findings = auditUsers(state, template)
	if len(findings) != 0 {
		t.Fatal("Expected 0 findings")
	}
}
