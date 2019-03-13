// +build integration

package main

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

var testDBObj backingStore

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

		tbs, err := getBackingStore("postgres", sqlURL)
		if err != nil {
			return testDBObj, err
		}
		testDBObj = tbs
	}
	return testDBObj, nil
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
