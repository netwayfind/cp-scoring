package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sumwonyuno/cp-scoring/model"
)

func main() {
	var server string

	flag.StringVar(&server, "server", "https://localhost:8443", "server URL")
	flag.Parse()

	ex, err := os.Executable()
	if err != nil {
		log.Fatal("ERROR: unable to get executable", err)
	}
	dir := filepath.Dir(ex)

	log.Println("Setting up data directory")
	dataDir := path.Join(dir, "data")
	err = os.MkdirAll(dataDir, 0700)
	if err != nil {
		log.Fatal("Unable to set up data directory;", err)
	}

	log.Println("Reading team key")
	teamKeyBytes, err := ioutil.ReadFile(path.Join(dir, "team.key"))
	if err != nil {
		log.Println("ERROR: cannot read team id file;", err)
		return
	}
	teamKey := string(teamKeyBytes)

	certs := x509.NewCertPool()

	log.Println("Reading server cert file")
	certBytes, err := ioutil.ReadFile(path.Join(dir, "server.crt"))
	if err != nil {
		log.Println("ERROR: cannot read server cert file;", err)
		return
	}
	certs.AppendCertsFromPEM(certBytes)
	tlsConfig := &tls.Config{
		RootCAs: certs,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	nextTime := time.Now()
	for {
		nextTime = nextTime.Add(time.Minute)
		saveState(dataDir)
		go sendState(dataDir, server, transport, teamKey)
		wait := time.Since(nextTime) * -1
		time.Sleep(wait)
	}
}

func saveState(dir string) {
	var state model.State

	// TODO: choose correct function based on OS
	log.Println("Getting state")
	state = getLinuxState()

	// convert to json bytes
	b, err := json.Marshal(state)
	if err != nil {
		log.Println("ERROR: marshalling state;", err)
		return
	}

	log.Println("Saving state")
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	filePath := path.Join(dir, timestamp)
	err = ioutil.WriteFile(filePath, b, 0600)
	if err != nil {
		log.Println("ERROR: saving state;", err)
		return
	}
}

func sendState(dir string, server string, transport *http.Transport, teamKey string) {
	url := server + "/audit"
	c := &http.Client{Transport: transport}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("ERROR: cannot read data directory;", err)
		return
	}
	for _, file := range files {
		fullPath := path.Join(dir, file.Name())
		log.Println("Found state file:", fullPath)
		stateBytes, err := ioutil.ReadFile(fullPath)
		if err != nil {
			log.Println("ERROR: unable to read state file;", err)
			log.Println("DELETING", fullPath)
			os.Remove(fullPath)
		} else {
			log.Println("Sending state to server", server)
			var submission model.StateSubmission
			submission.TeamKey = teamKey
			submission.StateBytes = stateBytes
			b, err := json.Marshal(submission)
			resp, err := c.Post(url, "application/json", bytes.NewBuffer(b))
			if err != nil {
				log.Println("ERROR:", err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("ERROR: unable to read server body")
			}
			log.Println(string(body))
			log.Println("DELETING", fullPath)
			os.Remove(fullPath)
		}
	}
}
