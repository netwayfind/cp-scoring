package main

import (
	"database/sql"
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

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS states(state varchar)")
	if err != nil {
		log.Fatal("ERROR: cannot create table states;", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("ERROR: cannot create table states;", err)
	}

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS templates(id INTEGER PRIMARY KEY, template VARCHAR NOT NULL)")
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

func dbSelectTemplates() (map[int64]string, error) {
	rows, err := db.Query("SELECT id, template FROM templates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int64
	var template string
	templates := make(map[int64]string)

	for rows.Next() {
		err = rows.Scan(&id, &template)
		if err != nil {
			return nil, err
		}
		templates[id] = template
	}

	return templates, nil
}

func dbSelectTemplate(id int64) (string, error) {
	var template string

	stmt, err := db.Prepare("SELECT template FROM templates where id=(?)")
	if err != nil {
		return template, err
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return template, err
	}

	for rows.Next() {
		err := rows.Scan(&template)
		if err != nil {
			return template, err
		}
		// only get first result
		break
	}

	return template, nil
}

func dbSelectTemplatesForHostname(hostname string) ([]model.Template, error) {
	stmt, err := db.Prepare("SELECT templates.template FROM templates, hosts, hosts_templates WHERE hosts.hostname=(?) AND hosts_templates.host_id=hosts.id AND hosts_templates.template_id=templates.id")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(hostname)
	if err != nil {
		return nil, err
	}

	var template model.Template
	for rows.Next() {
		err := rows.Scan(&template)
		if err != nil {
			return nil, err
		}
		log.Println(template)
	}
	return nil, nil
}

func dbInsertTemplate(template string) error {
	stmt, err := db.Prepare("INSERT INTO templates(template) VALUES(?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(template)
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

	var hostname string
	var os string
	count := 0
	for rows.Next() {
		err := rows.Scan(&hostname, &os)
		if err != nil {
			return host, err
		}
		// only get first result
		host.Hostname = hostname
		host.OS = os
		count++
		break
	}

	return host, nil
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
