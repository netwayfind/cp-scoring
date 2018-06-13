package model

import (
	"log"
	"os"
	"runtime"
	"time"
)

type State struct {
	Timestamp int64
	OS string
	Hostname string
	Users []string
	Groups map[string][]string
}

func GetNewStateTemplate() State {
	var state State
	var err error
	state.Timestamp = time.Now().Unix()
	state.OS = runtime.GOOS
	state.Hostname, err = os.Hostname()
	if err != nil {
		log.Println("ERROR: unable to get hostname;", err)
	}

	return state
}