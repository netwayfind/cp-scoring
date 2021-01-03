package main

import (
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("server")
	port := "8000"
	uiPath := "./ui/build"
	backingStoreStr := "postgres"
	dbURL := "postgres://postgres:password@localhost:5432?sslmode=disable"

	rand.Seed(time.Now().UTC().UnixNano())

	backingStore, err := getBackingStore(backingStoreStr, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	apiHandler := APIHandler{
		BackingStore: backingStore,
	}

	r := mux.NewRouter().StrictSlash(true)
	r.PathPrefix("/ui").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, uiPath+"/index.html")
	})
	r.PathPrefix("/static").Handler(http.FileServer(http.Dir(uiPath)))
	apiRouter := r.PathPrefix("/api").Subrouter()
	r.HandleFunc("/audit", apiHandler.audit).Methods("POST")
	hostTokenRouter := apiRouter.PathPrefix("/host-token").Subrouter()
	hostTokenRouter.HandleFunc("/request", apiHandler.requestHostToken).Methods("POST")
	hostTokenRouter.HandleFunc("/register", apiHandler.registerHostToken).Methods("POST")
	teamRouter := apiRouter.PathPrefix("/teams").Subrouter()
	teamRouter.HandleFunc("/", apiHandler.readTeams).Methods("GET")
	teamRouter.HandleFunc("/", apiHandler.createTeam).Methods("POST")
	teamRouter.HandleFunc("/{id:[0-9]+}", apiHandler.deleteTeam).Methods("DELETE")
	teamRouter.HandleFunc("/{id:[0-9]+}", apiHandler.readTeam).Methods("GET")
	teamRouter.HandleFunc("/{id:[0-9]+}", apiHandler.updateTeam).Methods("PUT")
	scenarioRouter := apiRouter.PathPrefix("/scenarios").Subrouter()
	scenarioRouter.HandleFunc("/", apiHandler.readScenarios).Methods("GET")
	scenarioRouter.HandleFunc("/", apiHandler.createScenario).Methods("POST")
	scenarioRouter.HandleFunc("/{id:[0-9]+}", apiHandler.deleteScenario).Methods("DELETE")
	scenarioRouter.HandleFunc("/{id:[0-9]+}", apiHandler.readScenario).Methods("GET")
	scenarioRouter.HandleFunc("/{id:[0-9]+}", apiHandler.updateScenario).Methods("PUT")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/answers", apiHandler.readScenarioAnswers).Methods("GET")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/answers", apiHandler.updateScenarioAnswers).Methods("PUT")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/checks", apiHandler.readScenarioChecks).Methods("GET")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/checks", apiHandler.updateScenarioChecks).Methods("PUT")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/config", apiHandler.readScenarioConfig).Methods("GET")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/report", apiHandler.readScenarioReport).Methods("GET")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/report/hostnames", apiHandler.readScenarioReportHostnames).Methods("GET")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/report/timeline", apiHandler.readScenarioReportTimeline).Methods("GET")
	scoreboardRouter := apiRouter.PathPrefix("/scoreboard").Subrouter()
	scoreboardRouter.HandleFunc("/scenarios", apiHandler.readScoreboardScenarios).Methods("GET")
	scoreboardRouter.HandleFunc("/scenarios/{id:[0-9]+}", apiHandler.readScoreboardForScenario).Methods("GET")

	log.Println("Ready to serve requests")
	addr := "0.0.0.0:" + port
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}

	err = http.Serve(l, r)
	if err != nil {
		log.Fatal("ERROR: cannot start server;", err)
	}

}
