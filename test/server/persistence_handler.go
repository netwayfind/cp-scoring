package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/netwayfind/cp-scoring/test/model"
)

type dbObj struct {
	dbConn *sql.DB
}

func (db dbObj) dbInit() {

	db.dbCreateTable("host_tokens", "CREATE TABLE IF NOT EXISTS host_tokens(host_token VARCHAR NOT NULL PRIMARY KEY, timestamp INTEGER NOT NULL, hostname VARCHAR NOT NULL, source VARCHAR NOT NULL)")
	db.dbCreateTable("teams", "CREATE TABLE IF NOT EXISTS teams(id BIGSERIAL PRIMARY KEY, name VARCHAR UNIQUE NOT NULL, poc VARCHAR NOT NULL, email VARCHAR NOT NULL, enabled BOOLEAN NOT NULL, key VARCHAR NOT NULL)")
	db.dbCreateTable("team_host_tokens", "CREATE TABLE IF NOT EXISTS team_host_tokens(team_id BIGSERIAL NOT NULL, host_token VARCHAR NOT NULL, timestamp INTEGER NOT NULL, FOREIGN KEY(team_id) REFERENCES teams(id), FOREIGN KEY(host_token) REFERENCES host_tokens(host_token))")
	db.dbCreateTable("scenarios", "CREATE TABLE IF NOT EXISTS scenarios(id BIGSERIAL PRIMARY KEY, name VARCHAR UNIQUE NOT NULL, description VARCHAR NOT NULL, enabled BOOLEAN NOT NULL)")
	db.dbCreateTable("scenario_hosts", "CREATE TABLE IF NOT EXISTS scenario_hosts(scenario_id BIGSERIAL NOT NULL, hostname VARCHAR NOT NULL, checks JSONB NOT NULL, answers JSONB NOT NULL, config JSONB NOT NULL, FOREIGN KEY(scenario_id) REFERENCES scenarios(id))")
	db.dbCreateTable("scoreboard", "CREATE TABLE IF NOT EXISTS scoreboard(scenario_id BIGSERIAL NOT NULL, team_id BIGSERIAL NOT NULL, hostname VARCHAR NOT NULL, score INTEGER NOT NULL, timestamp INTEGER NOT NULL, FOREIGN KEY(scenario_id) REFERENCES scenarios(id), FOREIGN KEY(team_id) REFERENCES teams(id))")
	db.dbCreateTable("audit_check_results", "CREATE TABLE IF NOT EXISTS audit_check_results(id BIGSERIAL NOT NULL PRIMARY KEY, scenario_id BIGSERIAL NOT NULL, team_id BIGSERIAL NOT NULL, host_token VARCHAR NOT NULL, timestamp_reported INTEGER NOT NULL, timestamp_received INTEGER NOT NULL, check_results JSONB NOT NULL, source VARCHAR NOT NULL, FOREIGN KEY(scenario_id) REFERENCES scenarios(id), FOREIGN KEY(team_id) REFERENCES teams(id), FOREIGN KEY(host_token) REFERENCES host_tokens(host_token))")
	db.dbCreateTable("audit_answer_results", "CREATE TABLE IF NOT EXISTS audit_answer_results(id BIGSERIAL NOT NULL PRIMARY KEY, scenario_id BIGSERIAL NOT NULL, team_id BIGSERIAL NOT NULL, host_token VARCHAR NOT NULL, timestamp INTEGER NOT NULL, audit_check_results_id BIGSERIAL NOT NULL, score INTEGER NOT NULL, answer_results JSONB NOT NULL, FOREIGN KEY(scenario_id) REFERENCES scenarios(id), FOREIGN KEY(team_id) REFERENCES teams(id), FOREIGN KEY(host_token) REFERENCES host_tokens(host_token), FOREIGN KEY(audit_check_results_id) REFERENCES audit_check_results(id))")

	log.Println("Finished setting up database")
}

func (db dbObj) dbClose() {
	db.dbConn.Close()
}

