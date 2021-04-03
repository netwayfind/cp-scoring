package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type hostWindows struct {
}

func (h hostWindows) copyTeamFiles() error {
	installPath := InstallPathWindows
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
	fileName = "Team Setup.bat"
	copyFile(filepath.Join(installPath, fileName), filepath.Join(currentDir, fileName))

	log.Println("Finished copying files")
	return nil
}

func (h hostWindows) install() error {
	installPath := InstallPathWindows

	// create installation folder
	err := os.MkdirAll(installPath, os.ModeDir)
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
	binFile := filepath.Join(installPath, FileAgentWindows)
	copyFile(ex, binFile)
	err = os.Chmod(binFile, 0755)
	if err != nil {
		log.Println("ERROR: unable to set file permissions")
		return err
	}

	// create Task Scheduler file
	log.Println("Creating Task Scheduler task")
	fileTaskSched := filepath.Join(installPath, "task.xml")
	err = ioutil.WriteFile(fileTaskSched, getScheduledTaskXML(), 0600)
	if err != nil {
		log.Println("ERROR: could not write task scheduler file")
		return err
	}
	// delete existing task and recreate
	exec.Command("C:\\Windows\\system32\\schtasks.exe", "/delete", "/F", "/tn", "cp-scoring").Run()
	err = exec.Command("C:\\Windows\\system32\\schtasks.exe", "/create", "/xml", fileTaskSched, "/tn", "cp-scoring").Run()
	if err != nil {
		log.Println("ERROR: unable to schedule task")
		return err
	}

	// create team setup
	log.Println("Creating team setup script")
	fileReg := filepath.Join(installPath, "Team Setup.bat")
	text := []byte("cd C:\\cp-scoring\r\ncp-scoring-agent-windows.exe -team_setup")
	err = ioutil.WriteFile(fileReg, text, 0600)
	if err != nil {
		log.Println("ERROR: unable to save team setup script")
		return err
	}

	log.Println("Finished install")
	return nil
}

func getScheduledTaskXML() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Date>2018-12-12T00:00:00.000</Date>
    <Author>WIN8\cyberpatriot</Author>
    <Description>cp-scoring. Do not delete or disable.</Description>
  </RegistrationInfo>
  <Triggers>
    <BootTrigger>
      <Enabled>true</Enabled>
    </BootTrigger>
  </Triggers>
  <Principals>
    <Principal id="Author">
      <UserId>S-1-5-18</UserId>
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>true</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>false</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>true</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>false</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>false</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>C:\cp-scoring\cp-scoring-agent-windows.exe</Command>
    </Exec>
  </Actions>
</Task>`)
}
