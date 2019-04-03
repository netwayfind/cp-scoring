package main

import (
	"errors"

	"github.com/sumwonyuno/cp-scoring/processing"

	"github.com/sumwonyuno/cp-scoring/model"
)

type backingStore interface {
	InsertState(timestamp int64, source string, hostToken string, state model.State) error
	SelectStates(hostToken string, timeStart int64, timeEnd int64) ([]model.State, error)
	SelectStateDiffs(hostToken string, timeStart int64, timeEnd int64) ([]processing.Change, error)
	SelectAdmins() ([]string, error)
	IsAdmin(username string) (bool, error)
	SelectAdminPasswordHash(username string) (string, error)
	InsertAdmin(username string, passwordHash string) error
	UpdateAdmin(username string, passwordHash string) error
	DeleteAdmin(username string) error
	SelectHosts() ([]model.Host, error)
	SelectHost(hostID uint64) (model.Host, error)
	SelectHostIDForHostname(hostname string) (uint64, error)
	InsertHost(host model.Host) (uint64, error)
	UpdateHost(hostID uint64, host model.Host) error
	DeleteHost(hostID uint64) error
	SelectTeams() ([]model.TeamSummary, error)
	SelectTeam(teamID uint64) (model.Team, error)
	SelectTeamIDFromHostToken(hostToken string) (uint64, error)
	SelectTeamIDForKey(teamKey string) (uint64, error)
	InsertTeam(team model.Team) (uint64, error)
	UpdateTeam(teamID uint64, team model.Team) error
	DeleteTeam(teamID uint64) error
	SelectTemplates() ([]model.Template, error)
	SelectTemplatesForHostname(scenarioID uint64, hostname string) ([]model.Template, error)
	SelectTemplate(templateID uint64) (model.Template, error)
	InsertTemplate(template model.Template) (uint64, error)
	UpdateTemplate(templateID uint64, template model.Template) error
	DeleteTemplate(templateID uint64) error
	SelectScenarios(onlyEnabled bool) ([]model.ScenarioSummary, error)
	SelectScenariosForHostname(hostname string) ([]uint64, error)
	SelectScenario(scenarioID uint64) (model.Scenario, error)
	InsertScenario(scenario model.Scenario) (uint64, error)
	UpdateScenario(scenarioID uint64, scenario model.Scenario) error
	DeleteScenario(scenarioID uint64) error
	SelectLatestScenarioScores(scenarioID uint64) ([]model.TeamScore, error)
	InsertScenarioReport(scenarioID uint64, hostToken string, report model.Report) error
	SelectScenarioReports(scenarioID uint64, hostToken string, timeStart int64, timeEnd int64) ([]model.Report, error)
	SelectScenarioReportDiffs(scenarioID uint64, hostToken string, timeStart int64, timeEnd int64) ([]processing.Change, error)
	InsertScenarioScore(score model.ScenarioHostScore) error
	SelectScenarioTimeline(scenarioID uint64, hostToken string) (model.ScenarioTimeline, error)
	SelectLatestScenarioReport(scenarioID uint64, hostToken string) (model.Report, error)
	SelectTeamScenarioHosts(teamID uint64) ([]model.ScenarioHosts, error)
	InsertHostToken(hostToken string, timestamp int64, hostname string, source string) error
	InsertTeamHostToken(teamID uint64, hostToken string, timestamp int64) error
	SelectHostTokens(teamID uint64, hostname string) ([]string, error)
}

func getBackingStore(store string, args ...string) (backingStore, error) {
	if store == "postgres" {
		db := dbObj{}
		dbConn, err := newPostgresDBConn(args)
		if err != nil {
			return nil, err
		}
		db.dbConn = dbConn
		db.dbInit()
		return db, nil
	}
	return nil, errors.New("Unsupported backing store " + store)
}
