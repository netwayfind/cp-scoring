package main

import (
	"io/ioutil"
	"os"
	"path"
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

func TestGetServerURL(t *testing.T) {
	// empty string
	_, err := getServerURL("")
	if err == nil {
		t.Fatal("Expected error for empty server file path")
	}

	// non-existent file
	_, err = getServerURL("not_here")
	if err == nil {
		t.Fatal("Expected error for non-existent server file path")
	}

	// temp dir
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// not a URL
	tempFile, err := ioutil.TempFile(tempDir, "")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	tempFile.Write([]byte("notaURL"))
	err = tempFile.Close()
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	_, err = getServerURL(tempFile.Name())
	if err == nil {
		t.Fatal("Expected error for not a URL")
	}

	// valid URL
	tempFile, err = ioutil.TempFile(tempDir, "")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	tempFile.Write([]byte("https://localhost:8443"))
	err = tempFile.Close()
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	u, err := getServerURL(tempFile.Name())
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	if u != "https://localhost:8443" {
		t.Fatal("Unexpected URL " + u)
	}
}

func TestCheckValidServerURL(t *testing.T) {
	// empty string
	err := checkValidServerURL("")
	if err == nil {
		t.Fatal("Expected error for empty server URL")
	}

	// bad string
	err = checkValidServerURL("bad")
	if err == nil {
		t.Fatal("Expected error for bad server URL")
	}

	// must be HTTPS
	err = checkValidServerURL("http://localhost:8443")
	if err == nil {
		t.Fatal("Expected error for unacceptable server URL")
	}

	// good URL
	err = checkValidServerURL("https://localhost:8443")
	if err != nil {
		t.Fatal("Expected no errors")
	}
}

func TestCreateLinkScoreboard(t *testing.T) {
	// temp dir
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// empty string
	err := createLinkScoreboard("https://localhost:8443", path.Join(tempDir, ""))
	if err == nil {
		t.Fatal("Expected error for empty string scoreboard path")
	}

	// valid
	err = createLinkScoreboard("https://localhost:8443", path.Join(tempDir, "scoreboard.html"))
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
}

func TestCreateLinkReport(t *testing.T) {
	// temp dir
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// empty path
	err := createLinkReport("https://localhost:8443", path.Join(tempDir, ""), "hostname", "")
	if err == nil {
		t.Fatal("Expected error for empty string report path")
	}

	// empty hostname
	err = createLinkReport("https://localhost:8443", path.Join(tempDir, "report.html"), "", "host_token")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}

	// empty host token
	err = createLinkReport("https://localhost:8443", path.Join(tempDir, "report.html"), "hostname", "")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}

	// non-empty host token
	err = createLinkReport("https://localhost:8443", path.Join(tempDir, "report.html"), "hostname", "host_token")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
}
