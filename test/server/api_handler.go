package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/netwayfind/cp-scoring/test/model"
)

// APIHandler asdf
type APIHandler struct {
	BackingStore backingStore
}

func getRequestID(r *http.Request) (uint64, error) {
	vars := mux.Vars(r)

	return strconv.ParseUint(vars["id"], 10, 64)
}

func getSourceIP(r *http.Request) string {
	conn := r.RemoteAddr
	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		conn = strings.Split(ips, ",")[0]
	}
	return strings.Split(conn, ":")[0]
}

func httpErrorBadRequest(w http.ResponseWriter) {
	msg := "ERROR: bad request;"
	http.Error(w, msg, http.StatusBadRequest)
}

func httpErrorDatabase(w http.ResponseWriter, err error) {
	msg := "ERROR: database query;"
	log.Println(msg, err)
	http.Error(w, msg, http.StatusInternalServerError)
}

func httpErrorInvalidID(w http.ResponseWriter) {
	msg := "ERROR: cannot parse scenario id;"
	log.Println(msg)
	http.Error(w, msg, http.StatusBadRequest)
}

func httpErrorMarshall(w http.ResponseWriter, err error) {
	msg := "ERROR: cannot marshall;"
	log.Println(msg, err)
	http.Error(w, msg, http.StatusInternalServerError)
}

func httpErrorNotFound(w http.ResponseWriter) {
	msg := "ERROR: not found;"
	log.Println(msg)
	http.Error(w, msg, http.StatusNotFound)
}

func httpErrorReadRequestBody(w http.ResponseWriter, err error) {
	msg := "ERROR: cannot read request body;"
	log.Println(msg, err)
	http.Error(w, msg, http.StatusInternalServerError)
}

func httpErrorUnmarshall(w http.ResponseWriter, err error) {
	msg := "ERROR: cannot unmarshall;"
	log.Println(msg, err)
	http.Error(w, msg, http.StatusBadRequest)
}

func readRequestBody(w http.ResponseWriter, r *http.Request, o interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpErrorReadRequestBody(w, err)
		return err
	}

	err = json.Unmarshal(body, &o)
	if err != nil {
		httpErrorUnmarshall(w, err)
		return err
	}

	return err
}

func sendResponse(w http.ResponseWriter, o interface{}) {
	b, err := json.Marshal(o)
	if err != nil {
		httpErrorMarshall(w, err)
		return
	}
	w.Write(b)
}

func (handler APIHandler) audit(w http.ResponseWriter, r *http.Request) {
	log.Println("audit")

	source := getSourceIP(r)
	timestamp := time.Now().Unix()

	var auditCheckResults model.AuditCheckResults
	err := readRequestBody(w, r, &auditCheckResults)
	if err != nil {
		return
	}

	scenario, err := handler.BackingStore.scenarioSelect(auditCheckResults.ScenarioID)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if scenario.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	answersMap, err := handler.BackingStore.scenarioAnswersSelectAll(auditCheckResults.ScenarioID)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	if len(auditCheckResults.HostToken) == 0 {
		httpErrorBadRequest(w)
		return
	}

	teamID, err := handler.BackingStore.hostTokenSelectTeamID(auditCheckResults.HostToken)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if teamID == 0 {
		httpErrorNotFound(w)
		return
	}

	checkResultsID, err := handler.BackingStore.auditCheckResultsInsert(auditCheckResults, teamID, timestamp, source)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	hostname, err := handler.BackingStore.hostTokenSelectHostname(auditCheckResults.HostToken)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if len(hostname) == 0 {
		httpErrorBadRequest(w)
		return
	}

	answers := answersMap[hostname]
	if answers == nil {
		httpErrorNotFound(w)
		return
	}

	if len(answers) != len(auditCheckResults.CheckResults) {
		httpErrorBadRequest(w)
		return
	}

	answerResults := make([]bool, len(answers))
	score := 0
	for i, answer := range answers {
		checkResult := auditCheckResults.CheckResults[i]
		if answer.Operator == model.OperatorTypeEqual {
			matched := answer.Value == checkResult
			answerResults[i] = matched
			if matched {
				score += answer.Points
			}
		}
	}

	auditAnswerResults := model.AuditAnswerResults{
		ScenarioID:     scenario.ID,
		TeamID:         teamID,
		HostToken:      auditCheckResults.HostToken,
		Timestamp:      auditCheckResults.Timestamp,
		CheckResultsID: checkResultsID,
		Score:          score,
		AnswerResults:  answerResults,
	}

	err = handler.BackingStore.auditAnswerResultsInsert(auditAnswerResults)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	err = handler.BackingStore.scoreboardUpdate(scenario.ID, teamID, hostname, score, auditCheckResults.Timestamp)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
}

func (handler APIHandler) requestHostToken(w http.ResponseWriter, r *http.Request) {
	log.Println("request host token")

	var hostTokenRequest model.HostTokenRequest
	err := readRequestBody(w, r, &hostTokenRequest)
	if err != nil {
		return
	}

	hostToken := randHexStr(16)
	hostname := hostTokenRequest.Hostname
	timestamp := time.Now().Unix()
	sourceIP := getSourceIP(r)
	err = handler.BackingStore.hostTokenInsert(hostToken, hostname, timestamp, sourceIP)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, hostToken)
}

