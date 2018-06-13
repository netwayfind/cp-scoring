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