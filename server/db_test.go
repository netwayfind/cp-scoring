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
	backingStore, err := getTestBackingStore()
	if err != nil {
		t.Fatal(err)
	}
	err = clearTables()
	if err != nil {
		t.Fatal(err)
	}

	state := model.State{Hostname: "test"}
	stateBytes, _ := json.Marshal(state)

	// state, no existing host token
	err = backingStore.InsertState(1000, "127.0.0.1", "host-token", stateBytes)
	if err == nil {
		t.Fatal("Expected error due to missing host token")
	}

	// state, existing host token
	backingStore.InsertHostToken("host-token", 1000, "host", "127.0.0.1")
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
	backingStore, err := getTestBackingStore()
	if err != nil {
		t.Fatal(err)
	}
	err = clearTables()
	if err != nil {
		t.Fatal(err)
	}

	// insert score, no existing scenario
	err = backingStore.InsertScenarioScore(model.ScenarioHostScore{})
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
	var readScenarioID int64
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
	backingStore, err := getTestBackingStore()
	if err != nil {
		t.Fatal(err)
	}
	err = clearTables()
	if err != nil {
		t.Fatal(err)
	}

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
	err = backingStore.InsertTeamHostToken(team1ID, "host1", "host1_1", 0)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertTeamHostToken(team1ID, "host1", "host1_2", 0)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertTeamHostToken(team1ID, "host2", "host2_1", 0)
	if err != nil {
		t.Fatal(err)
	}
	err = backingStore.InsertTeamHostToken(team2ID, "host1", "host1_3", 0)
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
	backingStore, err := getTestBackingStore()
	if err != nil {
		t.Fatal(err)
	}
	err = clearTables()
	if err != nil {
		t.Fatal(err)
	}

	// sample report
	findings := append(make([]model.Finding, 0), model.Finding{Show: true, Message: "test", Value: 1})
	report := model.Report{Timestamp: 1500, Findings: findings}

	// insert report without scenario and host token
	err = backingStore.InsertScenarioReport(-1, "host-token", report)
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
	var readScenarioID int64
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
	backingStore, err := getTestBackingStore()
	if err != nil {
		t.Fatal(err)
	}
	err = clearTables()
	if err != nil {
		t.Fatal(err)
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
	backingStore, err := getTestBackingStore()
	if err != nil {
		t.Fatal(err)
	}
	err = clearTables()
	if err != nil {
		t.Fatal(err)
	}

	// no existing data
	timeline, err := backingStore.SelectScenarioTimeline(-1, "host-token1")
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
