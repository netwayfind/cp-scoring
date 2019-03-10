package main

import (
	"errors"

	"github.com/sumwonyuno/cp-scoring/model"
)

type backingStore interface {
	InsertState(state string) error
	SelectAdmins() ([]string, error)
	IsAdmin(username string) (bool, error)
	SelectAdminPasswordHash(username string) (string, error)
	InsertAdmin(username string, passwordHash string) error
	UpdateAdmin(username string, passwordHash string) error
	DeleteAdmin(username string) error
	SelectHosts() ([]model.Host, error)
	SelectHost(hostID int64) (model.Host, error)
	SelectHostIDForHostname(hostname string) (int64, error)
	InsertHost(host model.Host) (int64, error)
	UpdateHost(hostID int64, host model.Host) error
	DeleteHost(hostID int64) error
	SelectTeams() ([]model.TeamSummary, error)
	SelectTeam(teamID int64) (model.Team, error)
	SelectTeamIDFromHostToken(hostToken string) (int64, error)
	SelectTeamIDForKey(teamKey string) (int64, error)
	InsertTeam(team model.Team) (int64, error)
	UpdateTeam(teamID int64, team model.Team) error
	DeleteTeam(teamID int64) error
	SelectTemplates() ([]model.Template, error)
	SelectTemplatesForHostname(scenarioID int64, hostname string) ([]model.Template, error)
	SelectTemplate(templateID int64) (model.Template, error)
	InsertTemplate(template model.Template) (int64, error)
	UpdateTemplate(templateID int64, template model.Template) error
	DeleteTemplate(templateID int64) error
	SelectScenarios(onlyEnabled bool) ([]model.ScenarioSummary, error)
	SelectScenariosForHostname(hostname string) ([]int64, error)
	SelectScenario(scenarioID int64) (model.Scenario, error)
	InsertScenario(scenario model.Scenario) (int64, error)
	UpdateScenario(scenarioID int64, scenario model.Scenario) error
	DeleteScenario(scenarioID int64) error
	SelectLatestScenarioScores(scenarioID int64) ([]model.TeamScore, error)
	InsertScenarioReport(scenarioID int64, hostToken string, report model.Report) error
	InsertScenarioScore(score model.ScenarioHostScore) error
	SelectScenarioTimeline(scenarioID int64, hostToken string) (model.ScenarioTimeline, error)
	SelectLatestScenarioReport(scenarioID int64, hostToken string) (model.Report, error)
	SelectTeamScenarioHosts(teamID int64) ([]model.ScenarioHosts, error)
	InsertHostToken(hostToken string, timestamp int64, hostname string, source string) error
	InsertTeamHostToken(teamID int64, hostname string, hostToken string, timestamp int64) error
	SelectHostTokens(teamID int64, hostname string) ([]string, error)
}

func getBackingStore(store string, args ...string) (backingStore, error) {
	if store == "sqlite" {
		db := dbObj{}
		dbConn, err := newSQLiteDBConn(args)
		if err != nil {
			return nil, err
		}
		db.dbConn = dbConn
		db.dbInit()
		return db, nil
	} else if store == "postgres" {
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
