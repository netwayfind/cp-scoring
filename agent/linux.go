// +build linux

package main

import (
	"io/ioutil"
	"log"
	"os/exec"
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

	groups := parseEtcGroup(bs)

	return groups
}

func getProcesses() []model.Process {
	out, err := exec.Command("/bin/ps", "-eo", "pid,user:32,command", "--sort=pid").Output()
	if err != nil {
		log.Fatal("ERROR: cannot get processes;", err)
	}

	processes := parseBinPs(out)

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
	// TCP connections
	bs, err := ioutil.ReadFile("/proc/net/tcp")
	if err != nil {
		log.Fatal("ERROR: cannot get tcp connections;", err)
	}
	tcpConns := parseProcNet("TCP", bs)

	// UDP connections
	bs, err = ioutil.ReadFile("/proc/net/udp")
	if err != nil {
		log.Fatal("ERROR: cannot get udp connections;", err)
	}
	udpConns := parseProcNet("UDP", bs)

	return append(tcpConns, udpConns...)
}
