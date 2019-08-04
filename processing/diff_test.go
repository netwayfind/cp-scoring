package processing

import (
	"testing"

	"github.com/cnf/structhash"

	"github.com/sumwonyuno/cp-scoring/model"
)

func TestGetReportEntries(t *testing.T) {
	// empty report
	report := model.Report{}
	entries := GetReportEntries(report)
	if len(entries) != 1 {
		t.Fatal("Unexpected entry count:", len(entries))
	}
	findings, present := entries["Findings"]
	if !present {
		t.Fatal("Expected findings key")
	}
	if len(findings) != 0 {
		t.Fatal("Unexpected findings count:", len(entries))
	}

	// sample report with findings
	finding := model.Finding{
		Show:    true,
		Value:   1,
		Message: "test",
	}
	report = model.Report{
		Timestamp: 15,
		Findings:  []model.Finding{finding},
	}
	entries = GetReportEntries(report)
	if len(entries) != 1 {
		t.Fatal("Unexpected entry count:", len(entries))
	}
	if len(entries) != 1 {
		t.Fatal("Unexpected entry count:", len(entries))
	}
	findings, present = entries["Findings"]
	if !present {
		t.Fatal("Expected findings key")
	}
	if len(findings) != 1 {
		t.Fatal("Unexpected findings count:", len(entries))
	}
	if _, present = findings[string(structhash.Sha1(finding, 1))]; !present {
		t.Fatal("Did not find expected entry")
	}
}

func TestGetStateEntries(t *testing.T) {
	// empty state
	state := model.State{}
	entries := GetStateEntries(state)
	if len(entries) != 5 {
		t.Fatal("Unexpected entry count:", len(entries))
	}
	users, present := entries["Users"]
	if !present {
		t.Fatal("Expected users key")
	}
	if len(users) != 0 {
		t.Fatal("Unexpected users count:", len(users))
	}
	groups, present := entries["Groups"]
	if !present {
		t.Fatal("Expected groups key")
	}
	if len(groups) != 0 {
		t.Fatal("Unexpected groups count:", len(groups))
	}
	software, present := entries["Software"]
	if !present {
		t.Fatal("Expected software key")
	}
	if len(software) != 0 {
		t.Fatal("Unexpected software count:", len(software))
	}
	processes, present := entries["Processes"]
	if !present {
		t.Fatal("Expected processes key")
	}
	if len(processes) != 0 {
		t.Fatal("Unexpected processes count:", len(processes))
	}
	conns, present := entries["Network Connections"]
	if !present {
		t.Fatal("Expected network connections key")
	}
	if len(conns) != 0 {
		t.Fatal("Unexpected network connections count:", len(conns))
	}

	// sample state with software
	sw := model.Software{
		Name:    "test-software",
		Version: "0.1.0",
	}
	state = model.State{
		Timestamp: 15,
		Software:  []model.Software{sw},
	}
	entries = GetStateEntries(state)
	if len(entries) != 5 {
		t.Fatal("Unexpected entry count:", len(entries))
	}
	users, present = entries["Users"]
	if !present {
		t.Fatal("Expected users key")
	}
	if len(users) != 0 {
		t.Fatal("Unexpected users count:", len(users))
	}
	groups, present = entries["Groups"]
	if !present {
		t.Fatal("Expected groups key")
	}
	if len(groups) != 0 {
		t.Fatal("Unexpected groups count:", len(groups))
	}
	software, present = entries["Software"]
	if !present {
		t.Fatal("Expected software key")
	}
	if len(software) != 1 {
		t.Fatal("Unexpected software count:", len(software))
	}
	_, present = software[string(structhash.Sha1(sw, 1))]
	if !present {
		t.Fatal("Did not find expected entry")
	}
	processes, present = entries["Processes"]
	if !present {
		t.Fatal("Expected processes key")
	}
	if len(processes) != 0 {
		t.Fatal("Unexpected processes count:", len(processes))
	}
	conns, present = entries["Network Connections"]
	if !present {
		t.Fatal("Expected network connections key")
	}
	if len(conns) != 0 {
		t.Fatal("Unexpected network connections count:", len(conns))
	}
}

