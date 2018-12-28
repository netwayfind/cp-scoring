package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func createTempDir(t *testing.T) string {
	dir := "test_temp"
	os.Mkdir(dir, 0700)
	tempDir, err := ioutil.TempDir(dir, "")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	return tempDir
}

func TestGetBackingStore(t *testing.T) {
	// empty string
	_, err := getBackingStore("")
	if err == nil {
		t.Fatal("Expected error for empty backing store string")
	}

	// invalid
	_, err = getBackingStore("invalid")
	if err == nil {
		t.Fatal("Expected error for invalid backing store string")
	}

	// temp dir
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// sqlite, missing args
	bs, err := getBackingStore("sqlite")
	if err == nil {
		t.Fatal("Expected error for missing args")
	}

	// sqlite, given args
	bs, err = getBackingStore("sqlite", tempDir)
	if err != nil {
		t.Fatal("Unexpected error for sqlite backing store;", err)
	}
	if bs == nil {
		t.Fatal("Unexpected nil backing store")
	}
}