func (db dbObj) dbCreateTable(name string, stmtStr string) {
	stmt, err := db.dbConn.Prepare(stmtStr)
	if err != nil {
		log.Fatal("ERROR: cannot create table "+name+";", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("ERROR: cannot create table "+name+";", err)
	}
}

func (db dbObj) dbDelete(stmtStr string, args ...interface{}) error {
	stmt, err := db.dbConn.Prepare(stmtStr)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}

func (db dbObj) dbInsert(stmtStr string, args ...interface{}) (uint64, error) {
	stmt, err := db.dbConn.Prepare(stmtStr)
	if err != nil {
		return 0, err
	}

	if strings.Contains(stmtStr, "RETURNING") {
		var id uint64
		err = stmt.QueryRow(args...).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (db dbObj) dbUpdate(stmtStr string, args ...interface{}) error {
	stmt, err := db.dbConn.Prepare(stmtStr)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(model.ErrorDBUpdateNoChange)
	}

	return nil
}

func (db dbObj) auditAnswerResultsInsert(results model.AuditAnswerResults) error {
	b, err := json.Marshal(results.AnswerResults)
	if err != nil {
		return err
	}
	_, err = db.dbInsert("INSERT INTO audit_answer_results(scenario_id, team_id, host_token, timestamp, audit_check_results_id, score, answer_results) VALUES($1, $2, $3, $4, $5, $6, $7)",
		results.ScenarioID, results.TeamID, results.HostToken, results.Timestamp, results.CheckResultsID, results.Score, b)
	return err
}

func (db dbObj) auditAnswerResultsSelectHostnames(scenarioID uint64, teamID uint64) ([]string, error) {
	rows, err := db.dbConn.Query("SELECT h.hostname FROM audit_answer_results a JOIN host_tokens h ON a.host_token=h.host_token GROUP BY h.hostname ORDER BY h.hostname ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hostnames := make([]string, 0)
	for rows.Next() {
		var hostname string
		err = rows.Scan(&hostname)
		if err != nil {
			return nil, err
		}
		hostnames = append(hostnames, hostname)
	}

	return hostnames, nil
}

func (db dbObj) auditAnswerResultsReport(scenarioID uint64, teamID uint64, hostname string) (model.Report, error) {
	var report model.Report

	rows, err := db.dbConn.Query("SELECT a.timestamp, a.answer_results FROM audit_answer_results a JOIN host_tokens h ON a.host_token=h.host_token WHERE a.scenario_id=$1 AND a.team_id=$2 AND h.hostname=$3 ORDER BY a.timestamp DESC LIMIT 1", scenarioID, teamID, hostname)
	if err != nil {
		return report, err
	}
	defer rows.Close()

	for rows.Next() {
		var timestamp int64
		var answerResultsBs []byte
		err = rows.Scan(&timestamp, &answerResultsBs)
		if err != nil {
			return report, err
		}
		var answerResults []model.AnswerResult
		err = json.Unmarshal(answerResultsBs, &answerResults)
		if err != nil {
			return report, err
		}
		report = model.Report{
			Timestamp:     timestamp,
			AnswerResults: answerResults,
		}

		// only get first result
		break
	}

	return report, nil
}

func (db dbObj) auditAnswerResultsReportTimeline(scenarioID uint64, teamID uint64, hostname string) ([]model.ReportTimeline, error) {
	rows, err := db.dbConn.Query("SELECT a.host_token, a.timestamp, a.score FROM audit_answer_results a JOIN host_tokens h ON a.host_token=h.host_token WHERE a.scenario_id=$1 AND a.team_id=$2 AND h.hostname=$3 ORDER BY h.hostname, a.timestamp ASC", scenarioID, teamID, hostname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hostTokenMap := make(map[string]int)
	timeline := make([]model.ReportTimeline, 0)
	for rows.Next() {
		var hostToken string
		var timestamp int64
		var score int
		err = rows.Scan(&hostToken, &timestamp, &score)
		if err != nil {
			return nil, err
		}
		hostIndex, present := hostTokenMap[hostToken]
		if !present {
			hostIndex = len(hostTokenMap)
			hostTokenMap[hostToken] = hostIndex
			hostTimeline := model.ReportTimeline{
				Timestamps: make([]int64, 0),
				Scores:     make([]int, 0),
			}
			timeline = append(timeline, hostTimeline)
		}
		currentTimeline := timeline[hostIndex]
		currentTimeline.Timestamps = append(currentTimeline.Timestamps, timestamp)
		currentTimeline.Scores = append(currentTimeline.Scores, score)
		timeline[hostIndex] = currentTimeline
	}

	return timeline, nil
}

func (db dbObj) auditCheckResultsInsert(results model.AuditCheckResults, teamID uint64, timestampProcessed int64, source string) (uint64, error) {
	b, err := json.Marshal(results.CheckResults)
	if err != nil {
		return 0, err
	}
	return db.dbInsert("INSERT INTO audit_check_results(scenario_id, team_id, host_token, timestamp_reported, timestamp_received, check_results, source) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		results.ScenarioID, teamID, results.HostToken, results.Timestamp, timestampProcessed, b, source)
}

func (db dbObj) hostTokenInsert(hostToken string, hostname string, timestamp int64, source string) error {
	_, err := db.dbInsert("INSERT INTO host_tokens(host_token, hostname, timestamp, source) VALUES($1, $2, $3, $4)", hostToken, hostname, timestamp, source)
	return err
}

func (db dbObj) hostTokenSelectHostname(hostToken string) (string, error) {
	var hostname string

	rows, err := db.dbConn.Query("SELECT hostname FROM host_tokens WHERE host_token=$1", hostToken)
	if err != nil {
		return hostname, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&hostname)
		if err != nil {
			return hostname, err
		}
		// only get first result
		break
	}

	return hostname, nil
}

func (db dbObj) hostTokenSelectTeamID(hostToken string) (uint64, error) {
	var teamID uint64

	rows, err := db.dbConn.Query("SELECT team_id FROM team_host_tokens WHERE host_token=$1", hostToken)
	if err != nil {
		return teamID, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&teamID)
		if err != nil {
			return teamID, err
		}
		// only get first result
		break
	}

	return teamID, nil
}

func (db dbObj) scenarioDelete(id uint64) error {
	// TODO: transaction
	err := db.scenarioHostsDelete(id)
	if err != nil {
		return err
	}
	return db.dbDelete("DELETE FROM scenarios where id=$1", id)
}

func (db dbObj) scenarioInsert(scenario model.Scenario) (model.Scenario, error) {
	id, err := db.dbInsert("INSERT INTO scenarios(name, description, enabled) VALUES($1, $2, $3) RETURNING id", scenario.Name, scenario.Description, scenario.Enabled)
	if err != nil {
		return model.Scenario{}, err
	}

	return db.scenarioSelect(id)
}

func (db dbObj) scenarioSelect(id uint64) (model.Scenario, error) {
	var scenario model.Scenario

	rows, err := db.dbConn.Query("SELECT id, name, description, enabled FROM scenarios WHERE id=$1", id)
	if err != nil {
		return scenario, err
	}
	defer rows.Close()

	for rows.Next() {
		scenario = model.Scenario{}
		err = rows.Scan(&scenario.ID, &scenario.Name, &scenario.Description, &scenario.Enabled)
		if err != nil {
			return scenario, err
		}
		// only get first result
		break
	}

	return scenario, nil
}

func (db dbObj) scenarioSelectAll() ([]model.ScenarioSummary, error) {
	rows, err := db.dbConn.Query("SELECT id, name, enabled FROM scenarios ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries := make([]model.ScenarioSummary, 0)

	for rows.Next() {
		summary := model.ScenarioSummary{}
		err = rows.Scan(&summary.ID, &summary.Name, &summary.Enabled)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (db dbObj) scenarioUpdate(id uint64, scenario model.Scenario) (model.Scenario, error) {
	enabled := 1
	if !scenario.Enabled {
		enabled = 0
	}

	err := db.dbUpdate("UPDATE scenarios SET name=$1, description=$2, enabled=$3 WHERE id=$4", scenario.Name, scenario.Description, enabled, id)
	if err != nil {
		return model.Scenario{}, err
	}

	return db.scenarioSelect(id)
}

func (db dbObj) scenarioHostsSelectAll(scenarioID uint64) (map[string]model.ScenarioHost, error) {
	rows, err := db.dbConn.Query("SELECT hostname, checks, answers, config FROM scenario_hosts WHERE scenario_id=$1", scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hostMap := make(map[string]model.ScenarioHost)
	for rows.Next() {
		var hostname string
		var checks []model.Action
		var checksBs []byte
		var answers []model.Answer
		var answersBs []byte
		var config []model.Action
		var configBs []byte
		err = rows.Scan(&hostname, &checksBs, &answersBs, &configBs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(checksBs, &checks)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(answersBs, &answers)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(configBs, &config)
		if err != nil {
			return nil, err
		}

		hostMap[hostname] = model.ScenarioHost{
			Checks:  checks,
			Answers: answers,
			Config:  config,
		}
	}

	return hostMap, nil
}

func (db dbObj) scenarioHostsSelectAnswers(scenarioID uint64, hostname string) ([]model.Answer, error) {
	rows, err := db.dbConn.Query("SELECT answers FROM scenario_hosts WHERE scenario_id=$1 AND hostname=$2", scenarioID, hostname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []model.Answer
	var answersBs []byte
	for rows.Next() {
		err = rows.Scan(&answersBs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(answersBs, &answers)
		if err != nil {
			return nil, err
		}
		break
	}

	return answers, nil
}

func (db dbObj) scenarioHostsSelectChecks(scenarioID uint64, hostname string) ([]model.Action, error) {
	rows, err := db.dbConn.Query("SELECT checks FROM scenario_hosts WHERE scenario_id=$1 AND hostname=$2", scenarioID, hostname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checks []model.Action
	var checksBs []byte
	for rows.Next() {
		err = rows.Scan(&checksBs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(checksBs, &checks)
		if err != nil {
			return nil, err
		}
		break
	}

	return checks, nil
}

func (db dbObj) scenarioHostsSelectConfig(scenarioID uint64, hostname string) ([]model.Action, error) {
	rows, err := db.dbConn.Query("SELECT config FROM scenario_hosts WHERE scenario_id=$1 AND hostname=$2", scenarioID, hostname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var config []model.Action
	var configBs []byte
	for rows.Next() {
		err = rows.Scan(&configBs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(configBs, &config)
		if err != nil {
			return nil, err
		}
		break
	}

	return config, nil
}

func (db dbObj) scenarioHostsDelete(scenarioID uint64) error {
	return db.dbDelete("DELETE FROM scenario_hosts WHERE scenario_id=$1", scenarioID)
}

func (db dbObj) scenarioHostsUpdate(scenarioID uint64, scenarioHosts map[string]model.ScenarioHost) error {
	// TODO: transaction
	err := db.scenarioHostsDelete(scenarioID)
	if err != nil {
		return err
	}

	for hostname, scenarioHost := range scenarioHosts {
		checksBs, err := json.Marshal(scenarioHost.Checks)
		if err != nil {
			return err
		}
		answersBs, err := json.Marshal(scenarioHost.Answers)
		if err != nil {
			return err
		}
		configBs, err := json.Marshal(scenarioHost.Config)
		if err != nil {
			return err
		}
		_, err = db.dbInsert("INSERT INTO scenario_hosts(scenario_id, hostname, checks, answers, config) VALUES ($1, $2, $3, $4, $5)", scenarioID, hostname, checksBs, answersBs, configBs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db dbObj) scoreboardSelectByScenarioID(scenarioID uint64) ([]model.ScenarioScore, error) {
	rows, err := db.dbConn.Query("SELECT t.name, s.hostname, s.score, s.timestamp FROM scoreboard s JOIN teams t ON s.team_id=t.id WHERE s.scenario_id=$1", scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scoreboard := make([]model.ScenarioScore, 0)
	for rows.Next() {
		var teamName string
		var hostname string
		var score int
		var timestamp int64
		err = rows.Scan(&teamName, &hostname, &score, &timestamp)
		if err != nil {
			return nil, err
		}
		scenarioScore := model.ScenarioScore{
			TeamName:  teamName,
			Hostname:  hostname,
			Score:     score,
			Timestamp: timestamp,
		}
		scoreboard = append(scoreboard, scenarioScore)
	}

	return scoreboard, nil
}

func (db dbObj) scoreboardSelectScenarios() ([]model.ScenarioSummary, error) {
	rows, err := db.dbConn.Query("SELECT id, name FROM scenarios WHERE enabled=true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scenarios := make([]model.ScenarioSummary, 0)
	for rows.Next() {
		var id uint64
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		scenario := model.ScenarioSummary{
			ID:      id,
			Name:    name,
			Enabled: true,
		}
		scenarios = append(scenarios, scenario)
	}

	return scenarios, nil
}

func (db dbObj) scoreboardUpdate(scenarioID uint64, teamID uint64, hostname string, score int, timestamp int64) error {
	// TODO: transaction
	err := db.dbDelete("DELETE FROM scoreboard WHERE scenario_id=$1 AND team_id=$2 AND hostname=$3", scenarioID, teamID, hostname)
	if err != nil {
		return err
	}

	_, err = db.dbInsert("INSERT INTO scoreboard(scenario_id, team_id, hostname, score, timestamp) VALUES($1, $2, $3, $4, $5)", scenarioID, teamID, hostname, score, timestamp)
	if err != nil {
		return err
	}
	return nil
}

func (db dbObj) teamDelete(id uint64) error {
	return db.dbDelete("DELETE FROM teams where id=$1", id)
}

func (db dbObj) teamInsert(team model.Team) (model.Team, error) {
	key := team.Key
	if len(key) == 0 {
		key = randHexStr(8)
	}
	enabled := 1
	if !team.Enabled {
		enabled = 0
	}
	id, err := db.dbInsert("INSERT INTO teams(name, poc, email, enabled, key) VALUES($1, $2, $3, $4, $5) RETURNING id", team.Name, team.POC, team.Email, enabled, key)
	if err != nil {
		return model.Team{}, err
	}

	return db.teamSelect(id)
}

func (db dbObj) teamSelect(id uint64) (model.Team, error) {
	var team model.Team

	rows, err := db.dbConn.Query("SELECT id, name, poc, email, enabled, key FROM teams WHERE id=$1", id)
	if err != nil {
		return team, err
	}
	defer rows.Close()

	for rows.Next() {
		team = model.Team{}
		err = rows.Scan(&team.ID, &team.Name, &team.POC, &team.Email, &team.Enabled, &team.Key)
		if err != nil {
			return team, err
		}
		// only get first result
		break
	}

	return team, nil
}

func (db dbObj) teamSelectByKey(key string) (model.Team, error) {
	var team model.Team

	rows, err := db.dbConn.Query("SELECT id, name, poc, email, enabled, key FROM teams WHERE key=$1", key)
	if err != nil {
		return team, err
	}
	defer rows.Close()

	for rows.Next() {
		team = model.Team{}
		err = rows.Scan(&team.ID, &team.Name, &team.POC, &team.Email, &team.Enabled, &team.Key)
		if err != nil {
			return team, err
		}
		// only get first result
		break
	}

	return team, nil
}

func (db dbObj) teamSelectAll() ([]model.TeamSummary, error) {
	rows, err := db.dbConn.Query("SELECT id, name, enabled FROM teams ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries := make([]model.TeamSummary, 0)

	for rows.Next() {
		summary := model.TeamSummary{}
		err = rows.Scan(&summary.ID, &summary.Name, &summary.Enabled)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (db dbObj) teamUpdate(id uint64, team model.Team) (model.Team, error) {
	enabled := 1
	if !team.Enabled {
		enabled = 0
	}

	err := db.dbUpdate("UPDATE teams SET name=$1, poc=$2, email=$3, enabled=$4, key=$5 WHERE id=$6", team.Name, team.POC, team.Email, enabled, team.Key, id)
	if err != nil {
		return model.Team{}, err
	}

	return db.teamSelect(id)
}

func (db dbObj) teamHostTokenInsert(teamID uint64, hostToken string, timestamp int64) error {
	_, err := db.dbInsert("INSERT INTO team_host_tokens(team_id, host_token, timestamp) VALUES($1, $2, $3)", teamID, hostToken, timestamp)
	return err
}
