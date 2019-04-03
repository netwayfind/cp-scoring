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
		configFileBytes, err := ioutil.ReadFile("cp-scoring.test.conf")
		if err != nil {
			log.Fatal("ERROR: unable to read test config file;", err)
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
	// insert sample hosts
	_, err = backingStore.InsertHost(model.Host{Hostname: "host1"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = backingStore.InsertHost(model.Host{Hostname: "host2"})
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
	// team 1: host 1 has 2 points (latest instance), host 2 has 1 points
	if scores[0].Score != 3 {
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

	// insert duplicate team host token
	err = backingStore.InsertTeamHostToken(team1ID, "host1_1", 100)
	if err != nil {
		t.Fatal(err)
	}

	// scores shouldn't change
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
	// this shouldn't change
	if scores[0].Score != 3 {
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

	// add duplicate host token
	err = backingStore.InsertTeamHostToken(team1ID, "host-token1a", 900)
	if err != nil {
		t.Fatal(err)
	}

	// same host tokens
	hostTokens, err = backingStore.SelectHostTokens(team1ID, "host1")
	if err != nil {
		t.Fatal(err)
	}
	if len(hostTokens) != 2 {
		t.Fatal("Unexpected host tokens count:", len(hostTokens))
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

func TestInsertScenario(t *testing.T) {
	backingStore := initBackingStore(t)

	scenarioID, err := backingStore.InsertScenario(model.Scenario{
		Name:        "Test scenario",
		Description: "description",
		Enabled:     true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// make sure scenario added
	rows, err := directDBConn.Query("SELECT * FROM scenarios")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	var readScenarioID uint64
	var name string
	var description string
	var enabled bool
	for rows.Next() {
		err = rows.Scan(&readScenarioID, &name, &description, &enabled)
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected row count:", counter)
	}
	if readScenarioID != scenarioID {
		t.Fatal("Unexpected scenario ID")
	}
	if name != "Test scenario" {
		t.Fatal("Unexpected scenario name")
	}
	if description != "description" {
		t.Fatal("Unexpected description")
	}
	if enabled != true {
		t.Fatal("Unexpected scenario enabled setting")
	}
}

func TestSelectScenario(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing scenario
	scenario, err := backingStore.SelectScenario(0)
	if err != nil {
		t.Fatal(err)
	}
	if scenario.ID != 0 {
		t.Fatal("Unexpected scenario ID")
	}

	// add scenarios
	// enabled
	scenario1 := model.Scenario{Name: "scenario1", Enabled: true}
	scenario1ID, err := backingStore.InsertScenario(scenario1)
	if err != nil {
		t.Fatal(err)
	}
	// disabled
	scenario2 := model.Scenario{Name: "scenario2", Enabled: false}
	scenario2ID, err := backingStore.InsertScenario(scenario2)
	if err != nil {
		t.Fatal(err)
	}

	// get enabled
	scenario, err = backingStore.SelectScenario(scenario1ID)
	if err != nil {
		t.Fatal(err)
	}
	if scenario.ID != scenario1ID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenario.Name != "scenario1" {
		t.Fatal("Unexpected scenario name")
	}
	if scenario.Enabled != true {
		t.Fatal("Unexpected scenario enabled setting")
	}
	if len(scenario.HostTemplates) != 0 {
		t.Fatal("Unexpected host template count:", len(scenario.HostTemplates))
	}

	// get disabled
	scenario, err = backingStore.SelectScenario(scenario2ID)
	if err != nil {
		t.Fatal(err)
	}
	if scenario.ID != scenario2ID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenario.Name != "scenario2" {
		t.Fatal("Unexpected scenario name")
	}
	if scenario.Enabled != false {
		t.Fatal("Unexpected scenario enabled setting")
	}
	if len(scenario.HostTemplates) != 0 {
		t.Fatal("Unexpected host template count:", len(scenario.HostTemplates))
	}

	// add host templates to scenario
	host := model.Host{Hostname: "host"}
	hostID, err := backingStore.InsertHost(host)
	if err != nil {
		t.Fatal(err)
	}
	template2 := model.Template{Name: "template2"}
	template2ID, err := backingStore.InsertTemplate(template2)
	if err != nil {
		t.Fatal(err)
	}
	template1 := model.Template{Name: "template1"}
	template1ID, err := backingStore.InsertTemplate(template1)
	if err != nil {
		t.Fatal(err)
	}
	scenario1.HostTemplates = map[uint64][]uint64{
		hostID: {template1ID, template2ID},
	}
	err = backingStore.UpdateScenario(scenario1ID, scenario1)
	if err != nil {
		t.Fatal(err)
	}

	// host templates should be present
	scenario, err = backingStore.SelectScenario(scenario1ID)
	if err != nil {
		t.Fatal(err)
	}
	if scenario.ID != scenario1ID {
		t.Fatal("Unexpected scenario ID")
	}
	if len(scenario.HostTemplates) != 1 {
		t.Fatal("Unexpected host template count:", len(scenario.HostTemplates))
	}
	templates, present := scenario.HostTemplates[hostID]
	if !present {
		t.Fatal("Expected host to be in host templates")
	}
	if len(templates) != 2 {
		t.Fatal("Unexpected template count")
	}
	// should be in template ID (insertion) order
	if templates[0] != template2ID {
		t.Fatal("Unexpected template ID")
	}
	if templates[1] != template1ID {
		t.Fatal("Unexpected template ID")
	}
}

func TestSelectScenarios(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing scenarios (only enabled is true)
	scenarios, err := backingStore.SelectScenarios(true)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarios) != 0 {
		t.Fatal("Unexpected scenario count:", len(scenarios))
	}
	// no existing scenarios (only enabled is false)
	scenarios, err = backingStore.SelectScenarios(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarios) != 0 {
		t.Fatal("Unexpected scenario count:", len(scenarios))
	}

	// insert sample scenarios
	scenario4ID, err := backingStore.InsertScenario(model.Scenario{Name: "scenario 4", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	scenario3ID, err := backingStore.InsertScenario(model.Scenario{Name: "scenario 3", Enabled: false})
	if err != nil {
		t.Fatal(err)
	}
	scenario2ID, err := backingStore.InsertScenario(model.Scenario{Name: "scenario 2", Enabled: false})
	if err != nil {
		t.Fatal(err)
	}
	scenario1ID, err := backingStore.InsertScenario(model.Scenario{Name: "scenario 1", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}

	// select all scenarios
	scenarios, err = backingStore.SelectScenarios(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarios) != 4 {
		t.Fatal("Unexpected scenario count:", len(scenarios))
	}
	// should be ordered by scenario name
	if scenarios[0].ID != scenario1ID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenarios[0].Name != "scenario 1" {
		t.Fatal("Unexpected scenario name")
	}
	if scenarios[1].ID != scenario2ID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenarios[1].Name != "scenario 2" {
		t.Fatal("Unexpected scenario name")
	}
	if scenarios[2].ID != scenario3ID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenarios[2].Name != "scenario 3" {
		t.Fatal("Unexpected scenario name")
	}
	if scenarios[3].ID != scenario4ID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenarios[3].Name != "scenario 4" {
		t.Fatal("Unexpected scenario name")
	}

	// select all scenarios (enabled only)
	scenarios, err = backingStore.SelectScenarios(true)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarios) != 2 {
		t.Fatal("Unexpected scenario count:", len(scenarios))
	}
	// should be ordered by scenario name
	if scenarios[0].ID != scenario1ID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenarios[0].Name != "scenario 1" {
		t.Fatal("Unexpected scenario name")
	}
	if scenarios[1].ID != scenario4ID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenarios[1].Name != "scenario 4" {
		t.Fatal("Unexpected scenario name")
	}
}

func TestUpdateScenario(t *testing.T) {
	backingStore := initBackingStore(t)

	// no existing scenario
	err := backingStore.UpdateScenario(0, model.Scenario{Name: "scenario"})
	if err != nil {
		t.Fatal(err)
	}

	// should not have created any scenarios
	scenarios, err := backingStore.SelectScenarios(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarios) != 0 {
		t.Fatal("Unexpected scenario count:", len(scenarios))
	}

	// add sample scenarios
	scenario1ID, err := backingStore.InsertScenario(model.Scenario{Name: "scenario1", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	scenario2ID, err := backingStore.InsertScenario(model.Scenario{Name: "scenario2", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}

	// update scenario
	err = backingStore.UpdateScenario(scenario2ID, model.Scenario{Name: "Scenario 2", Enabled: false})
	if err != nil {
		t.Fatal(err)
	}

	// check updates
	// scenario1
	scenario, err := backingStore.SelectScenario(scenario1ID)
	if err != nil {
		t.Fatal(err)
	}
	if scenario.ID != scenario1ID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenario.Name != "scenario1" {
		t.Fatal("Unexpected scenario name")
	}
	if scenario.Enabled != true {
		t.Fatal("Unexpected scenario enabled setting")
	}
	// Scenario 2
	scenario, err = backingStore.SelectScenario(scenario2ID)
	if err != nil {
		t.Fatal(err)
	}
	if scenario.ID != scenario2ID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenario.Name != "Scenario 2" {
		t.Fatal("Unexpected scenario name")
	}
	if scenario.Enabled != false {
		t.Fatal("Unexpected scenario enabled setting")
	}
}

func TestDeleteScenario(t *testing.T) {
	backingStore := initBackingStore(t)

	// no exsting scenario
	err := backingStore.DeleteScenario(0)
	if err != nil {
		t.Fatal(err)
	}

	// add sample scenario
	scenarioID, err := backingStore.InsertScenario(model.Scenario{Name: "scenario"})
	if err != nil {
		t.Fatal(err)
	}

	// make sure scenario exists
	scenario, err := backingStore.SelectScenario(scenarioID)
	if err != nil {
		t.Fatal(err)
	}
	if scenario.ID != scenarioID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenario.Name != "scenario" {
		t.Fatal("Unexpected scenario name")
	}

	// delete scenario
	err = backingStore.DeleteScenario(scenarioID)
	if err != nil {
		t.Fatal(err)
	}

	// make sure deleted
	scenarios, err := backingStore.SelectScenarios(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarios) != 0 {
		t.Fatal("Unexpected scenario count:", len(scenarios))
	}
}

func TestSelectScenariosForHostname(t *testing.T) {
	backingStore := initBackingStore(t)

	// non-existent hostname
	scenarioIDs, err := backingStore.SelectScenariosForHostname("hostname")
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarioIDs) != 0 {
		t.Fatal("Unexpected scenario count:", len(scenarioIDs))
	}

	// add host
	hostID, err := backingStore.InsertHost(model.Host{Hostname: "hostname"})
	if err != nil {
		t.Fatal(err)
	}

	// no scenarios yet
	scenarioIDs, err = backingStore.SelectScenariosForHostname("hostname")
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarioIDs) != 0 {
		t.Fatal("Unexpected scenario count:", len(scenarioIDs))
	}

	// add sample templates
	template1ID, err := backingStore.InsertTemplate(model.Template{Name: "template1"})
	if err != nil {
		t.Fatal(err)
	}
	template2ID, err := backingStore.InsertTemplate(model.Template{Name: "template2"})
	if err != nil {
		t.Fatal(err)
	}

	// add sample scenarios
	scenario1ID, err := backingStore.InsertScenario(model.Scenario{
		Name:    "scenario1",
		Enabled: false,
		HostTemplates: map[uint64][]uint64{
			hostID: []uint64{template1ID, template2ID},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	scenario2ID, err := backingStore.InsertScenario(model.Scenario{
		Name:    "scenario2",
		Enabled: true,
		HostTemplates: map[uint64][]uint64{
			hostID: []uint64{template1ID, template2ID},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// make sure scenario host templates exist
	rows, err := directDBConn.Query("SELECT * FROM hosts_templates ORDER BY scenario_id, host_id, template_id ASC")
	if err != nil {
		t.Fatal(err)
	}
	counter := 0
	var readScenarioID uint64
	var readHostID uint64
	var readTemplateID uint64
	for rows.Next() {
		rows.Scan(&readScenarioID, &readHostID, &readTemplateID)
		if counter == 0 {
			if readScenarioID != scenario1ID {
				t.Fatal("Unexpected scenario ID")
			}
			if readTemplateID != template1ID {
				t.Fatal("Unexpected template ID")
			}
		} else if counter == 1 {
			if readScenarioID != scenario1ID {
				t.Fatal("Unexpected scenario ID")
			}
			if readTemplateID != template2ID {
				t.Fatal("Unexpected template ID")
			}
		} else if counter == 2 {
			if readScenarioID != scenario2ID {
				t.Fatal("Unexpected scenario ID")
			}
			if readTemplateID != template1ID {
				t.Fatal("Unexpected template ID")
			}
		} else if counter == 3 {
			if readScenarioID != scenario2ID {
				t.Fatal("Unexpected scenario ID")
			}
			if readTemplateID != template2ID {
				t.Fatal("Unexpected template ID")
			}
		} else {
			t.Fatal("Unexpected row")
		}
		if readHostID != hostID {
			t.Fatal("Unexpected host ID")
		}
		counter++
	}
	if counter != 4 {
		t.Fatal("Unexpected host templates:", counter)
	}

	// check scenarios
	scenarioIDs, err = backingStore.SelectScenariosForHostname("hostname")
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarioIDs) != 1 {
		t.Fatal("Unexpected scenario count:", len(scenarioIDs))
	}
	// only scenario2 enabled
	if scenarioIDs[0] != scenario2ID {
		t.Fatal("Unexpected scenario ID")
	}

}

func TestSelectTemplatesForHostname(t *testing.T) {
	backingStore := initBackingStore(t)

	// nothing exists
	templates, err := backingStore.SelectTemplatesForHostname(0, "hostname")
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 0 {
		t.Fatal("Unexpected template count:", len(templates))
	}

	// add host
	host := model.Host{Hostname: "hostname"}
	hostID, err := backingStore.InsertHost(host)
	if err != nil {
		t.Fatal(err)
	}

	// add scenario, no templates
	scenario := model.Scenario{Name: "scenario", Enabled: true}
	scenarioID, err := backingStore.InsertScenario(scenario)
	if err != nil {
		t.Fatal(err)
	}

	// should not have any templates
	templates, err = backingStore.SelectTemplatesForHostname(scenarioID, host.Hostname)
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 0 {
		t.Fatal("Unexpected template count:", len(templates))
	}

	// add template to scenario
	templateID, err := backingStore.InsertTemplate(model.Template{Name: "template"})
	if err != nil {
		t.Fatal(err)
	}
	scenario.HostTemplates = map[uint64][]uint64{
		hostID: {templateID},
	}
	err = backingStore.UpdateScenario(scenarioID, scenario)
	if err != nil {
		t.Fatal(err)
	}

	// scenario should have host templates
	readScenario, err := backingStore.SelectScenario(scenarioID)
	if err != nil {
		t.Fatal(err)
	}
	if readScenario.ID != scenarioID {
		t.Fatal("Unexpected scenario ID")
	}
	readTemplates, present := readScenario.HostTemplates[hostID]
	if !present {
		t.Fatal("Expected ")
	}
	if len(readTemplates) != 1 {
		t.Fatal("Unexpected template count:", len(readTemplates))
	}
	if readTemplates[0] != templateID {
		t.Fatal("Unexpected template ID")
	}

	// should get template
	templates, err = backingStore.SelectTemplatesForHostname(scenarioID, host.Hostname)
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 1 {
		t.Fatal("Unexpected template count:", len(templates))
	}
	if templates[0].ID != templateID {
		t.Fatal("Unexpected template ID")
	}
	if templates[0].Name != "template" {
		t.Fatal("Unexpected template name")
	}

	// disable scenario
	scenario.Enabled = false
	err = backingStore.UpdateScenario(scenarioID, scenario)
	if err != nil {
		t.Fatal(err)
	}

	// should have no templates
	templates, err = backingStore.SelectTemplatesForHostname(scenarioID, host.Hostname)
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 0 {
		t.Fatal("Unexpected template count:", len(templates))
	}
}

func TestSelectTeamScenarioHosts(t *testing.T) {
	backingStore := initBackingStore(t)

	// nothing exists
	scenarioHosts, err := backingStore.SelectTeamScenarioHosts(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarioHosts) != 0 {
		t.Fatal("Unexpected scenario host count")
	}

	// add team
	team := model.Team{Name: "team", Enabled: true}
	teamID, err := backingStore.InsertTeam(team)
	if err != nil {
		t.Fatal(err)
	}

	// should still have no scenario hosts
	scenarioHosts, err = backingStore.SelectTeamScenarioHosts(teamID)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarioHosts) != 0 {
		t.Fatal("Unexpected scenario host count")
	}

	// add scenario
	scenario := model.Scenario{Name: "scenario", Enabled: true}
	scenarioID, err := backingStore.InsertScenario(scenario)
	if err != nil {
		t.Fatal(err)
	}

	// should still have no scenario hosts
	scenarioHosts, err = backingStore.SelectTeamScenarioHosts(teamID)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarioHosts) != 0 {
		t.Fatal("Unexpected scenario host count")
	}

	// add host
	host := model.Host{Hostname: "host"}
	hostID, err := backingStore.InsertHost(host)
	if err != nil {
		t.Fatal(err)
	}
	// add host token
	err = backingStore.InsertHostToken("host-token", 100, host.Hostname, "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	// assigned host token to team
	err = backingStore.InsertTeamHostToken(teamID, "host-token", 101)
	if err != nil {
		t.Fatal(err)
	}

	// check, not yet added to scenario
	scenarioHosts, err = backingStore.SelectTeamScenarioHosts(teamID)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarioHosts) != 0 {
		t.Fatal("Unexpected scenario host count")
	}

	// add host to scenario
	template := model.Template{Name: "template"}
	templateID, err := backingStore.InsertTemplate(template)
	if err != nil {
		t.Fatal(err)
	}
	scenario.HostTemplates = map[uint64][]uint64{
		hostID: {templateID},
	}
	err = backingStore.UpdateScenario(scenarioID, scenario)
	if err != nil {
		t.Fatal(err)
	}

	// should have team scenario hosts
	scenarioHosts, err = backingStore.SelectTeamScenarioHosts(teamID)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarioHosts) != 1 {
		t.Fatal("Unexpected scenario host count:", len(scenarioHosts))
	}
	if scenarioHosts[0].ScenarioID != scenarioID {
		t.Fatal("Unexpected scenario ID")
	}
	if scenarioHosts[0].ScenarioName != scenario.Name {
		t.Fatal("Unexpected scenario ID")
	}
	if len(scenarioHosts[0].Hosts) != 1 {
		t.Fatal("Unexpected scenario host count:", len(scenarioHosts[0].Hosts))
	}
	if scenarioHosts[0].Hosts[0].ID != hostID {
		t.Fatal("Unexpected host ID")
	}
	if scenarioHosts[0].Hosts[0].Hostname != host.Hostname {
		t.Fatal("Unexpected hostname")
	}

	// disable scenario
	scenario.Enabled = false
	err = backingStore.UpdateScenario(scenarioID, scenario)
	if err != nil {
		t.Fatal(err)
	}

	// should not have any scenario hosts
	scenarioHosts, err = backingStore.SelectTeamScenarioHosts(teamID)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarioHosts) != 0 {
		t.Fatal("Unexpected scenario host count:", len(scenarioHosts))
	}

	// re-enable scenario
	scenario.Enabled = true
	err = backingStore.UpdateScenario(scenarioID, scenario)
	if err != nil {
		t.Fatal(err)
	}

	// disable team
	team.Enabled = false
	err = backingStore.UpdateTeam(teamID, team)
	if err != nil {
		t.Fatal(err)
	}

	// should not have any scenario hosts
	scenarioHosts, err = backingStore.SelectTeamScenarioHosts(teamID)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarioHosts) != 0 {
		t.Fatal("Unexpected scenario host count:", len(scenarioHosts))
	}
}

func TestSelectScenarioReports(t *testing.T) {
	backingStore := initBackingStore(t)

	// nothing set up
	reports, err := backingStore.SelectScenarioReports(0, "host-token", 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 0 {
		t.Fatal("Unexpected report count:", len(reports))
	}

	// add scenario, host, and host token
	_, err = backingStore.InsertHost(model.Host{Hostname: "host1"})
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host-token", 0, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	scenarioID, err := backingStore.InsertScenario(model.Scenario{Name: "scenario1"})
	if err != nil {
		t.Fatal(err)
	}

	// should not have any reports
	reports, err = backingStore.SelectScenarioReports(scenarioID, "host-token", 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 0 {
		t.Fatal("Unexpected report count:", len(reports))
	}

	// add sample reports
	report1 := model.Report{
		Timestamp: 15,
		Findings: []model.Finding{
			model.Finding{
				Show:    false,
				Value:   0,
				Message: "no test",
			},
		},
	}
	report2 := model.Report{
		Timestamp: 30,
		Findings: []model.Finding{
			model.Finding{
				Show:    false,
				Value:   0,
				Message: "no test",
			},
		},
	}
	report3 := model.Report{
		Timestamp: 45,
		Findings: []model.Finding{
			model.Finding{
				Show:    true,
				Value:   1,
				Message: "test",
			},
		},
	}
	report4 := model.Report{
		Timestamp: 60,
		Findings: []model.Finding{
			model.Finding{
				Show:    true,
				Value:   1,
				Message: "test",
			},
		},
	}
	// insert out of order
	err = backingStore.InsertScenarioReport(scenarioID, "host-token", report3)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioReport(scenarioID, "host-token", report1)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioReport(scenarioID, "host-token", report4)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioReport(scenarioID, "host-token", report2)
	if err != nil {
		t.Fatal(err)
	}

	// should get reports in order
	reports, err = backingStore.SelectScenarioReports(scenarioID, "host-token", 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 4 {
		t.Fatal("Unexpected report count:", len(reports))
	}
	if reports[0].Timestamp != report1.Timestamp {
		t.Fatal("Unexpected report timestamp")
	}
	if len(reports[0].Findings) != len(report1.Findings) {
		t.Fatal("Unexpected report findings count")
	}
	if reports[0].Findings[0].Message != report1.Findings[0].Message {
		t.Fatal("Unexpected report finding message")
	}
	if reports[1].Timestamp != report2.Timestamp {
		t.Fatal("Unexpected report timestamp")
	}
	if len(reports[1].Findings) != len(report2.Findings) {
		t.Fatal("Unexpected report findings count")
	}
	if reports[1].Findings[0].Message != report2.Findings[0].Message {
		t.Fatal("Unexpected report finding message")
	}
	if reports[2].Timestamp != report3.Timestamp {
		t.Fatal("Unexpected report timestamp")
	}
	if len(reports[2].Findings) != len(report3.Findings) {
		t.Fatal("Unexpected report findings count")
	}
	if reports[2].Findings[0].Message != report3.Findings[0].Message {
		t.Fatal("Unexpected report finding message")
	}
	if reports[3].Timestamp != report4.Timestamp {
		t.Fatal("Unexpected report timestamp")
	}
	if len(reports[3].Findings) != len(report4.Findings) {
		t.Fatal("Unexpected report findings count")
	}
	if reports[3].Findings[0].Message != report4.Findings[0].Message {
		t.Fatal("Unexpected report finding message")
	}

	// should only get reports 2 and 3
	reports, err = backingStore.SelectScenarioReports(scenarioID, "host-token", 30, 45)
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 2 {
		t.Fatal("Unexpected report count:", len(reports))
	}
	if reports[0].Timestamp != report2.Timestamp {
		t.Fatal("Unexpected report timestamp")
	}
	if len(reports[0].Findings) != len(report2.Findings) {
		t.Fatal("Unexpected report findings count")
	}
	if reports[0].Findings[0].Message != report2.Findings[0].Message {
		t.Fatal("Unexpected report finding message")
	}
	if reports[1].Timestamp != report3.Timestamp {
		t.Fatal("Unexpected report timestamp")
	}
	if len(reports[1].Findings) != len(report3.Findings) {
		t.Fatal("Unexpected report findings count")
	}
	if reports[1].Findings[0].Message != report3.Findings[0].Message {
		t.Fatal("Unexpected report finding message")
	}
}

func TestSelectScenarioReportDiffs(t *testing.T) {
	backingStore := initBackingStore(t)

	// nothing set up
	diffs, err := backingStore.SelectScenarioReportDiffs(0, "host-token", 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(diffs) != 0 {
		t.Fatal("Unexpected report diff count:", len(diffs))
	}

	// add scenario, host, and host token
	_, err = backingStore.InsertHost(model.Host{Hostname: "host1"})
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertHostToken("host-token", 0, "host1", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	scenarioID, err := backingStore.InsertScenario(model.Scenario{Name: "scenario1"})
	if err != nil {
		t.Fatal(err)
	}

	// should not have any report diffs
	diffs, err = backingStore.SelectScenarioReportDiffs(scenarioID, "host-token", 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(diffs) != 0 {
		t.Fatal("Unexpected report diff count:", len(diffs))
	}

	// add sample reports
	report1 := model.Report{
		Timestamp: 15,
		Findings: []model.Finding{
			model.Finding{
				Show:    false,
				Value:   0,
				Message: "no test",
			},
		},
	}
	report2 := model.Report{
		Timestamp: 30,
		Findings: []model.Finding{
			model.Finding{
				Show:    false,
				Value:   0,
				Message: "no test",
			},
		},
	}
	report3 := model.Report{
		Timestamp: 45,
		Findings: []model.Finding{
			model.Finding{
				Show:    true,
				Value:   1,
				Message: "test",
			},
		},
	}
	report4 := model.Report{
		Timestamp: 60,
		Findings: []model.Finding{
			model.Finding{
				Show:    true,
				Value:   1,
				Message: "test",
			},
		},
	}
	// insert out of order
	err = backingStore.InsertScenarioReport(scenarioID, "host-token", report3)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioReport(scenarioID, "host-token", report1)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioReport(scenarioID, "host-token", report4)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertScenarioReport(scenarioID, "host-token", report2)
	if err != nil {
		t.Fatal(err)
	}

	// should get reports diffs in order
	diffs, err = backingStore.SelectScenarioReportDiffs(scenarioID, "host-token", 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(diffs) != 2 {
		t.Fatal("Unexpected report diff count:", len(diffs))
	}
	if diffs[0].Type != "Removed" {
		t.Fatal("Unexpected diff type")
	}
	if diffs[0].Key != "Findings" {
		t.Fatal("Unexpected diff key")
	}
	if diffs[0].Item != "{\"Value\":0,\"Show\":false,\"Message\":\"no test\"}" {
		t.Fatal("Unexpected diff item")
	}
	if diffs[1].Type != "Added" {
		t.Fatal("Unexpected diff type")
	}
	if diffs[1].Key != "Findings" {
		t.Fatal("Unexpected diff key")
	}
	if diffs[1].Item != "{\"Value\":1,\"Show\":true,\"Message\":\"test\"}" {
		t.Fatal("Unexpected diff item")
	}

	// nothing before report 1
	diffs, err = backingStore.SelectScenarioReportDiffs(scenarioID, "host-token", -100, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(diffs) != 0 {
		t.Fatal("Unexpected report diff count:", len(diffs))
	}

	// no diff between report 1 and 2
	diffs, err = backingStore.SelectScenarioReportDiffs(scenarioID, "host-token", 15, 30)
	if err != nil {
		t.Fatal(err)
	}
	if len(diffs) != 0 {
		t.Fatal("Unexpected report diff count:", len(diffs))
	}

	// diff between report 2 and 3
	diffs, err = backingStore.SelectScenarioReportDiffs(scenarioID, "host-token", 30, 45)
	if err != nil {
		t.Fatal(err)
	}
	if len(diffs) != 2 {
		t.Fatal("Unexpected report diff count:", len(diffs))
	}
	if diffs[0].Type != "Removed" {
		t.Fatal("Unexpected diff type")
	}
	if diffs[0].Key != "Findings" {
		t.Fatal("Unexpected diff key")
	}
	if diffs[0].Item != "{\"Value\":0,\"Show\":false,\"Message\":\"no test\"}" {
		t.Fatal("Unexpected diff item")
	}
	if diffs[1].Type != "Added" {
		t.Fatal("Unexpected diff type")
	}
	if diffs[1].Key != "Findings" {
		t.Fatal("Unexpected diff key")
	}
	if diffs[1].Item != "{\"Value\":1,\"Show\":true,\"Message\":\"test\"}" {
		t.Fatal("Unexpected diff item")
	}

	// no diff between report 3 and 4
	diffs, err = backingStore.SelectScenarioReportDiffs(scenarioID, "host-token", 45, 60)
	if err != nil {
		t.Fatal(err)
	}
	if len(diffs) != 0 {
		t.Fatal("Unexpected report diff count:", len(diffs))
	}

	// nothing after report 4
	diffs, err = backingStore.SelectScenarioReportDiffs(scenarioID, "host-token", 60, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(diffs) != 0 {
		t.Fatal("Unexpected report diff count:", len(diffs))
	}
}
