package main

import (
	"database/sql"
	"log"
	"github.com/sumwonyuno/cp-scoring/model"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var err error

func DBInit() {
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

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS template_assignment(template_id INTEGER NOT NULL, host_id INTEGER NOT NULL, FOREIGN KEY(template_id) REFERENCES templates(id), FOREIGN KEY(host_id) REFERENCES hosts(id))")
	if err != nil {
		log.Fatal("ERROR: cannot create table template_assignment;", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("ERROR: cannot create table template_assignment;", err)
	}

	log.Println("Finished setting up database")
}

func DBClose() {
	db.Close()
}

func DBSelectStates() {
	rows, err := db.Query("SELECT * FROM states")
	if err != nil {
		log.Println("ERROR: cannot select from table states;", err)
		return
	}
	defer rows.Close()

	var name string

	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			log.Println("ERROR: fetching row;", err)
			break
		}

		log.Print(name)
	}
}

func DBInsertState(state string) {
	stmt, err := db.Prepare("INSERT INTO states(state) VALUES(?)")
	if err != nil {
		log.Println("ERROR: cannot insert into table states;", err)
		return
	}
	_, err = stmt.Exec(state)
	if err != nil {
		log.Println("ERROR: cannot insert into table states;", err)
		return
	}
}

func DBSelectTemplates() map[int64]string {
	rows, err := db.Query("SELECT id, template FROM templates")
	if err != nil {
		log.Println("ERROR: cannot select from templates;", err)
		return nil
	}
	defer rows.Close()
	
	var id int64
	var template string
	templates := make(map[int64]string)

	for rows.Next() {
		err = rows.Scan(&id, &template)
		if err != nil {
			log.Println("ERROR: fetching row;", err)
			break
		}
		templates[id] = template
	}

	return templates
}

func DBSelectTemplate(id int64) string {
	template := "{}"

	stmt, err := db.Prepare("SELECT template FROM templates where id=(?)")
	if err != nil {
		log.Println("ERROR: cannot select from templates;", err)
		return template
	}
	rows, err := stmt.Query(id)
	if err != nil {
		log.Println("ERROR: cannot select from templates;", err)
		return template
	}

	for rows.Next() {
		err := rows.Scan(&template)
		if err != nil {
			log.Println("ERROR: fetching row;", err)
			return template
		}
		// only get first result
		break
	}

	return template
}

func DBSelectTemplatesForHostname(hostname string) []string {
	return nil
}

func DBInsertTemplate(template string) {
	stmt, err := db.Prepare("INSERT INTO templates(template) VALUES(?)")
	if err != nil {
		log.Println("ERROR: cannot insert into table templates;", err)
		return
	}
	_, err = stmt.Exec(template)
	if err != nil {
		log.Println("ERROR: cannot insert into table templates;", err)
		return
	}
}

func DBSelectHosts() map[int64]model.Host {
	rows, err := db.Query("SELECT id, hostname, os FROM hosts")
	if err != nil {
		log.Println("ERROR: cannot select from hosts;", err)
		return nil
	}
	defer rows.Close()
	
	var id int64
	var hostname string
	var os string
	hosts := make(map[int64]model.Host)

	for rows.Next() {
		err = rows.Scan(&id, &hostname, &os)
		if err != nil {
			log.Println("ERROR: fetching row;", err)
			break
		}
		var host model.Host
		host.Hostname = hostname
		host.OS = os
		hosts[id] = host
	}

	return hosts
}

func DBSelectHost(id int64) (model.Host, error) {
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
	for rows.Next() {
		err := rows.Scan(&hostname, &os)
		if err != nil {
			return host, err
		}
		// only get first result
		host.Hostname = hostname
		host.OS = os
		break
	}

	return host, nil
}

func DBInsertHost(host model.Host) {
	stmt, err := db.Prepare("INSERT INTO hosts(hostname, os) VALUES(?, ?)")
	if err != nil {
		log.Println("ERROR: cannot insert into table hosts;", err)
		return
	}
	_, err = stmt.Exec(host.Hostname, host.OS)
	if err != nil {
		log.Println("ERROR: cannot insert into table hosts;", err)
		return
	}
}
