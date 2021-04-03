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
	err := createLinkReport("https://localhost:8443", path.Join(tempDir, ""))
	if err == nil {
		t.Fatal("Expected error for empty string report path")
	}

	// valid
	err = createLinkReport("https://localhost:8443", path.Join(tempDir, "report.html"))
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
}

func TestReadServerPubKey(t *testing.T) {
	// temp dir
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// empty string
	_, err := readServerPubKey(path.Join(tempDir, ""))
	if err == nil {
		t.Fatal("Expected error for empty string public key path")
	}

	// non-existent file
	_, err = readServerPubKey(path.Join(tempDir, "notafile"))
	if err == nil {
		t.Fatal("Expected error for non-existent public key path")
	}

	// empty file
	tempFile, err := ioutil.TempFile(tempDir, "")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	tempFile.Write([]byte(""))
	err = tempFile.Close()
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	_, err = readServerPubKey(tempFile.Name())
	if err == nil {
		t.Fatal("Expected error for empty public key file")
	}

	// invalid file
	tempFile, err = ioutil.TempFile(tempDir, "")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	tempFile.Write([]byte("invalid"))
	err = tempFile.Close()
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	_, err = readServerPubKey(tempFile.Name())
	if err == nil {
		t.Fatal("Expected error for invalid public key file")
	}

	// valid file
	tempFile, err = ioutil.TempFile(tempDir, "")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	tempFile.Write([]byte(`-----BEGIN PGP PUBLIC KEY BLOCK-----

xsBNBFwhp0MBCAC9V4PbAgk7Q3YPWe7Iz5eJAL0fx2aY45kQkHQlYlKxO9989yyV
c/+/cEWyn6viL8lERB3n2rNuVKEzffalKZSNONhvTqGoHlkz0c/Ty4YF9zAtE5b3
cXJYSLpy1Cndz+ygeaV18CXck802yYmhHB3YXpx6EKhVwMk4SIu9yD+wH9GWCrly
pIOshCgpvCccVvgBGP8meB5Zbfdwgrpk4nMzRcB+aR6oUxGkKCX3HfNWtY5BlzFJ
yJRUNlvuyZjFNK3vyjQwp/yMalpDPpDeGCslDUC9kaywnXMkaM+4tNcCJVL5sVnP
qJjdyPpNCSmbgRRVSgcFYw79wO/LH4dX3khjABEBAAHNJGNwLXNjb3JpbmcgKHRl
c3QpIDx0ZXN0QGV4YW1wbGUuY29tPsLAYgQTAQgAFgUCXCGnQwkQjMvLR6ssDYEC
GwMCGQEAABNbCABROvfLvSp9LICnzXFlI5QqfJzpmKF17VD3DUvpPnU8KpGPyx6F
sOgJ/tNh9i2nVLwCEaGg72AzvGjprw6qzQEvf1LyRDoqUYSp/Snxsra++4s/AQct
KZrTcXvUmgnjfPly6g/QnkqGCW6ujx17gvfAOq4gCYlOfitUCsos1+WpB52IuS4f
hwtrR21O8DfCgfFk2kaBxmpqWK66Dvuel97YVHVjpZn3zoFLx5f67HiCJiqFrq1J
xOrlz07QG44crQEPr8Q59d8hXg/hY3nzQACzURGjQaytZCAVl9Gg9Y90iWbL5v2W
jLvcN/jcOBBKElbhGYh+N0IWA2odYw5kjv5wzsBNBFwhp0MBCADVg5JlbzipMfuJ
+BLn7VCDQDTn2kfYR+m6qxcWbUQFfOmHl/6rq/bfsXbUB95ulAJ5CqtOYDROGLZH
2yV/abu7hrsVMEoGjY6eofbNBLCJibo60WSYcFZiyp81xeI2Sjo7/UE6byvAvLTK
EfeebKaK8rS2u/oQVVBeij/Lf6ZIdvpaE7RK/eSJ3gnZ6q0yuUILyAnw3qYC5PBN
+lYMRWAURifvg5ZDv+PkoAAncYLq3OyCkRtSULx3bjf0/UtGmUm7Aais2/mV6wzU
muzSVYfKPYLZPOKu+2QpU+/M3B6kD2iMxK5vLO3DjKMdvVN2ddQk/fP2k9luIFfu
SD94I5N1ABEBAAHCwF8EGAEIABMFAlwhp0MJEIzLy0erLA2BAhsMAAA5QQgAK0kY
bqXJzTZeNRxWiMwjMjKeVRbotUDR41MUn4tmL7yaRzVnGz+hKai3VyCLYj3j+XE+
IXDr8C9U+4IrIYqly7YpqpOGAmybD8JoOBSausUd3GDQojeoRWI19Bt7d6fgxzQ7
CYO4FqztYxUrpJiKEdmyYO3uHb8fkh3W92aQhqDXXuSl6LHiEsbGVAwygdtlbPUx
E0XihGYOonfSg4mvQItWwCsMYdLiL17HBaKFPj+bu1MBwFylfvmeMEPJidDWePED
f9NqngjkT7E0IrHx4DSS546QWJvxXFfiLuXk4YShMhfHNkUZC46s43oKEgdTnFax
J4M2XZdqSA3PFVF70g==
=p0MY
-----END PGP PUBLIC KEY BLOCK-----`))
	err = tempFile.Close()
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	pk, err := readServerPubKey(tempFile.Name())
	if err != nil {
		t.Fatal("Unexpected error reading valid public key file;", err)
	}
	if pk == nil {
		t.Fatal("Expected entity list to not be nil")
	}
}

