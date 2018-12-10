package agent

import (
	"fmt"
	"os/exec"

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
	out, err := powershellCsv("Get-LocalUser", "Name,SID,Enabled,AccountExpires,PasswordLastSet,PasswordExpires")
	if err != nil {
		return nil, err
	}
	return parseWindowsUsers(out), nil

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
