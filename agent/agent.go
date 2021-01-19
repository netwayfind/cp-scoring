package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/netwayfind/cp-scoring/model"
)

const applicationJSON string = "application/json"
const exitCodeFail int = 1
const exitCodeSuccess int = 0
const fileNameHostToken string = "host_token"
const fileNameScenario string = "scenario"
const fileNameServer string = "server"
const fileNameTeamKey string = "team_key"

// to be set by build
var version string

func config(dirConfig string, hostname string) {
	log.Println("Running agent config")

	// don't override existing files
	serverURL, err := readServerURL(dirConfig)
	if err == nil || len(serverURL) > 0 {
		log.Fatalln("ERROR: server URL already set")
	}

	// ask for server URL
	log.Println("Enter server URL: ")
	_, err = fmt.Scan(&serverURL)
	if err != nil {
		log.Fatalln("Error asking for server URL;", err)
	}
	// remove trailing slash
	serverURL = strings.TrimRight(strings.TrimSpace(serverURL), ("/"))

	// ask for admin credentials
	var username string
	var password string
	log.Println("Enter username: ")
	_, err = fmt.Scan(&username)
	if err != nil {
		log.Fatalln("Error asking for username;", err)
	}
	log.Println("Enter password: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		log.Fatalln("Error asking for password;", err)
	}

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln("ERROR: unable to create cookie jar;", err)
	}

	c := &http.Client{
		Jar: cookieJar,
	}

	// test server URL
	resp, err := c.Get(serverURL + "/api/version")
	if err != nil {
		log.Fatalln("ERROR: unable to access server;", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalln("ERROR: unexpected server response: ", resp.StatusCode)
	}
	apiVersion, err := readBody(resp)
	if err != nil {
		log.Fatalln("ERROR: could not read server API version;", err)
	}
	if apiVersion != version {
		log.Fatalln("ERROR: could not verify server API version")
	}
	log.Println("Server checks passed")

	loginUser := model.LoginUser{
		Username: username,
		Password: password,
	}
	bs, err := json.Marshal(loginUser)
	if err != nil {
		log.Fatalln("ERROR: could not form login user request;", err)
	}

	// server admin login
	resp, err = c.Post(serverURL+"/api/login/", applicationJSON, bytes.NewBuffer(bs))
	if err != nil {
		log.Fatalln("ERROR: unable to access server;", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalln("ERROR: authentication failure")
	}
	log.Println("User authenticated")

	var scenarioID string
	// ask for scenario
	log.Println("Enter scenario: ")
	_, err = fmt.Scan(&scenarioID)
	if err != nil {
		log.Fatalln("Error asking for server URL;", err)
	}
	scenarioID = strings.TrimSpace(scenarioID)

	// get scenario config
	resp, err = c.Get(serverURL + "/api/scenarios/" + scenarioID + "/config?hostname=" + hostname)
	if err != nil {
		log.Fatalln("ERROR: unable to access server;", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalln("ERROR: cannot access scenario, status code: ", resp.StatusCode)
	}
	var config []model.Action
	err = json.NewDecoder(resp.Body).Decode(&config)
	if err != nil {
		log.Fatalln("ERROR: cannot read scenario config;", err)
	}
	executeConfig(config)

	log.Println("Saving config files")
	err = saveFile(dirConfig, fileNameServer, serverURL)
	if err != nil {
		log.Fatalln("ERROR: unable to save server URL;", err)
	}
	err = saveFile(dirConfig, fileNameScenario, scenarioID)
	if err != nil {
		log.Fatalln("ERROR: unable to save scenario;", err)
	}
}

func createDir(dir string) {
	// data directory
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		log.Fatalln("Unable to set up directory "+dir+";", err)
	}
}

func executeConfig(config []model.Action) {
	log.Println("Executing scenario config")

	for _, action := range config {
		log.Println(" - ", action.Type, ": ", action.Command, "[", strings.Join(action.Args, ","), "]")
		if action.Type != model.ActionTypeExec {
			continue
		}
		if &action.Command == nil || len(action.Command) == 0 {
			continue
		}
		_, err := exec.Command(action.Command, action.Args...).Output()
		if err != nil {
			log.Fatalln("ERROR: unable to execute config;", err)
		}
	}

	log.Println("Successfully applied config")
}

func executeScenarioChecks(serverURL string, scenarioID uint64, hostname string, hostToken string) {
	scenarioIDStr := strconv.FormatUint(scenarioID, 10)
	log.Println("Read scenario checks")
	resp, err := http.Get(serverURL + "/api/scenario-checks/" + scenarioIDStr + "?hostname=" + hostname)
	if err != nil {
		log.Println("ERROR: could not access server;", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Fatal("ERROR: could not get scenario checks")
		return
	}

	var checks []model.Action
	err = json.NewDecoder(resp.Body).Decode(&checks)

	log.Println("Executing scenario checks")
	checkResults := []string{}
	for _, check := range checks {
		var result string
		if check.Type == model.ActionTypeExec {
			if &check.Command == nil || len(check.Command) == 0 {
				result = "nope"
			} else {
				out, err := exec.Command(check.Command, check.Args...).Output()
				if err != nil {
					log.Fatal(err)
				}
				result = strings.TrimSpace(string(out))
			}
		} else if check.Type == model.ActionTypeFileExist {
			if _, err := os.Stat(check.Args[0]); err == nil {
				result = "true"

			} else {
				result = "false"
			}
		} else if check.Type == model.ActionTypeFileRegex {
			fp := check.Args[0]
			rgx := regexp.MustCompile(check.Args[1])
			contents, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}
			b := rgx.Match(contents)
			if b {
				result = "true"
			} else {
				result = "false"
			}
		} else if check.Type == model.ActionTypeFileValue {
			fp := check.Args[0]
			rgx := regexp.MustCompile(check.Args[1])
			contents, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}
			rrs := rgx.FindAllString(string(contents), -1)
			result = strconv.Itoa(len(rrs))
		}
		checkResults = append(checkResults, result)
	}
	auditCheckResults := model.AuditCheckResults{}
	auditCheckResults.ScenarioID = scenarioID
	auditCheckResults.HostToken = hostToken
	auditCheckResults.Timestamp = time.Now().Unix()
	auditCheckResults.CheckResults = checkResults

	log.Println("Sending results to server")
	body, err := json.Marshal(auditCheckResults)
	if err != nil {
		log.Fatal(err)
	}
	resp, err = http.Post(serverURL+"/api/audit/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.Status)
}

func install() {
	log.Println("Installing agent")
}

func pressEnterBeforeExit(code int) {
	log.Println("Press enter to exit")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	os.Exit(code)
}

func readBody(response *http.Response) (string, error) {
	var body string
	err := json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return "", err
	}
	return body, nil
}

