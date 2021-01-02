package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
	"github.com/netwayfind/cp-scoring/test/model"
)

type backingStore interface {
	auditAnswerResultsInsert(results model.AuditAnswerResults) error
	auditCheckResultsInsert(results model.AuditCheckResults, teamID uint64, timestamp int64, source string) (uint64, error)
	hostTokenInsert(hostToken string, hostname string, timestamp int64, source string) error
	hostTokenSelectHostname(hostToken string) (string, error)
	hostTokenSelectTeamID(hostToken string) (uint64, error)
	scenarioDelete(id uint64) error
	scenarioInsert(scenario model.Scenario) (model.Scenario, error)
	scenarioSelect(id uint64) (model.Scenario, error)
	scenarioSelectAll() ([]model.ScenarioSummary, error)
	scenarioUpdate(id uint64, scenario model.Scenario) (model.Scenario, error)
	scenarioChecksSelectAll(id uint64) (map[string][]model.Action, error)
	scenarioChecksDelete(id uint64) error
	scenarioChecksUpdate(id uint64, hostnameChecks map[string][]model.Action) error
	scenarioAnswersSelectAll(id uint64) (map[string][]model.Answer, error)
	scenarioAnswersDelete(id uint64) error
	scenarioAnswersUpdate(id uint64, answersMap map[string][]model.Answer) error
	scoreboardSelectByScenarioID(scenarioID uint64) ([]model.ScenarioScore, error)
	scoreboardSelectScenarios() ([]model.ScenarioSummary, error)
	scoreboardUpdate(scenarioID uint64, teamID uint64, hostname string, score int, timestamp int64) error
	teamDelete(id uint64) error
	teamInsert(team model.Team) (model.Team, error)
	teamSelect(id uint64) (model.Team, error)
	teamSelectByKey(key string) (model.Team, error)
	teamSelectAll() ([]model.TeamSummary, error)
	teamUpdate(id uint64, team model.Team) (model.Team, error)
	teamHostTokenInsert(teamID uint64, hostToken string, timestamp int64) error
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
