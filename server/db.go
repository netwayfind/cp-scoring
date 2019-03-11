package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/sumwonyuno/cp-scoring/model"
)

type dbObj struct {
	dbConn *sql.DB
}

func newPostgresDBConn(args []string) (*sql.DB, error) {
	// must have first argument as URL
	if len(args) < 1 {
		return nil, errors.New("ERROR: URL required")
	}
	connStr := args[0]
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to database")
	return dbConn, nil
}

func (db dbObj) dbInit() {
	db.createTable("states", "CREATE TABLE IF NOT EXISTS states(id BIGSERIAL PRIMARY KEY, timestamp INTEGER NOT NULL, source VARCHAR NOT NULL, host_token VARCHAR NOT NULL, state JSONB)")
	db.createTable("admins", "CREATE TABLE IF NOT EXISTS admins(username VARCHAR PRIMARY KEY, password_hash VARCHAR NOT NULL)")
	db.createTable("teams", "CREATE TABLE IF NOT EXISTS teams(id BIGSERIAL PRIMARY KEY, name VARCHAR NOT NULL, poc VARCHAR NOT NULL, email VARCHAR NOT NULL, enabled BOOLEAN NOT NULL, key VARCHAR NOT NULL)")
	db.createTable("templates", "CREATE TABLE IF NOT EXISTS templates(id BIGSERIAL PRIMARY KEY, name VARCHAR NOT NULL, state BYTEA NOT NULL)")
	db.createTable("hosts", "CREATE TABLE IF NOT EXISTS hosts(id BIGSERIAL PRIMARY KEY, hostname VARCHAR NOT NULL, os VARCHAR NOT NULL)")
	db.createTable("host_tokens", "CREATE TABLE IF NOT EXISTS host_tokens(host_token VARCHAR UNIQUE NOT NULL PRIMARY KEY, timestamp INTEGER NOT NULL, hostname VARCHAR UNIQUE NOT NULL, source VARCHAR NOT NULL)")
	db.createTable("team_host_tokens", "CREATE TABLE IF NOT EXISTS team_host_tokens(team_id INTEGER NOT NULL, hostname VARCHAR NOT NULL, host_token VARCHAR NOT NULL, timestamp INTEGER NOT NULL, FOREIGN KEY(team_id) REFERENCES teams(id), FOREIGN KEY(hostname) REFERENCES host_tokens(hostname), FOREIGN KEY(host_token) REFERENCES host_tokens(host_token))")
	db.createTable("scenarios", "CREATE TABLE IF NOT EXISTS scenarios(id BIGSERIAL PRIMARY KEY, name VARCHAR NOT NULL, description VARCHAR NOT NULL, enabled BOOLEAN NOT NULL)")
	db.createTable("hosts_templates", "CREATE TABLE IF NOT EXISTS hosts_templates(scenario_id INTEGER NOT NULL, host_id INTEGER NOT NULL, template_id INTEGER NOT NULL, FOREIGN KEY(scenario_id) REFERENCES scenarios(id), FOREIGN KEY(template_id) REFERENCES templates(id), FOREIGN KEY(host_id) REFERENCES hosts(id))")
	db.createTable("scores", "CREATE TABLE IF NOT EXISTS scores(scenario_id INTEGER NOT NULL, host_token VARCHAR NOT NULL, timestamp INTEGER NOT NULL, score INTEGER NOT NULL, FOREIGN KEY(scenario_id) REFERENCES scenarios(id), FOREIGN KEY(host_token) REFERENCES host_tokens(host_token))")
	db.createTable("reports", "CREATE TABLE IF NOT EXISTS reports(scenario_id INTEGER NOT NULL, host_token VARCHAR NOT NULL, timestamp INTEGER NOT NULL, report BYTEA NOT NULL, FOREIGN KEY(scenario_id) REFERENCES scenarios(id), FOREIGN KEY(host_token) REFERENCES host_tokens(host_token))")

	log.Println("Finished setting up database")
}

func (db dbObj) dbClose() {
	db.dbConn.Close()
}

