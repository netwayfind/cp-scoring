package processing

import (
	"encoding/json"
	"testing"

	"github.com/sumwonyuno/cp-scoring/model"
)

func TestDiffReport(t *testing.T) {
	// empty reports
	report1 := model.Report{}
	report2 := model.Report{}
	changes, err := DiffReports(report1, report2)
	if err != nil {
		t.Fatal(err)
	}
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
	changes, err = DiffReports(report1, report2)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) > 0 {
		t.Fatal("Unexpected changes")
	}

	// test add finding
	report1 = model.Report{
		Timestamp: 14,
		Findings:  []model.Finding{},
	}
	report2 = model.Report{
		Timestamp: 15,
		Findings: []model.Finding{
			model.Finding{
				Message: "Test message",
				Value:   1,
				Show:    true,
			},
		},
	}
	changes, err = DiffReports(report1, report2)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatal("Unexpected number of changes:", len(changes))
	}
	if changes[0].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[0].Key != "Findings" {
		t.Fatal("Unexpected change key")
	}
	if changes[0].Item != "{\"Value\":1,\"Show\":true,\"Message\":\"Test message\"}" {
		t.Fatal("Unexpected change item")
	}

	// test remove finding
	report1 = model.Report{
		Timestamp: 16,
		Findings: []model.Finding{
			model.Finding{
				Message: "Test message",
				Value:   1,
				Show:    true,
			},
		},
	}
	report2 = model.Report{
		Timestamp: 17,
		Findings:  []model.Finding{},
	}
	changes, err = DiffReports(report1, report2)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatal("Unexpected number of changes:", len(changes))
	}
	if changes[0].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[0].Key != "Findings" {
		t.Fatal("Unexpected change key")
	}
	if changes[0].Item != "{\"Value\":1,\"Show\":true,\"Message\":\"Test message\"}" {
		t.Fatal("Unexpected change item")
	}

	// test changed finding
	report1 = model.Report{
		Timestamp: 18,
		Findings: []model.Finding{
			model.Finding{
				Message: "Test message",
				Value:   1,
				Show:    true,
			},
		},
	}
	report2 = model.Report{
		Timestamp: 19,
		Findings: []model.Finding{
			model.Finding{
				Message: "Test message",
				Value:   0,
				Show:    false,
			},
		},
	}
	changes, err = DiffReports(report1, report2)
	if err != nil {
		t.Fatal(err)
	}
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
	if changes[0].Item != "{\"Value\":1,\"Show\":true,\"Message\":\"Test message\"}" {
		t.Fatal("Unexpected change item")
	}
	// then added
	if changes[1].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[1].Key != "Findings" {
		t.Fatal("Unexpected change key")
	}
	if changes[1].Item != "{\"Value\":0,\"Show\":false,\"Message\":\"Test message\"}" {
		t.Fatal("Unexpected change item")
	}
}

