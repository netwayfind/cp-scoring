package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sumwonyuno/cp-scoring/model"
)

var db *sql.DB
var err error

func dbInit(dir string) {
	db, err = sql.Open("sqlite3", path.Join(dir, "cp-scoring.db"))
	if err != nil {
		log.Fatal("ERROR: cannot open db file;", err)
	}
	log.Println("Connected to database")

	createTable("states", "CREATE TABLE IF NOT EXISTS states(state VARCHAR)")
	createTable("admins", "CREATE TABLE IF NOT EXISTS admins(username VARCHAR PRIMARY KEY, password_hash VARCHAR NOT NULL)")
	createTable("teams", "CREATE TABLE IF NOT EXISTS teams(id INTEGER PRIMARY KEY, name VARCHAR NOT NULL, poc VARCHAR NOT NULL, email VARCHAR NOT NULL, enabled BIT NOT NULL, key VARCHAR NOT NULL)")
	createTable("templates", "CREATE TABLE IF NOT EXISTS templates(id INTEGER PRIMARY KEY, name VARCHAR NOT NULL, template BLOB NOT NULL)")
	createTable("hosts", "CREATE TABLE IF NOT EXISTS hosts(id INTEGER PRIMARY KEY, hostname VARCHAR NOT NULL, os VARCHAR NOT NULL)")
	createTable("hosts_templates", "CREATE TABLE IF NOT EXISTS hosts_templates(scenario_id INTEGER NOT NULL, host_id INTEGER NOT NULL, template_id INTEGER NOT NULL, FOREIGN KEY(scenario_id) REFERENCES scenarios(id), FOREIGN KEY(template_id) REFERENCES templates(id), FOREIGN KEY(host_id) REFERENCES hosts(id))")
	createTable("scenarios", "CREATE TABLE IF NOT EXISTS scenarios(id INTEGER PRIMARY KEY, name VARCHAR NOT NULL, description VARCHAR NOT NULL, enabled BIT NOT NULL)")
	createTable("scores", "CREATE TABLE IF NOT EXISTS scores(scenario_id INTEGER NOT NULL, team_id INTEGER NOT NULL, host_id INTEGER NOT NULL, timestamp INTEGER NOT NULL, score INTEGER NOT NULL, FOREIGN KEY(scenario_id) REFERENCES scenarios(id), FOREIGN KEY(team_id) REFERENCES teams(id), FOREIGN KEY(host_id) REFERENCES hosts(id))")
	createTable("reports", "CREATE TABLE IF NOT EXISTS reports(scenario_id INTEGER NOT NULL, team_id INTEGER NOT NULL, host_id INTEGER NOT NULL, timestamp INTEGER NOT NULL, report BLOB NOT NULL, FOREIGN KEY(scenario_id) REFERENCES scenarios(id), FOREIGN KEY(team_id) REFERENCES teams(id), FOREIGN KEY(host_id) REFERENCES hosts(id))")

	log.Println("Finished setting up database")
}

func dbClose() {
	db.Close()
}

