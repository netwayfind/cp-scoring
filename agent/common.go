package agent

import (
	"log"
	"runtime"

	"github.com/sumwonyuno/cp-scoring/model"
)

func GetState(host model.CurrentHost) model.State {
	state := model.GetNewStateTemplate()
	errors := make([]string, 0)
	users, err := host.GetUsers()
	if err == nil {
		state.Users = users
	} else {
		errors = append(errors, "ERROR: unable to get users; "+err.Error())
	}
	groups, err := host.GetGroups()
	if err == nil {
		state.Groups = groups
	} else {
		errors = append(errors, "ERROR: cannot get groups; "+err.Error())
	}
	processes, err := host.GetProcesses()
	if err == nil {
		state.Processes = processes
	} else {
		errors = append(errors, "ERROR: cannot get processes; "+err.Error())
	}
	software, err := host.GetSoftware()
	if err == nil {
		state.Software = software
	} else {
		errors = append(errors, "ERROR: cannot get software; "+err.Error())
	}
	conns, err := host.GetNetworkConnections()
	if err == nil {
		state.NetworkConnections = conns
	} else {
		errors = append(errors, "ERROR: cannot get network connections; "+err.Error())
	}
	state.Errors = errors
	return state
}

func GetCurrentHost() model.CurrentHost {
	if runtime.GOOS == "linux" {
		return hostLinux{}
	} else if runtime.GOOS == "windows" {
		return hostWindows{}
	} else {
		log.Fatal("ERROR: unsupported platform: " + runtime.GOOS)
		return nil
	}
}
