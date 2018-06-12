package main

import (
	"encoding/json"
	"fmt"
	"github.com/sumwonyuno/cp-scoring/model"
)

func main() {
	state := GetLinuxState()

	// convert to json bytes
	b, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}

	// check can unmarshal
	var s model.State
	err = json.Unmarshal(b, &s)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}