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
		sendState(server, transport, teamKey)
		wait := time.Since(nextTime) * -1
		time.Sleep(wait)
	}
}

func sendState(server string, transport *http.Transport, teamKey string) {
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
	c := &http.Client{Transport: transport}
	log.Println("Sending state to server", server)
	resp, err := c.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}
