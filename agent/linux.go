// +build linux

package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/sumwonyuno/cp-scoring/model"
)

func getState() model.State {
	state := model.GetNewStateTemplate()
	state.Users = getUsers()
	state.Groups = getGroups()
	return state
}

func getUsers() []model.User {
	// use separate list to keep track of user names
	usernames := make([]string, 0)
	// use map to keep track of users (no guarantee user names in same order)
	// assume that all users exist in all files
	usersMap := make(map[string]model.User)

	// get user and uid
	bs, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		log.Fatal("ERROR: cannot get users info;", err)
	}
	for _, line := range strings.Split(string(bs), "\n") {
		tokens := strings.Split(line, ":")
		if len(tokens) != 7 {
			continue
		}
		username, id := tokens[0], tokens[2]
		var entry model.User
		entry.Name = username
		entry.ID = id
		entry.AccountPresent = true
		usernames = append(usernames, username)
		usersMap[username] = entry
	}

	// get other user information (sensitive)
	bs, err = ioutil.ReadFile("/etc/shadow")
	if err != nil {
		log.Fatal("ERROR: cannot get users info;", err)
	}
	for _, line := range strings.Split(string(bs), "\n") {
		tokens := strings.Split(line, ":")
		if len(tokens) != 9 {
			continue
		}
		username, passwordHash, unixDayPasswordLastChange, unixDayPasswordExpires, unixDayAccountDisabled := tokens[0], tokens[1], tokens[2], tokens[4], tokens[7]
		entry, present := usersMap[username]
		// user should be in /etc/passwd, but if not, create entry
		if !present {
			entry = model.User{}
			entry.Name = username
			// user doesn't exist
			entry.AccountPresent = false
			usernames = append(usernames, username)
			usersMap[username] = entry
		}
		// account active
		// user is locked if password hash entry starts with !
		// do not store the password hash
		// 33 (dec) is !
		if passwordHash[0] == 33 {
			entry.AccountActive = false
		} else {
			entry.AccountActive = true
		}
		// password last changed
		numDays, err := strconv.Atoi(unixDayPasswordLastChange)
		// should not have error
		if err == nil {
			// convert days to seconds (unix timestamp)
			entry.PasswordLastSet = int64(numDays * 86400)
		}
		// password expires
		if unixDayPasswordExpires == "99999" {
			entry.PasswordExpires = false
		} else {
			entry.PasswordExpires = true
		}
		// account expired
		if len(unixDayAccountDisabled) == 0 {
			entry.AccountExpires = false
		} else {
			entry.AccountExpires = true
		}

		// save changes
		usersMap[username] = entry
	}

	// turn users map into users array
	users := make([]model.User, len(usernames))
	for i, username := range usernames {
		users[i] = usersMap[username]
	}

	return users
}

func getGroups() map[string][]string {
	bs, err := ioutil.ReadFile("/etc/group")
	if err != nil {
		log.Fatal("ERROR: cannot get groups;", err)
	}

	groups := make(map[string][]string)
	for _, line := range strings.Split(string(bs), "\n") {
		tokens := strings.Split(line, ":")
		if len(tokens) != 4 {
			continue
		}
		group, membersStr := tokens[0], tokens[3]
		groups[group] = strings.Split(membersStr, ",")
	}

	return groups
}
