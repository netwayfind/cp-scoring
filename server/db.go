package main

import (
	"database/sql"
	"log"
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