func TestReadServerCert(t *testing.T) {
	// temp dir
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// empty string
	_, err := readServerCert(path.Join(tempDir, ""))
	if err == nil {
		t.Fatal("Expected error for empty string cert path")
	}

	// non-existent file
	_, err = readServerCert(path.Join(tempDir, "notafile"))
	if err == nil {
		t.Fatal("Expected error for non-existent cert path")
	}

	// empty file
	tempFile, err := ioutil.TempFile(tempDir, "")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	tempFile.Write([]byte(""))
	err = tempFile.Close()
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	_, err = readServerCert(tempFile.Name())
	if err == nil {
		t.Fatal("Expected error for empty cert file")
	}

	// invalid file
	tempFile, err = ioutil.TempFile(tempDir, "")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	tempFile.Write([]byte("invalid"))
	err = tempFile.Close()
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	_, err = readServerCert(tempFile.Name())
	if err == nil {
		t.Fatal("Expected error for invalid cert file")
	}

	// valid file
	tempFile, err = ioutil.TempFile(tempDir, "")
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	tempFile.Write([]byte(`-----BEGIN CERTIFICATE-----
MIIE+zCCAuOgAwIBAgIJAOSTiyZ4SQZlMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNV
BAMMCWxvY2FsaG9zdDAeFw0xODEyMjUwMzQzMDlaFw0xOTEyMjUwMzQzMDlaMBQx
EjAQBgNVBAMMCWxvY2FsaG9zdDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoC
ggIBALVGpeMe1ERPHVUeA53j4gRixB3SMJGdhN+rJyDuRltxt7wB00og8T4RpycJ
DsFDGrZ4JYSboGRl7jr9sNiTxxpOepgqXTN4dYDF6ks8KHhL3xTtDcX3gd1m2dmF
jKBlLruuAp4uc+wXjK3OmfTqpUsUN81AOkvtnl4IK5qBbfY/PsTXXUFmySG16FIj
jqZQJ0oVWr3XKn70Dr6VEUQ3lrjG/RBlMhcx7HwLEwFOgEAiEFeL4m3pMy1DKVeP
yLvPFeMkl/81oWCAoa0+14LOFwA3HF+kuXqziDNxpzB8EKhmzWuXbclsCCxCzrVR
7OGZ1qaix4xVcM02lY7e4EHXjriG00alYXuI6RGi05TYv9zx9ZYxFxqHMXwmam6F
8cc9NTx4vt/rHa8bFbd5fCqwUMIhbYpbQYqfkyefI2UgzeKzIdJhORcXDOrXgFfE
1+WJeaGAZ0/1oVTt2wJc3SaR5uW0+hHP1EF/tMMr6W5eZSX38ZbTr0ahwXjzc+sT
t+GLZ4SfnrJ+AEYuIxxT140583hj2hzQWB93qyut5V1fg1jqtg3WEQ1IaIU1LEQV
FG3b7zER1g9u44cLTy8pJ5nb+rGjX/wmMgJOOVaPe4jkEP4GEBfdIsQt4VzN1KH0
2ubfN8+8lxsOZlf2CD9kng/U7rrc94adPdNCqZ6WQl5ueCvPAgMBAAGjUDBOMB0G
A1UdDgQWBBQfUJpqwL5SXE78ZFR4QOkB+tTSKTAfBgNVHSMEGDAWgBQfUJpqwL5S
XE78ZFR4QOkB+tTSKTAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4ICAQAi
lrVdHeGvHnEHb4H/gyeiGpgOgh0Zj8WJL+vkfvb7ot403azAoOWs/eQIp89Bg5Xr
dRTJbWOSzyFyyiqmdc1QPP6Nasw5Q74riAM7hgz5nP+9hw0u2BAoADs7SVxO0Vtq
a5bVwaXv4r+IM0aSPrtzCKJucz4YQsS223rZdJ6BisWF8Yr5xLxwCrc1/LJjAEc+
xJWoZPrOTho7NP/Pp2eT0+LYD8qV5fsh2K8yWV7FbrwbKQOa2108DsGn0rLlBQLJ
4r2cJyHpyIgwfsMQQzgHQM1AtR1lu8KZ/fCJDO+AaiQDMMtt+ZPR3Io3DX8jEuZ5
mui8+1Pfj9yRhugbUAetKzPRG1Zm8jG+34ANTpqwk0HW5++muk/8Qwf5MDiDZjgm
Blmk6aGn0/U/NU+Xs0LoSBQ+C6EdiYrQVYSeN2yenzbHUOUFlbNEwdOuK+HXt0D1
vZua5sv5Xoda0uatUZc9JOYzTzsLA9SuT5XeNMczghmaEdE+++0car7sdVK+kGOL
aKTdRLhIFr6ZtucQrhkGMFwlyp88FVbVvrF0Uuewnq775XG8qKbInd58dVATTqUf
34uPYoMQGgMlXgu5aJnsZKKevWw0XWHihOKT2i6aqP4ifVUQs6XK7bkz+3tO8pzM
KlpJNAp2G84XZShvm2/w7kyRn78xKuqY4zc+1eh+WQ==
-----END CERTIFICATE-----`))
	err = tempFile.Close()
	if err != nil {
		t.Fatal("Unexpected error;", err)
	}
	c, err := readServerCert(tempFile.Name())
	if err != nil {
		t.Fatal("Unexpected error reading valid cert file;", err)
	}
	if c == nil {
		t.Fatal("Expected cert to not be nil")
	}
	if len(c.TLSClientConfig.RootCAs.Subjects()) != 1 {
		t.Fatal("Expected 1 root CA subject")
	}
}
