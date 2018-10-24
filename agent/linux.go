// +build linux

package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"

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

func getUsers() []model.User {
	// get user and uid
	bs, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		log.Fatal("ERROR: cannot get users info;", err)
	}
	userMapEtcPasswd := parseEtcPasswd(bs)

	// get other user information (sensitive)
	bs, err = ioutil.ReadFile("/etc/shadow")
	if err != nil {
		log.Fatal("ERROR: cannot get users info;", err)
	}
	userMapEtcShadow := parseEtcShadow(bs)

	return mergeUserMaps(userMapEtcPasswd, userMapEtcShadow)
}

func getGroups() map[string][]string {
	bs, err := ioutil.ReadFile("/etc/group")
	if err != nil {
		log.Fatal("ERROR: cannot get groups;", err)
	}

	groups := make(map[string][]string)
	for _, line := range strings.Split(string(bs), "\n") {
		tokens := strings.Split(line, ":")
		if len(tokens) != 4 {
			continue
		}
		group, membersStr := tokens[0], tokens[3]
		groups[group] = strings.Split(membersStr, ",")
	}

	return groups
}

func getProcesses() []model.Process {
	out, err := exec.Command("/bin/ps", "-eo", "pid,user,command", "--sort=pid").Output()
	if err != nil {
		log.Fatal("ERROR: cannot get processes;", err)
	}

	processes := make([]model.Process, 0)
	var posUser int
	var posPID int
	var posCommand int
	for i, line := range strings.Split(string(out), "\n") {
		if len(line) == 0 {
			continue
		}

		// get positions of columns
		if i == 0 {
			// PID is kept at 0, column is right justified
			posUser = strings.Index(line, "USER")
			posCommand = strings.Index(line, "COMMAND")
			continue
		}

		var process model.Process
		process.PID, _ = strconv.ParseInt(strings.TrimSpace(line[posPID:posUser]), 10, 64)
		process.User = strings.TrimSpace(line[posUser:posCommand])
		process.CommandLine = strings.TrimSpace(line[posCommand:])
		processes = append(processes, process)
	}

	return processes
}

func getSoftware() []model.Software {
	software := make([]model.Software, 0)

	out, err := exec.Command("/usr/bin/apt", "list", "--installed").Output()
	if err != nil {
		log.Fatal("ERROR: unable to get software list;", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		if len(line) == 0 {
			continue
		}

		tokens := strings.Split(line, " ")
		if len(tokens) != 4 {
			continue
		}
		pkgStr, version := tokens[0], tokens[1]
		pkg := strings.Split(pkgStr, "/")[0]
		var sw model.Software
		sw.Name = pkg
		sw.Version = version
		software = append(software, sw)
	}

	return software
}

func getNetworkConnections() []model.NetworkConnection {
	conns := make([]model.NetworkConnection, 0)

	out, err := exec.Command("/bin/sh", "-c", "/bin/ss -anptu | column -t").Output()
	if err != nil {
		log.Fatal("ERROR: cannot get network connections;", err)
	}

	var posNetid int
	var posState int
	var posRecv int
	var posLocal int
	var posAddrPort1 int
	var posPeer int
	for i, line := range strings.Split(string(out), "\n") {
		if len(line) == 0 {
			continue
		}

		// get position of columns
		if i == 0 {
			posNetid = strings.Index(line, "Netid")
			posState = strings.Index(line, "State")
			posRecv = strings.Index(line, "Recv-Q")
			posLocal = strings.Index(line, "Local")
			posAddrPort1 = strings.Index(line, "Address:Port")
			posPeer = strings.Index(line, "Peer")
			continue
		}

		var conn model.NetworkConnection
		// some connections may not have an PID (e.g. time wait)
		if posPeer >= len(line) {
			posPeer = len(line)
			conn.PID = -1
		} else {
			pidStr := strings.Split(strings.TrimSpace(line[posPeer:]), ",")[1]
			conn.PID, err = strconv.ParseInt(strings.TrimSpace(pidStr[4:]), 10, 64)
			if err != nil {
				log.Fatal("ERROR: unable to parse PID from network connection;", err)
			}
		}
		conn.Protocol = strings.ToUpper(strings.TrimSpace(line[posNetid:posState]))
		conn.State = model.GetNetworkConnectionState(strings.TrimSpace(line[posState:posRecv]))
		localAddrStr := strings.TrimSpace(line[posLocal:posAddrPort1])
		lastColon := strings.LastIndex(localAddrStr, ":")
		if lastColon == -1 {
			conn.LocalAddress = localAddrStr
		} else {
			conn.LocalAddress = localAddrStr[0:lastColon]
			conn.LocalPort = localAddrStr[lastColon+1:]
		}
		remoteAddrStr := strings.TrimSpace(line[posAddrPort1:posPeer])
		lastColon = strings.LastIndex(remoteAddrStr, ":")
		if lastColon == -1 {
			conn.RemoteAddress = remoteAddrStr
		} else {
			conn.RemoteAddress = remoteAddrStr[0:lastColon]
			conn.RemotePort = remoteAddrStr[lastColon+1:]
		}
		conns = append(conns, conn)
	}

	return conns
}
