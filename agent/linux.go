package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type hostLinux struct {
}

func (h hostLinux) copyTeamFiles() error {
	installPath := InstallPathLinux
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("ERROR: cannot get current directory")
		return err
	}
	log.Println("Copying files to: " + currentDir)

	// readme
	fileName := FileReadmeHTML
	copyFile(filepath.Join(installPath, fileName), filepath.Join(currentDir, fileName))

	// team setup shortcut
	fileName = "team_setup.desktop"
	copyFile(filepath.Join(installPath, fileName), filepath.Join(currentDir, fileName))
	os.Chmod(filepath.Join(currentDir, fileName), 0755)

	log.Println("Finished copying files")
	return nil
}

func (h hostLinux) install() error {
	installPath := InstallPathLinux

	// create installation folder
	err := os.MkdirAll(installPath, 0755)
	if err != nil {
		log.Println("ERROR: unable to create installation folder " + installPath)
		return err
	}
	log.Println("Created installation folder: " + installPath)

	// copy agent
	log.Println("Copying this executable to installation folder")
	ex, err := os.Executable()
	if err != nil {
		log.Println("ERROR: unable to copy executable")
		return err
	}
	binFile := filepath.Join(installPath, FileAgentLinux)
	copyFile(ex, binFile)
	err = os.Chmod(binFile, 0755)
	if err != nil {
		log.Println("ERROR: unable to set file permissions")
		return err
	}

	// create service
	log.Println("Creating service")
	serviceFile := filepath.Join(installPath, "cp-scoring.service")
	err = ioutil.WriteFile(serviceFile, getSystemdScript(), 0755)
	if err != nil {
		log.Println("ERROR: unable to write service file")
		return err
	}
	// delete existing service and recreate
	exec.Command("/bin/systemctl", "disable", "cp-scoring.service").Run()
	err = exec.Command("/bin/systemctl", "enable", serviceFile).Run()
	if err != nil {
		log.Println("ERROR: unable to enable service")
		return err
	}

	// create team setup
	log.Println("Creating team setup script")
	// script
	fileReg := filepath.Join(installPath, "team_setup.sh")
	text := []byte("#!/bin/sh\ncd /opt/cp-scoring\nsudo ./cp-scoring-agent-linux -team_setup")
	err = ioutil.WriteFile(fileReg, text, 0755)
	if err != nil {
		log.Println("ERROR: unable to save team setup script")
		return err
	}
	// shortcut
	fileShortcut := filepath.Join(installPath, "team_setup.desktop")
	text = []byte("[Desktop Entry]\nEncoding=UTF-8\nVersion=1.0\nName[en_US]=Team Key Registration\nExec=/opt/cp-scoring/team_setup.sh\nTerminal=true\nType=Application")
	err = ioutil.WriteFile(fileShortcut, text, 0755)
	if err != nil {
		log.Println("ERROR: unable to save team setup shortcut")
		return err
	}

	log.Println("Finished install")
	return nil
}

func getSystemdScript() []byte {
	return []byte(`[Unit]
Description=cp-scoring

[Service]
User=root
Group=root
WorkingDirectory=/opt/cp-scoring
ExecStart=/opt/cp-scoring/cp-scoring-agent-linux
Restart=on-failure

[Install]
WantedBy=multi-user.target
Alias=cp-scoring.service
`)
}
