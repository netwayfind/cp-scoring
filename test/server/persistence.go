package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
	"github.com/netwayfind/cp-scoring/test/model"
)

type backingStore interface {
	scenarioDelete(id uint64) error
	scenarioInsert(scenario model.Scenario) (model.Scenario, error)
	scenarioSelect(id uint64) (model.Scenario, error)
	scenarioSelectAll() ([]model.ScenarioSummary, error)
	scenarioUpdate(id uint64, scenario model.Scenario) (model.Scenario, error)
	scenarioChecksSelectAll(id uint64) (map[string][]model.Action, error)
	scenarioChecksDelete(id uint64) error
	scenarioChecksUpdate(id uint64, hostnameChecks map[string][]model.Action) error
	teamDelete(id uint64) error
	teamInsert(team model.Team) (model.Team, error)
	teamSelect(id uint64) (model.Team, error)
	teamSelectAll() ([]model.TeamSummary, error)
	teamUpdate(id uint64, team model.Team) (model.Team, error)
}

func getBackingStore(store string, args ...string) (backingStore, error) {
	if store == "postgres" {
		// must have first argument as URL
		if len(args) < 1 {
			return nil, errors.New("ERROR: URL required")
		}
		dbConn, err := sql.Open("postgres", args[0])
		if err != nil {
			return nil, err
		}
		log.Println("New connection to database")

		db := dbObj{
			dbConn: dbConn,
		}
		db.dbInit()
		return db, nil
	}

	return nil, errors.New("Unsupported backing store " + store)
}
