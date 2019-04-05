package processing

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/sumwonyuno/cp-scoring/model"
)

type Change struct {
	Type string
	Key  string
	Item string
}

func Diff(previous map[string]map[string]bool, current map[string]map[string]bool) []Change {
	changes := make([]Change, 0)

	// check keys in previous and current
	// sort keys
	previousKeys := make([]string, len(current))
	index := 0
	for key := range previous {
		previousKeys[index] = key
		index++
	}
	sort.Strings(previousKeys)
	for _, key := range previousKeys {
		entries, present := current[key]
		if !present {
			change := Change{Type: "Not present", Key: key}
			changes = append(changes, change)
			continue
		}

		previousEntries := previous[key]
		// sort keys
		previousEntriesKeys := make([]string, len(previousEntries))
		index = 0
		for k := range previousEntries {
			previousEntriesKeys[index] = k
			index++
		}
		sort.Strings(previousEntriesKeys)

		// removed is in previous entries but not current entries
		for _, h := range previousEntriesKeys {
			if _, present := entries[h]; !present {
				change := Change{Type: "Removed", Key: key, Item: h}
				changes = append(changes, change)
			}
		}

		// sort keys
		entriesKeys := make([]string, len(entries))
		index = 0
		for k := range entries {
			entriesKeys[index] = k
			index++
		}
		sort.Strings(entriesKeys)

		// added is in current entries but not previous entries
		for _, h := range entriesKeys {
			if _, present := previousEntries[h]; !present {
				change := Change{Type: "Added", Key: key, Item: h}
				changes = append(changes, change)
			}
		}
	}
	// check for keys in current not in present
	// sort keys
	currentKeys := make([]string, len(current))
	index = 0
	for key := range current {
		currentKeys[index] = key
		index++
	}
	sort.Strings(currentKeys)
	for _, key := range currentKeys {
		if _, present := previous[key]; !present {
			change := Change{Type: "Not present", Key: key}
			changes = append(changes, change)
		}
	}

	return changes
}

func GetReportEntries(report model.Report) (map[string]map[string]bool, error) {
	entries := make(map[string]map[string]bool)

	// findings
	findings := make(map[string]bool)
	entries["Findings"] = findings
	for _, finding := range report.Findings {
		h, err := json.Marshal(finding)
		if err != nil {
			return nil, err
		}
		findings[string(h)] = true
	}

	return entries, nil
}

func GetStateEntries(state model.State) (map[string]map[string]bool, error) {
	entries := make(map[string]map[string]bool)

	// users
	users := make(map[string]bool)
	entries["Users"] = users
	for _, user := range state.Users {
		h, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}
		users[string(h)] = true
	}

	// groups
	groups := make(map[string]bool)
	entries["Groups"] = groups

	// use sorted key order
	groupNames := make([]string, len(state.Groups))
	index := 0
	for name := range state.Groups {
		groupNames[index] = name
		index++
	}
	sort.Strings(groupNames)

	for _, name := range groupNames {
		// group name
		groupKey := fmt.Sprintf("{\"Group\":\"%s\"}", name)
		groups[groupKey] = true

		// group members
		members := state.Groups[name]
		for _, member := range members {
			key := fmt.Sprintf("{\"Group\":\"%s\",\"Member\":\"%s\"}", name, member.Name)
			groups[key] = true
		}
	}

	// software
	software := make(map[string]bool)
	entries["Software"] = software
	for _, sw := range state.Software {
		h, err := json.Marshal(sw)
		if err != nil {
			return nil, err
		}
		software[string(h)] = true
	}

	// processes
	processes := make(map[string]bool)
	entries["Processes"] = processes
	for _, process := range state.Processes {
		h, err := json.Marshal(process)
		if err != nil {
			return nil, err
		}
		processes[string(h)] = true
	}

	// network connections
	conns := make(map[string]bool)
	entries["Network Connections"] = conns
	for _, conn := range state.NetworkConnections {
		h, err := json.Marshal(conn)
		if err != nil {
			return nil, err
		}
		conns[string(h)] = true
	}

	return entries, nil
}
