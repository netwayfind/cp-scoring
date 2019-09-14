package agent

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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

func copyAgentLinux(installPath string) {
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
	binFile := filepath.Join(installPath, "cp-scoring-agent-linux")
	fileOut, err := os.Create(binFile)
	if err != nil {
		log.Fatalln("Unable to open destination file;", err)
	}
	defer fileOut.Close()
	_, err = io.Copy(fileOut, fileIn)
	if err != nil {
		log.Fatalln("Unable to copy file;", err)
	}
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

	log.Println("Finished installing to " + installPath)
}
