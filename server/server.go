package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sumwonyuno/cp-scoring/processing"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/sumwonyuno/cp-scoring/auditor"
	"github.com/sumwonyuno/cp-scoring/model"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/openpgp"
)

const cookieName string = "cp-scoring"
const configFileName string = "cp-scoring.conf"

var version string

type theServer struct {
	userTokens   map[string]string
	backingStore backingStore
}

func (theServer theServer) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			msg := "Unauthorized request"
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		if username, found := theServer.userTokens[cookie.Value]; found {
			isAdmin, err := theServer.backingStore.IsAdmin(username)
			if err != nil {
				msg := "ERROR: unable to check if user is admin"
				log.Println(msg, err)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			if !isAdmin {
				// delete any existing sessions
				delete(theServer.userTokens, cookie.Value)
				msg := "Unauthorized request"
				http.Error(w, msg, http.StatusUnauthorized)
				log.Println(r.RemoteAddr, ",", r.URL, ",", msg)
				return
			}
			next.ServeHTTP(w, r)
		} else {
			msg := "Unauthorized request"
			http.Error(w, msg, http.StatusUnauthorized)
			log.Println(r.RemoteAddr, ",", r.URL, ",", msg)
		}
	})
}

func saveAuditRequest(w http.ResponseWriter, r *http.Request, dataDir string) {
	log.Println("Received audit request")

	// expecting HTTP POST
	if r.Method != "POST" {
		msg := "HTTP 405"
		log.Println(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		// don't send error back to client
		return
	}

	// save request to temp file
	// temp file is dataDir/<timestamp>_<tempname>
	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	outFile, err := ioutil.TempFile(dataDir, timestampStr+"_")
	outFile.Chmod(0600)
	defer outFile.Close()
	outFile.Write(body)
}

func (theServer theServer) audit(dataDir string, entities openpgp.EntityList) {
	for {
		err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			log.Println("Auditing file " + path)
			auditErr := theServer.auditFile(path, entities)
			if auditErr != nil {
				return auditErr
			}
			log.Println("DELETING " + path)
			return os.Remove(path)
		})
		if err != nil {
			log.Println("ERROR: unable to walk data directory;", err)
		}
		time.Sleep(10 * time.Second)
	}
}

func (theServer theServer) auditFile(fileStr string, entities openpgp.EntityList) error {
	fileBytes, err := ioutil.ReadFile(fileStr)
	if err != nil {
		log.Println("ERROR: unable to read file")
		return err
	}

	var stateSubmission model.StateSubmission
	err = json.Unmarshal(fileBytes, &stateSubmission)
	if err != nil {
		log.Println("ERROR: cannot unmarshal state submission;")
		// allow file to be deleted
		return nil
	}
	state, err := processing.FromBytes(stateSubmission.StateBytes, entities)
	if err != nil {
		log.Println("ERROR: cannot unmarshal state;")
		// allow file to be deleted
		return nil
	}

	log.Println("Saving state")
	err = theServer.backingStore.InsertState(string(fileBytes))
	if err != nil {
		log.Println("ERROR: cannot insert state;")
		return err
	}

	log.Println("Getting information")
	hostToken := stateSubmission.HostToken
	if len(hostToken) == 0 {
		log.Println("ERROR: received state submission without host token")
		// allow file to be deleted
		return nil
	}

	log.Println("Getting scenarios")
	scenarioIDs, err := theServer.backingStore.SelectScenariosForHostname(state.Hostname)
	if err != nil {
		log.Println("ERROR: cannot get scenario IDs;")
		return err
	}
	if len(scenarioIDs) == 0 {
		log.Println("ERROR: no scenarios found")
		// allow file to be deleted
		return nil
	}
	for _, scenarioID := range scenarioIDs {
		log.Println(fmt.Sprintf("Processing scenario ID: %d", scenarioID))

		log.Println("Getting scenario templates")
		templates, err := theServer.backingStore.SelectTemplatesForHostname(scenarioID, state.Hostname)
		if err != nil {
			log.Println("ERROR: cannot get templates;")
			return err
		}
		log.Println(fmt.Sprintf("Found template count: %d", len(templates)))
		if len(templates) == 0 {
			log.Println("ERROR: no templates found")
			// skip this scenario
			continue
		}

		log.Println("Auditing state")
		report := auditor.Audit(state, templates)

		log.Println("Saving report")
		report.Timestamp = state.Timestamp
		err = theServer.backingStore.InsertScenarioReport(scenarioID, hostToken, report)
		if err != nil {
			log.Println("ERROR: cannot insert report;")
			return err
		}

		log.Println("Saving score")
		var score int64
		for _, finding := range report.Findings {
			score += finding.Value
		}

		var scoreEntry model.ScenarioHostScore
		scoreEntry.ScenarioID = scenarioID
		scoreEntry.HostToken = hostToken
		scoreEntry.Timestamp = state.Timestamp
		scoreEntry.Score = score
		err = theServer.backingStore.InsertScenarioScore(scoreEntry)
		if err != nil {
			log.Println("ERROR: cannot insert scenario score;")
			return err
		}
	}

	log.Println("Received and saved")
	return nil
}

