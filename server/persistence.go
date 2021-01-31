package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
	"github.com/netwayfind/cp-scoring/model"
)

type backingStore interface {
	auditAnswerResultsInsert(results model.AuditAnswerResults) error
	auditAnswerResultsSelectHostnames(scenarioID uint64, teamID uint64) ([]string, error)
	auditAnswerResultsReport(scenarioID uint64, teamID uint64, hostname string) (model.Report, error)
	auditAnswerResultsReportTimeline(scenarioID uint64, teamID uint64, hostname string) ([]model.ReportTimeline, error)
	auditQueueDelete(ids uint64) error
	auditQueueInsert(entry model.AuditQueueEntry) error
	auditQueueSelectStatusReceived() ([]model.AuditQueueEntry, error)
	auditQueueUpdateStatusFailed(id uint64) error
	auditCheckResultsInsert(results model.AuditCheckResults, teamID uint64, timestamp int64, source string) (uint64, error)
	hostTokenInsert(hostToken string, hostname string, timestamp int64, source string) error
	hostTokenSelectHostname(hostToken string) (string, error)
	hostTokenSelectTeamID(hostToken string) (uint64, error)
	scenarioDelete(id uint64) error
	scenarioInsert(scenario model.Scenario) (model.Scenario, error)
	scenarioSelect(id uint64) (model.Scenario, error)
	scenarioSelectAll() ([]model.ScenarioSummary, error)
	scenarioUpdate(id uint64, scenario model.Scenario) (model.Scenario, error)
	scenarioHostsSelectAll(scenarioID uint64) (map[string]model.ScenarioHost, error)
	scenarioHostsSelectAnswers(scenarioID uint64, hostname string) ([]model.Answer, error)
	scenarioHostsSelectChecks(scenarioID uint64, hostname string) ([]model.Action, error)
	scenarioHostsSelectConfig(scenarioID uint64, hostname string) ([]model.Action, error)
	scenarioHostsDelete(scenarioID uint64) error
	scenarioHostsUpdate(scenarioID uint64, scenarioHosts map[string]model.ScenarioHost) error
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
	userDelete(id uint64) error
	userInsert(user model.User) (model.User, error)
	userSelect(id uint64) (model.User, error)
	userSelectByUsername(username string) (model.User, error)
	userSelectAll() ([]model.UserSummary, error)
	userUpdate(id uint64, user model.User) (model.User, error)
	userRolesDelete(id uint64) error
	userRolesSelect(id uint64) ([]model.Role, error)
	userRolesUpdate(id uint64, roles []model.Role) error
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
