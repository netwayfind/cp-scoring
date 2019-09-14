package agent

import (
	"testing"
)

func TestGetCurrentHost(t *testing.T) {
	host := GetCurrentHost()
	if host == nil {
		t.Fatal("Could not get current host")
	}
}

func TestCollectState(t *testing.T) {
	state := GetState()
	if len(state.Hostname) == 0 {
		t.Fatal("hostname not set")
	}
	if len(state.OS) == 0 {
		t.Fatal("OS not set")
	}
	if state.Timestamp == 0 {
		t.Fatal("timestamp not set")
	}
	if state.Users == nil && state.Errors == nil {
		t.Fatal("users not set, no errors")
	}
	if state.Groups == nil && state.Errors == nil {
		t.Fatal("groups not set, no errors")
	}
	if state.Processes == nil && state.Errors == nil {
		t.Fatal("processes not set, no errors")
	}
	if state.Software == nil && state.Errors == nil {
		t.Fatal("software not set, no errors")
	}
	if state.NetworkConnections == nil && state.Errors == nil {
		t.Fatal("network connections not set, no errors")
	}
	if state.ScheduledTasks == nil && state.Errors == nil {
		t.Fatal("scheduled tasks not set, no errors")
	}
	if state.WindowsFirewallProfiles == nil && state.Errors == nil {
		t.Fatal("Windows firewall profiles not set, no errors")
	}
	if state.WindowsFirewallRules == nil && state.Errors == nil {
		t.Fatal("Windows firewall rules not set, no errors")
	}
	if state.Errors == nil {
		t.Fatal("No errors set")
	}
}
