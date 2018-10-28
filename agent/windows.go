// +build windows

package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/sumwonyuno/cp-scoring/model"
)

func getState() model.State {
	state := model.GetNewStateTemplate()
	state.Users = getUsers()
	state.Groups = getGroups()
	state.Processes = getProcesses()
	state.Software = getSoftware()
	state.NetworkConnections = getNetworkConnections()
	return state
}

type userinfo struct {
	username string
}

func powershellCsv(command string, columns string) []byte {
	cmdStr := fmt.Sprintf("%s | Select-Object %s | ConvertTo-Csv -NoTypeInformation", command, columns)
	out, err := exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "-command", cmdStr).Output()
	if err != nil {
		log.Println("ERROR: unable to execute powershell command;", err)
	}
	return out
}

func getUsers() []model.User {
	out := powershellCsv("Get-LocalUser", "Name,SID,Enabled,AccountExpires,PasswordLastSet,PasswordExpires")
	return parseWindowsUsers(out)
}

func getGroups() map[string][]string {
	out := powershellCsv("Get-WmiObject -class Win32_GroupUser", "GroupComponent,PartComponent")
	return parseWindowsGroups(out)
}

func getProcesses() []model.Process {
	out := powershellCsv("Get-Process -IncludeUserName", "ID,UserName,Path")
	return parseWindowsProcesses(out)
}

func getSoftware() []model.Software {
	// check two locations for software in registry
	loc1 := powershellCsv("Get-ItemProperty HKLM:SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall", "DisplayName,DisplayVersion")
	sw1 := parseWindowsSoftware(loc1)
	loc2 := powershellCsv("Get-ItemProperty HKLM:SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall", "DisplayName,DisplayVersion")
	sw2 := parseWindowsSoftware(loc2)

	return append(sw1, sw2...)
}

func getNetworkConnections() []model.NetworkConnection {
	out := powershellCsv("Get-NetTCPConnection", "OwningProcess,State,LocalAddress,LocalPort,RemoteAddress,RemotePort")
	tcpConns := parseWindowsTCPNetConns(out)

	out = powershellCsv("Get-NetUDPEndpoint", "OwningProcess,LocalAddress,LocalPort")
	udpConns := parseWindowsUDPNetConns(out)

	return append(tcpConns, udpConns...)
}
