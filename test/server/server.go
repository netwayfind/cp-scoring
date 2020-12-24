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

	rand.Seed(time.Now().UTC().UnixNano())

	backingStore, err := getBackingStore("postgres", "postgres://postgres:password@localhost:5432?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	apiHandler := APIHandler{
		BackingStore: backingStore,
	}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/audit", apiHandler.audit).Methods("POST")
	hostTokenRouter := r.PathPrefix("/host-token").Subrouter()
	hostTokenRouter.HandleFunc("/", apiHandler.readNewHostToken).Methods("GET")
	hostTokenRouter.HandleFunc("/", apiHandler.registerHostToken).Methods("POST")
	teamRouter := r.PathPrefix("/teams").Subrouter()
	teamRouter.HandleFunc("/", apiHandler.readTeams).Methods("GET")
	teamRouter.HandleFunc("/", apiHandler.createTeam).Methods("POST")
	teamRouter.HandleFunc("/{id:[0-9]+}", apiHandler.readTeam).Methods("GET")
	teamRouter.HandleFunc("/{id:[0-9]+}", apiHandler.updateTeam).Methods("PUT")
	scenarioRouter := r.PathPrefix("/scenarios").Subrouter()
	scenarioRouter.HandleFunc("/", apiHandler.readScenarios).Methods("GET")
	scenarioRouter.HandleFunc("/", apiHandler.createScenario).Methods("POST")
	scenarioRouter.HandleFunc("/{id:[0-9]+}", apiHandler.readScenario).Methods("GET")
	scenarioRouter.HandleFunc("/{id:[0-9]+}", apiHandler.updateScenario).Methods("PUT")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/checks", apiHandler.readScenarioChecks).Methods("GET")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/checks", apiHandler.updateScenarioChecks).Methods("PUT")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/config", apiHandler.readScenarioConfig).Methods("GET")

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
