package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
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
	"github.com/netwayfind/cp-scoring/processing"
	"golang.org/x/crypto/openpgp"

	_ "golang.org/x/crypto/ripemd160"
)

const applicationJSON string = "application/json"
const applicationOctetStream string = "application/octet-stream"
const exitCodeFail int = 1
const exitCodeSuccess int = 0
const fileNameHostToken string = "host_token"
const fileNameScenario string = "scenario"
const fileNameServer string = "server"
const fileNameServerPubKey string = "server.pub"
const fileNameTeamKey string = "team_key"

// to be set by build
var version string

func config(dirWork string, dirConfig string, hostname string) {
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

	// get server public key
	serverPubKey, _ := readServerPubKey(dirConfig)
	if serverPubKey == nil {
		log.Println("Downloading server public key")
		resp, err = c.Get(serverURL + "/public/" + fileNameServerPubKey)
		if err != nil {
			log.Fatalln("ERROR: unable to retrieve server public key;", err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("ERROR: could not download server public key: %d", resp.StatusCode)
		}
		pubKeyBs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("ERROR: unable to download server public key;", err)
		}
		err = saveFile(dirConfig, fileNameServerPubKey, string(pubKeyBs))
		if err != nil {
			log.Fatalln("ERROR: unable to save server public key;", err)
		}
	}

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

	writeReadmeHTML(dirWork, serverURL)
}

func copyTeamFiles(dirWork string) {
	host, err := getCurrentHost()
	if err != nil {
		log.Fatalln("ERROR: could not get current host;", err)
	}
	err = host.copyTeamFiles()
	if err != nil {
		log.Fatalln("ERROR: could not copy team files;", err)
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
		cmd := exec.Command(action.Command, action.Args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println("ERROR: unable to get stdout;", err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Println("ERROR: unable to get stderr;", err)
		}
		go func() {
			multiReader := io.MultiReader(stdout, stderr)
			scanner := bufio.NewScanner(multiReader)
			for scanner.Scan() {
				line := scanner.Text()
				log.Println(line)
			}
		}()

		err = cmd.Run()
		if err != nil {
			log.Println("Unable to execute config action;", err)
		}
	}

	log.Println("Applied config. Check log output.")
}

func getScenarioChecks(serverURL string, scenarioID uint64, hostname string, lastModified string) ([]model.Action, string, error) {
	log.Println("Read scenario checks")

	scenarioIDStr := strconv.FormatUint(scenarioID, 10)
	url := serverURL + "/api/scenario-checks/" + scenarioIDStr + "?hostname=" + hostname
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("If-Modified-Since", lastModified)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: could not access server;", err)
		return nil, "", err
	}

	var checks []model.Action

	if resp.StatusCode == 200 {
		log.Println("Scenario host checks updated")
		lastModified = resp.Header.Get("Last-Modified")
		err = json.NewDecoder(resp.Body).Decode(&checks)
		if err != nil {
			log.Println("ERROR: could not read scenario checks")
			return nil, "", err
		}
	} else if resp.StatusCode == 304 {
		// scenario checks not modified
	} else {
		return nil, "", fmt.Errorf("ERROR: could not get scenario checks: %d", resp.StatusCode)
	}

	return checks, lastModified, nil
}

