// +build windows

package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/sumwonyuno/cp-scoring/model"
	"golang.org/x/sys/windows/registry"
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
	out, err := exec.Command("C:\\Windows\\System32\\wbem\\WMIC.exe", "path", "win32_groupuser").Output()
	if err != nil {
		log.Fatal("ERROR: unable to get group users;", err)
	}

	groups := make(map[string][]string)
	var posGroupComponent int
	var posPartComponent int
	for i, line := range strings.Split(string(out), "\r\n") {
		if len(line) <= 1 {
			continue
		}

		// find positions of columns
		if i == 0 {
			// assume these exist
			posGroupComponent = strings.Index(line, "GroupComponent")
			posPartComponent = strings.Index(line, "PartComponent")
			continue
		}

		// parse out group and member
		groupComponentStr := strings.TrimSpace(line[posGroupComponent:posPartComponent])
		groupComponentStr = strings.Split(groupComponentStr, ",")[1]
		group := groupComponentStr[6 : len(groupComponentStr)-1]
		partComponentStr := strings.TrimSpace(line[posPartComponent:])
		partComponentStr = strings.Split(partComponentStr, ",")[1]
		member := partComponentStr[6 : len(partComponentStr)-1]
		g, present := groups[group]
		if !present {
			g = make([]string, 0)
		}
		g = append(g, member)
		groups[group] = g
	}

	return groups
}

func getProcesses() []model.Process {
	out := powershellCsv("Get-Process -IncludeUserName", "ID,UserName,Path")
	return parseWindowsProcesses(out)
}

func getSoftware() []model.Software {
	// based on add/remove programs
	// first location
	rights := uint32(registry.QUERY_VALUE | registry.ENUMERATE_SUB_KEYS | registry.QUERY_VALUE)
	loc1Path := "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall"
	loc1, err := registry.OpenKey(registry.LOCAL_MACHINE, loc1Path, rights)
	if err != nil {
		log.Fatal("ERROR: unable to get software;", err)
	}
	defer loc1.Close()
	// second location
	loc2Path := "SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall"
	loc2, err := registry.OpenKey(registry.LOCAL_MACHINE, loc2Path, rights)
	if err != nil {
		log.Fatal("ERROR: unable to get software;", err)
	}

	software := make([]model.Software, 0)

	// first location
	subkeys, err := loc1.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	for _, subkey := range subkeys {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, loc1Path+"\\"+subkey, rights)
		if err != nil {
			log.Fatal("ERROR: cannot read key;", err)
		}
		vn, err := key.ReadValueNames(-1)
		if err != nil {
			log.Fatal("ERROR: cannot read key value names;", err)
		}
		if len(vn) == 0 {
			continue
		}
		name, _, err := key.GetStringValue("DisplayName")
		if err != nil {
			continue
		}
		ver, _, err := key.GetStringValue("DisplayVersion")
		if err != nil {
			continue
		}
		var sw model.Software
		sw.Name = name
		sw.Version = ver
		software = append(software, sw)
	}
	// second location
	subkeys, err = loc2.ReadSubKeyNames(-1)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	for _, subkey := range subkeys {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, loc2Path+"\\"+subkey, rights)
		if err != nil {
			log.Fatal("ERROR: cannot read key;", err)
		}
		vn, _ := key.ReadValueNames(-1)
		if err != nil {
			log.Fatal("ERROR: cannot read key value names;", err)
		}
		if len(vn) == 0 {
			continue
		}
		name, _, err := key.GetStringValue("DisplayName")
		if err != nil {
			continue
		}
		ver, _, err := key.GetStringValue("DisplayVersion")
		if err != nil {
			continue
		}
		var sw model.Software
		sw.Name = name
		sw.Version = ver
		software = append(software, sw)
	}

	return software
}

func getNetworkConnections() []model.NetworkConnection {
	out := powershellCsv("Get-NetTCPConnection", "OwningProcess,State,LocalAddress,LocalPort,RemoteAddress,RemotePort")
	tcpConns := parseWindowsTCPNetConns(out)

	out = powershellCsv("Get-NetUDPEndpoint", "OwningProcess,LocalAddress,LocalPort")
	udpConns := parseWindowsUDPNetConns(out)

	return append(tcpConns, udpConns...)
}