func TestDiffForReports(t *testing.T) {
	// empty reports
	report1 := model.Report{}
	report2 := model.Report{}
	report1Entries := GetReportEntries(report1)
	report2Entries := GetReportEntries(report2)
	changes := Diff(report1Entries, report2Entries)
	if len(changes) > 0 {
		t.Fatal("Unexpected changes")
	}

	// no finding reports
	report1 = model.Report{
		Timestamp: 14,
	}
	report2 = model.Report{
		Timestamp: 15,
	}
	report1Entries = GetReportEntries(report1)
	report2Entries = GetReportEntries(report2)
	changes = Diff(report1Entries, report2Entries)
	if len(changes) > 0 {
		t.Fatal("Unexpected changes")
	}

	// test add finding
	report1 = model.Report{
		Timestamp: 14,
		Findings:  []model.Finding{},
	}
	finding := model.Finding{
		Message: "Test message",
		Value:   1,
		Show:    true,
	}
	report2 = model.Report{
		Timestamp: 15,
		Findings:  []model.Finding{finding},
	}
	report1Entries = GetReportEntries(report1)
	report2Entries = GetReportEntries(report2)
	changes = Diff(report1Entries, report2Entries)
	if len(changes) != 1 {
		t.Fatal("Unexpected number of changes:", len(changes))
	}
	if changes[0].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[0].Key != "Findings" {
		t.Fatal("Unexpected change key")
	}
	if changes[0].Item != finding {
		t.Fatal("Unexpected change item")
	}

	// test remove finding
	finding = model.Finding{
		Message: "Test message",
		Value:   1,
		Show:    true,
	}
	report1 = model.Report{
		Timestamp: 16,
		Findings:  []model.Finding{finding},
	}
	report2 = model.Report{
		Timestamp: 17,
		Findings:  []model.Finding{},
	}
	report1Entries = GetReportEntries(report1)
	report2Entries = GetReportEntries(report2)
	changes = Diff(report1Entries, report2Entries)
	if len(changes) != 1 {
		t.Fatal("Unexpected number of changes:", len(changes))
	}
	if changes[0].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[0].Key != "Findings" {
		t.Fatal("Unexpected change key")
	}
	if changes[0].Item != finding {
		t.Fatal("Unexpected change item")
	}

	// test changed finding
	finding1 := model.Finding{
		Message: "Test message",
		Value:   1,
		Show:    true,
	}
	report1 = model.Report{
		Timestamp: 18,
		Findings:  []model.Finding{finding1},
	}
	finding2 := model.Finding{
		Message: "Test message",
		Value:   0,
		Show:    false,
	}
	report2 = model.Report{
		Timestamp: 19,
		Findings:  []model.Finding{finding2},
	}
	report1Entries = GetReportEntries(report1)
	report2Entries = GetReportEntries(report2)
	changes = Diff(report1Entries, report2Entries)
	if len(changes) != 2 {
		t.Fatal("Unexpected number of changes:", len(changes))
	}
	// removed should be first
	if changes[0].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[0].Key != "Findings" {
		t.Fatal("Unexpected change key")
	}
	if changes[0].Item != finding1 {
		t.Fatal("Unexpected change item")
	}
	// then added
	if changes[1].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[1].Key != "Findings" {
		t.Fatal("Unexpected change key")
	}
	if changes[1].Item != finding2 {
		t.Fatal("Unexpected change item")
	}
}