func readHostToken(dirData string) (string, error) {
	fileHostToken := path.Join(dirData, fileNameHostToken)
	bs, err := ioutil.ReadFile(fileHostToken)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func readScenarioID(dirConfig string) (uint64, error) {
	fileScenario := path.Join(dirConfig, fileNameScenario)
	bs, err := ioutil.ReadFile(fileScenario)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(string(bs), 10, 64)
}

func readServerURL(dirConfig string) (string, error) {
	fileServer := path.Join(dirConfig, fileNameServer)
	bs, err := ioutil.ReadFile(fileServer)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func readTeamKey(dirData string) (string, error) {
	fileTeamKey := path.Join(dirData, fileNameTeamKey)
	bs, err := ioutil.ReadFile(fileTeamKey)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func requestHostToken(dirData string, serverURL string, hostname string) (string, error) {
	log.Println("Requesting host token")
	hostTokenRequest := model.HostTokenRequest{
		Hostname: hostname,
	}
	hostTokenRequestBs, err := json.Marshal(hostTokenRequest)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(serverURL+"/api/host-token/request", applicationJSON, bytes.NewBuffer(hostTokenRequestBs))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("Could not request host token, unexpected status code")
	}

	return readBody(resp)
}

func saveFile(dir string, fileName string, content string) error {
	file := path.Join(dir, fileName)
	return ioutil.WriteFile(file, []byte(content), 0400)
}

func teamSetup(dirData string, serverURL string) {
	log.Println("Running team setup")

	// team key exists
	teamKey, err := readTeamKey(dirData)
	if err == nil || len(teamKey) > 0 {
		log.Println("Team key already set")
		pressEnterBeforeExit(exitCodeSuccess)
	}

	// no host token yet
	hostToken, err := readHostToken(dirData)
	if err != nil || len(hostToken) == 0 {
		log.Println("Cannot register, agent not running or unable to access scoring server. Try again later.")
		pressEnterBeforeExit(exitCodeFail)
	}

	for {
		// ask for team key
		log.Println("Enter team key: ")
		_, err := fmt.Scan(&teamKey)
		if err != nil {
			log.Println("Error asking for team key;", err)
			pressEnterBeforeExit(exitCodeFail)
		}

		// register team key with host token
		c := http.Client{}
		data := model.HostTokenRegistration{
			HostToken: hostToken,
			TeamKey:   teamKey,
		}
		bs, err := json.Marshal(data)
		if err != nil {
			log.Println("ERROR: could not form host token registration request;", err)
			pressEnterBeforeExit(exitCodeFail)
		}
		r, err := c.Post(serverURL+"/api/host-token/register",
			"application/json", bytes.NewBuffer(bs))
		if err != nil {
			log.Println("ERROR: unable to POST team key (try again later);", err)
			pressEnterBeforeExit(exitCodeFail)
		}
		if r.StatusCode == http.StatusOK {
			log.Println("Team key registered with host token.")
			break
		} else if r.StatusCode == http.StatusUnauthorized {
			log.Println("Team key rejected. Try again.")
			continue
		} else {
			log.Printf("ERROR: Unexpected status code from server: %d", r.StatusCode)
			continue
		}
	}
	err = saveFile(dirData, fileNameTeamKey, teamKey)
	if err != nil {
		log.Println("ERROR: cannot save team key;", err)
		pressEnterBeforeExit(exitCodeFail)
	}
}

func main() {
	// set seed
	rand.Seed(time.Now().UTC().UnixNano())

	// default path
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("ERROR: unable to get executable", err)
	}
	dirWork := filepath.Dir(ex)

	// program arguments
	var askConfig bool
	var askInstall bool
	var askTeamSetup bool
	var askVersion bool
	flag.StringVar(&dirWork, "dir_work", dirWork, "working directory path")
	flag.BoolVar(&askConfig, "config", false, "run config")
	flag.BoolVar(&askInstall, "install", false, "run install")
	flag.BoolVar(&askTeamSetup, "team_setup", false, "team setup")
	flag.BoolVar(&askVersion, "version", false, "get version number")
	flag.Parse()

	// version
	if askVersion {
		log.Println("Version: " + version)
		os.Exit(exitCodeSuccess)
	}

	dirConfig := path.Join(dirWork, "config")
	dirData := path.Join(dirWork, "data")

	createDir(dirConfig)
	createDir(dirData)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("ERROR: could not get hostname", err)
	}

	// install
	if askInstall {
		install()
		os.Exit(exitCodeSuccess)
	}

	// config
	if askConfig {
		config(dirConfig, hostname)
		os.Exit(exitCodeSuccess)
	}

	serverURL, err := readServerURL(dirConfig)
	if err != nil {
		log.Println("Error reading server URL;", err)
		pressEnterBeforeExit(exitCodeFail)
	}

	// team setup
	if askTeamSetup {
		teamSetup(dirData, serverURL)
		pressEnterBeforeExit(exitCodeSuccess)
	}

	scenarioID, err := readScenarioID(dirConfig)
	if err != nil {
		log.Fatal("ERROR: unable to read scenario file;", err)
	}
	log.Println("scenario: ", scenarioID)

	var wg sync.WaitGroup

	// run scenario checks
	wg.Add(1)
	go func() {
		nextTime := time.Now()
		hostToken, _ := readHostToken(dirData)
		for {
			if len(hostToken) == 0 {
				hostToken, err = requestHostToken(dirData, serverURL, hostname)
				if err != nil {
					log.Println("ERROR: could not get host token;", err)
				} else {
					log.Println("Saving host token")
					err = saveFile(dirData, fileNameHostToken, hostToken)
					if err != nil {
						log.Println("ERROR: unable to save host token;", err)
					}
				}
			}
			if len(hostToken) > 0 {
				executeScenarioChecks(serverURL, scenarioID, hostname, hostToken)
			}
			nextTime = nextTime.Add(time.Minute)
			wait := time.Since(nextTime) * -1
			time.Sleep(wait)
		}
	}()

	// log.Println("host token: " + hostToken)
	// log.Println("team key: " + teamKey)

	// rtk := model.HostTokenRegistration{
	// 	HostToken: hostToken,
	// 	TeamKey:   teamKey,
	// }
	// rtkBs, err := json.Marshal(rtk)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// rrrr, err := http.Post(serverURL+"/api/host-token/register", applicationJSON, bytes.NewBuffer(rtkBs))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if rrrr.StatusCode != 200 {
	// 	log.Fatal("Could not register host token")
	// }

	wg.Wait()
}
