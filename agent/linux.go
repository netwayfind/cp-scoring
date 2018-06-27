package main

import (
	"os/exec"
	"strings"

	"github.com/sumwonyuno/cp-scoring/model"
)

func getLinuxState() model.State {
	state := model.GetNewStateTemplate()
	state.Users = getUsersLinux()
	state.Groups = getGroupsLinux()
	return state
}

func getUsersLinux() []string {
	out, err := exec.Command("sh", "-c", "getent passwd | cut -d: -f1").Output()
	if err != nil {
		panic(err)

	}
	lines := strings.Split(string(out), "\n")
	// remove last line if empty
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}

func getGroupsLinux() map[string][]string {
	out, err := exec.Command("sh", "-c", "getent group | awk -F':' '{print $1 \":\" $4}'").Output()
	if err != nil {
		panic(err)
	}

	results := make(map[string][]string)
	lines := strings.Split(string(out), "\n")
	for i := range lines {
		entry := strings.Split(lines[i], ":")
		group := entry[0]
		if len(group) == 0 {
			continue
		}
		var members []string
		if len(entry) > 1 {
			members = strings.Split(entry[1], ",")
		}

		results[group] = members
	}

	return results
}
