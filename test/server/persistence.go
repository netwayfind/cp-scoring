package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
	"github.com/netwayfind/cp-scoring/test/model"
)

type backingStore interface {
	scenarioInsert(scenario model.Scenario) (uint64, error)
	scenarioSelect(id uint64) (model.Scenario, error)
	scenarioSelectAll() ([]model.Scenario, error)
	scenarioUpdate(id uint64, scenario model.Scenario) error
	scenarioChecksSelectAll(id uint64) (map[string][]model.Action, error)
	scenarioChecksUpdate(id uint64, hostnameChecks map[string][]model.Action) error
	teamInsert(team model.Team) (uint64, error)
	teamSelect(id uint64) (model.Team, error)
	teamSelectAll() ([]model.Team, error)
	teamUpdate(id uint64, team model.Team) error
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
