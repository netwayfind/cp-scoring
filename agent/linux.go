// +build linux

package main

import (
	"io/ioutil"
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
	out, err := exec.Command("/usr/bin/apt", "list", "--installed").Output()
	if err != nil {
		log.Fatal("ERROR: unable to get software list;", err)
	}

	software := parseAptListInstalled(out)

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
