package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sumwonyuno/cp-scoring/agent"
	"github.com/sumwonyuno/cp-scoring/model"
	"github.com/sumwonyuno/cp-scoring/processing"
	"golang.org/x/crypto/openpgp"

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
	serverPubKeyFileURL := serverURL + "/public/server.pub"
	serverCrtFileURL := serverURL + "/public/server.crt"

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

func getHostToken(hostTokenURL string, hostTokenFile string, hostname string, transport *http.Transport) (string, error) {
	// if host token file doesn't exist, get new host token and save it
	if _, err := os.Stat(hostTokenFile); os.IsNotExist(err) {
		log.Println("Host token not found")
		c := http.Client{Transport: transport}
		url := hostTokenURL + "?hostname=" + hostname
		r, err := c.Get(url)
		if err != nil {
			return "", err
		}
		if r.StatusCode != 200 {
			return "", errors.New("Unexpected status code from server: " + r.Status)
		}
		defer r.Body.Close()
		tokenBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		if len(tokenBytes) == 0 {
			return "", errors.New("Empty host token")
		}
		err = ioutil.WriteFile(hostTokenFile, tokenBytes, 0400)
		if err != nil {
			log.Fatalln("ERROR: unable to save team key;", err)
		}
		log.Println("Saved host token")
	}

	// get host token file from file
	tokenBytes, err := ioutil.ReadFile(hostTokenFile)
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
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

func createLinkReport(serverURL string, linkReport string, hostname string, hostToken string) {
	url := serverURL + "/token/team"
	var s string
	// if have host token, then connected
	if len(hostToken) > 0 {
		s = "<html><head><body>" +
			"<form method=\"POST\" action=" + url + ">" +
			"<input name=\"host_token\" hidden value=\"" + hostToken + "\"/>" +
			"<input name=\"hostname\" hidden value=\"" + hostname + "\"/>" +
			"<label id=\"team_key\">Enter team key:</label>" +
			"<input name=\"team_key\" />" +
			"<button type=\"submit\">Submit</button>" +
			"</form></body></html>"
	} else {
		// refresh page every 30 seconds
		s = "<html><head><meta http-equiv=\"refresh\" content=\"30\"/></head><body>" +
			"Cannot connect to server. Check Internet connection and try again later." +
			"<p />" +
			"Last check: " + time.Now().Format("Jan 2 2006 15:04:05 MST") +
			"</body></html>"
	}
	err := ioutil.WriteFile(linkReport, []byte(s), 0644)
	if err != nil {
		log.Fatalln("ERROR: unable to save report link file;", err)
	}
	log.Println("Created report link file")
}

func installThis() {
	host := agent.GetCurrentHost()
	host.Install()
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
	hostTokenFile := path.Join(dir, "host_token")
	linkScoreboard := path.Join(dir, "scoreboard.html")
	linkReport := path.Join(dir, "report.html")
	dataDir := path.Join(dir, "data")

	var serverURL string
	var install bool

	flag.StringVar(&serverURL, "server", "", "server URL")
	flag.BoolVar(&install, "install", false, "install to this computer")
	flag.Parse()

	if install {
		installThis()
		os.Exit(0)
	}

	hostname, _ := os.Hostname()

	if len(serverURL) > 0 {
		// remove trailing slash
		serverURL = strings.TrimRight(serverURL, "/")
		downloadServerFiles(serverURL, serverURLFile, serverPubFile, serverCrtFile)
		createLinkScoreboard(serverURL, linkScoreboard)
		createLinkReport(serverURL, linkReport, hostname, "")
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

	// host token
	hostTokenURL := serverURL + "/token/host"
	hostToken, err := getHostToken(hostTokenURL, hostTokenFile, hostname, transport)
	if err != nil {
		log.Println("ERROR: getting host token;", err)
	}
	createLinkReport(serverURL, linkReport, hostname, hostToken)

	// main loop
	nextTime := time.Now()
	for {
		nextTime = nextTime.Add(time.Minute)
		saveState(dataDir, entities)
		go sendState(dataDir, serverURL, transport, hostToken)
		wait := time.Since(nextTime) * -1
		time.Sleep(wait)
	}
}

func saveState(dir string, entities []*openpgp.Entity) {
	log.Println("Getting state")
	state := agent.GetState()

	bs, err := processing.ToBytes(state, entities)
	if err != nil {
		log.Println("ERROR: unable to convert state to bytes;", err)
		return
	}

	log.Println("Saving state")
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	filePath := path.Join(dir, timestamp+".enc")
	err = ioutil.WriteFile(filePath, bs, 0600)
	if err != nil {
		log.Println("ERROR: saving state;", err)
		return
	}
}

func sendState(dir string, server string, transport *http.Transport, hostToken string) {
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
			submission.HostToken = hostToken
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
