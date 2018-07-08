package main

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sumwonyuno/cp-scoring/model"
)

var db *sql.DB
var err error

func dbInit() {
	db, err = sql.Open("sqlite3", "cp-scoring.db")
	if err != nil {
		log.Fatal("ERROR: cannot open db file;", err)
	}
	log.Println("Connected to database")

	createTable("states", "CREATE TABLE IF NOT EXISTS states(state VARCHAR)")
	createTable("teams", "CREATE TABLE IF NOT EXISTS teams(id INTEGER PRIMARY KEY, name VARCHAR NOT NULL, poc VARCHAR NOT NULL, email VARCHAR NOT NULL, enabled BIT NOT NULL)")
	createTable("teams_tokens", "CREATE TABLE IF NOT EXISTS teams_tokens(team_id INTEGER NOT NULL, token VARCHAR NOT NULL, FOREIGN KEY(team_id) REFERENCES teams(id))")
	createTable("templates", "CREATE TABLE IF NOT EXISTS templates(id INTEGER PRIMARY KEY, template BLOB NOT NULL)")
	createTable("hosts", "CREATE TABLE IF NOT EXISTS hosts(id INTEGER PRIMARY KEY, hostname VARCHAR NOT NULL, os VARCHAR NOT NULL)")
	createTable("host_templates", "CREATE TABLE IF NOT EXISTS hosts_templates(host_id INTEGER NOT NULL, template_id INTEGER NOT NULL, FOREIGN KEY(template_id) REFERENCES templates(id), FOREIGN KEY(host_id) REFERENCES hosts(id))")
	createTable("scenarios", "CREATE TABLE IF NOT EXISTS scenarios(id INTEGER PRIMARY KEY, name VARCHAR NOT NULL, description VARCHAR NOT NULL, enabled BIT NOT NULL)")

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

func dbInsert(stmtStr string, args ...interface{}) error {
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
	return dbInsert("INSERT INTO states(state) VALUES(?)", state)
}

func dbSelectTeams() ([]model.TeamSummary, error) {
	rows, err := db.Query("SELECT id, name FROM teams")
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

func dbSelectTeam(id int64) (model.Team, error) {
	var team model.Team

	rows, err := db.Query("SELECT name, poc, email, enabled FROM teams where id=(?)", id)
	if err != nil {
		return team, err
	}
	defer rows.Close()

	var name string
	var poc string
	var email string
	var enabled bool
	for rows.Next() {
		err := rows.Scan(&name, &poc, &email, &enabled)
		if err != nil {
			return team, err
		}
		team.ID = id
		team.Name = name
		team.POC = poc
		team.Email = email
		team.Enabled = enabled
		// only get first result
		break
	}

	return team, nil
}

func dbDeleteTeam(id int64) error {
	return dbDelete("DELETE FROM teams where id=(?)", id)
}

func dbInsertTeam(team model.Team) error {
	return dbInsert("INSERT INTO teams(name, poc, email, enabled) VALUES(?, ?, ?, ?)", team.Name, team.POC, team.Email, team.Enabled)
}

func dbUpdateTeam(id int64, team model.Team) error {
	return dbUpdate("UPDATE teams SET name=(?), poc=(?), email=(?), enabled=(?) WHERE id=(?)", team.Name, team.POC, team.Email, team.Enabled, id)
}

func dbSelectTemplates() ([]map[int64]model.Template, error) {
	rows, err := db.Query("SELECT id, template FROM templates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int64
	var templateBytes []byte
	templates := make([]map[int64]model.Template, 0)

	for rows.Next() {
		err = rows.Scan(&id, &templateBytes)
		if err != nil {
			return nil, err
		}
		var template model.Template
		err = json.Unmarshal(templateBytes, &template)
		if err != nil {
			continue
		}
		entry := make(map[int64]model.Template)
		entry[id] = template
		templates = append(templates, entry)
	}

	return templates, nil
}

func dbSelectTemplate(id int64) (model.Template, error) {
	var template model.Template
	var templateBytes []byte

	rows, err := db.Query("SELECT template FROM templates where id=(?)", id)
	if err != nil {
		return template, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&templateBytes)
		if err != nil {
			return template, err
		}
		err = json.Unmarshal(templateBytes, &template)
		if err != nil {
			return template, err
		}
		// only get first result
		break
	}

	return template, nil
}

func dbDeleteTemplate(id int64) error {
	return dbDelete("DELETE FROM templates where id=(?)", id)
}

func dbSelectTemplatesForHostname(hostname string) ([]model.Template, error) {
	rows, err := db.Query("SELECT templates.template FROM templates, hosts, hosts_templates WHERE hosts.hostname=(?) AND hosts_templates.host_id=hosts.id AND hosts_templates.template_id=templates.id", hostname)
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

func dbInsertTemplate(template model.Template) error {
	b, err := json.Marshal(template)
	if err != nil {
		return err
	}
	return dbInsert("INSERT INTO templates(template) VALUES(?)", b)
}

func dbUpdateTemplate(id int64, template model.Template) error {
	b, err := json.Marshal(template)
	if err != nil {
		return err
	}
	return dbUpdate("UPDATE templates SET template=(?) WHERE id=(?)", b, id)
}

func dbSelectHosts() ([]model.Host, error) {
	rows, err := db.Query("SELECT id, hostname, os FROM hosts")
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

	stmt, err := db.Prepare("SELECT hostname, os FROM hosts where id=(?)")
	if err != nil {
		return host, err
	}
	rows, err := stmt.Query(id)
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
	return dbInsert("INSERT INTO hosts(hostname, os) VALUES(?, ?)", host.Hostname, host.OS)
}

func dbUpdateHost(id int64, host model.Host) error {
	return dbUpdate("UPDATE hosts SET hostname=(?),os=(?) WHERE id=(?)", host.Hostname, host.OS, id)
}

func dbSelectHostsTemplates() ([]model.HostsTemplates, error) {
	rows, err := db.Query("SELECT host_id, template_id FROM hosts_templates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templateID int64
	var hostID int64
	hostsTemplates := make([]model.HostsTemplates, 0)

	for rows.Next() {
		err = rows.Scan(&hostID, &templateID)
		if err != nil {
			return nil, err
		}
		var entry model.HostsTemplates
		entry.HostID = hostID
		entry.TemplateID = templateID
		hostsTemplates = append(hostsTemplates, entry)
	}

	return hostsTemplates, nil
}

func dbInsertHostsTemplates(hostID int64, templateID int64) error {
	return dbInsert("INSERT INTO hosts_templates(host_id, template_id) VALUES(?, ?)", hostID, templateID)
}

func dbDeleteHostsTemplates(hostID int64, templateID int64) error {
	return dbDelete("DELETE FROM hosts_templates WHERE host_id=(?) AND template_id=(?)", hostID, templateID)
}

func dbSelectScenarios() ([]model.Scenario, error) {
	rows, err := db.Query("SELECT id, name, description, enabled FROM scenarios")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int64
	var name string
	var description string
	var enabled bool
	scenarios := make([]model.Scenario, 0)

	for rows.Next() {
		err = rows.Scan(&id, &name, &description, &enabled)
		if err != nil {
			return nil, err
		}
		var scenario model.Scenario
		scenario.ID = id
		scenario.Name = name
		scenario.Description = description
		scenario.Enabled = enabled
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
		// only get first result
		scenario.ID = id
		scenario.Name = name
		scenario.Description = description
		scenario.Enabled = enabled
		break
	}

	return scenario, nil
}

func dbDeleteScenario(id int64) error {
	return dbDelete("DELETE FROM scenarios where id=(?)", id)
}

func dbInsertScenario(scenario model.Scenario) error {
	return dbInsert("INSERT INTO scenarios(name, description, enabled) VALUES(?, ?, ?)", scenario.Name, scenario.Description, scenario.Enabled)
}

func dbUpdateScenario(id int64, scenario model.Scenario) error {
	return dbUpdate("UPDATE scenarios SET name=(?), description=(?), enabled=(?) WHERE id=(?)", scenario.Name, scenario.Description, scenario.Enabled, id)
}
