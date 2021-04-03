package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
)

func copyFile(srcPath string, dstPath string) {
	src, err := os.Open(srcPath)
	if err != nil {
		log.Fatalln("Unable to open source file;", err)
	}
	defer src.Close()
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Fatalln("Unable to open destination file;", err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		log.Fatalln("Unable to copy file;", err)
	}
}

func createDir(dir string) {
	// data directory
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		log.Fatalln("Unable to set up directory "+dir+";", err)
	}
}

func getCurrentHost() (currentHost, error) {
	if runtime.GOOS == "linux" {
		return hostLinux{}, nil
	} else if runtime.GOOS == "windows" {
		return hostWindows{}, nil
	}
	return nil, errors.New("ERROR: unsupported platform: " + runtime.GOOS)

}

func writeReadmeHTML(dir string, serverURL string) error {
	outFile := path.Join(dir, "README.html")
	log.Println("Creating " + outFile)
	url := serverURL + "/ui/team-dashboard"
	s := "<html><head><meta http-equiv=\"refresh\" content=\"0; url=" + url + "\"></head><body><a href=\"" + url + "\">Team Dashboard</a></body></html>"
	err := ioutil.WriteFile(outFile, []byte(s), 0644)
	if err != nil {
		log.Println("ERROR: unable to save " + outFile)
		return err
	}
	log.Println("Saved " + outFile)
	return nil
}