func (theServer theServer) getHosts(w http.ResponseWriter, r *http.Request) {
	log.Println("get all hosts")

	// get all hosts
	hosts, err := theServer.backingStore.SelectHosts()
	if err != nil {
		msg := "ERROR: cannot retrieve hosts;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	b, err := json.Marshal(hosts)
	if err != nil {
		msg := "ERROR: cannot marshal hosts;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(b)
}

func (theServer theServer) getHost(w http.ResponseWriter, r *http.Request) {
	log.Println("get a host")

	// parse out int64 id
	// remove /hosts/ from URL
	id, err := strconv.ParseInt(r.URL.Path[7:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse host id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)
	host, err := theServer.backingStore.SelectHost(id)
	if err != nil {
		msg := "ERROR: cannot retrieve host;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	out, err := json.Marshal(host)
	if err != nil {
		msg := "ERROR: cannot marshal host;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(out)
}

func (theServer theServer) newHost(w http.ResponseWriter, r *http.Request) {
	log.Println("new host")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var host model.Host
	err = json.Unmarshal(body, &host)
	if err != nil {
		msg := "ERROR: cannot unmarshal host;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	id, err := theServer.backingStore.InsertHost(host)
	if err != nil {
		msg := "ERROR: cannot insert host;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	idStr := strconv.FormatInt(id, 10)

	// new host
	log.Println("Saved host " + idStr)
	w.Write([]byte(idStr))
}

func (theServer theServer) editHost(w http.ResponseWriter, r *http.Request) {
	log.Println("edit host")

	// parse out int64 id
	// remove /hosts/ from URL
	id, err := strconv.ParseInt(r.URL.Path[7:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse host id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var host model.Host
	err = json.Unmarshal(body, &host)
	if err != nil {
		msg := "ERROR: cannot unmarshal host;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	err = theServer.backingStore.UpdateHost(id, host)
	if err != nil {
		msg := "ERROR: cannot update host;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	msg := "Updated host"
	log.Println(msg)
	w.Write([]byte(msg))
}

func (theServer theServer) deleteHost(w http.ResponseWriter, r *http.Request) {
	log.Println("delete host")

	// parse out int64 id
	// remove /hosts/ from URL
	id, err := strconv.ParseInt(r.URL.Path[7:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse host id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)
	err = theServer.backingStore.DeleteHost(id)
	if err != nil {
		msg := "ERROR: cannot delete host;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
}

func (theServer theServer) getTeams(w http.ResponseWriter, r *http.Request) {
	log.Println("get all teams")

	// get all teams
	teams, err := theServer.backingStore.SelectTeams()
	if err != nil {
		msg := "ERROR: cannot retrieve teams;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	b, err := json.Marshal(teams)
	if err != nil {
		msg := "ERROR: cannot marshal teams;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(b)
}

func (theServer theServer) getTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("get a team")

	// parse out int64 id
	// remove /teams/ from URL
	id, err := strconv.ParseInt(r.URL.Path[7:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse team id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)
	team, err := theServer.backingStore.SelectTeam(id)
	if err != nil {
		msg := "ERROR: cannot retrieve team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	out, err := json.Marshal(team)
	if err != nil {
		msg := "ERROR: cannot marshal team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(out)
}

func (theServer theServer) newTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("new team")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var team model.Team
	err = json.Unmarshal(body, &team)
	if err != nil {
		msg := "ERROR: cannot unmarshal team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	id, err := theServer.backingStore.InsertTeam(team)
	if err != nil {
		msg := "ERROR: cannot insert team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	idStr := strconv.FormatInt(id, 10)

	// new team
	log.Println("Saved team " + idStr)
	w.Write([]byte(idStr))
}

func (theServer theServer) editTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("edit team")

	// parse out int64 id
	// remove /teams/ from URL
	id, err := strconv.ParseInt(r.URL.Path[7:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse team id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var team model.Team
	err = json.Unmarshal(body, &team)
	if err != nil {
		msg := "ERROR: cannot unmarshal team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	err = theServer.backingStore.UpdateTeam(id, team)
	if err != nil {
		msg := "ERROR: cannot update team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	msg := "Updated team"
	log.Println(msg)
	w.Write([]byte(msg))
}

func (theServer theServer) deleteTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("delete team")

	// parse out int64 id
	// remove /teams/ from URL
	id, err := strconv.ParseInt(r.URL.Path[7:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse team id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)
	err = theServer.backingStore.DeleteTeam(id)
	if err != nil {
		msg := "ERROR: cannot delete team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
}

func (theServer theServer) getTemplates(w http.ResponseWriter, r *http.Request) {
	log.Println("get all templates")

	// get all templates
	templates, err := theServer.backingStore.SelectTemplates()
	if err != nil {
		msg := "ERROR: cannot retrieve templates;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	b, err := json.Marshal(templates)
	if err != nil {
		msg := "ERROR: cannot marshal templates;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(b)
}

func (theServer theServer) getTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("get a template")

	// parse out int64 id
	// remove /templates/ from URL
	id, err := strconv.ParseInt(r.URL.Path[11:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse template id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)
	template, err := theServer.backingStore.SelectTemplate(id)
	if err != nil {
		msg := "ERROR: cannot retrieve template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	out, err := json.Marshal(template)
	if err != nil {
		msg := "ERROR: cannot marshal template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(out)
}

func (theServer theServer) newTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("new template")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var template model.Template
	err = json.Unmarshal(body, &template)
	if err != nil {
		msg := "ERROR: cannot unmarshal template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	id, err := theServer.backingStore.InsertTemplate(template)
	if err != nil {
		msg := "ERROR: cannot insert template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	idStr := strconv.FormatInt(id, 10)

	// new template
	log.Println("Saved template " + idStr)
	w.Write([]byte(idStr))
}

func (theServer theServer) editTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("edit template")

	// parse out int64 id
	// remove /templates/ from URL
	id, err := strconv.ParseInt(r.URL.Path[11:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse template id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var template model.Template
	err = json.Unmarshal(body, &template)
	if err != nil {
		msg := "ERROR: cannot unmarshal template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	err = theServer.backingStore.UpdateTemplate(id, template)
	if err != nil {
		msg := "ERROR: cannot update template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	msg := "Updated template"
	log.Println(msg)
	w.Write([]byte(msg))
}

func (theServer theServer) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("delete template")

	// parse out int64 id
	// remove /templates/ from URL
	id, err := strconv.ParseInt(r.URL.Path[11:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse template id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)
	err = theServer.backingStore.DeleteTemplate(id)
	if err != nil {
		msg := "ERROR: cannot delete template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
}

func (theServer theServer) getScenarios(w http.ResponseWriter, r *http.Request) {
	log.Println("get all scenarios")

	// get all scenarios, even not enabled
	scenarios, err := theServer.backingStore.SelectScenarios(false)
	if err != nil {
		msg := "ERROR: cannot retrieve scenarios;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	b, err := json.Marshal(scenarios)
	if err != nil {
		msg := "ERROR: cannot marshal scenarios;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(b)
}

func (theServer theServer) getScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("get a scenario")

	// parse out int64 id
	// remove /scenarios/ from URL
	id, err := strconv.ParseInt(r.URL.Path[11:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse scenario id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)
	scenario, err := theServer.backingStore.SelectScenario(id)
	if err != nil {
		msg := "ERROR: cannot retrieve scenario;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	out, err := json.Marshal(scenario)
	if err != nil {
		msg := "ERROR: cannot marshal scenario;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(out)
}

func (theServer theServer) newScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("new scenario")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var scenario model.Scenario
	err = json.Unmarshal(body, &scenario)
	if err != nil {
		msg := "ERROR: cannot unmarshal scenario;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	id, err := theServer.backingStore.InsertScenario(scenario)
	if err != nil {
		msg := "ERROR: cannot insert scenario;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	idStr := strconv.FormatInt(id, 10)

	// new scenario
	log.Println("Saved scenario " + idStr)
	w.Write([]byte(idStr))
}

func (theServer theServer) editScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("edit scenario")

	// parse out int64 id
	// remove /scenarios/ from URL
	id, err := strconv.ParseInt(r.URL.Path[11:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse scenario id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var scenario model.Scenario
	err = json.Unmarshal(body, &scenario)
	if err != nil {
		msg := "ERROR: cannot unmarshal scenario;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	err = theServer.backingStore.UpdateScenario(id, scenario)
	if err != nil {
		msg := "ERROR: cannot update scenario;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	msg := "Updated scenario"
	log.Println(msg)
	w.Write([]byte(msg))
}

func (theServer theServer) deleteScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("delete scenario")

	// parse out int64 id
	// remove /scenarios/ from URL
	id, err := strconv.ParseInt(r.URL.Path[11:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse scenario id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)
	err = theServer.backingStore.DeleteScenario(id)
	if err != nil {
		msg := "ERROR: cannot delete scenario;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
}

func (theServer theServer) getScenariosForScoreboard(w http.ResponseWriter, r *http.Request) {
	log.Println("get scenarios for scoreboard")

	// get scenarios, only enabled
	scenarios, err := theServer.backingStore.SelectScenarios(true)
	if err != nil {
		msg := "ERROR: cannot retrieve scenarios;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	b, err := json.Marshal(scenarios)
	if err != nil {
		msg := "ERROR: cannot marshal scenarios;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(b)
}

func (theServer theServer) getScenarioScores(w http.ResponseWriter, r *http.Request) {
	log.Println("get scenario scores")

	// parse out int64 id
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse scenario id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(fmt.Sprintf("Scenario ID: %d", id))
	scores, err := theServer.backingStore.SelectLatestScenarioScores(id)
	if err != nil {
		msg := "ERROR: cannot retrieve scenario scores;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	out, err := json.Marshal(scores)
	if err != nil {
		msg := "ERROR: cannot marshal scenario scores;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(out)
}

func (theServer theServer) getScenarioScoresTimeline(w http.ResponseWriter, r *http.Request) {
	log.Println("get scenario timeline for team")

	// parse out int64 id
	vars := mux.Vars(r)

	scenarioID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse scenario id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(fmt.Sprintf("Scenario ID: %d", scenarioID))
	teamKey := r.FormValue("team_key")
	teamID, err := theServer.backingStore.SelectTeamIDForKey(teamKey)
	if err != nil {
		msg := "ERROR: cannot retrieve team id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(fmt.Sprintf("Team ID: %d", teamID))
	hostname := r.FormValue("hostname")
	hostTokens, err := theServer.backingStore.SelectHostTokens(teamID, hostname)
	if err != nil {
		msg := "ERROR: cannot retrieve host token;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	timelines := make([]model.ScenarioTimeline, 0)
	for _, hostToken := range hostTokens {
		timeline, err := theServer.backingStore.SelectScenarioTimeline(scenarioID, hostToken)
		if err != nil {
			msg := "ERROR: cannot retrieve scenario timeline for team;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}
		timelines = append(timelines, timeline)
	}
	out, err := json.Marshal(timelines)
	if err != nil {
		msg := "ERROR: cannot marshal scenario timeline for team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(out)
}

func (theServer theServer) getScenarioScoresReport(w http.ResponseWriter, r *http.Request) {
	log.Println("get scenario report for team")

	// parse out int64 id
	vars := mux.Vars(r)

	scenarioID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse scenario id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(fmt.Sprintf("Scenario ID: %d", scenarioID))
	teamKey := r.FormValue("team_key")
	teamID, err := theServer.backingStore.SelectTeamIDForKey(teamKey)
	if err != nil {
		msg := "ERROR: cannot retrieve team id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(fmt.Sprintf("Team ID: %d", teamID))
	hostname := r.FormValue("hostname")
	hostTokens, err := theServer.backingStore.SelectHostTokens(teamID, hostname)
	if err != nil {
		msg := "ERROR: cannot retrieve host token;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	// only take latest host token
	hostToken := hostTokens[len(hostTokens)-1]
	report, err := theServer.backingStore.SelectLatestScenarioReport(scenarioID, hostToken)
	if err != nil {
		msg := "ERROR: cannot retrieve scenario report for team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	// take out findings to not show
	findingsToShow := report.Findings[:0]
	for _, finding := range report.Findings {
		if finding.Show {
			findingsToShow = append(findingsToShow, finding)
		}
	}
	report.Findings = findingsToShow
	out, err := json.Marshal(report)
	if err != nil {
		msg := "ERROR: cannot marshal scenario report for team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(out)
}

func (theServer theServer) getTeamScenarioHosts(w http.ResponseWriter, r *http.Request) {
	log.Println("get team scenario hosts")

	teamKey := r.FormValue("team_key")
	teamID, err := theServer.backingStore.SelectTeamIDForKey(teamKey)
	if err != nil {
		msg := "ERROR: cannot retrieve team id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(fmt.Sprintf("Team ID: %d", teamID))
	scenarioHosts, err := theServer.backingStore.SelectTeamScenarioHosts(teamID)
	if err != nil {
		msg := "ERROR: cannot retrieve team scenario hosts;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	out, err := json.Marshal(scenarioHosts)
	if err != nil {
		msg := "ERROR: cannot marshal team scenario hosts;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(out)
}

func createEncryptionKeys(filePGPPub string, filePGPPriv string) {
	log.Println("Creating openpgp files")
	pubKey, privKey, err := processing.NewPubPrivKeys()
	if err != nil {
		log.Println("ERROR: cannot get openpgp entity;", err)
		return
	}

	log.Println("Writing openpgp public key")
	err = ioutil.WriteFile(filePGPPub, pubKey, 0600)
	if err != nil {
		log.Println("ERROR: cannot write public key file;", err)
		return
	}

	log.Println("Writing openpgp private key")
	err = ioutil.WriteFile(filePGPPriv, privKey, 0600)
	if err != nil {
		log.Println("ERROR: cannot write private key file;", err)
		return
	}
}

func passwordHash(cleartext string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(cleartext), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func checkPasswordHash(cleartext string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(cleartext))
	if err != nil {
		return false
	}
	return true
}

func (theServer theServer) getAdmins(w http.ResponseWriter, r *http.Request) {
	log.Println("get all admins")

	// get all admins
	admins, err := theServer.backingStore.SelectAdmins()
	if err != nil {
		msg := "ERROR: cannot retrieve admins;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	b, err := json.Marshal(admins)
	if err != nil {
		msg := "ERROR: cannot marshal admins;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(b)
}

func getRandStr() string {
	randKey := securecookie.GenerateRandomKey(32)
	// make sure URL safe, no padding
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randKey)
}

func (theServer theServer) authAdmin(w http.ResponseWriter, r *http.Request) {
	log.Println("Authenticating admin")

	r.ParseForm()
	username := r.Form.Get("username")
	log.Println("username: " + username)
	password := r.Form.Get("password")

	storedPasswordHash, err := theServer.backingStore.SelectAdminPasswordHash(username)
	if err != nil {
		msg := "ERROR: cannot authenticate admin;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	if checkPasswordHash(password, storedPasswordHash) {
		log.Println("User authentication successful")

		value := getRandStr()
		theServer.userTokens[value] = username

		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    value,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		return
	}
	msg := "User authenticated failed"
	log.Println(msg)
	http.Error(w, msg, http.StatusUnauthorized)
}

type credentials struct {
	Username string
	Password string
}

func (theServer theServer) newAdmin(w http.ResponseWriter, r *http.Request) {
	log.Println("new admin")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var creds credentials
	err = json.Unmarshal(body, &creds)
	if err != nil {
		msg := "ERROR: cannot unmarshal credentials;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	if len(creds.Username) == 0 || len(creds.Password) == 0 {
		msg := "ERROR: invalid username or password;"
		log.Println(msg)
		w.Write([]byte(msg))
		return
	}

	log.Println("username: " + creds.Username)
	passwordHash, err := passwordHash(creds.Password)
	if err != nil {
		msg := "ERROR: cannot generate password hash;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	err = theServer.backingStore.InsertAdmin(creds.Username, passwordHash)
	if err != nil {
		msg := "ERROR: cannot insert admin;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	log.Println("Admin added")
}

func (theServer theServer) editAdmin(w http.ResponseWriter, r *http.Request) {
	log.Println("editing admin")

	// parse out int64 id
	vars := mux.Vars(r)

	username := vars["username"]
	log.Println("username: " + username)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var creds credentials
	err = json.Unmarshal(body, &creds)
	if err != nil {
		msg := "ERROR: cannot unmarshal credentials;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	if len(username) == 0 || len(creds.Password) == 0 {
		msg := "ERROR: invalid username or password;"
		log.Println(msg)
		w.Write([]byte(msg))
		return
	}

	passwordHash, err := passwordHash(creds.Password)
	if err != nil {
		msg := "ERROR: cannot generate password hash;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	err = theServer.backingStore.UpdateAdmin(username, passwordHash)
	if err != nil {
		msg := "ERROR: cannot update admin;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	log.Println("Admin edited")
}

func (theServer theServer) deleteAdmin(w http.ResponseWriter, r *http.Request) {
	log.Println("deleting admin")

	vars := mux.Vars(r)

	username := vars["username"]
	log.Println("username: " + username)
	err := theServer.backingStore.DeleteAdmin(username)
	if err != nil {
		msg := "ERROR: cannot delete admin;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	log.Println("Admin deleted")
}

func (theServer theServer) logoutAdmin(w http.ResponseWriter, r *http.Request) {
	log.Println("logout request")
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return
	}
	if user, found := theServer.userTokens[cookie.Value]; found {
		log.Println("Logging out user " + user)
		delete(theServer.userTokens, cookie.Value)
	}
}

func (theServer theServer) getNewHostToken(w http.ResponseWriter, r *http.Request) {
	log.Println("new host token")

	r.ParseForm()
	hostname := r.Form.Get("hostname")
	hostname = strings.TrimSpace(hostname)
	if len(hostname) == 0 {
		log.Println("Empty hostname")
		return
	}

	// record request
	timestamp := time.Now().Unix()
	source := r.RemoteAddr
	token := getRandStr()
	err := theServer.backingStore.InsertHostToken(token, timestamp, hostname, source)
	if err != nil {
		msg := "ERROR: unable to get host token;"
		log.Println(msg, err)
		return
	}

	w.Write([]byte(token))
}

func (theServer theServer) postTeamHostToken(w http.ResponseWriter, r *http.Request) {
	log.Println("team host token")

	timestamp := time.Now().Unix()
	r.ParseForm()
	hostToken := r.Form.Get("host_token")
	if len(hostToken) == 0 {
		http.Error(w, "Host token missing", http.StatusBadRequest)
		return
	}
	hostname := r.Form.Get("hostname")
	if len(hostname) == 0 {
		http.Error(w, "Hostname missing", http.StatusBadRequest)
		return
	}
	teamKey := r.Form.Get("team_key")
	if len(teamKey) == 0 {
		http.Error(w, "Team key missing", http.StatusBadRequest)
		return
	}
	teamID, err := theServer.backingStore.SelectTeamIDForKey(teamKey)
	if err != nil {
		log.Println("Could not get team id;", err)
		http.Error(w, "Team not found", http.StatusBadRequest)
		return
	}

	err = theServer.backingStore.InsertTeamHostToken(teamID, hostname, hostToken, timestamp)
	if err != nil {
		log.Println("ERROR: unable to insert team host token;", err)
		http.Error(w, "Internal server error. Try again later", http.StatusInternalServerError)
		return
	}

	// redirect to team's score page
	url := "https://" + r.Host + "/ui/report?team_key=" + teamKey
	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("ERROR: unable to get executable", err)
	}
	workDir := filepath.Dir(ex)
	publicDir := path.Join(workDir, "public")
	privateDir := path.Join(workDir, "private")
	dataDir := path.Join(workDir, "data")
	configFile := path.Join(workDir, configFileName)

	var askVersion bool

	flag.BoolVar(&askVersion, "version", false, "get version number")
	flag.Parse()

	if askVersion {
		log.Println("Version: " + version)
		os.Exit(0)
	}

	// if config file doesn't exist, generate default config file
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Printf("Creating default config file %s\n", configFile)
		w, err := os.Create(configFile)
		if err != nil {
			log.Fatal("ERROR: could not create default config file;", err)
		}
		w.Chmod(0600)
		fmt.Fprintf(w, "key %s\n", path.Join(privateDir, "server.key"))
		fmt.Fprintf(w, "cert %s\n", path.Join(publicDir, "server.crt"))
		fmt.Fprintf(w, "port %d\n", 8443)
		fmt.Fprintf(w, "sql_url %s\n", "postgres://localhost")
	}

	// read config file
	log.Printf("Reading config file %s\n", configFile)
	configFileBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("ERROR: unable to read config file;", err)
	}
	var fileKey string
	var fileCert string
	var port string
	var sqlURL string
	for _, line := range strings.Split(string(configFileBytes), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.Split(line, " ")
		if tokens[0] == "key" {
			fileKey = tokens[1]
		} else if tokens[0] == "cert" {
			fileCert = tokens[1]
		} else if tokens[0] == "port" {
			port = tokens[1]
		} else if tokens[0] == "sql_url" {
			sqlURL = tokens[1]
		} else {
			log.Fatalf("ERROR: unknown config file setting %s\n", tokens[0])
		}
	}

	err = os.MkdirAll(publicDir, 0700)
	if err != nil {
		log.Fatalln("ERROR: unable to create public directory;", err)
	}
	err = os.MkdirAll(privateDir, 0700)
	if err != nil {
		log.Fatalln("ERROR: unable to create private directory;", err)
	}
	err = os.MkdirAll(dataDir, 0700)
	if err != nil {
		log.Fatalln("ERROR: unable to create data directory;", err)
	}

	log.Println("Using server key file: " + fileKey)
	log.Println("Using server cert file: " + fileCert)
	log.Println("Using network port: " + port)
	log.Println("Using SQL URL: " + sqlURL)

	theServer := theServer{}
	theServer.userTokens = make(map[string]string)
	// theServer.backingStore, err = getBackingStore("sqlite", workDir)
	theServer.backingStore, err = getBackingStore("postgres", sqlURL)
	if err != nil {
		log.Fatal("ERROR: Unable to get backing store;", err)
	}

	// generate default admin if no admins
	admins, err := theServer.backingStore.SelectAdmins()
	if err != nil {
		log.Fatal("Could not get admin list;", err)
	}
	if len(admins) == 0 {
		log.Println("Creating default admin")
		// default credentials
		passwordHash, err := passwordHash("admin")
		if err != nil {
			log.Fatal("ERROR: cannot generate password hash;", err)
		}
		theServer.backingStore.InsertAdmin("admin", passwordHash)
	}

	filePGPPub := path.Join(publicDir, "server.pub")
	filePGPPriv := path.Join(privateDir, "server.priv")
	if _, err := os.Stat(filePGPPriv); os.IsNotExist(err) {
		createEncryptionKeys(filePGPPub, filePGPPriv)
	}

	log.Println("Reading server openpgp private key file")
	privKeyFile, err := os.Open(filePGPPriv)
	if err != nil {
		log.Println("ERROR: cannot read server openpgp private key file;", err)
		return
	}
	entities, err := openpgp.ReadArmoredKeyRing(privKeyFile)
	if err != nil {
		log.Println("ERROR: cannot read entity;", err)
		return
	}
	// process audit requests asynchronously
	go theServer.audit(dataDir, entities)

	r := mux.NewRouter()
	r.Handle("", http.RedirectHandler("/ui/", http.StatusMovedPermanently))
	r.Handle("/", http.RedirectHandler("/ui/", http.StatusMovedPermanently))
	r.PathPrefix("/ui").Handler(http.FileServer(http.Dir(workDir)))

	r.PathPrefix("/public").Handler(http.FileServer(http.Dir(workDir)))

	r.HandleFunc("/audit", func(w http.ResponseWriter, r *http.Request) {
		saveAuditRequest(w, r, dataDir)
	}).Methods("POST")
	r.HandleFunc("/login", theServer.authAdmin).Methods("POST")
	r.HandleFunc("/logout", theServer.logoutAdmin).Methods("DELETE")
	adminRouter := r.PathPrefix("/admins").Subrouter()
	adminRouter.Use(theServer.middleware)
	adminRouter.HandleFunc("", theServer.getAdmins).Methods("GET")
	adminRouter.HandleFunc("/", theServer.getAdmins).Methods("GET")
	adminRouter.HandleFunc("", theServer.newAdmin).Methods("POST")
	adminRouter.HandleFunc("/", theServer.newAdmin).Methods("POST")
	adminRouter.HandleFunc("/{username}", theServer.editAdmin).Methods("POST")
	adminRouter.HandleFunc("/{username}", theServer.deleteAdmin).Methods("DELETE")
	templatesRouter := r.PathPrefix("/templates").Subrouter()
	templatesRouter.Use(theServer.middleware)
	templatesRouter.HandleFunc("", theServer.getTemplates).Methods("GET")
	templatesRouter.HandleFunc("/", theServer.getTemplates).Methods("GET")
	templatesRouter.HandleFunc("", theServer.newTemplate).Methods("POST")
	templatesRouter.HandleFunc("/", theServer.newTemplate).Methods("POST")
	templatesRouter.HandleFunc("/{id:[0-9]+}", theServer.getTemplate).Methods("GET")
	templatesRouter.HandleFunc("/{id:[0-9]+}", theServer.editTemplate).Methods("POST")
	templatesRouter.HandleFunc("/{id:[0-9]+}", theServer.deleteTemplate).Methods("DELETE")
	hostsRouter := r.PathPrefix("/hosts").Subrouter()
	hostsRouter.Use(theServer.middleware)
	hostsRouter.HandleFunc("", theServer.getHosts).Methods("GET")
	hostsRouter.HandleFunc("/", theServer.getHosts).Methods("GET")
	hostsRouter.HandleFunc("", theServer.newHost).Methods("POST")
	hostsRouter.HandleFunc("/", theServer.newHost).Methods("POST")
	hostsRouter.HandleFunc("/{id:[0-9]+}", theServer.getHost).Methods("GET")
	hostsRouter.HandleFunc("/{id:[0-9]+}", theServer.editHost).Methods("POST")
	hostsRouter.HandleFunc("/{id:[0-9]+}", theServer.deleteHost).Methods("DELETE")
	scenariosRouter := r.PathPrefix("/scenarios").Subrouter()
	scenariosRouter.Use(theServer.middleware)
	scenariosRouter.HandleFunc("", theServer.getScenarios).Methods("GET")
	scenariosRouter.HandleFunc("/", theServer.getScenarios).Methods("GET")
	scenariosRouter.HandleFunc("", theServer.newScenario).Methods("POST")
	scenariosRouter.HandleFunc("/", theServer.newScenario).Methods("POST")
	scenariosRouter.HandleFunc("/{id:[0-9]+}", theServer.getScenario).Methods("GET")
	scenariosRouter.HandleFunc("/{id:[0-9]+}", theServer.editScenario).Methods("POST")
	scenariosRouter.HandleFunc("/{id:[0-9]+}", theServer.deleteScenario).Methods("DELETE")
	scoresRouter := r.PathPrefix("/scores").Subrouter()
	// no auth
	scoresRouter.HandleFunc("/scenarios", theServer.getScenariosForScoreboard).Methods("GET")
	scoresRouter.HandleFunc("/scenario/{id:[0-9]+}", theServer.getScenarioScores).Methods("GET")
	reportRouter := r.PathPrefix("/reports").Subrouter()
	// using team key as auth
	reportRouter.HandleFunc("", theServer.getTeamScenarioHosts).Methods("GET")
	reportRouter.HandleFunc("/", theServer.getTeamScenarioHosts).Methods("GET")
	reportRouter.HandleFunc("/scenario/{id:[0-9]+}", theServer.getScenarioScoresReport).Methods("GET")
	reportRouter.HandleFunc("/scenario/{id:[0-9]+}/timeline", theServer.getScenarioScoresTimeline).Methods("GET")
	teamsRouter := r.PathPrefix("/teams").Subrouter()
	teamsRouter.Use(theServer.middleware)
	teamsRouter.HandleFunc("", theServer.getTeams).Methods("GET")
	teamsRouter.HandleFunc("/", theServer.getTeams).Methods("GET")
	teamsRouter.HandleFunc("", theServer.newTeam).Methods("POST")
	teamsRouter.HandleFunc("/", theServer.newTeam).Methods("POST")
	teamsRouter.HandleFunc("/{id:[0-9]+}", theServer.getTeam).Methods("GET")
	teamsRouter.HandleFunc("/{id:[0-9]+}", theServer.editTeam).Methods("POST")
	teamsRouter.HandleFunc("/{id:[0-9]+}", theServer.deleteTeam).Methods("DELETE")
	// no auth
	tokenRouter := r.PathPrefix("/token").Subrouter()
	tokenRouter.HandleFunc("/host", theServer.getNewHostToken).Methods("GET")
	tokenRouter.HandleFunc("/team", theServer.postTeamHostToken).Methods("POST")

	log.Println("Ready to serve requests")
	addr := "0.0.0.0:" + port
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	err = http.ServeTLS(l, r, fileCert, fileKey)
	if err != nil {
		log.Fatal("ERROR: cannot start server;", err)
	}
}
