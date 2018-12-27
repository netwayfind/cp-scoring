package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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

type authenticationMiddleware struct {
	tokenUsers map[string]string
}

func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			msg := "Unauthorized request"
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		if username, found := amw.tokenUsers[cookie.Value]; found {
			isAdmin, err := dbIsAdmin(username)
			if err != nil {
				msg := "ERROR: unable to check if user is admin"
				log.Println(msg, err)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			if !isAdmin {
				// delete any existing sessions
				delete(amw.tokenUsers, cookie.Value)
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

func audit(dataDir string, entities openpgp.EntityList) {
	for {
		err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			log.Println("Auditing file " + path)
			auditErr := auditFile(path, entities)
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

func auditFile(fileStr string, entities openpgp.EntityList) error {
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
	err = dbInsertState(string(fileBytes))
	if err != nil {
		log.Println("ERROR: cannot insert state;")
		return err
	}

	log.Println("Getting information")
	hostToken := stateSubmission.HostToken
	teamID, err := dbSelectTeamIDFromHostToken(hostToken)
	if err != nil {
		log.Println("ERROR: cannot get team id;")
		return err
	}
	log.Println(fmt.Sprintf("Team ID: %d", teamID))
	hostID, err := dbSelectHostIDForHostname(state.Hostname)
	if err != nil {
		log.Println("ERROR: cannot get host id;")
		return err
	}
	log.Println(fmt.Sprintf("Host ID: %d", hostID))

	log.Println("Getting scenarios")
	scenarioIDs, err := dbSelectScenariosForHostname(state.Hostname)
	if err != nil {
		log.Println("ERROR: cannot get scenario IDs;")
		return err
	}
	if len(scenarioIDs) == 0 {
		log.Println("ERROR: no scenarios found")
		return nil
	}
	for _, scenarioID := range scenarioIDs {
		log.Println(fmt.Sprintf("Processing scenario ID: %d", scenarioID))

		log.Println("Getting scenario templates")
		templates, err := dbSelectTemplatesForHostname(scenarioID, state.Hostname)
		if err != nil {
			log.Println("ERROR: cannot get templates;")
			return err
		}
		log.Println(fmt.Sprintf("Found template count: %d", len(templates)))
		if len(templates) == 0 {
			log.Println("ERROR: no templates found")
			return nil
		}

		log.Println("Auditing state")
		report := auditor.Audit(state, templates)

		log.Println("Saving report")
		report.Timestamp = state.Timestamp
		err = dbInsertScenarioReport(scenarioID, teamID, hostID, report)
		if err != nil {
			log.Println("ERROR: cannot insert report;")
			return err
		}

		log.Println("Saving score")
		var score int64
		for _, finding := range report.Findings {
			score += finding.Value
		}

		var scoreEntry model.ScenarioScore
		scoreEntry.ScenarioID = scenarioID
		scoreEntry.TeamID = teamID
		scoreEntry.HostID = hostID
		scoreEntry.Timestamp = state.Timestamp
		scoreEntry.Score = score
		err = dbInsertScenarioScore(scoreEntry)
		if err != nil {
			log.Println("ERROR: cannot insert scenario score;")
			return err
		}
	}

	log.Println("Received and saved")
	return nil
}

func getHosts(w http.ResponseWriter, r *http.Request) {
	log.Println("get all hosts")

	// get all hosts
	hosts, err := dbSelectHosts()
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

func getHost(w http.ResponseWriter, r *http.Request) {
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
	host, err := dbSelectHost(id)
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

func newHost(w http.ResponseWriter, r *http.Request) {
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

	err = dbInsertHost(host)
	if err != nil {
		msg := "ERROR: cannot insert host;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	// new host
	msg := "Saved host"
	log.Println(msg)
	w.Write([]byte(msg))
}

func editHost(w http.ResponseWriter, r *http.Request) {
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

	err = dbUpdateHost(id, host)
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

func deleteHost(w http.ResponseWriter, r *http.Request) {
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
	err = dbDeleteHost(id)
	if err != nil {
		msg := "ERROR: cannot delete host;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
}

func getTeams(w http.ResponseWriter, r *http.Request) {
	log.Println("get all teams")

	// get all teams
	teams, err := dbSelectTeams()
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

func getTeam(w http.ResponseWriter, r *http.Request) {
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
	team, err := dbSelectTeam(id)
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

func newTeam(w http.ResponseWriter, r *http.Request) {
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

	err = dbInsertTeam(team)
	if err != nil {
		msg := "ERROR: cannot insert team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	// new team
	msg := "Saved team"
	log.Println(msg)
	w.Write([]byte(msg))
}

func editTeam(w http.ResponseWriter, r *http.Request) {
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

	err = dbUpdateTeam(id, team)
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

func deleteTeam(w http.ResponseWriter, r *http.Request) {
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
	err = dbDeleteTeam(id)
	if err != nil {
		msg := "ERROR: cannot delete team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
}

func getTemplates(w http.ResponseWriter, r *http.Request) {
	log.Println("get all templates")

	// get all templates
	templates, err := dbSelectTemplates()
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

func getTemplate(w http.ResponseWriter, r *http.Request) {
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
	template, err := dbSelectTemplate(id)
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

func newTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("new template")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var templateEntry model.TemplateEntry
	err = json.Unmarshal(body, &templateEntry)
	if err != nil {
		msg := "ERROR: cannot unmarshal template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	err = dbInsertTemplate(templateEntry)
	if err != nil {
		msg := "ERROR: cannot insert template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	// new template
	msg := "Saved template"
	log.Println(msg)
	w.Write([]byte(msg))
}

func editTemplate(w http.ResponseWriter, r *http.Request) {
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

	var templateEntry model.TemplateEntry
	err = json.Unmarshal(body, &templateEntry)
	if err != nil {
		msg := "ERROR: cannot unmarshal template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	err = dbUpdateTemplate(id, templateEntry)
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

func deleteTemplate(w http.ResponseWriter, r *http.Request) {
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
	err = dbDeleteTemplate(id)
	if err != nil {
		msg := "ERROR: cannot delete template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
}

func getScenarios(w http.ResponseWriter, r *http.Request) {
	log.Println("get all scenarios")

	// get all scenarios, even not enabled
	scenarios, err := dbSelectScenarios(false)
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

func getScenario(w http.ResponseWriter, r *http.Request) {
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
	scenario, err := dbSelectScenario(id)
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

func newScenario(w http.ResponseWriter, r *http.Request) {
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

	err = dbInsertScenario(scenario)
	if err != nil {
		msg := "ERROR: cannot insert scenario;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	// new scenario
	msg := "Saved scenario"
	log.Println(msg)
	w.Write([]byte(msg))
}

func editScenario(w http.ResponseWriter, r *http.Request) {
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

	err = dbUpdateScenario(id, scenario)
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

func deleteScenario(w http.ResponseWriter, r *http.Request) {
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
	err = dbDeleteScenario(id)
	if err != nil {
		msg := "ERROR: cannot delete scenario;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
}

func getScenariosForScoreboard(w http.ResponseWriter, r *http.Request) {
	log.Println("get scenarios for scoreboard")

	// get scenarios, only enabled
	scenarios, err := dbSelectScenarios(true)
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

func getScenarioScores(w http.ResponseWriter, r *http.Request) {
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
	log.Println(id)
	scores, err := dbSelectScenarioLatestScores(id)
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

func getScenarioScoresTimeline(w http.ResponseWriter, r *http.Request) {
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
	teamID, err := dbSelectTeamIDForKey(teamKey)
	if err != nil {
		msg := "ERROR: cannot retrieve team id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(fmt.Sprintf("Team ID: %d", teamID))
	hostname := r.FormValue("hostname")
	hostID, err := dbSelectHostIDForHostname(hostname)
	if err != nil {
		msg := "ERROR: cannot retrieve host id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(fmt.Sprintf("Host ID: %d", hostID))
	timeline, err := dbSelectScenarioTimeline(scenarioID, teamID, hostID)
	if err != nil {
		msg := "ERROR: cannot retrieve scenario timeline for team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	out, err := json.Marshal(timeline)
	if err != nil {
		msg := "ERROR: cannot marshal scenario timeline for team;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(out)
}

func getScenarioScoresReport(w http.ResponseWriter, r *http.Request) {
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
	teamID, err := dbSelectTeamIDForKey(teamKey)
	if err != nil {
		msg := "ERROR: cannot retrieve team id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(fmt.Sprintf("Team ID: %d", teamID))
	hostname := r.FormValue("hostname")
	hostID, err := dbSelectHostIDForHostname(hostname)
	if err != nil {
		msg := "ERROR: cannot retrieve host id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(fmt.Sprintf("Host ID: %d", hostID))
	report, err := dbSelectLatestScenarioReport(scenarioID, teamID, hostID)
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

func getTeamScenarioHosts(w http.ResponseWriter, r *http.Request) {
	log.Println("get team scenario hosts")

	teamKey := r.FormValue("team_key")
	teamID, err := dbSelectTeamIDForKey(teamKey)
	if err != nil {
		msg := "ERROR: cannot retrieve team id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(fmt.Sprintf("Team ID: %d", teamID))
	scenarioHosts, err := dbSelectTeamScenarioHosts(teamID)
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

func getAdmins(w http.ResponseWriter, r *http.Request) {
	log.Println("get all admins")

	// get all admins
	admins, err := dbSelectAdmins()
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
	return base64.StdEncoding.EncodeToString(randKey)
}

func authAdmin(w http.ResponseWriter, r *http.Request, amw authenticationMiddleware) {
	log.Println("Authenticating admin")

	r.ParseForm()
	username := r.Form.Get("username")
	log.Println("username: " + username)
	password := r.Form.Get("password")

	storedPasswordHash, err := dbSelectAdminPasswordHash(username)
	if err != nil {
		msg := "ERROR: cannot authenticate admin;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	if checkPasswordHash(password, storedPasswordHash) {
		log.Println("User authentication successful")

		value := getRandStr()
		amw.tokenUsers[value] = username

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

func newAdmin(w http.ResponseWriter, r *http.Request) {
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
	err = dbInsertAdmin(creds.Username, passwordHash)
	if err != nil {
		msg := "ERROR: cannot insert admin;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	log.Println("Admin added")
}

func editAdmin(w http.ResponseWriter, r *http.Request) {
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

	err = dbUpdateAdmin(username, passwordHash)
	if err != nil {
		msg := "ERROR: cannot update admin;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	log.Println("Admin edited")
}

func deleteAdmin(w http.ResponseWriter, r *http.Request) {
	log.Println("deleting admin")

	vars := mux.Vars(r)

	username := vars["username"]
	log.Println("username: " + username)
	err = dbDeleteAdmin(username)
	if err != nil {
		msg := "ERROR: cannot delete admin;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	log.Println("Admin deleted")
}

func logoutAdmin(w http.ResponseWriter, r *http.Request, amw authenticationMiddleware) {
	log.Println("logout request")
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return
	}
	if user, found := amw.tokenUsers[cookie.Value]; found {
		log.Println("Logging out user " + user)
		delete(amw.tokenUsers, cookie.Value)
	}
}

func getNewHostToken(w http.ResponseWriter, r *http.Request) {
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
	err := dbInsertHostToken(token, timestamp, hostname, source)
	if err != nil {
		msg := "ERROR: unable to get host token;"
		log.Println(msg, err)
		return
	}

	w.Write([]byte(token))
}

func postTeamHostToken(w http.ResponseWriter, r *http.Request) {
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
	hostID, err := dbSelectHostIDForHostname(hostname)
	if err != nil {
		log.Println("Could not get host id;", err)
		http.Error(w, "Host not found", http.StatusBadRequest)
		return
	}
	teamKey := r.Form.Get("team_key")
	if len(teamKey) == 0 {
		http.Error(w, "Team key missing", http.StatusBadRequest)
		return
	}
	teamID, err := dbSelectTeamIDForKey(teamKey)
	if err != nil {
		log.Println("Could not get team id;", err)
		http.Error(w, "Team not found", http.StatusBadRequest)
		return
	}

	err = dbInsertTeamHostToken(teamID, hostID, hostToken, timestamp)
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
	dir := filepath.Dir(ex)

	dbInit(dir)

	// generate default admin if no admins
	admins, err := dbSelectAdmins()
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
		dbInsertAdmin("admin", passwordHash)
	}

	publicDir := path.Join(dir, "public")
	privateDir := path.Join(dir, "private")
	dataDir := path.Join(dir, "data")
	err = os.MkdirAll(publicDir, 0700)
	err = os.MkdirAll(privateDir, 0700)
	err = os.MkdirAll(dataDir, 0700)

	var fileKey string
	var fileCert string
	var port int

	flag.StringVar(&fileKey, "key", path.Join(privateDir, "server.key"), "server key")
	flag.StringVar(&fileCert, "cert", path.Join(publicDir, "server.crt"), "server cert")
	flag.IntVar(&port, "port", 8443, "port")
	flag.Parse()

	log.Println("Using server key file: " + fileKey)
	log.Println("Using server cert file: " + fileCert)
	log.Println("Using network port: " + strconv.Itoa(port))

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
	go audit(dataDir, entities)

	authenticator := authenticationMiddleware{}
	authenticator.tokenUsers = make(map[string]string)

	r := mux.NewRouter()
	r.Handle("", http.RedirectHandler("/ui/", http.StatusMovedPermanently))
	r.Handle("/", http.RedirectHandler("/ui/", http.StatusMovedPermanently))
	r.PathPrefix("/ui").Handler(http.FileServer(http.Dir(dir)))

	r.PathPrefix("/public").Handler(http.FileServer(http.Dir(dir)))

	r.HandleFunc("/audit", func(w http.ResponseWriter, r *http.Request) {
		saveAuditRequest(w, r, dataDir)
	}).Methods("POST")
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		authAdmin(w, r, authenticator)
	}).Methods("POST")
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logoutAdmin(w, r, authenticator)
	}).Methods("DELETE")
	adminRouter := r.PathPrefix("/admins").Subrouter()
	adminRouter.Use(authenticator.Middleware)
	adminRouter.HandleFunc("", getAdmins).Methods("GET")
	adminRouter.HandleFunc("/", getAdmins).Methods("GET")
	adminRouter.HandleFunc("", newAdmin).Methods("POST")
	adminRouter.HandleFunc("/", newAdmin).Methods("POST")
	adminRouter.HandleFunc("/{username}", editAdmin).Methods("POST")
	adminRouter.HandleFunc("/{username}", deleteAdmin).Methods("DELETE")
	templatesRouter := r.PathPrefix("/templates").Subrouter()
	templatesRouter.Use(authenticator.Middleware)
	templatesRouter.HandleFunc("", getTemplates).Methods("GET")
	templatesRouter.HandleFunc("/", getTemplates).Methods("GET")
	templatesRouter.HandleFunc("", newTemplate).Methods("POST")
	templatesRouter.HandleFunc("/", newTemplate).Methods("POST")
	templatesRouter.HandleFunc("/{id:[0-9]+}", getTemplate).Methods("GET")
	templatesRouter.HandleFunc("/{id:[0-9]+}", editTemplate).Methods("POST")
	templatesRouter.HandleFunc("/{id:[0-9]+}", deleteTemplate).Methods("DELETE")
	hostsRouter := r.PathPrefix("/hosts").Subrouter()
	hostsRouter.Use(authenticator.Middleware)
	hostsRouter.HandleFunc("", getHosts).Methods("GET")
	hostsRouter.HandleFunc("/", getHosts).Methods("GET")
	hostsRouter.HandleFunc("", newHost).Methods("POST")
	hostsRouter.HandleFunc("/", newHost).Methods("POST")
	hostsRouter.HandleFunc("/{id:[0-9]+}", getHost).Methods("GET")
	hostsRouter.HandleFunc("/{id:[0-9]+}", editHost).Methods("POST")
	hostsRouter.HandleFunc("/{id:[0-9]+}", deleteHost).Methods("DELETE")
	scenariosRouter := r.PathPrefix("/scenarios").Subrouter()
	scenariosRouter.Use(authenticator.Middleware)
	scenariosRouter.HandleFunc("", getScenarios).Methods("GET")
	scenariosRouter.HandleFunc("/", getScenarios).Methods("GET")
	scenariosRouter.HandleFunc("", newScenario).Methods("POST")
	scenariosRouter.HandleFunc("/", newScenario).Methods("POST")
	scenariosRouter.HandleFunc("/{id:[0-9]+}", getScenario).Methods("GET")
	scenariosRouter.HandleFunc("/{id:[0-9]+}", editScenario).Methods("POST")
	scenariosRouter.HandleFunc("/{id:[0-9]+}", deleteScenario).Methods("DELETE")
	scoresRouter := r.PathPrefix("/scores").Subrouter()
	// no auth
	scoresRouter.HandleFunc("/scenarios", getScenariosForScoreboard).Methods("GET")
	scoresRouter.HandleFunc("/scenario/{id:[0-9]+}", getScenarioScores).Methods("GET")
	reportRouter := r.PathPrefix("/reports").Subrouter()
	// using team key as auth
	reportRouter.HandleFunc("", getTeamScenarioHosts).Methods("GET")
	reportRouter.HandleFunc("/", getTeamScenarioHosts).Methods("GET")
	reportRouter.HandleFunc("/scenario/{id:[0-9]+}", getScenarioScoresReport).Methods("GET")
	reportRouter.HandleFunc("/scenario/{id:[0-9]+}/timeline", getScenarioScoresTimeline).Methods("GET")
	teamsRouter := r.PathPrefix("/teams").Subrouter()
	teamsRouter.Use(authenticator.Middleware)
	teamsRouter.HandleFunc("", getTeams).Methods("GET")
	teamsRouter.HandleFunc("/", getTeams).Methods("GET")
	teamsRouter.HandleFunc("", newTeam).Methods("POST")
	teamsRouter.HandleFunc("/", newTeam).Methods("POST")
	teamsRouter.HandleFunc("/{id:[0-9]+}", getTeam).Methods("GET")
	teamsRouter.HandleFunc("/{id:[0-9]+}", editTeam).Methods("POST")
	teamsRouter.HandleFunc("/{id:[0-9]+}", deleteTeam).Methods("DELETE")
	// no auth
	tokenRouter := r.PathPrefix("/token").Subrouter()
	tokenRouter.HandleFunc("/host", getNewHostToken).Methods("GET")
	tokenRouter.HandleFunc("/team", postTeamHostToken).Methods("POST")

	log.Println("Ready to serve requests")
	addr := ":" + strconv.Itoa(port)
	err = http.ListenAndServeTLS(addr, fileCert, fileKey, r)
	if err != nil {
		log.Println("ERROR: cannot start server;", err)
	}
}