func executeScenarioChecks(scenarioID uint64, hostToken string, checks []model.Action, lastModified string, outputDir string, tempDir string, entities []*openpgp.Entity) {
	log.Println("Executing scenario checks")
	checkResults := []string{}
	for _, check := range checks {
		var result string
		if check.Type == model.ActionTypeExec {
			if &check.Command == nil || len(check.Command) == 0 {
				result = "invalid command"
			} else {
				cmd := exec.Command(check.Command, check.Args...)
				cmd.Dir = tempDir
				out, err := cmd.Output()
				if err != nil {
					result = "could not execute command"
				}
				result = strings.TrimSpace(string(out))
			}
		} else if check.Type == model.ActionTypeFileContains {
			if len(check.Args) == 2 {
				f, err := os.Open(check.Args[0])
				if err != nil {
					result = "could not read file"
				} else {
					scanner := bufio.NewScanner(f)
					bs := []byte(check.Args[1])
					log.Println(bs)
					result = "false"
					for scanner.Scan() {
						lineBs := scanner.Bytes()
						log.Println(lineBs)
						if bytes.Contains(lineBs, bs) {
							result = "true"
							break
						}
					}
					log.Println(result)
				}
			}
		} else if check.Type == model.ActionTypeFileExist {
			if len(check.Args) == 1 {
				if _, err := os.Stat(check.Args[0]); err == nil {
					result = "true"

				} else {
					result = "false"
				}
			}
		} else if check.Type == model.ActionTypeFileRegex {
			if len(check.Args) == 2 {
				fp := check.Args[0]
				rgx := regexp.MustCompile(check.Args[1])
				contents, err := ioutil.ReadFile(fp)
				if err != nil {
					result = "could not read file"
				} else {
					b := rgx.MatchString(string(contents))
					if b {
						result = "true"
					} else {
						result = "false"
					}
				}
			}
		} else if check.Type == model.ActionTypeFileValue {
			if len(check.Args) == 2 {
				fp := check.Args[0]
				rgx := regexp.MustCompile(check.Args[1])
				contents, err := ioutil.ReadFile(fp)
				if err != nil {
					result = "could not read file"
				} else {
					rrs := rgx.FindAllString(string(contents), -1)
					result = strconv.Itoa(len(rrs))
				}
			}
		}
		checkResults = append(checkResults, result)
	}
	auditCheckResults := model.AuditCheckResults{}
	auditCheckResults.ScenarioID = scenarioID
	auditCheckResults.HostToken = hostToken
	auditCheckResults.Timestamp = time.Now().Unix()
	auditCheckResults.CheckResults = checkResults
	auditCheckResults.ChecksLastModified = lastModified

	// save results
	bs, err := processing.ToBytes(auditCheckResults, entities)
	if err != nil {
		log.Println("ERROR: could not prepare saving results to file;", err)
	} else {
		fileName := strconv.FormatInt(auditCheckResults.Timestamp, 10)
		saveFile(outputDir, fileName, string(bs))
	}
}

func executeSubmitScenarioCheckResults(serverURL string, outputDir string) {
	files, err := ioutil.ReadDir(outputDir)
	if err != nil {
		log.Println("ERROR: cannot read results directory;", err)
		return
	}
	if len(files) == 0 {
		return
	}

	log.Println("Sending results to server")

	for _, resultFile := range files {
		filePath := path.Join(outputDir, resultFile.Name())
		bs, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Println("ERROR: unable to read results file;", err)
			continue
		}

		resp, err := http.Post(serverURL+"/api/audit/", applicationOctetStream, bytes.NewBuffer(bs))
		if err != nil {
			log.Println("ERROR: unable to send results file;", err)
			break
		}
		if resp.StatusCode == http.StatusBadRequest {
			log.Println("SERVER REJECTED")
		}
		log.Println("DELETING", filePath)
		os.Remove(filePath)
	}
}

