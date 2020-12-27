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
	db.dbCreateTable("teams", "CREATE TABLE IF NOT EXISTS teams(id BIGSERIAL PRIMARY KEY, name VARCHAR UNIQUE NOT NULL, poc VARCHAR NOT NULL, email VARCHAR NOT NULL, enabled BOOLEAN NOT NULL, key VARCHAR NOT NULL)")
	db.dbCreateTable("scenarios", "CREATE TABLE IF NOT EXISTS scenarios(id BIGSERIAL PRIMARY KEY, name VARCHAR UNIQUE NOT NULL, description VARCHAR NOT NULL, enabled BOOLEAN NOT NULL)")
	db.dbCreateTable("scenario_checks", "CREATE TABLE IF NOT EXISTS scenario_checks(scenario_id BIGSERIAL NOT NULL, checks JSONB NOT NULL, FOREIGN KEY(scenario_id) REFERENCES scenarios(id))")

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

func (db dbObj) scenarioChecksSelectAll(id uint64) (map[string][]model.Action, error) {
	rows, err := db.dbConn.Query("SELECT checks FROM scenario_checks WHERE scenario_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hostnameChecks map[string][]model.Action
	var hostnameChecksBs []byte

	for rows.Next() {
		err = rows.Scan(&hostnameChecksBs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(hostnameChecksBs, &hostnameChecks)
		if err != nil {
			return nil, err
		}
		break
	}

	if hostnameChecks == nil {
		hostnameChecks = make(map[string][]model.Action)
	}

	return hostnameChecks, nil
}

func (db dbObj) scenarioChecksUpdate(id uint64, hostnameChecks map[string][]model.Action) error {
	// TODO: transaction
	err := db.dbDelete("DELETE FROM scenario_checks WHERE scenario_id=$1", id)
	if err != nil {
		return err
	}

	b, err := json.Marshal(hostnameChecks)
	if err != nil {
		return err
	}
	_, err = db.dbInsert("INSERT INTO scenario_checks(scenario_id, checks) VALUES ($1, $2)", id, b)
	if err != nil {
		return err
	}

	return nil
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
