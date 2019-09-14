package model

import (
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type CurrentHost interface {
	GetUsers() ([]User, error)
	GetGroups() (map[string][]GroupMember, error)
	GetProcesses() ([]Process, error)
	GetSoftware() ([]Software, error)
	GetNetworkConnections() ([]NetworkConnection, error)
	GetScheduledTasks() ([]ScheduledTask, error)
	GetWindowsFirewallProfiles() ([]WindowsFirewallProfile, error)
	Install()
}

type Host struct {
	ID       uint64
	Hostname string
	OS       string
}

type StateSubmission struct {
	HostToken  string
	StateBytes []byte
}

type State struct {
	Timestamp          int64
	OS                 string
	Hostname           string
	Errors             []string
	Users              []User
	Groups             map[string][]GroupMember
	Processes          []Process
	Software           []Software
	NetworkConnections []NetworkConnection
	ScheduledTasks     []ScheduledTask
	WindowsFirewall    []WindowsFirewallProfile
}

type ObjectState string

const (
	ObjectStateAdd    ObjectState = "Add"
	ObjectStateKeep   ObjectState = "Keep"
	ObjectStateRemove ObjectState = "Remove"
)

type User struct {
	Name            string
	ID              string
	ObjectState     ObjectState
	AccountActive   bool
	AccountExpires  bool
	PasswordLastSet int64
	PasswordExpires bool
}

type GroupMember struct {
	Name        string
	ObjectState ObjectState
}

type Process struct {
	ObjectState ObjectState
	PID         int64
	User        string
	CommandLine string
}

type Software struct {
	Name        string
	Version     string
	ObjectState ObjectState
}

type ScheduledTask struct {
	Name        string
	Path        string
	Enabled     bool
	ObjectState ObjectState
}

type WindowsFirewallProfile struct {
	Name                  string
	Enabled               bool
	DefaultInboundAction  string
	DefaultOutboundAction string
}

type NetworkConnectionState string

const (
	NetworkConnectionClosed      NetworkConnectionState = "CLOSED"
	NetworkConnectionCloseWait   NetworkConnectionState = "CLOSE_WAIT"
	NetworkConnectionClosing     NetworkConnectionState = "CLOSING"
	NetworkConnectionDeleteTcb   NetworkConnectionState = "DELETE_TCB"
	NetworkConnectionEstablished NetworkConnectionState = "ESTABLISHED"
	NetworkConnectionFinWait1    NetworkConnectionState = "FIN_WAIT1"
	NetworkConnectionFinWait2    NetworkConnectionState = "FIN_WAIT2"
	NetworkConnectionLastAck     NetworkConnectionState = "LAST_ACK"
	NetworkConnectionListen      NetworkConnectionState = "LISTEN"
	NetworkConnectionSynReceived NetworkConnectionState = "SYN_RECV"
	NetworkConnectionSynSent     NetworkConnectionState = "SYN_SENT"
	NetworkConnectionTimeWait    NetworkConnectionState = "TIME_WAIT"
	NetworkConnectionUnconn      NetworkConnectionState = "UNCONN"
	NetworkConnectionUnknown     NetworkConnectionState = "UNKNOWN"
)

func GetNetworkConnectionStateLinux(hex string) NetworkConnectionState {
	switch hex {
	case "01":
		return NetworkConnectionEstablished
	case "02":
		return NetworkConnectionSynSent
	case "03":
		return NetworkConnectionSynReceived
	case "04":
		return NetworkConnectionFinWait1
	case "05":
		return NetworkConnectionFinWait2
	case "06":
		return NetworkConnectionTimeWait
	case "07":
		return NetworkConnectionClosed
	case "08":
		return NetworkConnectionCloseWait
	case "09":
		return NetworkConnectionLastAck
	case "0A":
		return NetworkConnectionListen
	case "0B":
		return NetworkConnectionClosing
	case "0C":
		return NetworkConnectionSynReceived
	default:
		return NetworkConnectionUnknown
	}
}

func GetNetworkConnectionState(stateStr string) NetworkConnectionState {
	// narrow down possible state strings
	stateStr = strings.ToLower(stateStr)
	stateStr = strings.Replace(stateStr, "_", "", -1)
	stateStr = strings.Replace(stateStr, "-", "", -1)
	switch stateStr {
	case "closed":
		return NetworkConnectionClosed
	case "closewait":
		return NetworkConnectionCloseWait
	case "closing":
		return NetworkConnectionClosing
	case "deletetcb":
		return NetworkConnectionDeleteTcb
	case "estab":
		return NetworkConnectionEstablished
	case "established":
		return NetworkConnectionEstablished
	case "finwait1":
		return NetworkConnectionFinWait1
	case "finwait2":
		return NetworkConnectionFinWait2
	case "lastack":
		return NetworkConnectionLastAck
	case "listen":
		return NetworkConnectionListen
	case "listening":
		return NetworkConnectionListen
	case "synrecv":
		return NetworkConnectionSynReceived
	case "synreceived":
		return NetworkConnectionSynReceived
	case "synsent":
		return NetworkConnectionSynSent
	case "timewait":
		return NetworkConnectionTimeWait
	case "unconn":
		return NetworkConnectionUnconn
	case "unknown":
		return NetworkConnectionUnknown
	default:
		return NetworkConnectionUnknown
	}
}

type NetworkConnection struct {
	Protocol      string
	PID           int64
	LocalAddress  string
	LocalPort     string
	RemoteAddress string
	RemotePort    string
	State         NetworkConnectionState
	ObjectState   ObjectState
}

func GetNewStateTemplate() State {
	var state State
	var err error
	state.Timestamp = time.Now().Unix()
	state.OS = runtime.GOOS
	state.Hostname, err = os.Hostname()
	if err != nil {
		log.Println("ERROR: unable to get hostname;", err)
	}

	return state
}
