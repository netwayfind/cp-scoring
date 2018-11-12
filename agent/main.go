package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sumwonyuno/cp-scoring/model"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"

	_ "golang.org/x/crypto/ripemd160"
)

func getServerURL(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("ERROR: cannot read from server URL file")
	}
	serverURL := strings.TrimSpace(string(b))
	checkValidServerURL(serverURL)
	return serverURL
}

func checkValidServerURL(serverURL string) {
	_, err := url.ParseRequestURI(serverURL)
	if err != nil {
		log.Fatalln("ERROR: could not parse server URL;", err)
	}
	// probably not a https:// URL
	if len(serverURL) <= 8 || serverURL[0:8] != "https://" {
		log.Fatalln("ERROR: not a HTTPS URL: " + serverURL)
	}
}

func downloadServerFiles(serverURL string, serverURLFile string, serverPubKeyFile string, serverCrtFile string) {
	checkValidServerURL(serverURL)

	// don't override existing files
	if _, err := os.Stat(serverURLFile); !os.IsNotExist(err) {
		log.Fatalln("ERROR: server URL already set")
	}
	if _, err := os.Stat(serverPubKeyFile); !os.IsNotExist(err) {
		log.Fatalln("ERROR: server public key already set")
	}
	if _, err := os.Stat(serverCrtFile); !os.IsNotExist(err) {
		log.Fatalln("ERROR: server certificate already set")
	}

	// expected URLs
	serverPubKeyFileURL := serverURL + "public/server.pub"
	serverCrtFileURL := serverURL + "public/server.crt"

	// do insecure fetch, as need to get cert to check later...
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	log.Println("Retrieving files from server")

	r, err := client.Get(serverPubKeyFileURL)
	if err != nil {
		log.Fatalln("ERROR: unable to get server public key;", err)
	}
	defer r.Body.Close()
	serverPubKeyFileBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("ERROR: unable to read server public key response;", err)
	}
	log.Println("Fetched server public key")

	r, err = client.Get(serverCrtFileURL)
	if err != nil {
		log.Fatalln("ERROR: unabel to get server certificate")
	}
	defer r.Body.Close()
	serverCrtFileBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("ERROR: unable to read server public key response;", err)
	}
	log.Println("Fetch server certificate")

	// if here, must be OK
	// write to files, only readable by this process
	log.Println("Saving server files")
	err = ioutil.WriteFile(serverURLFile, []byte(serverURL), 0400)
	if err != nil {
		log.Fatalln("ERROR: unable to save server URL;", err)
	}
	err = ioutil.WriteFile(serverPubKeyFile, serverPubKeyFileBytes, 0400)
	if err != nil {
		log.Fatalln("ERROR: unable to save server public key to file;", err)
	}
	err = ioutil.WriteFile(serverCrtFile, serverCrtFileBytes, 0400)
	if err != nil {
		log.Fatalln("ERROR: unable to save server certificate to file;", err)
	}
}

