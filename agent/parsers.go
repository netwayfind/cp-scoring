package main

import (
	"strconv"
	"strings"

	"github.com/sumwonyuno/cp-scoring/model"
)

func parseEtcPasswd(bs []byte) map[string]model.User {
	users := make(map[string]model.User)
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
		users[username] = entry
	}
	return users
}

func parseEtcShadow(bs []byte) map[string]model.User {
	users := make(map[string]model.User)
	for _, line := range strings.Split(string(bs), "\n") {
		tokens := strings.Split(line, ":")
		if len(tokens) != 9 {
			continue
		}
		username, passwordHash, unixDayPasswordLastChange, unixDayPasswordExpires, unixDayAccountDisabled := tokens[0], tokens[1], tokens[2], tokens[4], tokens[7]
		entry := model.User{}
		entry.Name = username

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
		users[username] = entry
	}
	return users
}

func mergeUserMaps(userMapEtcPasswd map[string]model.User, userMapEtcShadow map[string]model.User) []model.User {
	users := make([]model.User, 0)

	// assume users exist in both /etc/passwd and /etc/shadow
	for username, entry := range userMapEtcPasswd {
		entryShadow := userMapEtcShadow[username]
		// set settings from shadow into entry
		entry.AccountActive = entryShadow.AccountActive
		entry.AccountExpires = entryShadow.AccountExpires
		entry.PasswordLastSet = entryShadow.PasswordLastSet
		entry.PasswordExpires = entryShadow.PasswordExpires

		users = append(users, entry)
	}

	return users
}
