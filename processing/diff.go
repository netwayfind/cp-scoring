package processing

import (
	"sort"

	"github.com/cnf/structhash"
	"github.com/netwayfind/cp-scoring/model"
)

type Change struct {
	Type string
	Key  string
	Item interface{}
}

type DocumentDiff struct {
	Timestamp int64
	Changes   []Change
}

type groupMemberEntry struct {
	Group  string
	Member string
}

func Diff(previous map[string]map[string]interface{}, current map[string]map[string]interface{}) []Change {
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
				change := Change{Type: "Removed", Key: key, Item: previousEntries[h]}
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
				change := Change{Type: "Added", Key: key, Item: entries[h]}
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

func GetReportEntries(report model.Report) map[string]map[string]interface{} {
	entries := make(map[string]map[string]interface{})

	// findings
	findings := make(map[string]interface{})
	entries["Findings"] = findings
	for _, finding := range report.Findings {
		h := string(structhash.Sha1(finding, 1))
		findings[h] = finding
	}

	return entries
}

func GetStateEntries(state model.State) map[string]map[string]interface{} {
	entries := make(map[string]map[string]interface{})

	// users
	users := make(map[string]interface{})
	entries["Users"] = users
	for _, user := range state.Users {
		h := string(structhash.Sha1(user, 1))
		users[h] = user
	}

	// groups
	groups := make(map[string]interface{})
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
		h := string(structhash.Sha1(name, 1))
		groups[h] = name

		// group members
		members := state.Groups[name]
		for _, member := range members {
			entry := groupMemberEntry{Group: name, Member: member.Name}
			h := string(structhash.Sha1(entry, 1))
			groups[h] = entry
		}
	}

	// software
	software := make(map[string]interface{})
	entries["Software"] = software
	for _, sw := range state.Software {
		h := string(structhash.Sha1(sw, 1))
		software[h] = sw
	}

	// processes
	processes := make(map[string]interface{})
	entries["Processes"] = processes
	for _, process := range state.Processes {
		h := string(structhash.Sha1(process, 1))
		processes[h] = process
	}

	// network connections
	conns := make(map[string]interface{})
	entries["Network Connections"] = conns
	for _, conn := range state.NetworkConnections {
		h := string(structhash.Sha1(conn, 1))
		conns[h] = conn
	}

	return entries
}