func (handler APIHandler) registerHostToken(w http.ResponseWriter, r *http.Request) {
	log.Println("register host token")

	var hostTokenRegistration model.HostTokenRegistration
	err := readRequestBody(w, r, &hostTokenRegistration)
	if err != nil {
		return
	}

	team, err := handler.BackingStore.teamSelectByKey(hostTokenRegistration.TeamKey)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if team.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	hostToken := hostTokenRegistration.HostToken
	timestamp := time.Now().Unix()

	err = handler.BackingStore.teamHostTokenInsert(team.ID, hostToken, timestamp)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
}

func (handler APIHandler) createScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("create scenario")

	var scenario model.Scenario
	err := readRequestBody(w, r, &scenario)
	if err != nil {
		return
	}

	s, err := handler.BackingStore.scenarioInsert(scenario)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) deleteScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("delete scenario")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	team, err := handler.BackingStore.scenarioSelect(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if team.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	err = handler.BackingStore.scenarioDelete(id)
	if err != nil {
		httpErrorDatabase(w, err)
	}
}

func (handler APIHandler) readScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenario")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	s, err := handler.BackingStore.scenarioSelect(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if s.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) readScenarios(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenarios")

	s, err := handler.BackingStore.scenarioSelectAll()
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) updateScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("update scenario")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	var scenario model.Scenario
	err = readRequestBody(w, r, &scenario)
	if err != nil {
		return
	}

	s, err := handler.BackingStore.scenarioUpdate(id, scenario)
	if err != nil {
		if err.Error() == model.ErrorDBUpdateNoChange {
			httpErrorNotFound(w)
			return
		}
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) readScenarioAnswers(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenario answers")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	s, err := handler.BackingStore.scenarioAnswersSelectAll(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) deleteScenarioAnswers(w http.ResponseWriter, r *http.Request) {
	log.Println("delete scenario answers")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	err = handler.BackingStore.scenarioAnswersDelete(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
}

func (handler APIHandler) updateScenarioAnswers(w http.ResponseWriter, r *http.Request) {
	log.Println("update scenario answers")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	var answerMap map[string][]model.Answer
	err = readRequestBody(w, r, &answerMap)
	if err != nil {
		return
	}

	err = handler.BackingStore.scenarioAnswersUpdate(id, answerMap)
	if err != nil {
		if err.Error() == model.ErrorDBUpdateNoChange {
			httpErrorNotFound(w)
			return
		}
		httpErrorDatabase(w, err)
		return
	}
}

func (handler APIHandler) readScenarioChecks(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenario checks")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	s, err := handler.BackingStore.scenarioChecksSelectAll(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) updateScenarioChecks(w http.ResponseWriter, r *http.Request) {
	log.Println("update scenario checks")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	var hostnameChecks map[string][]model.Action
	err = readRequestBody(w, r, &hostnameChecks)
	if err != nil {
		return
	}

	err = handler.BackingStore.scenarioChecksUpdate(id, hostnameChecks)
	if err != nil {
		if err.Error() == model.ErrorDBUpdateNoChange {
			httpErrorNotFound(w)
			return
		}
		httpErrorDatabase(w, err)
		return
	}
}

func (handler APIHandler) readScenarioConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenario config")
}

func (handler APIHandler) readScoreboardForScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("read scoreboard for scenarios")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	s, err := handler.BackingStore.scoreboardSelectByScenarioID(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) readScoreboardScenarios(w http.ResponseWriter, r *http.Request) {
	log.Println("read scoreboard scenarios")

	s, err := handler.BackingStore.scoreboardSelectScenarios()
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) createTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("create team")

	var team model.Team
	err := readRequestBody(w, r, &team)
	if err != nil {
		return
	}

	t, err := handler.BackingStore.teamInsert(team)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, t)
}

func (handler APIHandler) deleteTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("delete team")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	team, err := handler.BackingStore.teamSelect(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if team.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	err = handler.BackingStore.teamDelete(id)
	if err != nil {
		httpErrorDatabase(w, err)
	}
}

func (handler APIHandler) readTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("read team")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	t, err := handler.BackingStore.teamSelect(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if t.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	sendResponse(w, t)
}

func (handler APIHandler) readTeams(w http.ResponseWriter, r *http.Request) {
	log.Println("read teams")

	t, err := handler.BackingStore.teamSelectAll()
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, t)
}

func (handler APIHandler) updateTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("update team")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}
	log.Println(id)

	var team model.Team
	err = readRequestBody(w, r, &team)
	if err != nil {
		return
	}

	s, err := handler.BackingStore.teamUpdate(id, team)
	if err != nil {
		if err.Error() == model.ErrorDBUpdateNoChange {
			httpErrorNotFound(w)
			return
		}
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}