func TestDiffForStates(t *testing.T) {
	// empty reports
	state1 := model.State{}
	state2 := model.State{}
	state1Entries := GetStateEntries(state1)
	state2Entries := GetStateEntries(state2)
	changes := Diff(state1Entries, state2Entries)
	if len(changes) > 0 {
		t.Fatal("Unexpected changes")
	}

	// empty state entries
	state1 = model.State{
		Timestamp: 14,
	}
	state2 = model.State{
		Timestamp: 15,
	}
	state1Entries = GetStateEntries(state1)
	state2Entries = GetStateEntries(state2)
	changes = Diff(state1Entries, state2Entries)
	if len(changes) > 0 {
		t.Fatal("Unexpected changes")
	}

	// test add entries
	state1 = model.State{
		Timestamp:          14,
		Users:              []model.User{},
		Groups:             map[string][]model.GroupMember{},
		Software:           []model.Software{},
		Processes:          []model.Process{},
		NetworkConnections: []model.NetworkConnection{},
	}
	state2 = model.State{
		Timestamp: 15,
		Users: []model.User{
			model.User{
				Name:          "bob",
				AccountActive: true,
			},
		},
		Groups: map[string][]model.GroupMember{
			"Users": []model.GroupMember{
				model.GroupMember{
					Name: "bob",
				},
			},
		},
		Software: []model.Software{
			model.Software{
				Name:    "test-software",
				Version: "0.1.0",
			},
		},
		Processes: []model.Process{
			model.Process{
				PID:         5,
				User:        "user",
				CommandLine: "cmd 1 2 3",
			},
		},
		NetworkConnections: []model.NetworkConnection{
			model.NetworkConnection{
				LocalAddress: "127.0.0.1",
				LocalPort:    "80",
				State:        model.NetworkConnectionListen,
			},
		},
	}
	state1Entries = GetStateEntries(state1)
	state2Entries = GetStateEntries(state2)
	changes = Diff(state1Entries, state2Entries)
	if len(changes) != 6 {
		t.Fatal("Unexpected number of changes:", len(changes))
	}
	// group members
	if changes[0].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[0].Key != "Groups" {
		t.Fatal("Unexpected change key")
	}
	entry := groupMemberEntry{
		Group:  "Users",
		Member: "bob",
	}
	if changes[0].Item != entry {
		t.Fatal("Unexpected change item")
	}
	// group
	if changes[1].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[1].Key != "Groups" {
		t.Fatal("Unexpected change key")
	}
	if changes[1].Item != "Users" {
		t.Fatal("Unexpected change item")
	}
	// network connections
	if changes[2].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[2].Key != "Network Connections" {
		t.Fatal("Unexpected change key")
	}
	if changes[2].Item != state2.NetworkConnections[0] {
		t.Fatal("Unexpected change item")
	}
	// processes
	if changes[3].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[3].Key != "Processes" {
		t.Fatal("Unexpected change key")
	}
	if changes[3].Item != state2.Processes[0] {
		t.Fatal("Unexpected change item")
	}
	// software
	if changes[4].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[4].Key != "Software" {
		t.Fatal("Unexpected change key")
	}
	if changes[4].Item != state2.Software[0] {
		t.Fatal("Unexpected change item")
	}
	// users
	if changes[5].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[5].Key != "Users" {
		t.Fatal("Unexpected change key")
	}
	if changes[5].Item != state2.Users[0] {
		t.Fatal("Unexpected change item")
	}

	// test remove entries
	state1 = model.State{
		Timestamp: 16,
		Users: []model.User{
			model.User{
				Name:          "bob",
				AccountActive: true,
			},
		},
		Groups: map[string][]model.GroupMember{
			"Users": []model.GroupMember{
				model.GroupMember{
					Name: "bob",
				},
			},
			"Empty": []model.GroupMember{},
		},
		Software: []model.Software{
			model.Software{
				Name:    "test-software",
				Version: "0.1.0",
			},
		},
		Processes: []model.Process{
			model.Process{
				PID:         5,
				User:        "user",
				CommandLine: "cmd 1 2 3",
			},
		},
		NetworkConnections: []model.NetworkConnection{
			model.NetworkConnection{
				LocalAddress: "127.0.0.1",
				LocalPort:    "80",
				State:        model.NetworkConnectionListen,
			},
		},
	}
	state2 = model.State{
		Timestamp: 17,
		Users:     []model.User{},
		Groups: map[string][]model.GroupMember{
			"Users": []model.GroupMember{},
		},
		Software:           []model.Software{},
		Processes:          []model.Process{},
		NetworkConnections: []model.NetworkConnection{},
	}
	state1Entries = GetStateEntries(state1)
	state2Entries = GetStateEntries(state2)
	changes = Diff(state1Entries, state2Entries)
	if len(changes) != 6 {
		t.Fatal("Unexpected number of changes:", len(changes))
	}
	// groups (Users)
	if changes[0].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[0].Key != "Groups" {
		t.Fatal("Unexpected change key")
	}
	item := groupMemberEntry{
		Group:  "Users",
		Member: "bob",
	}
	if changes[0].Item != item {
		t.Fatal("Unexpected change item")
	}
	// groups (Empty)
	if changes[1].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[1].Key != "Groups" {
		t.Fatal("Unexpected change key")
	}
	if changes[1].Item != "Empty" {
		t.Fatal("Unexpected change item")
	}
	// network connections
	if changes[2].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[2].Key != "Network Connections" {
		t.Fatal("Unexpected change key")
	}
	if changes[2].Item != state1.NetworkConnections[0] {
		t.Fatal("Unexpected change item")
	}
	// processes
	if changes[3].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[3].Key != "Processes" {
		t.Fatal("Unexpected change key")
	}
	if changes[3].Item != state1.Processes[0] {
		t.Fatal("Unexpected change item")
	}
	// software
	if changes[4].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[4].Key != "Software" {
		t.Fatal("Unexpected change key")
	}
	if changes[4].Item != state1.Software[0] {
		t.Fatal("Unexpected change item")
	}
	// users
	if changes[5].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[5].Key != "Users" {
		t.Fatal("Unexpected change key")
	}
	if changes[5].Item != state1.Users[0] {
		t.Fatal("Unexpected change item")
	}

	// test changed entries
	state1 = model.State{
		Timestamp: 18,
		Users: []model.User{
			model.User{
				Name:          "alice",
				AccountActive: true,
			},
			model.User{
				Name:          "bob",
				AccountActive: true,
			},
		},
		Groups: map[string][]model.GroupMember{
			"Users": []model.GroupMember{
				model.GroupMember{
					Name: "bob",
				},
			},
		},
		Software: []model.Software{
			model.Software{
				Name:    "test-software",
				Version: "0.1.0",
			},
		},
		Processes: []model.Process{
			model.Process{
				PID:         5,
				User:        "user",
				CommandLine: "cmd 1 2 3",
			},
		},
		NetworkConnections: []model.NetworkConnection{
			model.NetworkConnection{
				LocalAddress:  "127.0.0.1",
				LocalPort:     "45678",
				RemoteAddress: "192.168.1.1",
				RemotePort:    "443",
				State:         model.NetworkConnectionEstablished,
			},
		},
	}
	state2 = model.State{
		Timestamp: 19,
		Users: []model.User{
			model.User{
				Name:          "alice",
				AccountActive: true,
			},
			model.User{
				Name:          "bob",
				AccountActive: false,
			},
		},
		Groups: map[string][]model.GroupMember{
			"Users": []model.GroupMember{
				model.GroupMember{
					Name: "alice",
				},
				model.GroupMember{
					Name: "bob",
				},
			},
		},
		Software: []model.Software{
			model.Software{
				Name:    "test-software",
				Version: "0.2.0",
			},
		},
		Processes: []model.Process{
			model.Process{
				PID:         6,
				User:        "user",
				CommandLine: "cmd 1 2 3",
			},
		},
		NetworkConnections: []model.NetworkConnection{
			model.NetworkConnection{
				LocalAddress:  "127.0.0.1",
				LocalPort:     "46000",
				RemoteAddress: "192.168.1.1",
				RemotePort:    "443",
				State:         model.NetworkConnectionEstablished,
			},
		},
	}
	state1Entries = GetStateEntries(state1)
	state2Entries = GetStateEntries(state2)
	changes = Diff(state1Entries, state2Entries)
	if len(changes) != 9 {
		t.Fatal("Unexpected number of changes:", len(changes))
	}
	// groups
	if changes[0].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[0].Key != "Groups" {
		t.Fatal("Unexpected change key")
	}
	item = groupMemberEntry{
		Group:  "Users",
		Member: "alice",
	}
	if changes[0].Item != item {
		t.Fatal("Unexpected change item")
	}
	// network connections
	// removed should be first
	if changes[1].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[1].Key != "Network Connections" {
		t.Fatal("Unexpected change key")
	}
	if changes[1].Item != state1.NetworkConnections[0] {
		t.Fatal("Unexpected change item")
	}
	// then added
	if changes[2].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[2].Key != "Network Connections" {
		t.Fatal("Unexpected change key")
	}
	if changes[2].Item != state2.NetworkConnections[0] {
		t.Fatal("Unexpected change item")
	}
	// processes
	// removed should be first
	if changes[3].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[3].Key != "Processes" {
		t.Fatal("Unexpected change key")
	}
	if changes[3].Item != state1.Processes[0] {
		t.Fatal("Unexpected change item")
	}
	// then added
	if changes[4].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[4].Key != "Processes" {
		t.Fatal("Unexpected change key")
	}
	if changes[4].Item != state2.Processes[0] {
		t.Fatal("Unexpected change item")
	}
	// software
	// removed should be first
	if changes[5].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[5].Key != "Software" {
		t.Fatal("Unexpected change key")
	}
	if changes[5].Item != state1.Software[0] {
		t.Fatal("Unexpected change item")
	}
	// then added
	if changes[6].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[6].Key != "Software" {
		t.Fatal("Unexpected change key")
	}
	if changes[6].Item != state2.Software[0] {
		t.Fatal("Unexpected change item")
	}
	// users
	// removed should be first
	if changes[7].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[7].Key != "Users" {
		t.Fatal("Unexpected change key")
	}
	if changes[7].Item != state1.Users[1] {
		t.Fatal("Unexpected change item")
	}
	// then added
	if changes[8].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[8].Key != "Users" {
		t.Fatal("Unexpected change key")
	}
	if changes[8].Item != state2.Users[1] {
		t.Fatal("Unexpected change item")
	}
}
