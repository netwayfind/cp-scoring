package agent

import (
	"io/ioutil"
	"os/exec"

	"github.com/sumwonyuno/cp-scoring/model"
)

type hostLinux struct {
}

func (h hostLinux) GetUsers() ([]model.User, error) {
	// get user and uid
	bs, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return nil, err
	}
	userMapEtcPasswd := parseEtcPasswd(bs)

	// get other user information (sensitive)
	bs, err = ioutil.ReadFile("/etc/shadow")
	if err != nil {
		return nil, err
	}
	userMapEtcShadow := parseEtcShadow(bs)

	return mergeUserMaps(userMapEtcPasswd, userMapEtcShadow), nil
}

func (h hostLinux) GetGroups() (map[string][]string, error) {
	bs, err := ioutil.ReadFile("/etc/group")
	if err != nil {
		return nil, err
	}

	groups := parseEtcGroup(bs)

	return groups, nil
}

func (h hostLinux) GetProcesses() ([]model.Process, error) {
	out, err := exec.Command("/bin/ps", "-eo", "pid,user:32,command", "--sort=pid").Output()
	if err != nil {
		return nil, err
	}

	processes := parseBinPs(out)

	return processes, nil
}

func (h hostLinux) GetSoftware() ([]model.Software, error) {
	out, err := exec.Command("/usr/bin/apt", "list", "--installed").Output()
	if err != nil {
		return nil, err
	}

	software := parseAptListInstalled(out)

	return software, nil
}

func (h hostLinux) GetNetworkConnections() ([]model.NetworkConnection, error) {
	// TCP connections
	bs, err := ioutil.ReadFile("/proc/net/tcp")
	if err != nil {
		return nil, err
	}
	tcpConns := parseProcNet("TCP", bs)

	// UDP connections
	bs, err = ioutil.ReadFile("/proc/net/udp")
	if err != nil {
		return nil, err
	}
	udpConns := parseProcNet("UDP", bs)

	return append(tcpConns, udpConns...), err
}
