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

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS states(state VARCHAR)")
	if err != nil {
		log.Fatal("ERROR: cannot create table states;", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("ERROR: cannot create table states;", err)
	}

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS teams(id INTEGER PRIMARY KEY, name VARCHAR NOT NULL, poc VARCHAR NOT NULL, email VARCHAR NOT NULL, enabled BIT NOT NULL)")
	if err != nil {
		log.Fatal("ERROR: cannot create table teams;", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("ERROR: cannot create table teams;", err)
	}

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS teams_tokens(team_id INTEGER NOT NULL, token VARCHAR NOT NULL, FOREIGN KEY(team_id) REFERENCES teams(id))")
	if err != nil {
		log.Fatal("ERROR: cannot create table teams_tokens;", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("ERROR: cannot create table teams_tokens;", err)
	}

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS templates(id INTEGER PRIMARY KEY, template BLOB NOT NULL)")
	if err != nil {
		log.Fatal("ERROR: cannot create table templates;", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("ERROR: cannot create table templates;", err)
	}

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS hosts(id INTEGER PRIMARY KEY, hostname VARCHAR NOT NULL, os VARCHAR NOT NULL)")
	if err != nil {
		log.Fatal("ERROR: cannot create table hosts;", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("ERROR: cannot create table hosts;", err)
	}

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS hosts_templates(host_id INTEGER NOT NULL, template_id INTEGER NOT NULL, FOREIGN KEY(template_id) REFERENCES templates(id), FOREIGN KEY(host_id) REFERENCES hosts(id))")
	if err != nil {
		log.Fatal("ERROR: cannot create table hosts_templates;", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("ERROR: cannot create table hosts_templates;", err)
	}

	log.Println("Finished setting up database")
}

func dbClose() {
	db.Close()
}

func dbInsertState(state string) error {
	stmt, err := db.Prepare("INSERT INTO states(state) VALUES(?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(state)
	if err != nil {
		return err
	}

	return nil
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
	stmt, err := db.Prepare("DELETE FROM teams where id=(?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func dbInsertTeam(team model.Team) error {
	stmt, err := db.Prepare("INSERT INTO teams(name, poc, email, enabled) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(team.Name, team.POC, team.Email, team.Enabled)
	if err != nil {
		return err
	}

	return nil
}

func dbUpdateTeam(id int64, team model.Team) error {
	stmt, err := db.Prepare("UPDATE teams SET name=(?), poc=(?), email=(?), enabled=(?) WHERE id=(?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(team.Name, team.POC, team.Email, team.Enabled, id)
	if err != nil {
		return err
	}

	return nil
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
	stmt, err := db.Prepare("DELETE FROM templates where id=(?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
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
	stmt, err := db.Prepare("INSERT INTO templates(template) VALUES(?)")
	if err != nil {
		return err
	}
	b, err := json.Marshal(template)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(b)
	if err != nil {
		return err
	}

	return nil
}

func dbUpdateTemplate(id int64, template model.Template) error {
	stmt, err := db.Prepare("UPDATE templates SET template=(?) WHERE id=(?)")
	if err != nil {
		return err
	}
	b, err := json.Marshal(template)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(b, id)
	if err != nil {
		return err
	}

	return nil
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
	stmt, err := db.Prepare("DELETE FROM hosts where id=(?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func dbInsertHost(host model.Host) error {
	stmt, err := db.Prepare("INSERT INTO hosts(hostname, os) VALUES(?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(host.Hostname, host.OS)
	if err != nil {
		return err
	}

	return nil
}

func dbUpdateHost(id int64, host model.Host) error {
	stmt, err := db.Prepare("UPDATE hosts SET hostname=(?),os=(?) WHERE id=(?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(host.Hostname, host.OS, id)
	if err != nil {
		return err
	}

	return nil
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
	stmt, err := db.Prepare("INSERT INTO hosts_templates(host_id, template_id) VALUES(?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(hostID, templateID)
	if err != nil {
		return err
	}

	return nil
}

func dbDeleteHostsTemplates(hostID int64, templateID int64) error {
	stmt, err := db.Prepare("DELETE FROM hosts_templates WHERE host_id=(?) AND template_id=(?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(hostID, templateID)
	if err != nil {
		return err
	}

	return nil
}
