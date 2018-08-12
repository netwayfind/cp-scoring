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
	"runtime"
	"strconv"
	"time"

	"github.com/sumwonyuno/cp-scoring/model"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"

	_ "golang.org/x/crypto/ripemd160"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("ERROR: unable to get executable", err)
	}
	dir := filepath.Dir(ex)

	var serverURL string
	defaultURL := "https://localhost:8443"

	flag.StringVar(&serverURL, "server", defaultURL, "server URL")
	flag.Parse()

	serverURLFile := path.Join(dir, "server")
	serverPubFile := path.Join(dir, "server.pub")
	serverCrtFile := path.Join(dir, "server.crt")
	teamKeyFile := path.Join(dir, "team.key")

	// read server URL from saved file
	if _, err := os.Stat(serverURLFile); os.IsExist(err) {
		log.Println("Found server URL file")
		b, err := ioutil.ReadFile(serverURLFile)
		if err != nil {
			log.Fatal("ERROR: cannot read from server URL file")
		}
		serverURL = string(b)
	} else {
		// not set up, save server URL
		log.Println("Saving server URL")
		serverFileText := []byte(serverURL)
		err = ioutil.WriteFile(serverURLFile, serverFileText, 0600)
		if err != nil {
			log.Fatal("ERROR: cannot write to server URL file;", err)
		}
	}
	log.Println("Server URL: " + serverURL)

	log.Println("Setting up data directory")
	dataDir := path.Join(dir, "data")
	err = os.MkdirAll(dataDir, 0700)
	if err != nil {
		log.Fatal("Unable to set up data directory;", err)
	}

	log.Println("Reading team key")
	teamKeyBytes, err := ioutil.ReadFile(teamKeyFile)
	if err != nil {
		log.Println("ERROR: cannot read team id file;", err)
		return
	}
	teamKey := string(teamKeyBytes)

	certs := x509.NewCertPool()

	log.Println("Reading server cert file")
	certBytes, err := ioutil.ReadFile(serverCrtFile)
	if err != nil {
		log.Println("ERROR: cannot read server cert file;", err)
		return
	}
	certs.AppendCertsFromPEM(certBytes)
	tlsConfig := &tls.Config{
		RootCAs: certs,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	log.Println("Reading server openpgp public key file")
	pubKeyFile, err := os.Open(serverPubFile)
	if err != nil {
		log.Println("ERROR: cannot read server openpgp public key file;", err)
		return
	}
	defer pubKeyFile.Close()
	entities, err := openpgp.ReadArmoredKeyRing(pubKeyFile)
	if err != nil {
		log.Println("ERROR: cannot read entity;", err)
		return
	}

	nextTime := time.Now()
	for {
		nextTime = nextTime.Add(time.Minute)
		saveState(dataDir, entities)
		go sendState(dataDir, serverURL, transport, teamKey)
		wait := time.Since(nextTime) * -1
		time.Sleep(wait)
	}
}

func encryptBytes(theBytes []byte, entities []*openpgp.Entity) ([]byte, error) {
	log.Println("Encrypting bytes")
	encbuf := bytes.NewBuffer(nil)
	writer, err := armor.Encode(encbuf, openpgp.PublicKeyType, nil)
	if err != nil {
		return nil, err
	}

	plaintext, err := openpgp.Encrypt(writer, entities, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	plaintext.Write(theBytes)
	plaintext.Close()
	writer.Close()

	return encbuf.Bytes(), nil
}

func saveState(dir string, entities []*openpgp.Entity) {
	var state model.State

	// TODO: choose correct function based on OS
	log.Println("Getting state")
	if runtime.GOOS == "linux" {
		state = getState()
	} else if runtime.GOOS == "windows" {
		state = getState()
	} else {
		log.Fatal("ERROR: unsupported platform: " + runtime.GOOS)
	}

	// convert to json bytes
	b, err := json.Marshal(state)
	if err != nil {
		log.Println("ERROR: marshalling state;", err)
		return
	}
	encryptedBytes, err := encryptBytes(b, entities)
	if err != nil {
		log.Println("ERROR: cannot encrypt state bytes;", err)
		return
	}

	log.Println("Saving state")
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	filePath := path.Join(dir, timestamp+".enc")
	err = ioutil.WriteFile(filePath, encryptedBytes, 0600)
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