func askForTeam(teamKeyFile string) {
	// don't override existing team key
	if _, err := os.Stat(teamKeyFile); !os.IsNotExist(err) {
		log.Fatalln("ERROR: team key already set")
	}

	reader := bufio.NewReader(os.Stdin)
	log.Print("Enter team key: ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)

	err := ioutil.WriteFile(teamKeyFile, []byte(key), 0400)
	if err != nil {
		log.Fatalln("ERROR: unable to save team key;", err)
	}
	log.Println("Saved team key")
}

func createLinkScoreboard(serverURL string, linkScoreboard string) {
	url := serverURL + "/ui/scoreboard"
	s := "<html><head><meta http-equiv=\"refresh\" content=\"0; url=" + url + "\"></head><body><a href=\"" + url + "\">Scoreboard</a></body></html>"
	err := ioutil.WriteFile(linkScoreboard, []byte(s), 0644)
	if err != nil {
		log.Fatalln("ERROR: unable to save scoreboard link file;", err)
	}
	log.Println("Created scoreboard link file")
}

func createLinkReport(serverURL string, linkReport string, teamKey string) {
	url := serverURL + "/ui/report?team_key=" + teamKey
	s := "<html><head><meta http-equiv=\"refresh\" content=\"0; url=" + url + "\"></head><body><a href=\"" + url + "\">Reports</a></body></html>"
	err := ioutil.WriteFile(linkReport, []byte(s), 0644)
	if err != nil {
		log.Fatalln("ERROR: unable to save report link file;", err)
	}
	log.Println("Created report link file")
}

func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalln("ERROR: unable to get executable", err)
	}
	dir := filepath.Dir(ex)

	serverURLFile := path.Join(dir, "server")
	serverPubFile := path.Join(dir, "server.pub")
	serverCrtFile := path.Join(dir, "server.crt")
	teamKeyFile := path.Join(dir, "team.key")
	linkScoreboard := path.Join(dir, "scoreboard.html")
	linkReport := path.Join(dir, "report.html")
	dataDir := path.Join(dir, "data")

	var serverURL string
	var askTeam bool

	flag.StringVar(&serverURL, "server", "", "server URL")
	flag.BoolVar(&askTeam, "ask_team", false, "ask for team key")
	flag.Parse()

	if askTeam {
		askForTeam(teamKeyFile)
		os.Exit(0)
	}

	if len(serverURL) > 0 {
		downloadServerFiles(serverURL, serverURLFile, serverPubFile, serverCrtFile)
		createLinkScoreboard(serverURL, linkScoreboard)
		os.Exit(0)
	}

	// data directory
	err = os.MkdirAll(dataDir, 0700)
	if err != nil {
		log.Fatalln("Unable to set up data directory;", err)
	}

	// these files must exist
	if _, err := os.Stat(serverURLFile); os.IsNotExist(err) {
		log.Fatalln("ERROR: server file not found")
	}
	if _, err := os.Stat(serverPubFile); os.IsNotExist(err) {
		log.Fatalln("ERROR: server public key not found")
	}
	if _, err := os.Stat(serverCrtFile); os.IsNotExist(err) {
		log.Fatalln("ERROR: server certificate not found")
	}

	// get values from files
	// server URL
	serverURL = getServerURL(serverURLFile)

	// server public key
	pubKeyFile, err := os.Open(serverPubFile)
	if err != nil {
		log.Fatalln("ERROR: cannot read server openpgp public key file;", err)
	}
	defer pubKeyFile.Close()
	entities, err := openpgp.ReadArmoredKeyRing(pubKeyFile)
	if err != nil {
		log.Fatalln("ERROR: cannot read entity;", err)
	}

	// server certificate
	certs := x509.NewCertPool()
	certBytes, err := ioutil.ReadFile(serverCrtFile)
	if err != nil {
		log.Fatalln("ERROR: cannot read server cert file;", err)
	}
	certs.AppendCertsFromPEM(certBytes)
	tlsConfig := &tls.Config{
		RootCAs: certs,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	teamKey := ""

	// main loop
	nextTime := time.Now()
	for {
		nextTime = nextTime.Add(time.Minute)
		saveState(dataDir, entities)
		// if team key not set, check for it
		if len(teamKey) == 0 {
			log.Println("Looking for team key")
			teamKeyBytes, err := ioutil.ReadFile(teamKeyFile)
			if err == nil && len(teamKeyBytes) > 0 {
				log.Println("Found team key")
				teamKey = strings.TrimSpace(string(teamKeyBytes))
			}
			createLinkReport(serverURL, linkReport, teamKey)
		}
		// only send if have team key
		if len(teamKey) > 0 {
			go sendState(dataDir, serverURL, transport, teamKey)
		}
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
			if resp.StatusCode == 200 {
				log.Println("DELETING", fullPath)
				os.Remove(fullPath)
			}
		}
	}
}
