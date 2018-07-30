package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/sumwonyuno/cp-scoring/model"
)

func main() {
	var server string

	flag.StringVar(&server, "server", "http://localhost:8080", "server URL")
	flag.Parse()

	ex, err := os.Executable()
	if err != nil {
		log.Fatal("ERROR: unable to get executable", err)
	}
	dir := filepath.Dir(ex)

	teamKeyBytes, err := ioutil.ReadFile(path.Join(dir, "team.key"))
	if err != nil {
		log.Println("ERROR: cannot read team id file;", err)
		return
	}
	teamKey := string(teamKeyBytes)

	log.Println("Sending state to server", server)

	nextTime := time.Now()
	for {
		nextTime = nextTime.Add(time.Minute)
		sendState(server, teamKey)
		wait := time.Since(nextTime) * -1
		time.Sleep(wait)
	}
}

func sendState(server string, teamKey string) {
	var state model.State

	// TODO: choose correct function based on OS
	state = getLinuxState()
	state.TeamKey = teamKey

	// convert to json bytes
	b, err := json.Marshal(state)
	if err != nil {
		log.Println("ERROR: marshalling state;", err)
		return
	}

	url := server + "/audit"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}
