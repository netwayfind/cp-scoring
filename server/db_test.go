// +build integration

package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/sumwonyuno/cp-scoring/model"
)

var testDBObj backingStore
var directDBConn *sql.DB

func getTestBackingStore() (backingStore, error) {
	if testDBObj == nil {
		configFileBytes, err := ioutil.ReadFile("cp-config.test.conf")
		if err != nil {
			log.Fatal("ERROR: unable to read config file;", err)
		}
		var sqlURL string
		for _, line := range strings.Split(string(configFileBytes), "\n") {
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			if strings.HasPrefix(line, "#") {
				continue
			}

			tokens := strings.Split(line, " ")
			if tokens[0] == "sql_url" {
				sqlURL = tokens[1]
			} else {
				log.Fatalf("ERROR: unknown config file setting %s\n", tokens[0])
			}
		}

		dbc, err := sql.Open("postgres", sqlURL)
		if err != nil {
			return testDBObj, err
		}
		directDBConn = dbc

		err = dropTables()
		if err != nil {
			return testDBObj, err
		}

		tbs, err := getBackingStore("postgres", sqlURL)
		if err != nil {
			return testDBObj, err
		}
		testDBObj = tbs
	}
	return testDBObj, nil
}

func dropTables() error {
	tx, err := directDBConn.Begin()
	if err != nil {
		return err
	}
	tx.Exec("DROP TABLE IF EXISTS reports")
	tx.Exec("DROP TABLE IF EXISTS scores")
	tx.Exec("DROP TABLE IF EXISTS hosts_templates")
	tx.Exec("DROP TABLE IF EXISTS scenarios")
	tx.Exec("DROP TABLE IF EXISTS states")
	tx.Exec("DROP TABLE IF EXISTS team_host_tokens")
	tx.Exec("DROP TABLE IF EXISTS host_tokens")
	tx.Exec("DROP TABLE IF EXISTS hosts")
	tx.Exec("DROP TABLE IF EXISTS templates")
	tx.Exec("DROP TABLE IF EXISTS teams")
	tx.Exec("DROP TABLE IF EXISTS admins")
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func clearTables() error {
	tx, err := directDBConn.Begin()
	if err != nil {
		return err
	}
	tx.Exec("TRUNCATE TABLE reports CASCADE")
	tx.Exec("TRUNCATE TABLE scores CASCADE")
	tx.Exec("TRUNCATE TABLE hosts_templates CASCADE")
	tx.Exec("TRUNCATE TABLE scenarios CASCADE")
	tx.Exec("TRUNCATE TABLE states CASCADE")
	tx.Exec("TRUNCATE TABLE team_host_tokens CASCADE")
	tx.Exec("TRUNCATE TABLE host_tokens CASCADE")
	tx.Exec("TRUNCATE TABLE hosts CASCADE")
	tx.Exec("TRUNCATE TABLE templates CASCADE")
	tx.Exec("TRUNCATE TABLE teams CASCADE")
	tx.Exec("TRUNCATE TABLE admins CASCADE")
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func initBackingStore(t *testing.T) backingStore {
	backingStore, err := getTestBackingStore()
	if err != nil {
		t.Fatal(err)
	}
	err = clearTables()
	if err != nil {
		t.Fatal(err)
	}
	return backingStore
}

func TestGetPostgresBackingStore(t *testing.T) {
	backingStore, err := getTestBackingStore()
	if err != nil {
		t.Fatal(err)
	}
	if backingStore == nil {
		t.Fatal("Expected postgres backing store to not be nil")
	}
}

func TestInsertState(t *testing.T) {
	backingStore := initBackingStore(t)

	state := model.State{Hostname: "test"}
	stateBytes, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}

	// state, no existing host token
	err = backingStore.InsertState(1000, "127.0.0.1", "host-token", stateBytes)
	if err == nil {
		t.Fatal("Expected error due to missing host token")
	}

	// state, existing host token
	err = backingStore.InsertHostToken("host-token", 1000, "host", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertState(1000, "127.0.0.1", "host-token", stateBytes)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := directDBConn.Query("SELECT * FROM states")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	var timestamp int64
	var source string
	var hostToken string
	var readStateBytes []byte
	var readState model.State
	for rows.Next() {
		err = rows.Scan(&timestamp, &source, &hostToken, &readStateBytes)
		if err != nil {
			t.Fatal(err)
		}
		if timestamp != 1000 {
			t.Fatal("Unexpected timestamp value")
		}
		if source != "127.0.0.1" {
			t.Fatal("Unexpected source value")
		}
		if hostToken != "host-token" {
			t.Fatal("Unexpected host token value")
		}
		// check bytes can be turned into State
		err = json.Unmarshal(readStateBytes, &readState)
		if err != nil {
			t.Fatal(err)
		}
		if readState.Hostname != "test" {
			t.Fatal("Unexpected hostname from state")
		}
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected number of rows:", counter)
	}
}

func TestInsertScenarioScore(t *testing.T) {
	backingStore := initBackingStore(t)

	// insert score, no existing scenario
	err := backingStore.InsertScenarioScore(model.ScenarioHostScore{})
	if err == nil {
		t.Fatal("Expected error due to missing scenario")
	}

	// insert sample scenario
	scenarioID, err := backingStore.InsertScenario(model.Scenario{})
	if err != nil {
		t.Fatal(err)
	}

	// insert sample scores, no existing host token
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		Timestamp:  1000,
		HostToken:  "host1",
		ScenarioID: scenarioID,
		Score:      1,
	})
	if err == nil {
		t.Fatal("Expected error due to missing host token")
	}

	// insert sample host token
	err = backingStore.InsertHostToken("host1", 0, "host", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	// insert sample score
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		Timestamp:  1000,
		HostToken:  "host1",
		ScenarioID: scenarioID,
		Score:      1,
	})
	if err != nil {
		t.Fatal(err)
	}

	rows, err := directDBConn.Query("SELECT * from scores")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	var readScenarioID uint64
	var readHostToken string
	var readTimestamp int64
	var readScore int64
	for rows.Next() {
		err = rows.Scan(&readScenarioID, &readHostToken, &readTimestamp, &readScore)
		if err != nil {
			t.Fatal(err)
		}
		if readScenarioID != scenarioID {
			t.Fatal("Unexpected scenario ID")
		}
		if readHostToken != "host1" {
			t.Fatal("Unexpected host token")
		}
		if readTimestamp != 1000 {
			t.Fatal("Unexpected timestamp")
		}
		if readScore != 1 {
			t.Fatal("Unexpected score")
		}
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected number of rows:", counter)
	}
}

