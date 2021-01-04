package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/netwayfind/cp-scoring/test/model"
)

func main() {
	var scenarioID uint64
	flag.Uint64Var(&scenarioID, "scenario", 0, "scenario")
	flag.Parse()

	log.Println("scenario: ", scenarioID)
	scenarioIDStr := strconv.FormatUint(scenarioID, 10)

	contentType := "application/json"

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	hostTokenRequest := model.HostTokenRequest{
		Hostname: hostname,
	}
	hostTokenRequestBs, err := json.Marshal(hostTokenRequest)
	if err != nil {
		log.Fatal(err)
	}

	rrrr, err := http.Post("http://localhost:8000/api/host-token/request", contentType, bytes.NewBuffer(hostTokenRequestBs))
	if err != nil {
		log.Fatal(err)
	}
	if rrrr.StatusCode != 200 {
		log.Fatal("Could not request host token")
	}
	var hostToken string
	err = json.NewDecoder(rrrr.Body).Decode(&hostToken)
	if err != nil {
		log.Fatal(err)
	}

	teamKey := "55555555"

	log.Println("host token: " + hostToken)
	log.Println("team key: " + teamKey)

	rtk := model.HostTokenRegistration{
		HostToken: hostToken,
		TeamKey:   teamKey,
	}
	rtkBs, err := json.Marshal(rtk)
	if err != nil {
		log.Fatal(err)
	}

	rrrr, err = http.Post("http://localhost:8000/api/host-token/register", contentType, bytes.NewBuffer(rtkBs))
	if err != nil {
		log.Fatal(err)
	}
	if rrrr.StatusCode != 200 {
		log.Fatal("Could not register host token")
	}

	x, err := http.Get("http://localhost:8000/api/scenarios/" + scenarioIDStr)
	if err != nil {
		log.Fatal(err)
	}
	y := model.Scenario{}
	err = json.NewDecoder(x.Body).Decode(&y)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(y.Name)

	x, err = http.Get("http://localhost:8000/api/scenarios/" + scenarioIDStr + "/checks?hostname=" + hostname)
	if err != nil {
		log.Fatal(err)
	}
	var yy []model.Action
	err = json.NewDecoder(x.Body).Decode(&yy)

	checkResults := []string{}
	for _, v := range yy {
		log.Println(v)
		log.Println(" - ", v.Type, ": ", v.Command)
		var r string
		if v.Type == model.ActionTypeExec {
			if &v.Command == nil || len(v.Command) == 0 {
				r = "nope"
			} else {
				out, err := exec.Command(v.Command, v.Args...).Output()
				if err != nil {
					log.Fatal(err)
				}
				r = strings.TrimSpace(string(out))
			}
		} else if v.Type == model.ActionTypeFileExist {
			if _, err := os.Stat(v.Args[0]); err == nil {
				r = "true"

			} else {
				r = "false"
			}
		} else if v.Type == model.ActionTypeFileRegex {
			fp := v.Args[0]
			rgx := regexp.MustCompile(v.Args[1])
			contents, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}
			b := rgx.Match(contents)
			if b {
				r = "true"
			} else {
				r = "false"
			}
		} else if v.Type == model.ActionTypeFileValue {
			fp := v.Args[0]
			rgx := regexp.MustCompile(v.Args[1])
			contents, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}
			rrs := rgx.FindAllString(string(contents), -1)
			r = strconv.Itoa(len(rrs))
		}
		checkResults = append(checkResults, r)
	}
	auditCheckResults := model.AuditCheckResults{}
	auditCheckResults.ScenarioID = scenarioID
	auditCheckResults.HostToken = hostToken
	auditCheckResults.Timestamp = time.Now().Unix()
	auditCheckResults.CheckResults = checkResults

	log.Println("to server")
	body, err := json.Marshal(auditCheckResults)
	if err != nil {
		log.Fatal(err)
	}
	x, err = http.Post("http://localhost:8000/audit", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(x.Status)

	log.Println("Done")

}
