package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sumwonyuno/cp-scoring/auditor"
	"github.com/sumwonyuno/cp-scoring/model"
)

func audit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		msg := "HTTP 405"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var state model.State
	err = json.Unmarshal(body, &state)
	if err != nil {
		msg := "ERROR: cannot unmarshal state;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	log.Println("Saving state")
	err = dbInsertState(string(body))
	if err != nil {
		msg := "ERROR: cannot insert state;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	log.Println("Auditing state")
	templates, err := dbSelectTemplatesForHostname(state.Hostname)
	if err != nil {
		msg := "ERROR: cannot get templates;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	report := auditor.Audit(state, templates)
	hostID, err := dbSelectHostIDForHostname(state.Hostname)
	if err != nil {
		msg := "ERROR: cannot get host id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	report.HostID = hostID
	teamID, err := dbSelectTeamIDForKey(state.TeamKey)
	if err != nil {
		msg := "ERROR: cannot get team id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	report.TeamID = teamID
	log.Println(report)

	response := "Received and saved"
	log.Println(response)
	w.Write([]byte(response))
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

	// get all scenarios
	scenarios, err := dbSelectScenarios()
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

func main() {
	dbInit()

	r := mux.NewRouter()
	r.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("./ui/"))))

	r.HandleFunc("/audit", audit)
	templatesRouter := r.PathPrefix("/templates").Subrouter()
	templatesRouter.HandleFunc("", getTemplates).Methods("GET")
	templatesRouter.HandleFunc("/", getTemplates).Methods("GET")
	templatesRouter.HandleFunc("", newTemplate).Methods("POST")
	templatesRouter.HandleFunc("/", newTemplate).Methods("POST")
	templatesRouter.HandleFunc("/{id:[0-9]+}", getTemplate).Methods("GET")
	templatesRouter.HandleFunc("/{id:[0-9]+}", editTemplate).Methods("POST")
	templatesRouter.HandleFunc("/{id:[0-9]+}", deleteTemplate).Methods("DELETE")
	hostsRouter := r.PathPrefix("/hosts").Subrouter()
	hostsRouter.HandleFunc("", getHosts).Methods("GET")
	hostsRouter.HandleFunc("/", getHosts).Methods("GET")
	hostsRouter.HandleFunc("", newHost).Methods("POST")
	hostsRouter.HandleFunc("/", newHost).Methods("POST")
	hostsRouter.HandleFunc("/{id:[0-9]+}", getHost).Methods("GET")
	hostsRouter.HandleFunc("/{id:[0-9]+}", editHost).Methods("POST")
	hostsRouter.HandleFunc("/{id:[0-9]+}", deleteHost).Methods("DELETE")
	scenariosRouter := r.PathPrefix("/scenarios").Subrouter()
	scenariosRouter.HandleFunc("", getScenarios).Methods("GET")
	scenariosRouter.HandleFunc("/", getScenarios).Methods("GET")
	scenariosRouter.HandleFunc("", newScenario).Methods("POST")
	scenariosRouter.HandleFunc("/", newScenario).Methods("POST")
	scenariosRouter.HandleFunc("/{id:[0-9]+}", getScenario).Methods("GET")
	scenariosRouter.HandleFunc("/{id:[0-9]+}", editScenario).Methods("POST")
	scenariosRouter.HandleFunc("/{id:[0-9]+}", deleteScenario).Methods("DELETE")
	teamsRouter := r.PathPrefix("/teams").Subrouter()
	teamsRouter.HandleFunc("", getTeams).Methods("GET")
	teamsRouter.HandleFunc("/", getTeams).Methods("GET")
	teamsRouter.HandleFunc("", newTeam).Methods("POST")
	teamsRouter.HandleFunc("/", newTeam).Methods("POST")
	teamsRouter.HandleFunc("/{id:[0-9]+}", getTeam).Methods("GET")
	teamsRouter.HandleFunc("/{id:[0-9]+}", editTeam).Methods("POST")
	teamsRouter.HandleFunc("/{id:[0-9]+}", deleteTeam).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}
