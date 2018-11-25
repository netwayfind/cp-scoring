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
	// template account state add
	templateUser := model.User{Name: "user1", AccountState: model.ObjectStateAdd}
	// account not present
	finding := auditUserAccountState(templateUser, false)
	checkFinding(t, finding, false, 0, "User not added: user1")
	// account present
	finding = auditUserAccountState(templateUser, true)
	checkFinding(t, finding, true, 1, "User added: user1")

	// template account state keep
	templateUser = model.User{Name: "user1", AccountState: model.ObjectStateKeep}
	// account not present
	finding = auditUserAccountState(templateUser, false)
	checkFinding(t, finding, true, -1, "User not present: user1")
	// account present
	finding = auditUserAccountState(templateUser, true)
	checkFinding(t, finding, false, 0, "User present: user1")

	// template account state remove
	templateUser = model.User{Name: "user1", AccountState: model.ObjectStateRemove}
	// account not present
	finding = auditUserAccountState(templateUser, false)
	checkFinding(t, finding, true, 1, "User removed: user1")
	// account present
	finding = auditUserAccountState(templateUser, true)
	checkFinding(t, finding, false, 0, "User not removed: user1")
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
