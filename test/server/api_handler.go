package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
	var result model.ScenarioHostResult
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)

	answers := make([]model.Answer, 1)
	answers[0] = model.Answer{
		Operator: model.OperatorTypeEqual,
		Value:    "1000",
	}

	match := false

	if len(answers) == len(result.Findings) {
		for i, answer := range answers {
			finding := result.Findings[i]
			if answer.Operator == model.OperatorTypeEqual {
				log.Println(answer.Value)
				log.Println(finding)
				if answer.Value == finding {
					match = true
				}
			}
		}
	}

	log.Println(match)
}

func (handler APIHandler) readNewHostToken(w http.ResponseWriter, r *http.Request) {
	log.Println("read new host token")

	x := randHexStr(16)
	w.Write([]byte(x))
}

func (handler APIHandler) registerHostToken(w http.ResponseWriter, r *http.Request) {
	log.Println("register host token")
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
