// +build windows

package main

import (
	"bytes"
	"encoding/csv"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

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

func getUsers() []model.User {
	users := make([]model.User, 0)

	// get users and info
	out, err := exec.Command("C:\\Windows\\System32\\wbem\\WMIC.exe", "UserAccount", "get", "Name,SID").Output()
	if err != nil {
		log.Fatal("ERROR: cannot get users info;", err)
	}
	var posName int
	var posSID int
	for i, line := range strings.Split(string(out), "\r\n") {
		if len(line) <= 1 {
			continue
		}

		var user model.User

		// parse header for positions
		if i == 0 {
			// assume all exist
			posName = strings.Index(line, "Name")
			posSID = strings.Index(line, "SID")
			continue
		}

		user.AccountPresent = true
		user.Name = strings.TrimSpace(line[posName:posSID])
		user.ID = strings.TrimSpace(line[posSID:])

		// use net user to get additional information
		out, err := exec.Command("C:\\Windows\\System32\\net.exe", "user", user.Name).Output()
		if err != nil {
			log.Println("ERROR: cannot get user info;", err)
			continue
		}
		for _, line := range strings.Split(string(out), "\r\n") {
			if len(line) == 0 {
				continue
			}
			// assume always fixed position
			value := strings.TrimSpace(line[29:])
			if strings.HasPrefix(line, "Account active") {
				if value == "Yes" {
					user.AccountActive = true
				} else {
					user.AccountActive = false
				}
			} else if strings.HasPrefix(line, "Account expires") {
				if value == "Never" {
					user.AccountExpires = false
				} else {
					user.AccountExpires = true
				}
			} else if strings.HasPrefix(line, "Password last set") {
				// add timezone to value
				timezone, _ := time.Now().Zone()
				value = value + " " + timezone
				layout := "1/2/2006 3:04:05 PM MST"
				t, err := time.Parse(layout, value)
				if err != nil {
					log.Println("ERROR: cannot parse date time string;", err)
				}
				user.PasswordLastSet = t.Unix()
			} else if strings.HasPrefix(line, "Password expires") {
				if value == "Never" {
					user.PasswordExpires = false
				} else {
					user.PasswordExpires = true
				}
			}
			continue
		}

		users = append(users, user)
	}
	return users
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
	// get user and PID
	pidUserMap := make(map[int64]string)
	out, err := exec.Command("C:\\Windows\\System32\\tasklist.exe", "/V", "/FO:csv").Output()
	if err != nil {
		log.Fatal("ERROR: cannot get process list;", err)
	}

	r := csv.NewReader(bytes.NewReader(out))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal("ERROR: cannot read process list csv; ", err)
	}
	for _, tokens := range records {
		pid, _ := strconv.ParseInt(tokens[1], 10, 64)
		user := tokens[6]
		// local account, remove hostname
		hostname, _ := os.Hostname()
		if strings.Index(user, hostname) != -1 {
			user = user[len(hostname)+1:]
		}
		pidUserMap[pid] = user
	}

	// get PID and command line
	processes := make([]model.Process, 0)
	var posCommandLine int
	var posProcessID int
	out, err = exec.Command("C:\\Windows\\System32\\wbem\\WMIC.exe", "process", "get", "CommandLine,ProcessId").Output()
	if err != nil {
		log.Fatal("ERROR: cannot get process list;", err)
	}

	for i, line := range strings.Split(string(out), "\r\n") {
		if len(line) <= 1 {
			continue
		}

		// get position of columns
		if i == 0 {
			posCommandLine = strings.Index(line, "CommandLine")
			posProcessID = strings.Index(line, "ProcessId")
			continue
		}

		var process model.Process
		process.CommandLine = strings.TrimSpace(line[posCommandLine:posProcessID])
		pid, _ := strconv.ParseInt(strings.TrimSpace(line[posProcessID:]), 10, 64)
		process.PID = pid
		user, _ := pidUserMap[pid]
		process.User = user
		processes = append(processes, process)
	}

	return processes
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
	out, err := exec.Command("C:\\Windows\\System32\\NETSTAT.exe", "-nao").Output()
	if err != nil {
		log.Fatal("ERROR: cannot get network connections;", err)
	}

	conns := make([]model.NetworkConnection, 0)
	for i, line := range strings.Split(string(out), "\r\n") {
		if len(line) <= 1 {
			continue
		}

		// skip header lines
		if i < 4 {
			continue
		}

		tokens := strings.Fields(line)

		var conn model.NetworkConnection
		conn.Protocol = strings.ToUpper(tokens[0])
		localAddrStr := tokens[1]
		lastColon := strings.LastIndex(localAddrStr, ":")
		if lastColon == -1 {
			conn.LocalAddress = localAddrStr
		} else {
			conn.LocalAddress = localAddrStr[0:lastColon]
			conn.LocalPort = localAddrStr[lastColon+1:]
		}
		remoteAddrStr := tokens[2]
		lastColon = strings.LastIndex(remoteAddrStr, ":")
		if lastColon == -1 {
			conn.RemoteAddress = remoteAddrStr
		} else {
			conn.RemoteAddress = remoteAddrStr[0:lastColon]
			conn.RemotePort = remoteAddrStr[lastColon+1:]
		}
		// five tokens has all columns available
		// four tokens is probably missing state (e.g. UDP)
		if len(tokens) == 5 {
			conn.State = model.GetNetworkConnectionState(tokens[3])
		}
		pidStr := tokens[len(tokens)-1]
		pid, err := strconv.ParseInt(pidStr, 10, 64)
		if err != nil {
			log.Fatal("ERROR: cannot parse PID;", err)
		}
		conn.PID = pid
	}

	return conns
}
