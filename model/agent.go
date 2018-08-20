package model

import (
	"log"
	"os"
	"runtime"
	"time"
)

type Host struct {
	ID       int64
	Hostname string
	OS       string
}

type StateSubmission struct {
	TeamKey    string
	StateBytes []byte
}

type State struct {
	Timestamp int64
	OS        string
	Hostname  string
	Users     []User
	Groups    map[string][]string
	Processes []Process
	Software  []Software
}

type User struct {
	Name            string
	ID              string
	AccountPresent  bool
	AccountActive   bool
	AccountExpires  bool
	PasswordLastSet int64
	PasswordExpires bool
}

type Process struct {
	PID         int64
	User        string
	CommandLine string
}

type Software struct {
	Name    string
	Version string
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