func TestDiffState(t *testing.T) {
	// empty reports
	state1 := model.State{}
	state2 := model.State{}
	changes, err := DiffStates(state1, state2)
	if err != nil {
		t.Fatal(err)
	}
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
	changes, err = DiffStates(state1, state2)
	if err != nil {
		t.Fatal(err)
	}
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
	changes, err = DiffStates(state1, state2)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 6 {
		t.Fatal("Unexpected number of changes:", len(changes))
	}
	// group name
	if changes[0].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[0].Key != "Groups" {
		t.Fatal("Unexpected change key")
	}
	if changes[0].Item != "{\"Group\":\"Users\",\"Member\":\"bob\"}" {
		t.Fatal("Unexpected change item")
	}
	// group members
	if changes[1].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[1].Key != "Groups" {
		t.Fatal("Unexpected change key")
	}
	if changes[1].Item != "{\"Group\":\"Users\"}" {
		t.Fatal("Unexpected change item")
	}
	// network connections
	if changes[2].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[2].Key != "Network Connections" {
		t.Fatal("Unexpected change key")
	}
	expected, err := json.Marshal(state2.NetworkConnections[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[2].Item != string(expected) {
		t.Fatal("Unexpected change item")
	}
	// processes
	if changes[3].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[3].Key != "Processes" {
		t.Fatal("Unexpected change key")
	}
	expected, err = json.Marshal(state2.Processes[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[3].Item != string(expected) {
		t.Fatal("Unexpected change item")
	}
	// software
	if changes[4].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[4].Key != "Software" {
		t.Fatal("Unexpected change key")
	}
	expected, err = json.Marshal(state2.Software[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[4].Item != string(expected) {
		t.Fatal("Unexpected change item")
	}
	// users
	if changes[5].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[5].Key != "Users" {
		t.Fatal("Unexpected change key")
	}
	expected, err = json.Marshal(state2.Users[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[5].Item != string(expected) {
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
	changes, err = DiffStates(state1, state2)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 6 {
		t.Fatal("Unexpected number of changes:", len(changes))
	}
	// groups (Empty)
	if changes[0].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[0].Key != "Groups" {
		t.Fatal("Unexpected change key")
	}
	if changes[0].Item != "{\"Group\":\"Empty\"}" {
		t.Fatal("Unexpected change item")
	}
	// groups (Users)
	if changes[1].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[1].Key != "Groups" {
		t.Fatal("Unexpected change key")
	}
	if changes[1].Item != "{\"Group\":\"Users\",\"Member\":\"bob\"}" {
		t.Fatal("Unexpected change item")
	}
	// network connections
	if changes[2].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[2].Key != "Network Connections" {
		t.Fatal("Unexpected change key")
	}
	expected, err = json.Marshal(state1.NetworkConnections[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[2].Item != string(expected) {
		t.Fatal("Unexpected change item")
	}
	// processes
	if changes[3].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[3].Key != "Processes" {
		t.Fatal("Unexpected change key")
	}
	expected, err = json.Marshal(state1.Processes[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[3].Item != string(expected) {
		t.Fatal("Unexpected change item")
	}
	// software
	if changes[4].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[4].Key != "Software" {
		t.Fatal("Unexpected change key")
	}
	expected, err = json.Marshal(state1.Software[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[4].Item != string(expected) {
		t.Fatal("Unexpected change item")
	}
	// users
	if changes[5].Type != "Removed" {
		t.Fatal("Unexpected change type")
	}
	if changes[5].Key != "Users" {
		t.Fatal("Unexpected change key")
	}
	expected, err = json.Marshal(state1.Users[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[5].Item != string(expected) {
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
	changes, err = DiffStates(state1, state2)
	if err != nil {
		t.Fatal(err)
	}
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
	if changes[0].Item != "{\"Group\":\"Users\",\"Member\":\"alice\"}" {
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
	expected, err = json.Marshal(state1.NetworkConnections[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[1].Item != string(expected) {
		t.Fatal("Unexpected change item")
	}
	// then added
	if changes[2].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[2].Key != "Network Connections" {
		t.Fatal("Unexpected change key")
	}
	expected, err = json.Marshal(state2.NetworkConnections[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[2].Item != string(expected) {
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
	expected, err = json.Marshal(state1.Processes[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[3].Item != string(expected) {
		t.Fatal("Unexpected change item")
	}
	// then added
	if changes[4].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[4].Key != "Processes" {
		t.Fatal("Unexpected change key")
	}
	expected, err = json.Marshal(state2.Processes[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[4].Item != string(expected) {
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
	expected, err = json.Marshal(state1.Software[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[5].Item != string(expected) {
		t.Fatal("Unexpected change item")
	}
	// then added
	if changes[6].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[6].Key != "Software" {
		t.Fatal("Unexpected change key")
	}
	expected, err = json.Marshal(state2.Software[0])
	if err != nil {
		t.Fatal(err)
	}
	if changes[6].Item != string(expected) {
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
	expected, err = json.Marshal(state1.Users[1])
	if err != nil {
		t.Fatal(err)
	}
	if changes[7].Item != string(expected) {
		t.Fatal("Unexpected change item")
	}
	// then added
	if changes[8].Type != "Added" {
		t.Fatal("Unexpected change type")
	}
	if changes[8].Key != "Users" {
		t.Fatal("Unexpected change key")
	}
	expected, err = json.Marshal(state2.Users[1])
	if err != nil {
		t.Fatal(err)
	}
	if changes[8].Item != string(expected) {
		t.Fatal("Unexpected change item")
	}
}
