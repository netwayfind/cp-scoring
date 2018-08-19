// +build windows

package main

import (
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/sumwonyuno/cp-scoring/model"
)

func getState() model.State {
	state := model.GetNewStateTemplate()
	state.Users = getUsers()
	state.Groups = getGroups()
	return state
}

type userinfo struct {
	username string
}

func getUsers() []model.User {
	users := make([]model.User, 0)

	// get users and info
	out, err := exec.Command("wmic", "UserAccount", "get", "Name,SID").Output()
	if err != nil {
		log.Fatal("ERROR: cannot get users info;", err)
	}
	var posName int
	var posSID int
	for i, line := range strings.Split(string(out), "\r\n") {
		if len(line) <= 1 {
			continue
		}

		var user model.User

		// parse header for positions
		if i == 0 {
			// assume all exist
			posName = strings.Index(line, "Name")
			posSID = strings.Index(line, "SID")
			continue
		}

		user.AccountPresent = true
		user.Name = strings.TrimSpace(line[posName:posSID])
		user.ID = strings.TrimSpace(line[posSID:])

		// use net user to get additional information
		out, err := exec.Command("net", "user", user.Name).Output()
		if err != nil {
			log.Println("ERROR: cannot get user info;", err)
			continue
		}
		for _, line := range strings.Split(string(out), "\r\n") {
			if len(line) == 0 {
				continue
			}
			// assume always fixed position
			value := strings.TrimSpace(line[29:])
			if strings.HasPrefix(line, "Account active") {
				if value == "Yes" {
					user.AccountActive = true
				} else {
					user.AccountActive = false
				}
			} else if strings.HasPrefix(line, "Account expires") {
				if value == "Never" {
					user.AccountExpires = false
				} else {
					user.AccountExpires = true
				}
			} else if strings.HasPrefix(line, "Password last set") {
				// add timezone to value
				timezone, _ := time.Now().Zone()
				value = value + " " + timezone
				layout := "1/2/2006 3:04:05 PM MST"
				t, err := time.Parse(layout, value)
				if err != nil {
					log.Println("ERROR: cannot parse date time string;", err)
				}
				user.PasswordLastSet = t.Unix()
			} else if strings.HasPrefix(line, "Password expires") {
				if value == "Never" {
					user.PasswordExpires = false
				} else {
					user.PasswordExpires = true
				}
			}
			continue
		}

		users = append(users, user)
	}
	return users
}

func getGroups() map[string][]string {
	out, err := exec.Command("wmic", "path", "win32_groupuser").Output()
	if err != nil {
		log.Fatal("ERROR: unable to get group users;", err)
	}

	groups := make(map[string][]string)
	var posGroupComponent int
	var posPartComponent int
	for i, line := range strings.Split(string(out), "\r\n") {
		if len(line) <= 1 {
			continue
		}

		// find positions of columns
		if i == 0 {
			// assume these exist
			posGroupComponent = strings.Index(line, "GroupComponent")
			posPartComponent = strings.Index(line, "PartComponent")
			continue
		}

		// parse out group and member
		groupComponentStr := strings.TrimSpace(line[posGroupComponent:posPartComponent])
		groupComponentStr = strings.Split(groupComponentStr, ",")[1]
		group := groupComponentStr[6 : len(groupComponentStr)-1]
		partComponentStr := strings.TrimSpace(line[posPartComponent:])
		partComponentStr = strings.Split(partComponentStr, ",")[1]
		member := partComponentStr[6 : len(partComponentStr)-1]
		g, present := groups[group]
		if !present {
			g = make([]string, 0)
		}
		g = append(g, member)
		groups[group] = g
	}

	return groups
}