func TestSelectLatestScenarioScores(t *testing.T) {
	backingStore := initBackingStore(t)

	// insert sample scenario
	scenarioID, err := backingStore.InsertScenario(model.Scenario{})
	if err != nil {
		t.Fatal(err)
	}
	// insert sample host tokens
	err = backingStore.InsertHostToken("host1_1", 0, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host1_2", 0, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host1_3", 0, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host2_1", 0, "host2", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	// insert sample teams
	team1ID, err := backingStore.InsertTeam(model.Team{Name: "Team 1"})
	if err != nil {
		t.Fatal(err)
	}
	team2ID, err := backingStore.InsertTeam(model.Team{Name: "Team 2"})
	if err != nil {
		t.Fatal(err)
	}
	// insert sample team, host token mapping
	err = backingStore.InsertTeamHostToken(team1ID, "host1_1", 0)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertTeamHostToken(team1ID, "host1_2", 0)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertTeamHostToken(team1ID, "host2_1", 0)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertTeamHostToken(team2ID, "host1_3", 0)
	if err != nil {
		t.Fatal(err)
	}

	// no existing scores
	scores, err := backingStore.SelectLatestScenarioScores(scenarioID)
	if err != nil {
		t.Fatal(err)
	}
	if len(scores) != 0 {
		t.Fatal("Expected no scores")
	}

	// insert sample scores
	// team 1 has 2 hosts, host1 has 2 instances, host2 has 1 instance
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		Timestamp:  1000,
		HostToken:  "host1_1",
		ScenarioID: scenarioID,
		Score:      1,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		Timestamp:  1005,
		HostToken:  "host1_1",
		ScenarioID: scenarioID,
		Score:      2,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		Timestamp:  1010,
		HostToken:  "host1_2",
		ScenarioID: scenarioID,
		Score:      2,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		Timestamp:  1002,
		HostToken:  "host2_1",
		ScenarioID: scenarioID,
		Score:      1,
	})
	if err != nil {
		t.Fatal(err)
	}
	// team 2 just has host1
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		Timestamp:  1040,
		HostToken:  "host1_3",
		ScenarioID: scenarioID,
		Score:      6,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		Timestamp:  1039,
		HostToken:  "host1_3",
		ScenarioID: scenarioID,
		Score:      7,
	})
	if err != nil {
		t.Fatal(err)
	}

	scores, err = backingStore.SelectLatestScenarioScores(scenarioID)
	if err != nil {
		t.Fatal(err)
	}
	// should just be 2 latest scores, one for each team
	if len(scores) != 2 {
		t.Fatal("Unexpected number of scores:", len(scores))
	}
	// should be ordered by team name
	if scores[0].TeamName != "Team 1" {
		t.Fatal("Unexpected team name")
	}
	if scores[0].Timestamp != 1010 {
		t.Fatal("Unexpected timestamp")
	}
	if scores[0].Score != 5 {
		t.Fatal("Unexpected score", scores[0].Score)
	}
	if scores[1].TeamName != "Team 2" {
		t.Fatal("Unexpected team name")
	}
	if scores[1].Timestamp != 1040 {
		t.Fatal("Unexpected timestamp")
	}
	if scores[1].Score != 6 {
		t.Fatal("Unexpected score")
	}
}

func TestInsertScenarioReport(t *testing.T) {
	backingStore := initBackingStore(t)

	// sample report
	findings := append(make([]model.Finding, 0), model.Finding{Show: true, Message: "test", Value: 1})
	report := model.Report{Timestamp: 1500, Findings: findings}

	// insert report without scenario and host token
	err := backingStore.InsertScenarioReport(0, "host-token", report)
	if err == nil {
		t.Fatal("Expected error")
	}

	// insert sample scenario and host token
	scenarioID, err := backingStore.InsertScenario(model.Scenario{})
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host-token", 0, "host1", "127.0.0.1")

	// insert sample report
	err = backingStore.InsertScenarioReport(scenarioID, "host-token", report)
	if err != nil {
		t.Fatal(err)
	}

	// check report inserted
	rows, err := directDBConn.Query("SELECT * from reports")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	var readScenarioID uint64
	var readHostToken string
	var readTimestamp int64
	var readReportBytes []byte
	var readReport model.Report
	for rows.Next() {
		err = rows.Scan(&readScenarioID, &readHostToken, &readTimestamp, &readReportBytes)
		if err != nil {
			t.Fatal(err)
		}
		if readScenarioID != scenarioID {
			t.Fatal("Unexpected scenario ID")
		}
		if readHostToken != "host-token" {
			t.Fatal("Unexpected host token")
		}
		if readTimestamp != 1500 {
			t.Fatal("Unexpected timestamp")
		}
		// check bytes can be turned into Report
		err = json.Unmarshal(readReportBytes, &readReport)
		if err != nil {
			t.Fatal(err)
		}
		if readReport.Timestamp != 1500 {
			t.Fatal("Unexpected timestamp from report")
		}
		if len(readReport.Findings) != 1 {
			t.Fatal("Unexpected number of findings in report")
		}
		finding := readReport.Findings[0]
		if finding.Show != true {
			t.Fatal("Unexpected finding show setting")
		}
		if finding.Message != "test" {
			t.Fatal("Unexpected finding message")
		}
		if finding.Value != 1 {
			t.Fatal("Unexpected finding value")
		}
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected number of rows:", counter)
	}
}

