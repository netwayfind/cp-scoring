package agent

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/netwayfind/cp-scoring/model"
)

type hostLinux struct {
}

func (h hostLinux) GetUsers() ([]model.User, error) {
	// get user and uid
	bs, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return nil, err
	}
	usersEtcPasswd := parseEtcPasswd(bs)

	// get other user information (sensitive)
	bs, err = ioutil.ReadFile("/etc/shadow")
	if err != nil {
		return nil, err
	}
	usersEtcShadow := parseEtcShadow(bs)

	return mergeLinuxUsers(usersEtcPasswd, usersEtcShadow), nil
}

func (h hostLinux) GetGroups() (map[string][]model.GroupMember, error) {
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
	// IPv4
	bs, err := ioutil.ReadFile("/proc/net/tcp")
	if err != nil {
		return nil, err
	}
	// IPv6
	tcpConns := parseProcNet("TCP", bs)
	bs, err = ioutil.ReadFile("/proc/net/tcp6")
	if err != nil {
		return nil, err
	}
	tcp6Conns := parseProcNet6("TCP", bs)

	// UDP connections
	// IPv4
	bs, err = ioutil.ReadFile("/proc/net/udp")
	if err != nil {
		return nil, err
	}
	udpConns := parseProcNet("UDP", bs)
	bs, err = ioutil.ReadFile("/proc/net/udp6")
	if err != nil {
		return nil, err
	}
	udp6Conns := parseProcNet6("UDP", bs)

	conns := append(tcpConns, tcp6Conns...)
	conns = append(conns, udpConns...)
	conns = append(conns, udp6Conns...)
	return conns, nil
}

func (h hostLinux) GetScheduledTasks() ([]model.ScheduledTask, error) {
	// not implemented yet
	return make([]model.ScheduledTask, 0), nil
}

func (h hostLinux) GetWindowsFirewallProfiles() ([]model.WindowsFirewallProfile, error) {
	// no Windows firewall
	return make([]model.WindowsFirewallProfile, 0), nil
}

func (h hostLinux) GetWindowsFirewallRules() ([]model.WindowsFirewallRule, error) {
	// no Windows firewall
	return make([]model.WindowsFirewallRule, 0), nil
}

func (h hostLinux) GetWindowsSettings() ([]model.WindowsSetting, error) {
	// not applicalbe
	return make([]model.WindowsSetting, 0), nil
}

func copyAgentLinux(installPath string) {
	log.Println("Copying this executable to installation folder")
	ex, err := os.Executable()
	if err != nil {
		log.Println("Unable to get this executable path;", err)
	}
	binFile := filepath.Join(installPath, "cp-scoring-agent-linux")
	copyFile(ex, binFile)
	err = os.Chmod(binFile, 0755)
	if err != nil {
		log.Fatalln("Unable to set permissions;", err)
	}
}

func getSystemdScript() []byte {
	return []byte(`[Unit]
Description=cp-scoring

[Service]
User=root
Group=root
WorkingDirectory=/opt/cp-scoring
ExecStart=/opt/cp-scoring/cp-scoring-agent-linux
Restart=on-failure

[Install]
WantedBy=multi-user.target
Alias=cp-scoring.service
`)
}

func createService(installPath string) {
	log.Println("Creating service")
	serviceFile := filepath.Join(installPath, "cp-scoring.service")
	err := ioutil.WriteFile(serviceFile, getSystemdScript(), 0755)
	if err != nil {
		log.Fatalln("Could not write Task Scheduler file")
	}
	// delete existing service and recreate
	exec.Command("/bin/systemctl", "disable", "cp-scoring.service").Run()
	err = exec.Command("/bin/systemctl", "enable", serviceFile).Run()
	if err != nil {
		log.Fatalln("Unable to enable service;", err)
	}
}

func createTeamKeyRegistrationLinux(installPath string) {
	log.Println("Creating team key registration")
	// script
	fileReg := filepath.Join(installPath, "teamkeyregistration.sh")
	text := []byte("#!/bin/sh\ncd /opt/cp-scoring\nsudo ./cp-scoring-agent-linux -teamKey")
	err := ioutil.WriteFile(fileReg, text, 0755)
	if err != nil {
		log.Fatalln("Could not write team key registration file")
	}
	// shortcut
	fileShortcut := filepath.Join(installPath, "teamkeyregistration.desktop")
	text = []byte("[Desktop Entry]\nEncoding=UTF-8\nVersion=1.0\nName[en_US]=Team Key Registration\nExec=/opt/cp-scoring/teamkeyregistration.sh\nTerminal=true\nType=Application")
	err = ioutil.WriteFile(fileShortcut, text, 0755)
	if err != nil {
		log.Fatalln("Could not write team key registration shortcut file")
	}
}

func (h hostLinux) Install() {
	installPath := "/opt/cp-scoring"

	// create installation folder
	err := os.MkdirAll(installPath, 0755)
	if err != nil {
		log.Fatalln("ERROR: cannot create installation folder;", err)
	}
	log.Println("Created installation folder: " + installPath)

	// copy agent
	copyAgentLinux(installPath)

	// create service
	createService(installPath)

	// create team key registration
	createTeamKeyRegistrationLinux(installPath)

	log.Println("Finished installing to " + installPath)
}

func (h hostLinux) CopyFiles() {
	installPath := "/opt/cp-scoring"
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalln("ERROR: cannot get current directory;", err)
	}
	log.Println("Copying files to: " + currentDir)

	// report
	copyFile(filepath.Join(installPath, "report.html"), filepath.Join(currentDir, "report.html"))

	// scoreboard
	copyFile(filepath.Join(installPath, "scoreboard.html"), filepath.Join(currentDir, "scoreboard.html"))

	// readme
	copyFile(filepath.Join(installPath, "README.html"), filepath.Join(currentDir, "README.html"))

	// team key registration shortcut
	copyFile(filepath.Join(installPath, "teamkeyregistration.desktop"), filepath.Join(currentDir, "teamkeyregistration.desktop"))
	os.Chmod(filepath.Join(currentDir, "teamkeyregistration.desktop"), 0755)

	log.Println("Finished copying files")
}