func install() {
	log.Println("Installing agent")

	host, err := getCurrentHost()
	if err != nil {
		log.Fatalln("ERROR: could not get current host;", err)
	}
	err = host.install()
	if err != nil {
		log.Fatalln("ERROR: could not install;", err)
	}
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

func readServerPubKey(dirConfig string) (openpgp.EntityList, error) {
	pubKeyFile, err := os.Open(path.Join(dirConfig, fileNameServerPubKey))
	if err != nil {
		return nil, err
	}
	defer pubKeyFile.Close()
	entities, err := openpgp.ReadArmoredKeyRing(pubKeyFile)
	if err != nil {
		return nil, err
	}
	return entities, nil
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

func requestHostToken(dirData string, serverURL string, scenarioID uint64, hostname string) (string, error) {
	log.Println("Requesting host token")
	hostTokenRequest := model.HostTokenRequest{
		ScenarioID: scenarioID,
		Hostname:   hostname,
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

	log.Println("Team setup complete")
	pressEnterBeforeExit(exitCodeSuccess)
}

func main() {
	// set seed
	rand.Seed(time.Now().UTC().UnixNano())

	// default path
	ex, err := os.Executable()
	if err != nil {
		log.Fatalln("ERROR: unable to get executable", err)
	}
	dirWork := filepath.Dir(ex)

	// program arguments
	var askConfig bool
	var askCopyFiles bool
	var askInstall bool
	var askTeamSetup bool
	var askVersion bool
	flag.StringVar(&dirWork, "dir_work", dirWork, "working directory path")
	flag.BoolVar(&askConfig, "config", false, "run config")
	flag.BoolVar(&askCopyFiles, "copy_files", false, "copy team files to current directory")
	flag.BoolVar(&askInstall, "install", false, "run install")
	flag.BoolVar(&askTeamSetup, "team_setup", false, "team setup")
	flag.BoolVar(&askVersion, "version", false, "get version number")
	flag.Parse()

	// version
	if askVersion {
		log.Println("Version: " + version)
		os.Exit(exitCodeSuccess)
	}

	// install
	if askInstall {
		install()
		os.Exit(exitCodeSuccess)
	}

	// copy files
	if askCopyFiles {
		copyTeamFiles(dirWork)
		os.Exit(exitCodeSuccess)
	}

	dirConfig := path.Join(dirWork, "config")
	dirData := path.Join(dirWork, "data")
	dirTemp := path.Join(dirData, "temp")
	dirResults := path.Join(dirData, "results")

	createDir(dirConfig)
	createDir(dirData)
	createDir(dirTemp)
	createDir(dirResults)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("ERROR: could not get hostname", err)
	}

	// config
	if askConfig {
		config(dirWork, dirConfig, hostname)
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
		log.Fatalln("ERROR: unable to read scenario file;", err)
	}
	log.Println("scenario: ", scenarioID)

	// get server public key
	entities, err := readServerPubKey(dirConfig)
	if err != nil {
		log.Fatalln("ERROR: could not read server public key; ", err)
	}

	var wg sync.WaitGroup

	// run scenario checks
	wg.Add(1)
	go func() {
		nextTime := time.Now()
		hostToken, _ := readHostToken(dirData)
		teamKey := ""
		lastModified := "Thu, 01 Jan 1970 00:00:00 GMT"
		var checks []model.Action
		for {
			if len(hostToken) == 0 {
				hostToken, err = requestHostToken(dirData, serverURL, scenarioID, hostname)
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
				// make sure team key registered before doing scenario checks
				if len(teamKey) == 0 {
					teamKey, _ = readTeamKey(dirData)
				}
				if len(teamKey) > 0 {
					checks2, lastModified2, err := getScenarioChecks(serverURL, scenarioID, hostname, lastModified)
					if err != nil {
						log.Println("ERROR: unable to get checks;", err)
					}
					if checks2 != nil {
						checks = checks2
						lastModified = lastModified2
					}
					if checks != nil {
						executeScenarioChecks(scenarioID, hostToken, checks, lastModified, dirResults, dirTemp, entities)
					}
				}
			}
			nextTime = nextTime.Add(time.Minute)
			wait := time.Since(nextTime) * -1
			time.Sleep(wait)
		}
	}()

	// flush scenario check results
	wg.Add(1)
	go func() {
		nextTime := time.Now()
		for {
			executeSubmitScenarioCheckResults(serverURL, dirResults)
			nextTime = nextTime.Add(5 * time.Second)
			wait := time.Since(nextTime) * -1
			time.Sleep(wait)
		}
	}()

	wg.Wait()
}
