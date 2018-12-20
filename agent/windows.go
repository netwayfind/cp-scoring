package agent

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sumwonyuno/cp-scoring/model"
)

type hostWindows struct {
	PowerShellVersion string
}

func powershellCsv(command string, columns string) ([]byte, error) {
	cmdStr := fmt.Sprintf("%s | Select-Object %s | ConvertTo-Csv -NoTypeInformation", command, columns)
	out, err := exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "-command", cmdStr).Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (host hostWindows) GetUsers() ([]model.User, error) {
	// powershell 5.1 required
	if host.PowerShellVersion == "5.1" {
		out, err := powershellCsv("Get-LocalUser", "Name,SID,Enabled,AccountExpires,PasswordLastSet,PasswordExpires")
		if err != nil {
			return nil, err
		}
		return parseWindowsUsersGetLocalUser(out), nil
	}
	// otherwise, use older command
	out, err := powershellCsv("Get-WmiObject -class Win32_UserAccount", "Name,SID")
	if err != nil {
		return nil, err
	}
	parsedUsers := parseWindowsUsersWin32UserAccount(out)
	users := make([]model.User, 0)
	// need to get further info
	for _, user := range parsedUsers {
		out, err := exec.Command("C:\\Windows\\System32\\net.exe", "user", user.Name).Output()
		if err != nil {
			return users, err
		}
		// merge user info
		users = append(users, mergeNetUser(user, parseWindowsNetUser(out)))
	}
	return users, nil
}

func mergeNetUser(original model.User, new model.User) model.User {
	original.AccountActive = new.AccountActive
	original.AccountExpires = new.AccountExpires
	original.PasswordExpires = new.PasswordExpires
	original.PasswordLastSet = new.PasswordLastSet
	return original
}

func getPowerShellVersion() (string, error) {
	out, err := powershellCsv("Get-Host", "Version")
	if err != nil {
		return "", err
	}
	return parsePowerShellVersion(out), nil
}

func (host hostWindows) GetGroups() (map[string][]string, error) {
	out, err := powershellCsv("Get-WmiObject -class Win32_GroupUser", "GroupComponent,PartComponent")
	if err != nil {
		return nil, err
	}
	return parseWindowsGroups(out), nil
}

func (host hostWindows) GetProcesses() ([]model.Process, error) {
	out, err := powershellCsv("Get-Process -IncludeUserName", "ID,UserName,Path")
	if err != nil {
		return nil, err
	}
	return parseWindowsProcesses(out), nil
}

func (host hostWindows) GetSoftware() ([]model.Software, error) {
	// check two locations for software in registry
	loc1, err := powershellCsv("Get-ItemProperty HKLM:SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall", "DisplayName,DisplayVersion")
	if err != nil {
		return nil, err
	}
	sw1 := parseWindowsSoftware(loc1)
	loc2, err := powershellCsv("Get-ItemProperty HKLM:SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall", "DisplayName,DisplayVersion")
	if err != nil {
		return nil, err
	}
	sw2 := parseWindowsSoftware(loc2)

	return append(sw1, sw2...), nil
}

func (host hostWindows) GetNetworkConnections() ([]model.NetworkConnection, error) {
	out, err := powershellCsv("Get-NetTCPConnection", "OwningProcess,State,LocalAddress,LocalPort,RemoteAddress,RemotePort")
	if err != nil {
		return nil, err
	}
	tcpConns := parseWindowsTCPNetConns(out)

	out, err = powershellCsv("Get-NetUDPEndpoint", "OwningProcess,LocalAddress,LocalPort")
	if err != nil {
		return nil, err
	}
	udpConns := parseWindowsUDPNetConns(out)

	return append(tcpConns, udpConns...), nil
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
      <UserId>S-1-5-19</UserId>
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

func copyAgent(installPath string) {
	log.Println("Copying this executable to installation folder")
	ex, err := os.Executable()
	if err != nil {
		log.Println("Unable to get this executable path;", err)
	}
	fileIn, err := os.Open(ex)
	if err != nil {
		log.Fatalln("Unable to open self file;", err)
	}
	defer fileIn.Close()
	fileExe := filepath.Join(installPath, "cp-scoring-agent-windows.exe")
	fileOut, err := os.Create(fileExe)
	if err != nil {
		log.Fatalln("Unable to open destination file;", err)
	}
	defer fileOut.Close()
	_, err = io.Copy(fileOut, fileIn)
	if err != nil {
		log.Fatalln("Unable to copy file;", err)
	}
}

func createScheduledTask(installPath string) {
	log.Println("Creating Task Scheduler task")
	fileTaskSched := filepath.Join(installPath, "task.xml")
	err := ioutil.WriteFile(fileTaskSched, getScheduledTaskXML(), 0600)
	if err != nil {
		log.Fatalln("Could not write Task Scheduler file")
	}
	// delete existing task and recreate
	exec.Command("C:\\Windows\\system32\\schtasks.exe", "/delete", "/F", "/tn", "cp-scoring").Run()
	err = exec.Command("C:\\Windows\\system32\\schtasks.exe", "/create", "/xml", fileTaskSched, "/tn", "cp-scoring").Run()
	if err != nil {
		log.Fatalln("Unable to load task;", err)
	}
}

func (host hostWindows) Install() {
	installPath := "C:\\cp-scoring"

	// create installation folder
	os.MkdirAll(installPath, os.ModeDir)
	log.Println("Created installation folder: " + installPath)

	// copy agent
	copyAgent(installPath)

	// create Task Scheduler file
	createScheduledTask(installPath)

	log.Println("Finished installing to " + installPath)
}
