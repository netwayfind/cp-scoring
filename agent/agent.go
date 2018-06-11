package main

import (
	"encoding/json"
	"fmt"
)

type State struct {
	Users []string
	Groups map[string][]string
}

func main() {
	state := GetLinuxState()

	// convert to json bytes
	b, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}

	// check can unmarshal
	var s State
	err = json.Unmarshal(b, &s)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}