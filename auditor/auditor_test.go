package auditor

import (
	"testing"

	"github.com/sumwonyuno/cp-scoring/model"
)

func checkFinding(t *testing.T, finding model.Finding, show bool, value int64, message string) {
	t.Log("---------------")
	t.Log(show)
	t.Log(finding.Show)
	t.Log(value)
	t.Log(finding.Value)
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
	// template account state not set
	templateUser := model.User{Name: "user1"}
	// account not present
	finding := auditUserAccountState(templateUser, false)
	checkFinding(t, finding, false, 0, "Unknown user account state: user1")
	// account present
	finding = auditUserAccountState(templateUser, true)
	checkFinding(t, finding, false, 0, "Unknown user account state: user1")

	// template account state add
	templateUser = model.User{Name: "user1", AccountState: model.ObjectStateAdd}
	// account not present
	finding = auditUserAccountState(templateUser, false)
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

func TestAuditUserActive(t *testing.T) {
	// template account active, user active
	templateUser := model.User{Name: "user1", AccountActive: true}
	user := model.User{Name: "user1", AccountActive: true}
	finding := auditUserAccountActive(templateUser, user)
	checkFinding(t, finding, true, 1, "User active: user1")

	// template account active, user not active
	templateUser = model.User{Name: "user1", AccountActive: true}
	user = model.User{Name: "user1", AccountActive: false}
	finding = auditUserAccountActive(templateUser, user)
	checkFinding(t, finding, true, -1, "User not active: user1")

	// template account not active, user active
	templateUser = model.User{Name: "user1", AccountActive: false}
	user = model.User{Name: "user1", AccountActive: true}
	finding = auditUserAccountActive(templateUser, user)
	checkFinding(t, finding, false, 0, "User active: user1")

	// template account not active, user not active
	templateUser = model.User{Name: "user1", AccountActive: false}
	user = model.User{Name: "user1", AccountActive: false}
	finding = auditUserAccountActive(templateUser, user)
	checkFinding(t, finding, true, 1, "User not active: user1")
}

func TestAuditUserExpire(t *testing.T) {
	// template account not expire, user not expire
	templateUser := model.User{Name: "user1", AccountExpires: false}
	user := model.User{Name: "user1", AccountExpires: false}
	finding := auditUserAccountExpire(templateUser, user)
	checkFinding(t, finding, false, 0, "User account does not expire: user1")

	// template account not expire, user expire
	templateUser = model.User{Name: "user1", AccountExpires: false}
	user = model.User{Name: "user1", AccountExpires: true}
	finding = auditUserAccountExpire(templateUser, user)
	checkFinding(t, finding, true, -1, "User account expires: user1")

	// template account expire, user not expire
	templateUser = model.User{Name: "user1", AccountExpires: true}
	user = model.User{Name: "user1", AccountExpires: false}
	finding = auditUserAccountExpire(templateUser, user)
	checkFinding(t, finding, false, 0, "User account does not expire: user1")

	// template account expire, user expire
	templateUser = model.User{Name: "user1", AccountExpires: true}
	user = model.User{Name: "user1", AccountExpires: true}
	finding = auditUserAccountExpire(templateUser, user)
	checkFinding(t, finding, true, 1, "User account expires: user1")
}

func TestAuditUserPasswordExpire(t *testing.T) {
	// temple password not expire, password not expire
	templateUser := model.User{Name: "user1", PasswordExpires: false}
	user := model.User{Name: "user1", PasswordExpires: false}
	finding := auditUserPasswordExpire(templateUser, user)
	checkFinding(t, finding, false, 0, "User password does not expire: user1")

	// temple password not expire, password expire
	templateUser = model.User{Name: "user1", PasswordExpires: false}
	user = model.User{Name: "user1", PasswordExpires: true}
	finding = auditUserPasswordExpire(templateUser, user)
	checkFinding(t, finding, true, -1, "User password expires: user1")

	// temple password expire, password not expire
	templateUser = model.User{Name: "user1", PasswordExpires: true}
	user = model.User{Name: "user1", PasswordExpires: false}
	finding = auditUserPasswordExpire(templateUser, user)
	checkFinding(t, finding, false, 0, "User password does not expire: user1")

	// temple password expire, password expire
	templateUser = model.User{Name: "user1", PasswordExpires: true}
	user = model.User{Name: "user1", PasswordExpires: true}
	finding = auditUserPasswordExpire(templateUser, user)
	checkFinding(t, finding, true, 1, "User password expires: user1")
}

func TestAuditUserPasswordChange(t *testing.T) {
	// password not changed
	templateUser := model.User{Name: "user1", PasswordLastSet: 10}
	user := model.User{Name: "user1", PasswordLastSet: 0}
	finding := auditUserPasswordChange(templateUser, user)
	checkFinding(t, finding, false, 0, "User password not changed: user1")

	// password changed
	templateUser = model.User{Name: "user1", PasswordLastSet: 10}
	user = model.User{Name: "user1", PasswordLastSet: 20}
	finding = auditUserPasswordChange(templateUser, user)
	checkFinding(t, finding, true, 1, "User password changed: user1")
}

func TestAuditUsers(t *testing.T) {
	state := model.State{}

	// no users in template, no users in state
	state.Users = make([]model.User, 0)
	template := model.Template{}
	findings := auditUsers(state, template)
	if len(findings) != 0 {
		t.Fatal("Expected 0 findings")
	}

	// no users in template, user in state
	state.Users = make([]model.User, 0)
	user := model.User{Name: "user1"}
	state.Users = append(state.Users, user)
	template = model.Template{}
	findings = auditUsers(state, template)
	if len(findings) != 0 {
		t.Fatal("Expected 0 findings")
	}

	// user to add in template, user in not state
	state.Users = make([]model.User, 0)
	template = model.Template{}
	template.Users = make([]model.User, 0)
	tUser := model.User{Name: "user1", AccountState: model.ObjectStateAdd}
	template.Users = append(template.Users, tUser)
	findings = auditUsers(state, template)
	if len(findings) != 1 {
		t.Fatal("Expected 1 findings")
	}
	checkFinding(t, findings[0], false, 0, "User not added: user1")

	// user to add in template, user in state
	state.Users = make([]model.User, 0)
	user = model.User{Name: "user1", AccountActive: true}
	state.Users = append(state.Users, user)
	template = model.Template{}
	template.Users = make([]model.User, 0)
	tUser = model.User{Name: "user1", AccountState: model.ObjectStateAdd, AccountActive: true}
	template.Users = append(template.Users, tUser)
	findings = auditUsers(state, template)
	if len(findings) != 5 {
		t.Fatal("Expected 5 findings")
	}
	checkFinding(t, findings[0], true, 1, "User added: user1")
	checkFinding(t, findings[1], true, 1, "User active: user1")
	checkFinding(t, findings[2], false, 0, "User account does not expire: user1")
	checkFinding(t, findings[3], false, 0, "User password does not expire: user1")
	checkFinding(t, findings[4], false, 0, "User password not changed: user1")

	// user to keep in template, user in not state
	state.Users = make([]model.User, 0)
	template = model.Template{}
	template.Users = make([]model.User, 0)
	tUser = model.User{Name: "user1", AccountState: model.ObjectStateKeep}
	template.Users = append(template.Users, tUser)
	findings = auditUsers(state, template)
	if len(findings) != 1 {
		t.Fatal("Expected 1 findings")
	}
	checkFinding(t, findings[0], true, -1, "User not present: user1")

	// user to keep in template, user in state
	state.Users = make([]model.User, 0)
	user = model.User{Name: "user1", AccountActive: true}
	state.Users = append(state.Users, user)
	template = model.Template{}
	template.Users = make([]model.User, 0)
	tUser = model.User{Name: "user1", AccountState: model.ObjectStateKeep, AccountActive: true}
	template.Users = append(template.Users, tUser)
	findings = auditUsers(state, template)
	if len(findings) != 5 {
		t.Fatal("Expected 5 findings")
	}
	// don't give indication that user should be kept
	checkFinding(t, findings[0], false, 0, "User present: user1")
	checkFinding(t, findings[1], false, 0, "User active: user1")
	// these are OK to show if user are kept
	checkFinding(t, findings[2], false, 0, "User account does not expire: user1")
	checkFinding(t, findings[3], false, 0, "User password does not expire: user1")
	checkFinding(t, findings[4], false, 0, "User password not changed: user1")

	// user to remove in template, user in not state
	state.Users = make([]model.User, 0)
	template = model.Template{}
	template.Users = make([]model.User, 0)
	tUser = model.User{Name: "user1", AccountState: model.ObjectStateRemove}
	template.Users = append(template.Users, tUser)
	findings = auditUsers(state, template)
	if len(findings) != 1 {
		t.Fatal("Expected 1 findings")
	}
	checkFinding(t, findings[0], true, 1, "User removed: user1")

	// user to remove in template, user in state
	state.Users = make([]model.User, 0)
	user = model.User{Name: "user1", AccountActive: true}
	state.Users = append(state.Users, user)
	template = model.Template{}
	template.Users = make([]model.User, 0)
	tUser = model.User{Name: "user1", AccountState: model.ObjectStateRemove, AccountActive: true}
	template.Users = append(template.Users, tUser)
	findings = auditUsers(state, template)
	if len(findings) != 1 {
		t.Fatal("Expected 1 findings")
	}
	checkFinding(t, findings[0], false, 0, "User not removed: user1")
}
