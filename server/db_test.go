// +build integration

package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
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

		tbs, err := getBackingStore("postgres", sqlURL)
		if err != nil {
			return testDBObj, err
		}
		testDBObj = tbs
	}
	return testDBObj, nil
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
		log.Print(err)
		t.Fatal("Unexpected error")
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
			log.Fatal(err)
		}
		if readState.Hostname != "test" {
			log.Fatal("Unexpected hostname from state")
		}
		counter++
	}
	if counter != 1 {
		t.Fatal("Unexpected number of rows:", counter)
	}
}