func createTable(name string, stmtStr string) {
	stmt, err := db.Prepare(stmtStr)
	if err != nil {
		log.Fatal("ERROR: cannot create table "+name+";", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("ERROR: cannot create table "+name+";", err)
	}
}

func dbInsert(stmtStr string, args ...interface{}) (int64, error) {
	stmt, err := db.Prepare(stmtStr)
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func dbDelete(stmtStr string, args ...interface{}) error {
	stmt, err := db.Prepare(stmtStr)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}

func dbUpdate(stmtStr string, args ...interface{}) error {
	stmt, err := db.Prepare(stmtStr)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}

func dbInsertState(state string) error {
	_, err := dbInsert("INSERT INTO states(state) VALUES(?)", state)
	return err
}

func dbSelectAdmins() ([]string, error) {
	var admins []string
	rows, err := db.Query("SELECT DISTINCT username FROM admins")
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

func dbIsAdmin(username string) (bool, error) {
	rows, err := db.Query("SELECT username FROM admins WHERE username=(?)", username)
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

func dbSelectAdminPasswordHash(username string) (string, error) {
	rows, err := db.Query("SELECT password_hash FROM admins WHERE username=(?)", username)
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

func dbDeleteAdmin(username string) error {
	return dbDelete("DELETE FROM admins where username=(?)", username)
}

func dbUpdateAdmin(username string, passwordHash string) error {
	return dbUpdate("UPDATE admins SET password_hash=(?) WHERE username=(?)", passwordHash, username)
}

func dbInsertAdmin(username string, passwordHash string) error {
	_, err := dbInsert("INSERT INTO admins(username, password_hash) VALUES(?, ?)", username, passwordHash)
	return err
}

func dbSelectTeams() ([]model.TeamSummary, error) {
	rows, err := db.Query("SELECT id, name FROM teams ORDER BY name ASC")
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

func dbSelectHostIDForHostname(hostname string) (int64, error) {
	var id int64 = -1

	rows, err := db.Query("SELECT id FROM hosts WHERE hostname=(?)", hostname)
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

func dbSelectTeamIDForKey(key string) (int64, error) {
	var id int64 = -1
	key = strings.TrimSpace(key)

	rows, err := db.Query("SELECT id FROM teams WHERE key=(?) AND enabled=1", key)
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

func dbSelectTeam(id int64) (model.Team, error) {
	var team model.Team

	rows, err := db.Query("SELECT name, poc, email, enabled, key FROM teams where id=(?)", id)
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

func dbDeleteTeam(id int64) error {
	return dbDelete("DELETE FROM teams where id=(?)", id)
}

func dbInsertTeam(team model.Team) error {
	_, err := dbInsert("INSERT INTO teams(name, poc, email, enabled, key) VALUES(?, ?, ?, ?, ?)", team.Name, team.POC, team.Email, team.Enabled, team.Key)
	return err
}

func dbUpdateTeam(id int64, team model.Team) error {
	return dbUpdate("UPDATE teams SET name=(?), poc=(?), email=(?), enabled=(?), key=(?) WHERE id=(?)", team.Name, team.POC, team.Email, team.Enabled, team.Key, id)
}

func dbSelectTemplates() ([]model.TemplateEntry, error) {
	rows, err := db.Query("SELECT id, name, template FROM templates ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int64
	var name string
	var templateBytes []byte
	templateEntries := make([]model.TemplateEntry, 0)

	for rows.Next() {
		err = rows.Scan(&id, &name, &templateBytes)
		if err != nil {
			return nil, err
		}
		var template model.Template
		err = json.Unmarshal(templateBytes, &template)
		if err != nil {
			continue
		}
		var entry model.TemplateEntry
		entry.ID = id
		entry.Name = name
		entry.Template = template
		templateEntries = append(templateEntries, entry)
	}

	return templateEntries, nil
}

func dbSelectTemplate(id int64) (model.TemplateEntry, error) {
	var templateEntry model.TemplateEntry
	var template model.Template
	var name string
	var templateBytes []byte

	rows, err := db.Query("SELECT name, template FROM templates where id=(?)", id)
	if err != nil {
		return templateEntry, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &templateBytes)
		if err != nil {
			return templateEntry, err
		}
		err = json.Unmarshal(templateBytes, &template)
		if err != nil {
			return templateEntry, err
		}
		// only get first result
		templateEntry.ID = id
		templateEntry.Name = name
		templateEntry.Template = template
		break
	}

	return templateEntry, nil
}

func dbDeleteTemplate(id int64) error {
	return dbDelete("DELETE FROM templates where id=(?)", id)
}

func dbSelectScenariosForHostname(hostname string) ([]int64, error) {
	rows, err := db.Query("SELECT scenarios.id FROM hosts, hosts_templates, scenarios WHERE hosts.hostname=(?) AND hosts_templates.host_id=hosts.id AND hosts_templates.scenario_id=scenarios.id AND scenarios.enabled=1", hostname)
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

func dbSelectTemplatesForHostname(scenarioID int64, hostname string) ([]model.Template, error) {
	rows, err := db.Query("SELECT templates.template FROM templates, hosts, hosts_templates WHERE hosts.hostname=(?) AND hosts_templates.scenario_id=(?) AND hosts_templates.host_id=hosts.id AND hosts_templates.template_id=templates.id", hostname, scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := make([]model.Template, 0)
	var templateBytes []byte
	for rows.Next() {
		err := rows.Scan(&templateBytes)
		if err != nil {
			return nil, err
		}
		var template model.Template
		err = json.Unmarshal(templateBytes, &template)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	return templates, nil
}

func dbInsertTemplate(templateEntry model.TemplateEntry) error {
	b, err := json.Marshal(templateEntry.Template)
	if err != nil {
		return err
	}
	_, err = dbInsert("INSERT INTO templates(name, template) VALUES(?,?)", templateEntry.Name, b)
	return err
}

func dbUpdateTemplate(id int64, templateEntry model.TemplateEntry) error {
	b, err := json.Marshal(templateEntry.Template)
	if err != nil {
		return err
	}
	return dbUpdate("UPDATE templates SET name=(?), template=(?) WHERE id=(?)", templateEntry.Name, b, id)
}

func dbSelectHosts() ([]model.Host, error) {
	rows, err := db.Query("SELECT id, hostname, os FROM hosts ORDER BY hostname ASC")
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

func dbSelectHost(id int64) (model.Host, error) {
	var host model.Host

	rows, err := db.Query("SELECT hostname, os FROM hosts where id=(?)", id)
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

func dbDeleteHost(id int64) error {
	return dbDelete("DELETE FROM hosts where id=(?)", id)
}

func dbInsertHost(host model.Host) error {
	_, err := dbInsert("INSERT INTO hosts(hostname, os) VALUES(?, ?)", host.Hostname, host.OS)
	return err
}

func dbUpdateHost(id int64, host model.Host) error {
	return dbUpdate("UPDATE hosts SET hostname=(?),os=(?) WHERE id=(?)", host.Hostname, host.OS, id)
}

func dbSelectScenarioHostTemplates(scenarioID int64) (map[int64][]int64, error) {
	rows, err := db.Query("SELECT host_id, template_id FROM hosts_templates WHERE scenario_id=(?)", scenarioID)
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

func dbSelectScenarios(onlyEnabled bool) ([]model.ScenarioSummary, error) {
	var q strings.Builder

	q.WriteString("SELECT id, name FROM scenarios")
	if onlyEnabled {
		// enabled is a boolean stored as a bit
		q.WriteString(" WHERE enabled=1")
	}
	q.WriteString(" ORDER BY name ASC")
	rows, err := db.Query(q.String())
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

func dbSelectScenario(id int64) (model.Scenario, error) {
	var scenario model.Scenario

	stmt, err := db.Prepare("SELECT name, description, enabled FROM scenarios where id=(?)")
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
		hostTemplates, err := dbSelectScenarioHostTemplates(id)
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

func dbDeleteScenarioHostTemplates(scenarioID int64) error {
	return dbDelete("DELETE FROM hosts_templates WHERE scenario_id=(?)", scenarioID)
}

func dbDeleteScenario(id int64) error {
	err := dbDelete("DELETE FROM scenarios WHERE id=(?)", id)
	if err != nil {
		return err
	}
	return dbDeleteScenarioHostTemplates(id)
}

func dbInsertScenarioHostTemplates(id int64, scenario model.Scenario) error {
	for hostID, templateIDs := range scenario.HostTemplates {
		for _, templateID := range templateIDs {
			_, err = dbInsert("INSERT INTO hosts_templates(scenario_id, host_id, template_id) VALUES(?, ?, ?)", id, hostID, templateID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dbInsertScenario(scenario model.Scenario) error {
	id, err := dbInsert("INSERT INTO scenarios(name, description, enabled) VALUES(?, ?, ?)", scenario.Name, scenario.Description, scenario.Enabled)
	if err != nil {
		return err
	}

	return dbInsertScenarioHostTemplates(id, scenario)
}

func dbUpdateScenario(id int64, scenario model.Scenario) error {
	err := dbUpdate("UPDATE scenarios SET name=(?), description=(?), enabled=(?) WHERE id=(?)", scenario.Name, scenario.Description, scenario.Enabled, id)
	if err != nil {
		return err
	}
	err = dbDeleteScenarioHostTemplates(id)
	if err != nil {
		return err
	}
	return dbInsertScenarioHostTemplates(id, scenario)
}

func dbSelectScenarioLatestScores(scenarioID int64) ([]model.ScenarioLatestScore, error) {
	rows, err := db.Query("SELECT teams.name, scores.timestamp, scores.score FROM scores, teams WHERE scenario_id=(?) AND scores.team_id=teams.id GROUP BY scores.team_id,scores.host_id ORDER BY teams.name ASC,max(scores.timestamp) DESC", scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teamName string
	var timestamp int64
	var score int64
	scores := make([]model.ScenarioLatestScore, 0)

	for rows.Next() {
		err = rows.Scan(&teamName, &timestamp, &score)
		if err != nil {
			return nil, err
		}
		var entry model.ScenarioLatestScore
		entry.TeamName = teamName
		entry.Timestamp = timestamp
		entry.Score = score
		scores = append(scores, entry)
	}

	// combine scores into total scores
	totalScoresMap := make(map[string]model.ScenarioLatestScore)
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
	totalScores := make([]model.ScenarioLatestScore, len(teamNames))
	for i, teamName := range teamNames {
		storedScore, _ := totalScoresMap[teamName]
		totalScores[i] = storedScore
	}

	return totalScores, nil
}

func dbSelectScenarioScores(scenarioID int64, teamID int64) ([]model.ScenarioScore, error) {
	rows, err := db.Query("SELECT host_id, timestamp, score FROM scores WHERE scenario_id=(?) AND team_id=(?) ORDER BY timestamp DESC", scenarioID, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scores := make([]model.ScenarioScore, 0)
	var hostID int64
	var timestamp int64
	var score int64

	for rows.Next() {
		err := rows.Scan(&hostID, &timestamp, &score)
		if err != nil {
			return nil, err
		}
		var entry model.ScenarioScore
		entry.ScenarioID = scenarioID
		entry.TeamID = teamID
		entry.HostID = hostID
		entry.Timestamp = timestamp
		entry.Score = score
		scores = append(scores, entry)
	}

	return scores, nil
}

func dbSelectScenarioTimeline(scenarioID int64, teamID int64, hostID int64) (model.ScenarioTimeline, error) {
	var timeline model.ScenarioTimeline
	timeline.Timestamps = make([]int64, 0)
	timeline.Scores = make([]int64, 0)

	rows, err := db.Query("SELECT timestamp, score FROM scores WHERE scenario_id=(?) AND team_id=(?) AND host_id=(?) ORDER BY timestamp ASC", scenarioID, teamID, hostID)
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

func dbInsertScenarioScore(entry model.ScenarioScore) error {
	_, err := dbInsert("INSERT INTO scores(scenario_id, team_id, host_id, timestamp, score) VALUES(?, ?, ?, ?, ?)", entry.ScenarioID, entry.TeamID, entry.HostID, entry.Timestamp, entry.Score)
	if err != nil {
		return err
	}

	return nil
}

func dbSelectLatestScenarioReport(scenarioID int64, teamID int64, hostID int64) (model.Report, error) {
	var report model.Report
	rows, err := db.Query("SELECT report FROM reports WHERE scenario_id=(?) AND team_id=(?) and host_id=(?) GROUP BY timestamp ORDER BY timestamp DESC", scenarioID, teamID, hostID)
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

func dbInsertScenarioReport(scenarioID int64, teamID int64, hostID int64, entry model.Report) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = dbInsert("INSERT INTO reports(scenario_id, team_id, host_id, timestamp, report) VALUES(?, ?, ?, ?, ?)", scenarioID, teamID, hostID, entry.Timestamp, b)
	if err != nil {
		return err
	}

	return nil
}