func TestSelectLatestScenarioReport(t *testing.T) {
	backingStore := initBackingStore(t)

	// insert sample scenario
	scenarioID, err := backingStore.InsertScenario(model.Scenario{})
	if err != nil {
		t.Fatal(err)
	}
	// insert sample host tokens
	err = backingStore.InsertHostToken("host-token1", 0, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host-token2", 0, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	// no existing reports
	report, err := backingStore.SelectLatestScenarioReport(scenarioID, "host-token")
	if err != nil {
		t.Fatal(err)
	}
	// should be empty report
	if report.Timestamp != 0 {
		t.Fatal("Unexpected report timestamp:", report.Timestamp)
	}
	if len(report.Findings) != 0 {
		t.Fatal("Expected no report findings")
	}

	// insert multiple reports
	findings1 := append(make([]model.Finding, 0), model.Finding{Show: true, Message: "test", Value: 1})
	report1a := model.Report{Timestamp: 1001, Findings: findings1}
	report1b := model.Report{Timestamp: 1000, Findings: findings1}
	findings2 := append(make([]model.Finding, 0), model.Finding{Show: false, Message: "test2", Value: 0})
	report2 := model.Report{Timestamp: 1200, Findings: findings2}

	err = backingStore.InsertScenarioReport(scenarioID, "host-token1", report1a)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioReport(scenarioID, "host-token1", report1b)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioReport(scenarioID, "host-token2", report2)
	if err != nil {
		t.Fatal(err)
	}

	// check first host token
	report, err = backingStore.SelectLatestScenarioReport(scenarioID, "host-token1")
	if err != nil {
		t.Fatal(err)
	}
	// should be latest report
	if report.Timestamp != 1001 {
		t.Fatal("Unexpected report timestamp:", report.Timestamp)
	}
	if len(report.Findings) != 1 {
		t.Fatal("Expected 1 report findings")
	}
	finding := report.Findings[0]
	if finding.Show != true {
		t.Fatal("Unexpected finding show setting")
	}
	if finding.Message != "test" {
		t.Fatal("Unexpected finding message")
	}
	if finding.Value != 1 {
		t.Fatal("Unexpected finding value")
	}

	// check second host token
	report, err = backingStore.SelectLatestScenarioReport(scenarioID, "host-token2")
	if err != nil {
		t.Fatal(err)
	}
	// should be latest report
	if report.Timestamp != 1200 {
		t.Fatal("Unexpected report timestamp:", report.Timestamp)
	}
	if len(report.Findings) != 1 {
		t.Fatal("Expected 1 report findings")
	}
	finding = report.Findings[0]
	if finding.Show != false {
		t.Fatal("Unexpected finding show setting")
	}
	if finding.Message != "test2" {
		t.Fatal("Unexpected finding message")
	}
	if finding.Value != 0 {
		t.Fatal("Unexpected finding value")
	}
}

func TestSelectScenarioTimeline(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing data
	timeline, err := backingStore.SelectScenarioTimeline(0, "host-token1")
	if err != nil {
		t.Fatal(err)
	}
	if len(timeline.Timestamps) != 0 {
		t.Fatal("Unexpected number of timestamps:", len(timeline.Timestamps))
	}
	if len(timeline.Scores) != 0 {
		t.Fatal("Unexpected number of scores:", len(timeline.Scores))
	}

	// insert sample scenario
	scenarioID, err := backingStore.InsertScenario(model.Scenario{})
	if err != nil {
		t.Fatal(err)
	}
	// insert sample host tokens
	err = backingStore.InsertHostToken("host-token1", 0, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host-token2", 0, "host2", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	// no existing scores
	timeline, err = backingStore.SelectScenarioTimeline(scenarioID, "host-token1")
	if err != nil {
		t.Fatal(err)
	}
	if len(timeline.Timestamps) != 0 {
		t.Fatal("Unexpected number of timestamps:", len(timeline.Timestamps))
	}
	if len(timeline.Scores) != 0 {
		t.Fatal("Unexpected number of scores:", len(timeline.Scores))
	}

	// insert sample scores
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		HostToken:  "host-token1",
		ScenarioID: scenarioID,
		Score:      14,
		Timestamp:  1300,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		HostToken:  "host-token1",
		ScenarioID: scenarioID,
		Score:      14,
		Timestamp:  1360,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		HostToken:  "host-token1",
		ScenarioID: scenarioID,
		Score:      15,
		Timestamp:  1420,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		HostToken:  "host-token2",
		ScenarioID: scenarioID,
		Score:      7,
		Timestamp:  521,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{
		HostToken:  "host-token2",
		ScenarioID: scenarioID,
		Score:      7,
		Timestamp:  520,
	})
	if err != nil {
		t.Fatal(err)
	}

	// check host token 1
	timeline, err = backingStore.SelectScenarioTimeline(scenarioID, "host-token1")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(timeline.Timestamps, []int64{1300, 1360, 1420}) {
		t.Fatal("Unexpected timestamp values")
	}
	if !reflect.DeepEqual(timeline.Scores, []int64{14, 14, 15}) {
		t.Fatal("Unexpected score values")
	}

	// check host token 2
	timeline, err = backingStore.SelectScenarioTimeline(scenarioID, "host-token2")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(timeline.Timestamps, []int64{520, 521}) {
		t.Fatal("Unexpected timestamp values")
	}
	if !reflect.DeepEqual(timeline.Scores, []int64{7, 7}) {
		t.Fatal("Unexpected timestamp values")
	}
}

func TestInsertAdmin(t *testing.T) {
	backingStore := initBackingStore(t)

	// insert sample admin
	err := backingStore.InsertAdmin("admin", "hash")
	if err != nil {
		t.Fatal(err)
	}

	rows, err := directDBConn.Query("SELECT * FROM admins")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	var username string
	var passwordHash string
	for rows.Next() {
		err = rows.Scan(&username, &passwordHash)
		if err != nil {
			t.Fatal(err)
		}
		if username != "admin" {
			t.Fatal("Unexpected username")
		}
		if passwordHash != "hash" {
			t.Fatal("Unexpected password hash")
		}
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected number of rows:", counter)
	}

	// insert admin with same name
	err = backingStore.InsertAdmin("admin", "hash")
	if err == nil {
		t.Fatal("Expected error for inserting duplicate admin")
	}
	// should still be 1 user
	rows, err = directDBConn.Query("SELECT * FROM admins")
	if err != nil {
		t.Fatal(err)
	}
	counter = 0
	for rows.Next() {
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected number of rows:", counter)
	}
}

func TestIsAdmin(t *testing.T) {
	backingStore := initBackingStore(t)

	// does not exist yet
	present, err := backingStore.IsAdmin("admin")
	if err != nil {
		t.Fatal(err)
	}
	if present {
		t.Fatal("user should not exist yet")
	}

	// add user
	err = backingStore.InsertAdmin("admin", "hash")
	if err != nil {
		t.Fatal(err)
	}

	// should exist
	present, err = backingStore.IsAdmin("admin")
	if err != nil {
		t.Fatal(err)
	}
	if !present {
		t.Fatal("user should exist")
	}
}

func TestSelectAdminPasswordHash(t *testing.T) {
	backingStore := initBackingStore(t)

	// user does not exist yet
	passwordHash, err := backingStore.SelectAdminPasswordHash("admin")
	if err != nil {
		t.Fatal(err)
	}

	// add user and hash
	err = backingStore.InsertAdmin("admin", "hash")
	if err != nil {
		t.Fatal(err)
	}

	// user should exist
	passwordHash, err = backingStore.SelectAdminPasswordHash("admin")
	if err != nil {
		t.Fatal(err)
	}
	if passwordHash != "hash" {
		t.Fatal("Unexpected password hash")
	}
}

func TestSelectAdmins(t *testing.T) {
	backingStore := initBackingStore(t)

	// no admins yet
	admins, err := backingStore.SelectAdmins()
	if err != nil {
		t.Fatal(err)
	}
	if len(admins) != 0 {
		t.Fatal("Unexpected admin count:", len(admins))
	}

	// add sample admins
	err = backingStore.InsertAdmin("admin2", "hash")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertAdmin("admin1", "hash")
	if err != nil {
		t.Fatal(err)
	}

	// should be in sorted order
	admins, err = backingStore.SelectAdmins()
	if err != nil {
		t.Fatal(err)
	}
	if len(admins) != 2 {
		t.Fatal("Unexpected admin count:", len(admins))
	}
	if admins[0] != "admin1" {
		t.Fatal("Unexpected admin")
	}
	if admins[1] != "admin2" {
		t.Fatal("Unexpected admin")
	}
}

func TestUpdateAdmin(t *testing.T) {
	backingStore := initBackingStore(t)

	// update user that does not exist, no errors
	err := backingStore.UpdateAdmin("admin1", "hash")
	if err != nil {
		t.Fatal(err)
	}

	// should not have created user
	rows, err := directDBConn.Query("SELECT * from admins")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	for rows.Next() {
		counter++
	}
	if counter != 0 {
		t.Fatal("Unexpected number of rows:", counter)
	}

	// create sample admin users
	err = backingStore.InsertAdmin("admin1", "hash")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertAdmin("admin2", "hash")
	if err != nil {
		t.Fatal(err)
	}

	// update password hash
	err = backingStore.UpdateAdmin("admin1", "hashhash")
	if err != nil {
		t.Fatal(err)
	}

	// check updates
	passwordHash, err := backingStore.SelectAdminPasswordHash("admin1")
	if err != nil {
		t.Fatal(err)
	}
	if passwordHash != "hashhash" {
		t.Fatal("Unexpected password hash")
	}
	passwordHash, err = backingStore.SelectAdminPasswordHash("admin2")
	if err != nil {
		t.Fatal(err)
	}
	if passwordHash != "hash" {
		t.Fatal("Unexpected password hash")
	}
}

func TestDeleteAdmin(t *testing.T) {
	backingStore := initBackingStore(t)

	// test delete user does not exist yet
	err := backingStore.DeleteAdmin("admin1")
	if err != nil {
		t.Fatal(err)
	}

	// create sample admin users
	err = backingStore.InsertAdmin("admin1", "hash")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertAdmin("admin2", "hash")
	if err != nil {
		t.Fatal(err)
	}

	// should be present
	admins, err := backingStore.SelectAdmins()
	if err != nil {
		t.Fatal(err)
	}
	if len(admins) != 2 {
		t.Fatal("Unexpected number of admins:", len(admins))
	}

	// test delete
	err = backingStore.DeleteAdmin("admin1")
	if err != nil {
		t.Fatal(err)
	}

	// should only be admin2
	admins, err = backingStore.SelectAdmins()
	if err != nil {
		t.Fatal(err)
	}
	if len(admins) != 1 {
		t.Fatal("Unexpected number of admins:", len(admins))
	}
	if admins[0] != "admin2" {
		t.Fatal("Unexpected admin")
	}
}

func TestInsertTeam(t *testing.T) {
	backingStore := initBackingStore(t)

	// sample team
	teamID, err := backingStore.InsertTeam(model.Team{
		Name:    "team1",
		POC:     "person1",
		Email:   "poc@example.com",
		Enabled: true,
		Key:     "12345",
	})

	// check in database
	rows, err := directDBConn.Query("SELECT * FROM teams")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	var readTeamID uint64
	var name string
	var poc string
	var email string
	var enabled bool
	var key string
	if rows.Next() {
		counter++
		err = rows.Scan(&readTeamID, &name, &poc, &email, &enabled, &key)
		if err != nil {
			t.Fatal(err)
		}
		if readTeamID != teamID {
			t.Fatal("Unexpected read team ID")
		}
		if name != "team1" {
			t.Fatal("Unexpected team name")
		}
		if poc != "person1" {
			t.Fatal("Unexpected team POC")
		}
		if email != "poc@example.com" {
			t.Fatal("Unexpected team email")
		}
		if enabled != true {
			t.Fatal("Unexpected team enabled setting")
		}
		if key != "12345" {
			t.Fatal("Unexpected team key")
		}
	}
	if counter != 1 {
		t.Fatal("Unexpected row count:", err)
	}

	// insert team with same name
	teamID, err = backingStore.InsertTeam(model.Team{Name: "team1"})
	if err == nil {
		t.Fatal("Expected error for team with same name")
	}

	// should be no change
	rows, err = directDBConn.Query("SELECT * FROM teams")
	if err != nil {
		t.Fatal(err)
	}
	counter = 0
	if rows.Next() {
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected row count:", counter)
	}
}

func TestSelectTeam(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing team
	team, err := backingStore.SelectTeam(0)
	if err != nil {
		t.Fatal(err)
	}
	if team.ID != 0 {
		t.Fatal("Expected empty team")
	}

	// add sample team
	team1ID, err := backingStore.InsertTeam(model.Team{Name: "team1"})
	if err != nil {
		t.Fatal(err)
	}

	// check can get team
	team, err = backingStore.SelectTeam(team1ID)
	if err != nil {
		t.Fatal(err)
	}
	if team.ID != team1ID {
		t.Fatal("Unexpected team ID:", team.ID)
	}
	if team.Name != "team1" {
		t.Fatal("Unexpected team ID:", team.ID)
	}
}

func TestSelectTeams(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing teams
	teams, err := backingStore.SelectTeams()
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 0 {
		t.Fatal("Unexpected teams count:", len(teams))
	}

	// add sample teams
	team2ID, err := backingStore.InsertTeam(model.Team{Name: "team2"})
	if err != nil {
		t.Fatal(err)
	}
	team1ID, err := backingStore.InsertTeam(model.Team{Name: "team1"})
	if err != nil {
		t.Fatal(err)
	}

	// check teams
	teams, err = backingStore.SelectTeams()
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 2 {
		t.Fatal("Unexpected teams count:", len(teams))
	}
	// should be in team name order
	if teams[0].ID != team1ID {
		t.Fatal("Unexpected team ID")
	}
	if teams[0].Name != "team1" {
		t.Fatal("Unexpected team name")
	}
	if teams[1].ID != team2ID {
		t.Fatal("Unexpected team ID")
	}
	if teams[1].Name != "team2" {
		t.Fatal("Unexpected team name")
	}
}

func TestSelectTeamIDForKey(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing team
	_, err := backingStore.SelectTeamIDForKey("team1key")
	if err == nil {
		t.Fatal("Expected error")
	}

	// insert sample teams
	team1ID, err := backingStore.InsertTeam(model.Team{Name: "team1", Key: "team1key", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	_, err = backingStore.InsertTeam(model.Team{Name: "team2", Key: "team2key", Enabled: false})
	if err != nil {
		t.Fatal(err)
	}

	// check can get team
	// team enabled
	teamID, err := backingStore.SelectTeamIDForKey("team1key")
	if err != nil {
		t.Fatal(err)
	}
	if teamID != team1ID {
		t.Fatal("Unexpected team ID from team key")
	}
	// team disabled
	teamID, err = backingStore.SelectTeamIDForKey("team2key")
	if err == nil {
		t.Fatal("Expected error, team disabled")
	}
}

func TestUpdateTeam(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing team
	err := backingStore.UpdateTeam(0, model.Team{Name: "team1"})
	if err != nil {
		t.Fatal(err)
	}

	// should be no teams
	teams, err := backingStore.SelectTeams()
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 0 {
		t.Fatal("Unexpected team count:", len(teams))
	}

	// add sample teams
	team1ID, err := backingStore.InsertTeam(model.Team{Name: "team1", Key: "team1key", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	team2ID, err := backingStore.InsertTeam(model.Team{Name: "team2", Key: "team2key", Enabled: false})
	if err != nil {
		t.Fatal(err)
	}

	// update team 2
	err = backingStore.UpdateTeam(team2ID, model.Team{Name: "Team 2", Key: "keykey", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}

	// check teams
	team1, err := backingStore.SelectTeam(team1ID)
	if err != nil {
		t.Fatal(err)
	}
	if team1.ID != team1ID {
		t.Fatal("Unexpected team ID")
	}
	if team1.Key != "team1key" {
		t.Fatal("Unexpected team key")
	}
	if team1.Enabled != true {
		t.Fatal("Unexpected team enabled setting")
	}
	team2, err := backingStore.SelectTeam(team2ID)
	if err != nil {
		t.Fatal(err)
	}
	if team2.ID != team2ID {
		t.Fatal("Unexpected team ID")
	}
	if team2.Key != "keykey" {
		t.Fatal("Unexpected team key")
	}
	if team2.Enabled != true {
		t.Fatal("Unexpected team enabled setting")
	}
}

func TestDeleteTeam(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing team
	err := backingStore.DeleteTeam(0)
	if err != nil {
		t.Fatal(err)
	}

	// sample team
	team1ID, err := backingStore.InsertTeam(model.Team{Name: "team1"})
	if err != nil {
		t.Fatal(err)
	}

	// team should be there
	teams, err := backingStore.SelectTeams()
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 1 {
		t.Fatal("Unexpected number of teams:", len(teams))
	}
	if teams[0].ID != team1ID {
		t.Fatal("Unexpected team ID")
	}

	// delete team
	err = backingStore.DeleteTeam(team1ID)
	if err != nil {
		t.Fatal(err)
	}

	// team should not be there
	teams, err = backingStore.SelectTeams()
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 0 {
		t.Fatal("Unexpected number of teams:", len(teams))
	}
}

func TestInsertHost(t *testing.T) {
	backingStore := initBackingStore(t)

	hostID, err := backingStore.InsertHost(model.Host{
		Hostname: "hostname",
		OS:       "this OS",
	})
	if err != nil {
		t.Fatal(err)
	}

	rows, err := directDBConn.Query("SELECT * FROM hosts")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	var readHostID uint64
	var readHostname string
	var readOS string
	if rows.Next() {
		err = rows.Scan(&readHostID, &readHostname, &readOS)
		if err != nil {
			t.Fatal(err)
		}
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected row count:", counter)
	}
	if readHostID != hostID {
		t.Fatal("Unexpected host ID")
	}
	if readHostname != "hostname" {
		t.Fatal("Unexpected hostname")
	}
	if readOS != "this OS" {
		t.Fatal("Unexpected host OS")
	}

	// test inserting host with same hostname
	_, err = backingStore.InsertHost(model.Host{
		Hostname: "hostname",
		OS:       "this OS",
	})
	if err == nil {
		t.Fatal("Expected error for same hostname")
	}
}

func TestSelectHost(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing host
	host, err := backingStore.SelectHost(0)
	if err != nil {
		t.Fatal(err)
	}
	if host.ID != 0 {
		t.Fatal("Expected empty host")
	}

	// add sample host
	hostID, err := backingStore.InsertHost(model.Host{Hostname: "hostname", OS: "this OS"})
	if err != nil {
		t.Fatal(err)
	}

	// check host
	host, err = backingStore.SelectHost(hostID)
	if err != nil {
		t.Fatal(err)
	}
	if host.ID != hostID {
		t.Fatal("Unexpected host ID")
	}
	if host.Hostname != "hostname" {
		t.Fatal("Unexpected hostname")
	}
	if host.OS != "this OS" {
		t.Fatal("Unexpected host OS")
	}
}

func TestSelectHosts(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing hosts
	hosts, err := backingStore.SelectHosts()
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 0 {
		t.Fatal("Unexpected host count:", len(hosts))
	}

	// add sample hosts
	host2ID, err := backingStore.InsertHost(model.Host{Hostname: "host2", OS: "this OS"})
	if err != nil {
		t.Fatal(err)
	}
	host1ID, err := backingStore.InsertHost(model.Host{Hostname: "host1", OS: "this OS"})
	if err != nil {
		t.Fatal(err)
	}

	hosts, err = backingStore.SelectHosts()
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 2 {
		t.Fatal("Unexpected host count:", len(hosts))
	}
	// hosts should be sorted by hostname
	if hosts[0].ID != host1ID {
		t.Fatal("Unexpected host ID")
	}
	if hosts[0].Hostname != "host1" {
		t.Fatal("Unexpected hostname")
	}
	if hosts[0].OS != "this OS" {
		t.Fatal("Unexpected host OS")
	}
	if hosts[1].ID != host2ID {
		t.Fatal("Unexpected host ID")
	}
	if hosts[1].Hostname != "host2" {
		t.Fatal("Unexpected hostname")
	}
	if hosts[1].OS != "this OS" {
		t.Fatal("Unexpected host OS")
	}
}

func TestSelectHostIDForHostname(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing host
	hostID, err := backingStore.SelectHostIDForHostname("hostname")
	if err == nil {
		t.Fatal("Expected error")
	}

	// add sample host
	hostID, err = backingStore.InsertHost(model.Host{Hostname: "hostname"})
	if err != nil {
		t.Fatal(err)
	}

	// check exists
	readHostID, err := backingStore.SelectHostIDForHostname("hostname")
	if err != nil {
		t.Fatal(err)
	}
	if readHostID != hostID {
		t.Fatal("Unexpected host ID")
	}
}

func TestUpdateHost(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing host
	err := backingStore.UpdateHost(0, model.Host{Hostname: "host1"})
	if err != nil {
		t.Fatal(err)
	}

	// should still be no hosts
	hosts, err := backingStore.SelectHosts()
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 0 {
		t.Fatal("Unexpected number of hosts:", len(hosts))
	}

	// add sample hosts
	host1ID, err := backingStore.InsertHost(model.Host{Hostname: "host1", OS: "this OS"})
	if err != nil {
		t.Fatal(err)
	}
	host2ID, err := backingStore.InsertHost(model.Host{Hostname: "host2", OS: "this OS"})
	if err != nil {
		t.Fatal(err)
	}

	// update host2
	err = backingStore.UpdateHost(host2ID, model.Host{Hostname: "Host 2", OS: "this OS"})
	if err != nil {
		t.Fatal(err)
	}

	// check hosts
	host, err := backingStore.SelectHost(host1ID)
	if host.ID != host1ID {
		t.Fatal("Unexpected host ID")
	}
	if host.Hostname != "host1" {
		t.Fatal("Unexpected hostname")
	}
	if host.OS != "this OS" {
		t.Fatal("Unexpected host OS")
	}
	host, err = backingStore.SelectHost(host2ID)
	if host.ID != host2ID {
		t.Fatal("Unexpected host ID")
	}
	if host.Hostname != "Host 2" {
		t.Fatal("Unexpected hostname")
	}
	if host.OS != "this OS" {
		t.Fatal("Unexpected host OS")
	}
}

func TestDeleteHost(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing host
	err := backingStore.DeleteHost(0)
	if err != nil {
		t.Fatal()
	}

	// sample host
	hostID, err := backingStore.InsertHost(model.Host{Hostname: "hostname"})
	if err != nil {
		t.Fatal(err)
	}

	// host should be there
	hosts, err := backingStore.SelectHosts()
	if len(hosts) != 1 {
		t.Fatal("Unexpected number of hosts:", len(hosts))
	}
	if hosts[0].ID != hostID {
		t.Fatal("Unexpected host ID")
	}
	if hosts[0].Hostname != "hostname" {
		t.Fatal("Unexpected hostname")
	}

	// delete host
	err = backingStore.DeleteHost(hostID)
	if err != nil {
		t.Fatal(err)
	}

	// host shouldn't be there
	hosts, err = backingStore.SelectHosts()
	if len(hosts) != 0 {
		t.Fatal("Unexpected number of hosts:", len(hosts))
	}
}

func TestInsertHostToken(t *testing.T) {
	backingStore := initBackingStore(t)

	// insert sample host token
	err := backingStore.InsertHostToken("host-token", 1540, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	// check host token added
	rows, err := directDBConn.Query("SELECT * FROM host_tokens")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	var hostToken string
	var timestamp int64
	var hostname string
	var source string
	for rows.Next() {
		err = rows.Scan(&hostToken, &timestamp, &hostname, &source)
		if err != nil {
			t.Fatal(err)
		}
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected row count:", counter)
	}
	if hostToken != "host-token" {
		t.Fatal("Unexpected host token")
	}
	if timestamp != 1540 {
		t.Fatal("Unexpected timestamp")
	}
	if hostname != "host1" {
		t.Fatal("Unexpected hostname")
	}
	if source != "127.0.0.1" {
		t.Fatal("Unexpected source")
	}
}

func TestInsertTeamHostToken(t *testing.T) {
	backingStore := initBackingStore(t)

	// insert team host token without existing team or host token
	err := backingStore.InsertTeamHostToken(0, "host-token", 1300)
	if err == nil {
		t.Fatal("Expected error")
	}

	// should be no team host token
	rows, err := directDBConn.Query("SELECT * FROM team_host_tokens")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	for rows.Next() {
		counter++
	}
	if counter != 0 {
		t.Fatal("Unexpected row count:", err)
	}

	// add team
	teamID, err := backingStore.InsertTeam(model.Team{Name: "team"})
	if err != nil {
		t.Fatal(err)
	}

	// insert team host token without existing host token
	err = backingStore.InsertTeamHostToken(teamID, "host-token", 1300)
	if err == nil {
		t.Fatal("Expected error")
	}

	// should be no team host token
	rows, err = directDBConn.Query("SELECT * FROM team_host_tokens")
	if err != nil {
		t.Fatal(err)
	}
	counter = 0
	for rows.Next() {
		counter++
	}
	if counter != 0 {
		t.Fatal("Unexpected row count:", err)
	}

	// add host token
	err = backingStore.InsertHostToken("host-token", 1200, "hostname", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	// insert team host token
	err = backingStore.InsertTeamHostToken(teamID, "host-token", 1300)
	if err != nil {
		t.Fatal(err)
	}

	// team host token should be present
	rows, err = directDBConn.Query("SELECT * FROM team_host_tokens")
	if err != nil {
		t.Fatal(err)
	}
	counter = 0
	var readTeamID uint64
	var hostToken string
	var timestamp int64
	for rows.Next() {
		err = rows.Scan(&readTeamID, &hostToken, &timestamp)
		if err != nil {
			t.Fatal(err)
		}
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected row count:", err)
	}
	if readTeamID != teamID {
		t.Fatal("Unexpected team ID")
	}
	if hostToken != "host-token" {
		t.Fatal("Unexpected host token")
	}
	if timestamp != 1300 {
		t.Fatal("Unexpected timestamp")
	}
}

func TestSelectHostTokens(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing host tokens
	hostTokens, err := backingStore.SelectHostTokens(0, "host1")
	if err == nil {
		t.Fatal("Expected error")
	}

	// add sample teams
	team1ID, err := backingStore.InsertTeam(model.Team{Name: "team1"})
	if err != nil {
		t.Fatal(err)
	}
	team2ID, err := backingStore.InsertTeam(model.Team{Name: "team2"})
	if err != nil {
		t.Fatal(err)
	}
	// add sample host tokens
	err = backingStore.InsertHostToken("host-token1a", 251, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host-token1b", 250, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host-token2", 300, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host-token3", 250, "host2", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	// assign host tokens to teams
	err = backingStore.InsertTeamHostToken(team1ID, "host-token1a", 401)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertTeamHostToken(team1ID, "host-token1b", 400)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertTeamHostToken(team2ID, "host-token2", 400)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertTeamHostToken(team1ID, "host-token3", 400)
	if err != nil {
		t.Fatal(err)
	}

	// check host tokens assigned to teams
	hostTokens, err = backingStore.SelectHostTokens(team1ID, "host1")
	if err != nil {
		t.Fatal(err)
	}
	if len(hostTokens) != 2 {
		t.Fatal("Unexpected host tokens count:", len(hostTokens))
	}
	// should be in added timestamp ordered
	if hostTokens[0] != "host-token1b" {
		t.Fatal("Unexpected host token")
	}
	if hostTokens[1] != "host-token1a" {
		t.Fatal("Unexpected host token")
	}
	hostTokens, err = backingStore.SelectHostTokens(team1ID, "host2")
	if err != nil {
		t.Fatal(err)
	}
	if len(hostTokens) != 1 {
		t.Fatal("Unexpected host tokens count:", len(hostTokens))
	}
	if hostTokens[0] != "host-token3" {
		t.Fatal("Unexpected host token")
	}
	hostTokens, err = backingStore.SelectHostTokens(team2ID, "host1")
	if err != nil {
		t.Fatal(err)
	}
	if len(hostTokens) != 1 {
		t.Fatal("Unexpected host tokens count:", len(hostTokens))
	}
	if hostTokens[0] != "host-token2" {
		t.Fatal("Unexpected host token")
	}
	hostTokens, err = backingStore.SelectHostTokens(team2ID, "host2")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestSelectTeamIDFromHostToken(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing host token or team host token
	_, err := backingStore.SelectTeamIDFromHostToken("host-token")
	if err == nil {
		t.Fatal("Expected error")
	}

	// add sample teams
	team1ID, err := backingStore.InsertTeam(model.Team{Name: "team1"})
	if err != nil {
		t.Fatal(err)
	}
	team2ID, err := backingStore.InsertTeam(model.Team{Name: "team2"})
	if err != nil {
		t.Fatal(err)
	}
	// add host tokens
	err = backingStore.InsertHostToken("host-token1a", 750, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host-token1b", 751, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host-token2", 750, "host2", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	// set up team host tokens
	err = backingStore.InsertTeamHostToken(team1ID, "host-token1a", 900)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertTeamHostToken(team2ID, "host-token1b", 900)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertTeamHostToken(team1ID, "host-token2", 900)
	if err != nil {
		t.Fatal(err)
	}

	// check team ID for host token
	teamID, err := backingStore.SelectTeamIDFromHostToken("host-token1a")
	if err != nil {
		t.Fatal(err)
	}
	if teamID != team1ID {
		t.Fatal("Unexpected team ID")
	}
	teamID, err = backingStore.SelectTeamIDFromHostToken("host-token1b")
	if err != nil {
		t.Fatal(err)
	}
	if teamID != team2ID {
		t.Fatal("Unexpected team ID")
	}
	teamID, err = backingStore.SelectTeamIDFromHostToken("host-token2")
	if err != nil {
		t.Fatal(err)
	}
	if teamID != team1ID {
		t.Fatal("Unexpected team ID")
	}
}

func TestInsertTemplate(t *testing.T) {
	backingStore := initBackingStore(t)

	templateID, err := backingStore.InsertTemplate(model.Template{Name: "template1", State: model.State{Hostname: "host1"}})
	if err != nil {
		t.Fatal()
	}

	rows, err := directDBConn.Query("SELECT * FROM templates")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	var readTemplateID uint64
	var templateName string
	var templateStateBytes []byte
	var templateState model.State
	for rows.Next() {
		err = rows.Scan(&readTemplateID, &templateName, &templateStateBytes)
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected row count:", counter)
	}
	if readTemplateID != templateID {
		t.Fatal("Unexpected template ID")
	}
	if templateName != "template1" {
		t.Fatal("Unexpected template name")
	}
	err = json.Unmarshal(templateStateBytes, &templateState)
	if err != nil {
		t.Fatal(err)
	}
	if templateState.Hostname != "host1" {
		t.Fatal("Unexpected hostname")
	}

	// insert template with same name
	templateID, err = backingStore.InsertTemplate(model.Template{Name: "template1", State: model.State{Hostname: "host1"}})
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestSelectTemplate(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing template
	template, err := backingStore.SelectTemplate(0)
	if err != nil {
		t.Fatal(err)
	}
	if template.ID != 0 {
		t.Fatal("Unexpected template ID")
	}

	// insert template
	templateID, err := backingStore.InsertTemplate(model.Template{Name: "template1", State: model.State{Hostname: "host1"}})
	if err != nil {
		t.Fatal()
	}

	// check template
	template, err = backingStore.SelectTemplate(templateID)
	if err != nil {
		t.Fatal(err)
	}
	if template.ID != templateID {
		t.Fatal("Unexpected template ID")
	}
	if template.Name != "template1" {
		t.Fatal("Unexpected template name")
	}
	if template.State.Hostname != "host1" {
		t.Fatal("Unexpected template state content")
	}
}

func TestSelectTemplates(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing templates
	templates, err := backingStore.SelectTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 0 {
		t.Fatal("Unexpected template count:", len(templates))
	}

	// insert sample templates
	template2ID, err := backingStore.InsertTemplate(model.Template{Name: "template2"})
	if err != nil {
		t.Fatal(err)
	}
	template1ID, err := backingStore.InsertTemplate(model.Template{Name: "template1"})
	if err != nil {
		t.Fatal(err)
	}

	// check templates
	templates, err = backingStore.SelectTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 2 {
		t.Fatal("Unexpected template count:", len(templates))
	}
	// should be template name order
	if templates[0].ID != template1ID {
		t.Fatal("Unexpected template ID")
	}
	if templates[0].Name != "template1" {
		t.Fatal("Unexpected template ID")
	}
	if templates[1].ID != template2ID {
		t.Fatal("Unexpected template ID")
	}
	if templates[1].Name != "template2" {
		t.Fatal("Unexpected template ID")
	}
}

func TestUpdateTemplate(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing template
	err := backingStore.UpdateTemplate(0, model.Template{Name: "template1"})
	if err != nil {
		t.Fatal(err)
	}

	// should be no templates added
	templates, err := backingStore.SelectTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 0 {
		t.Fatal("Unexpected template count:", len(templates))
	}

	// insert sample templates
	template1ID, err := backingStore.InsertTemplate(model.Template{Name: "template1"})
	if err != nil {
		t.Fatal(err)
	}
	template2ID, err := backingStore.InsertTemplate(model.Template{Name: "template2"})
	if err != nil {
		t.Fatal(err)
	}

	// update template2
	err = backingStore.UpdateTemplate(template2ID, model.Template{Name: "Template 2"})
	if err != nil {
		t.Fatal(err)
	}

	// check templates
	// template1
	template, err := backingStore.SelectTemplate(template1ID)
	if err != nil {
		t.Fatal(err)
	}
	if template.ID != template1ID {
		t.Fatal("Unexpected template ID")
	}
	if template.Name != "template1" {
		t.Fatal("Unexpected template name")
	}
	// Template 2
	template, err = backingStore.SelectTemplate(template2ID)
	if err != nil {
		t.Fatal(err)
	}
	if template.ID != template2ID {
		t.Fatal("Unexpected template ID")
	}
	if template.Name != "Template 2" {
		t.Fatal("Unexpected template name")
	}
}

func TestDeleteTemplate(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing template
	err := backingStore.DeleteTemplate(0)
	if err != nil {
		t.Fatal(err)
	}

	// add sample template
	templateID, err := backingStore.InsertTemplate(model.Template{Name: "template"})
	if err != nil {
		t.Fatal(err)
	}

	// make sure template present
	template, err := backingStore.SelectTemplate(templateID)
	if err != nil {
		t.Fatal(err)
	}
	if template.ID != templateID {
		t.Fatal("Unexpected template ID")
	}
	if template.Name != "template" {
		t.Fatal("Unexpected template name")
	}

	// delete template
	err = backingStore.DeleteTemplate(templateID)
	if err != nil {
		t.Fatal(err)
	}

	// make sure deleted
	templates, err := backingStore.SelectTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 0 {
		t.Fatal("Unexpected template count:", len(templates))
	}
}