func (db dbObj) createTable(name string, stmtStr string) {
	stmt, err := db.dbConn.Prepare(stmtStr)
	if err != nil {
		log.Fatal("ERROR: cannot create table "+name+";", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("ERROR: cannot create table "+name+";", err)
	}
}

func (db dbObj) dbInsert(stmtStr string, args ...interface{}) (int64, error) {
	stmt, err := db.dbConn.Prepare(stmtStr)
	if err != nil {
		return -1, err
	}

	if strings.Contains(stmtStr, "RETURNING") {
		var id int64
		err = stmt.QueryRow(args...).Scan(&id)
		if err != nil {
			return -1, err
		}
		return id, nil
	} else {
		_, err := stmt.Exec(args...)
		if err != nil {
			return -1, err
		}
		return -1, nil
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

func (db dbObj) dbUpdate(stmtStr string, args ...interface{}) error {
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

func (db dbObj) InsertState(timestamp int64, source string, hostToken string, state []byte) error {
	_, err := db.dbInsert("INSERT INTO states(timestamp, source, host_token, state) VALUES($1, $2, $3, $4) RETURNING id", timestamp, source, hostToken, state)
	return err
}

func (db dbObj) SelectAdmins() ([]string, error) {
	var admins []string
	rows, err := db.dbConn.Query("SELECT DISTINCT username FROM admins")
	if err != nil {
		return admins, err
	}
	defer rows.Close()

	admins = make([]string, 0)
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return admins, err
		}
		admins = append(admins, username)
	}

	return admins, nil
}

func (db dbObj) IsAdmin(username string) (bool, error) {
	rows, err := db.dbConn.Query("SELECT username FROM admins WHERE username=$1", username)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var count int
	var u string
	for rows.Next() {
		err := rows.Scan(&u)
		if err != nil {
			return false, err
		}
		count++
	}

	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (db dbObj) SelectAdminPasswordHash(username string) (string, error) {
	rows, err := db.dbConn.Query("SELECT password_hash FROM admins WHERE username=$1", username)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var passwordHash string
	for rows.Next() {
		err := rows.Scan(&passwordHash)
		if err != nil {
			return "", err
		}
	}

	return passwordHash, nil
}

func (db dbObj) DeleteAdmin(username string) error {
	return db.dbDelete("DELETE FROM admins where username=$1", username)
}

func (db dbObj) UpdateAdmin(username string, passwordHash string) error {
	return db.dbUpdate("UPDATE admins SET password_hash=$1 WHERE username=$2", passwordHash, username)
}

func (db dbObj) InsertAdmin(username string, passwordHash string) error {
	_, err := db.dbInsert("INSERT INTO admins(username, password_hash) VALUES($1, $2)", username, passwordHash)
	return err
}

func (db dbObj) SelectTeams() ([]model.TeamSummary, error) {
	rows, err := db.dbConn.Query("SELECT id, name FROM teams ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int64
	var name string
	teams := make([]model.TeamSummary, 0)

	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		var entry model.TeamSummary
		entry.ID = id
		entry.Name = name
		teams = append(teams, entry)
	}

	return teams, nil
}

func (db dbObj) SelectHostIDForHostname(hostname string) (int64, error) {
	var id int64 = -1

	rows, err := db.dbConn.Query("SELECT id FROM hosts WHERE hostname=$1", hostname)
	if err != nil {
		return id, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return id, err
		}
		// only get first result
		break
	}

	// did not find any
	if id == -1 {
		return id, &errorStr{hostname + " hostname not found"}
	}

	return id, nil
}

type errorStr struct {
	error string
}

func (e *errorStr) Error() string {
	return e.error
}

func (db dbObj) SelectTeamIDForKey(key string) (int64, error) {
	var id int64 = -1
	key = strings.TrimSpace(key)

	rows, err := db.dbConn.Query("SELECT id FROM teams WHERE key=$1 AND enabled=TRUE", key)
	if err != nil {
		return id, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return id, err
		}
		// only get first result
		break
	}

	// did not find any
	if id == -1 {
		return id, &errorStr{key + " key not found"}
	}

	return id, err
}

func (db dbObj) SelectTeam(id int64) (model.Team, error) {
	var team model.Team

	rows, err := db.dbConn.Query("SELECT name, poc, email, enabled, key FROM teams where id=$1", id)
	if err != nil {
		return team, err
	}
	defer rows.Close()

	var name string
	var poc string
	var email string
	var enabled bool
	var key string
	for rows.Next() {
		err := rows.Scan(&name, &poc, &email, &enabled, &key)
		if err != nil {
			return team, err
		}
		team.ID = id
		team.Name = name
		team.POC = poc
		team.Email = email
		team.Enabled = enabled
		team.Key = key
		// only get first result
		break
	}

	return team, nil
}

func (db dbObj) DeleteTeam(id int64) error {
	return db.dbDelete("DELETE FROM teams where id=$1", id)
}

func (db dbObj) InsertTeam(team model.Team) (int64, error) {
	return db.dbInsert("INSERT INTO teams(name, poc, email, enabled, key) VALUES($1, $2, $3, $4, $5) RETURNING id", team.Name, team.POC, team.Email, team.Enabled, team.Key)
}

func (db dbObj) UpdateTeam(id int64, team model.Team) error {
	return db.dbUpdate("UPDATE teams SET name=$1, poc=$2, email=$3, enabled=$4, key=$5 WHERE id=$6", team.Name, team.POC, team.Email, team.Enabled, team.Key, id)
}

func (db dbObj) SelectTemplates() ([]model.Template, error) {
	rows, err := db.dbConn.Query("SELECT id, name, state FROM templates ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int64
	var name string
	var stateBytes []byte
	templates := make([]model.Template, 0)

	for rows.Next() {
		err = rows.Scan(&id, &name, &stateBytes)
		if err != nil {
			return nil, err
		}
		var state model.State
		err = json.Unmarshal(stateBytes, &state)
		if err != nil {
			continue
		}
		var template model.Template
		template.ID = id
		template.Name = name
		template.State = state
		templates = append(templates, template)
	}

	return templates, nil
}

func (db dbObj) SelectTemplate(id int64) (model.Template, error) {
	var template model.Template
	var name string
	var state model.State
	var stateBytes []byte

	rows, err := db.dbConn.Query("SELECT name, state FROM templates where id=$1", id)
	if err != nil {
		return template, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &stateBytes)
		if err != nil {
			return template, err
		}
		err = json.Unmarshal(stateBytes, &state)
		if err != nil {
			return template, err
		}
		// only get first result
		template.ID = id
		template.Name = name
		template.State = state
		break
	}

	return template, nil
}

func (db dbObj) DeleteTemplate(id int64) error {
	return db.dbDelete("DELETE FROM templates where id=$1", id)
}

func (db dbObj) SelectScenariosForHostname(hostname string) ([]int64, error) {
	rows, err := db.dbConn.Query("SELECT scenarios.id FROM hosts, hosts_templates, scenarios WHERE hosts.hostname=$1 AND hosts_templates.host_id=hosts.id AND hosts_templates.scenario_id=scenarios.id AND scenarios.enabled=TRUE", hostname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scenarioIDs := make([]int64, 0)
	var id int64
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		scenarioIDs = append(scenarioIDs, id)
	}
	return scenarioIDs, nil
}

func (db dbObj) SelectTemplatesForHostname(scenarioID int64, hostname string) ([]model.Template, error) {
	rows, err := db.dbConn.Query("SELECT templates.id, templates.name, templates.state FROM templates, hosts, hosts_templates WHERE hosts.hostname=$1 AND hosts_templates.scenario_id=$2 AND hosts_templates.host_id=hosts.id AND hosts_templates.template_id=templates.id", hostname, scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := make([]model.Template, 0)
	var id int64
	var name string
	var stateBytes []byte
	for rows.Next() {
		err := rows.Scan(&id, &name, &stateBytes)
		if err != nil {
			return nil, err
		}
		var state model.State
		err = json.Unmarshal(stateBytes, &state)
		if err != nil {
			return nil, err
		}
		var template model.Template
		template.ID = id
		template.Name = name
		template.State = state
		templates = append(templates, template)
	}
	return templates, nil
}

func (db dbObj) InsertTemplate(template model.Template) (int64, error) {
	b, err := json.Marshal(template.State)
	if err != nil {
		return -1, err
	}
	return db.dbInsert("INSERT INTO templates(name, state) VALUES($1, $2) RETURNING id", template.Name, b)
}

func (db dbObj) UpdateTemplate(id int64, template model.Template) error {
	b, err := json.Marshal(template.State)
	if err != nil {
		return err
	}
	return db.dbUpdate("UPDATE templates SET name=$1, state=$2 WHERE id=$3", template.Name, b, id)
}

func (db dbObj) SelectHosts() ([]model.Host, error) {
	rows, err := db.dbConn.Query("SELECT id, hostname, os FROM hosts ORDER BY hostname ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int64
	var hostname string
	var os string
	hosts := make([]model.Host, 0)

	for rows.Next() {
		err = rows.Scan(&id, &hostname, &os)
		if err != nil {
			return nil, err
		}
		var host model.Host
		host.ID = id
		host.Hostname = hostname
		host.OS = os
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func (db dbObj) SelectHost(id int64) (model.Host, error) {
	var host model.Host

	rows, err := db.dbConn.Query("SELECT hostname, os FROM hosts where id=$1", id)
	if err != nil {
		return host, err
	}
	defer rows.Close()

	var hostname string
	var os string
	count := 0
	for rows.Next() {
		err := rows.Scan(&hostname, &os)
		if err != nil {
			return host, err
		}
		// only get first result
		host.ID = id
		host.Hostname = hostname
		host.OS = os
		count++
		break
	}

	return host, nil
}

func (db dbObj) DeleteHost(id int64) error {
	return db.dbDelete("DELETE FROM hosts where id=$1", id)
}

func (db dbObj) InsertHost(host model.Host) (int64, error) {
	return db.dbInsert("INSERT INTO hosts(hostname, os) VALUES($1, $2) RETURNING id", host.Hostname, host.OS)
}

func (db dbObj) UpdateHost(id int64, host model.Host) error {
	return db.dbUpdate("UPDATE hosts SET hostname=$1,os=$2 WHERE id=$3", host.Hostname, host.OS, id)
}

func (db dbObj) SelectScenarioHostTemplates(scenarioID int64) (map[int64][]int64, error) {
	rows, err := db.dbConn.Query("SELECT host_id, template_id FROM hosts_templates WHERE scenario_id=$1", scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templateID int64
	var hostID int64
	hostTemplates := make(map[int64][]int64)

	for rows.Next() {
		err = rows.Scan(&hostID, &templateID)
		if err != nil {
			return nil, err
		}
		entry := hostTemplates[hostID]
		if entry == nil {
			entry = make([]int64, 1)
			entry[0] = templateID
			hostTemplates[hostID] = entry
		} else {
			hostTemplates[hostID] = append(entry, templateID)
		}
	}

	return hostTemplates, nil
}

func (db dbObj) SelectScenarios(onlyEnabled bool) ([]model.ScenarioSummary, error) {
	var q strings.Builder

	q.WriteString("SELECT id, name FROM scenarios")
	if onlyEnabled {
		q.WriteString(" WHERE enabled=TRUE")
	}
	q.WriteString(" ORDER BY name ASC")
	rows, err := db.dbConn.Query(q.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int64
	var name string
	scenarios := make([]model.ScenarioSummary, 0)

	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		var scenario model.ScenarioSummary
		scenario.ID = id
		scenario.Name = name
		scenarios = append(scenarios, scenario)
	}

	return scenarios, nil
}

func (db dbObj) SelectScenario(id int64) (model.Scenario, error) {
	var scenario model.Scenario

	stmt, err := db.dbConn.Prepare("SELECT name, description, enabled FROM scenarios where id=$1")
	if err != nil {
		return scenario, err
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return scenario, err
	}
	defer rows.Close()

	var name string
	var description string
	var enabled bool
	for rows.Next() {
		err := rows.Scan(&name, &description, &enabled)
		if err != nil {
			return scenario, err
		}
		hostTemplates, err := db.SelectScenarioHostTemplates(id)
		if err != nil {
			return scenario, err
		}
		// only get first result
		scenario.ID = id
		scenario.Name = name
		scenario.Description = description
		scenario.Enabled = enabled
		scenario.HostTemplates = hostTemplates
		break
	}

	return scenario, nil
}

func (db dbObj) DeleteScenarioHostTemplates(scenarioID int64) error {
	return db.dbDelete("DELETE FROM hosts_templates WHERE scenario_id=$1", scenarioID)
}

func (db dbObj) DeleteScenario(id int64) error {
	err := db.dbDelete("DELETE FROM scenarios WHERE id=$1", id)
	if err != nil {
		return err
	}
	return db.DeleteScenarioHostTemplates(id)
}

func (db dbObj) InsertScenarioHostTemplates(id int64, scenario model.Scenario) error {
	for hostID, templateIDs := range scenario.HostTemplates {
		for _, templateID := range templateIDs {
			_, err := db.dbInsert("INSERT INTO hosts_templates(scenario_id, host_id, template_id) VALUES($1, $2, $3)", id, hostID, templateID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (db dbObj) InsertScenario(scenario model.Scenario) (int64, error) {
	id, err := db.dbInsert("INSERT INTO scenarios(name, description, enabled) VALUES($1, $2, $3) RETURNING id", scenario.Name, scenario.Description, scenario.Enabled)
	if err != nil {
		return -1, err
	}

	return id, db.InsertScenarioHostTemplates(id, scenario)
}

func (db dbObj) UpdateScenario(id int64, scenario model.Scenario) error {
	enabled := 1
	if !scenario.Enabled {
		enabled = 0
	}
	err := db.dbUpdate("UPDATE scenarios SET name=$1, description=$2, enabled=$3 WHERE id=$4", scenario.Name, scenario.Description, enabled, id)
	if err != nil {
		return err
	}
	err = db.DeleteScenarioHostTemplates(id)
	if err != nil {
		return err
	}
	return db.InsertScenarioHostTemplates(id, scenario)
}

func (db dbObj) SelectLatestScenarioScores(scenarioID int64) ([]model.TeamScore, error) {
	rows, err := db.dbConn.Query("SELECT teams.name, scores.timestamp, scores.score FROM scores, teams, team_host_tokens WHERE scenario_id=$1 AND scores.host_token=team_host_tokens.host_token AND teams.id=team_host_tokens.team_id GROUP BY team_host_tokens.team_id,team_host_tokens.hostname ORDER BY teams.name ASC,max(scores.timestamp) DESC", scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teamName string
	var timestamp int64
	var score int64
	scores := make([]model.TeamScore, 0)

	for rows.Next() {
		err = rows.Scan(&teamName, &timestamp, &score)
		if err != nil {
			return nil, err
		}
		var entry model.TeamScore
		entry.TeamName = teamName
		entry.Timestamp = timestamp
		entry.Score = score
		scores = append(scores, entry)
	}

	// combine scores into total scores
	totalScoresMap := make(map[string]model.TeamScore)
	teamNames := make([]string, 0)
	for _, score := range scores {
		storedScore, present := totalScoresMap[score.TeamName]
		if !present {
			totalScoresMap[score.TeamName] = score
			storedScore = score
			teamNames = append(teamNames, score.TeamName)
			continue
		}
		storedScore.Score += score.Score
		if score.Timestamp > storedScore.Timestamp {
			storedScore.Timestamp = score.Timestamp
		}
		totalScoresMap[score.TeamName] = storedScore
	}
	totalScores := make([]model.TeamScore, len(teamNames))
	for i, teamName := range teamNames {
		storedScore, _ := totalScoresMap[teamName]
		totalScores[i] = storedScore
	}

	return totalScores, nil
}

func (db dbObj) SelectScenarioTimeline(scenarioID int64, hostToken string) (model.ScenarioTimeline, error) {
	var timeline model.ScenarioTimeline
	timeline.Timestamps = make([]int64, 0)
	timeline.Scores = make([]int64, 0)

	rows, err := db.dbConn.Query("SELECT timestamp, score FROM scores WHERE scenario_id=$1 AND host_token=$2 ORDER BY timestamp ASC", scenarioID, hostToken)
	if err != nil {
		return timeline, err
	}
	defer rows.Close()

	var timestamp int64
	var score int64

	for rows.Next() {
		err := rows.Scan(&timestamp, &score)
		if err != nil {
			return timeline, err
		}
		timeline.Timestamps = append(timeline.Timestamps, timestamp)
		timeline.Scores = append(timeline.Scores, score)
	}

	return timeline, nil
}

func (db dbObj) InsertScenarioScore(entry model.ScenarioHostScore) error {
	_, err := db.dbInsert("INSERT INTO scores(scenario_id, host_token, timestamp, score) VALUES($1, $2, $3, $4)", entry.ScenarioID, entry.HostToken, entry.Timestamp, entry.Score)
	if err != nil {
		return err
	}

	return nil
}

func (db dbObj) SelectLatestScenarioReport(scenarioID int64, hostToken string) (model.Report, error) {
	var report model.Report
	rows, err := db.dbConn.Query("SELECT report FROM reports WHERE scenario_id=$1 AND host_token=$2 GROUP BY timestamp ORDER BY timestamp DESC", scenarioID, hostToken)
	if err != nil {
		return report, err
	}
	defer rows.Close()

	var reportBytes []byte

	for rows.Next() {
		err := rows.Scan(&reportBytes)
		if err != nil {
			return report, err
		}
		json.Unmarshal(reportBytes, &report)
		break
	}

	return report, nil
}

func (db dbObj) InsertScenarioReport(scenarioID int64, hostToken string, entry model.Report) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = db.dbInsert("INSERT INTO reports(scenario_id, host_token, timestamp, report) VALUES($1, $2, $3, $4)", scenarioID, hostToken, entry.Timestamp, b)
	if err != nil {
		return err
	}

	return nil
}

func (db dbObj) SelectTeamScenarioHosts(teamID int64) ([]model.ScenarioHosts, error) {
	scenarioHosts := make([]model.ScenarioHosts, 0)
	rows, err := db.dbConn.Query("SELECT DISTINCT scenarios.name, scenarios.id, hosts.hostname, hosts.id, hosts.os FROM reports, scenarios, hosts, team_host_tokens WHERE team_host_tokens.team_id=$1 AND scenarios.enabled=TRUE AND reports.scenario_id=scenarios.id AND reports.host_token=team_host_tokens.host_token AND hosts.hostname=team_host_tokens.hostname", teamID)
	if err != nil {
		return scenarioHosts, err
	}
	defer rows.Close()

	var scenarioName string
	var scenarioID int64
	var hostname string
	var hostID int64
	var hostOS string

	// need to collect scenario to hosts mapping
	collectedHosts := make(map[int64][]model.Host)
	// cache scenario name
	scenarioNames := make(map[int64]string)

	for rows.Next() {
		err := rows.Scan(&scenarioName, &scenarioID, &hostname, &hostID, &hostOS)
		if err != nil {
			return scenarioHosts, err
		}
		hosts, present := collectedHosts[scenarioID]
		if !present {
			hosts = make([]model.Host, 0)
			collectedHosts[scenarioID] = hosts
		}
		scenarioNames[scenarioID] = scenarioName
		var host model.Host
		host.Hostname = hostname
		host.ID = hostID
		host.OS = hostOS
		collectedHosts[scenarioID] = append(hosts, host)
	}

	// create model instances
	for scenarioID, hosts := range collectedHosts {
		var sh model.ScenarioHosts
		sh.ScenarioID = scenarioID
		sh.Hosts = hosts

		sh.ScenarioName = scenarioNames[scenarioID]

		scenarioHosts = append(scenarioHosts, sh)
	}

	return scenarioHosts, nil
}

func (db dbObj) InsertHostToken(token string, timestamp int64, hostname string, source string) error {
	_, err := db.dbInsert("INSERT INTO host_tokens(host_token, timestamp, hostname, source) VALUES($1, $2, $3, $4)", token, timestamp, hostname, source)
	if err != nil {
		return err
	}

	return nil
}

func (db dbObj) SelectTeamIDFromHostToken(hostToken string) (int64, error) {
	rows, err := db.dbConn.Query("SELECT team_id from team_host_tokens WHERE host_token=$1", hostToken)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var teamID int64
	teamID = -1

	// should only be one match
	for rows.Next() {
		rows.Scan(&teamID)
		break
	}
	// should have changed
	if teamID == -1 {
		return teamID, errors.New("No team ID found")
	}

	return teamID, nil
}

func (db dbObj) InsertTeamHostToken(teamID int64, hostname string, hostToken string, timestamp int64) error {
	_, err := db.dbInsert("INSERT INTO team_host_tokens(team_id, hostname, host_token, timestamp) VALUES($1, $2, $3, $4)", teamID, hostname, hostToken, timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (db dbObj) SelectHostTokens(teamID int64, hostname string) ([]string, error) {
	hostTokens := make([]string, 0)

	rows, err := db.dbConn.Query("SELECT DISTINCT host_token from team_host_tokens WHERE team_id=$1 AND hostname=$2 ORDER BY timestamp ASC", teamID, hostname)
	if err != nil {
		return hostTokens, err
	}
	defer rows.Close()

	for rows.Next() {
		var hostToken string
		rows.Scan(&hostToken)
		hostTokens = append(hostTokens, hostToken)
	}
	if len(hostTokens) == 0 {
		return hostTokens, errors.New("No host token found")
	}

	return hostTokens, nil
}
