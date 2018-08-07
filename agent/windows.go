package main

import (
	"github.com/sumwonyuno/cp-scoring/model"
)

func getWindowsState() model.State {
	state := model.GetNewStateTemplate()
	state.Users = getUsersWindows()
	state.Groups = getGroupsWindows()
	return state
}

func getUsersWindows() []string {
	return nil
}

func getGroupsWindows() map[string][]string {
	return nil
}
