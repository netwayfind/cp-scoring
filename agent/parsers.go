package main

import (
	"bytes"
	"encoding/csv"
	"math"
	"strconv"
	"strings"
	"time"

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

func parseWindowsUsers(bs []byte) []model.User {
	users := make([]model.User, 0)
	c := csv.NewReader(bytes.NewReader(bs))
	records, err := c.ReadAll()
	if err != nil {
		return users
	}
	for i, row := range records {
		// header row
		if i == 0 {
			continue
		}
		// must have exactly 6 columns, or else ignore line
		if len(row) != 6 {
			continue
		}
		user := model.User{}
		user.AccountPresent = true
		// Name,SID,Enabled,AccountExpires,PasswordLastSet,PasswordExpires
		user.Name = row[0]
		user.ID = row[1]
		if row[2] == "True" {
			user.AccountActive = true
		} else if row[2] == "False" {
			user.AccountActive = false
		} else {
			user.AccountActive = false
		}
		if len(row[3]) == 0 {
			user.AccountExpires = false
		} else {
			user.AccountExpires = true
		}
		if len(row[4]) == 0 {
			// never expire is max value
			user.PasswordLastSet = math.MaxInt64
		} else {
			timezone, _ := time.Now().Zone()
			value := row[4] + " " + timezone
			layout := "1/2/2006 3:04:05 PM MST"
			t, err := time.Parse(layout, value)
			if err == nil {
				user.PasswordLastSet = t.Unix()
			}
		}
		if len(row[5]) == 0 {
			user.PasswordExpires = false
		} else {
			user.PasswordExpires = true
		}
		users = append(users, user)
	}

	return users
}
