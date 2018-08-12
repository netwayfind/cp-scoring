// +build windows

package main

import (
	"github.com/sumwonyuno/cp-scoring/model"
)

func getState() model.State {
	state := model.GetNewStateTemplate()
	state.Users = getUsers()
	state.Groups = getGroups()
	return state
}

func getUsers() []model.User {
	return nil
}

func getGroups() map[string][]string {
	return nil
}
