package main

import (
	"math"
	"os"
	"testing"
	"time"

	"github.com/sumwonyuno/cp-scoring/model"
)

func TestEmptyEtcPasswd(t *testing.T) {
	bs := []byte("")
	users := parseEtcPasswd(bs)
	if len(users) != 0 {
		t.Fatal("Parsed users out of empty string")
	}
}
func TestBadEtcPasswd(t *testing.T) {
	bs := []byte("bad")
	users := parseEtcPasswd(bs)
	if len(users) != 0 {
		t.Fatal("Parsed users out of bad string")
	}

	// cut off
	bs = []byte("root:x:0:0:roo")
	users = parseEtcPasswd(bs)
	if len(users) != 0 {
		t.Fatal("Parsed users out of incomplete string")
	}
}

func TestEtcPasswd(t *testing.T) {
	// one user
	bs := []byte("root:x:0:0:root:/root:/bin/bash")
	users := parseEtcPasswd(bs)
	if len(users) != 1 {
		t.Fatal("Did not parse expected user")
	}
	user := users["root"]
	if user.Name != "root" {
		t.Fatal("Unexpected user name")
	}
	if user.ID != "0" {
		t.Fatal("Unexpected user ID")
	}
	if user.AccountPresent != true {
		t.Fatal("Account should be present")
	}

	// two users
	bs = []byte("root:x:0:0:root:/root:/bin/bash\nuser:x:1000:1000:user:/home/user:/bin/bash")
	users = parseEtcPasswd(bs)
	if len(users) != 2 {
		t.Fatal("Did not parse expected users")
	}
	user1 := users["root"]
	if user1.Name != "root" {
		t.Fatal("Unexpected user name")
	}
	if user1.ID != "0" {
		t.Fatal("Unexpected user ID")
	}
	if user1.AccountPresent != true {
		t.Fatal("Account should be present")
	}
	user2 := users["user"]
	if user2.Name != "user" {
		t.Fatal("Unexpected user name")
	}
	if user2.ID != "1000" {
		t.Fatal("Unexpected user ID")
	}
	if user2.AccountPresent != true {
		t.Fatal("Account should be present")
	}
}

func TestEmptyEtcShadow(t *testing.T) {
	bs := []byte("")
	users := parseEtcShadow(bs)
	if len(users) != 0 {
		t.Fatal("Parsed users out of empty string")
	}
}

func TestBadEtcShadow(t *testing.T) {
	bs := []byte("no")
	users := parseEtcShadow(bs)
	if len(users) != 0 {
		t.Fatal("Parsed users out of bad string")
	}

	// cut off
	bs = []byte("root:*:1648")
	users = parseEtcShadow(bs)
	if len(users) != 0 {
		t.Fatal("Parsed users out of incomplete string")
	}
}

func TestEtcShadow(t *testing.T) {
	// one user
	bs := []byte("root:*:16482:0:99999:7:::")
	users := parseEtcShadow(bs)
	if len(users) != 1 {
		t.Fatal("Did not parse expected user")
	}
	user := users["root"]
	if user.AccountExpires != false {
		t.Fatal("Account is expected to expire")
	}
	if user.AccountActive != true {
		t.Fatal("Account is expected to be active")
	}
	if user.PasswordExpires != false {
		t.Fatal("Account password is expected to not expire")
	}
	// 16482 * 86400 (seconds)
	if user.PasswordLastSet != 1424044800 {
		t.Fatal("Account last password set not expected value")
	}

	// multiple users
	bs = []byte("user1:*:16000:0:99999:7:::\nuser2:$6$notahash:16482:0:123:7:::\nuser3:!$6$notahash:16482:0:99999:7:::\nuser4:$6$notahash:16482:0:99999:7::17000:")
	users = parseEtcShadow(bs)
	if len(users) != 4 {
		t.Fatal("Did not parse expected users")
	}
	user1 := users["user1"]
	if user1.AccountExpires != false {
		t.Fatal("Account is expected to not expire")
	}
	if user1.AccountActive != true {
		t.Fatal("Account is expected to be active")
	}
	if user1.PasswordExpires != false {
		t.Fatal("Account password is expected to not expire")
	}
	// 16000 * 86400 (seconds)
	if user1.PasswordLastSet != 1382400000 {
		t.Fatal("Account last password set not expected value")
	}
	user2 := users["user2"]
	if user2.AccountExpires != false {
		t.Fatal("Account is expected to not expire")
	}
	if user2.AccountActive != true {
		t.Fatal("Account is expected to be active")
	}
	if user2.PasswordExpires != true {
		t.Fatal("Account password is expected to expire")
	}
	// 16482 * 86400 (seconds)
	if user2.PasswordLastSet != 1424044800 {
		t.Fatal("Account last password set not expected value")
	}
	user3 := users["user3"]
	if user3.AccountExpires != false {
		t.Fatal("Account is expected to not expire")
	}
	if user3.AccountActive != false {
		t.Fatal("Account is expected to be active")
	}
	if user3.PasswordExpires != false {
		t.Fatal("Account password is expected to not expire")
	}
	// 16482 * 86400 (seconds)
	if user3.PasswordLastSet != 1424044800 {
		t.Fatal("Account last password set not expected value")
	}
	user4 := users["user4"]
	if user4.AccountExpires != true {
		t.Fatal("Account is expected to expire")
	}
	if user4.AccountActive != true {
		t.Fatal("Account is expected to be not active/locked")
	}
	if user4.PasswordExpires != false {
		t.Fatal("Account password is expected to not expire")
	}
	// 16482 * 86400 (seconds)
	if user4.PasswordLastSet != 1424044800 {
		t.Fatal("Account last password set not expected value")
	}
}

func TestMergeUserMaps(t *testing.T) {
	// empty maps
	usersPasswd := make(map[string]model.User)
	usersShadow := make(map[string]model.User)

	result := mergeUserMaps(usersPasswd, usersShadow)
	if len(result) != 0 {
		t.Fatal("Expected no users")
	}

	// expected to have same users in both /etc/passwd and /etc/shadow
	userPart1 := model.User{}
	userPart1.Name = "bob"
	userPart1.ID = "175"
	userPart1.AccountPresent = true
	usersPasswd[userPart1.Name] = userPart1

	userPart2 := model.User{}
	userPart2.Name = userPart1.Name
	userPart2.AccountActive = true
	userPart2.AccountExpires = false
	userPart2.PasswordExpires = true
	userPart2.PasswordLastSet = 1000
	usersShadow[userPart2.Name] = userPart2

	result = mergeUserMaps(usersPasswd, usersShadow)
	if len(result) != 1 {
		t.Fatal("Expected to have 1 user entry")
	}
	user := result[0]
	if user.Name != "bob" {
		t.Fatal("Expected user name is bob")
	}
	if user.ID != "175" {
		t.Fatal("Expected user ID is 175")
	}
	if user.AccountPresent != true {
		t.Fatal("Expected user is present")
	}
	if user.AccountExpires != false {
		t.Fatal("Expected user account not to expire")
	}
	if user.AccountActive != true {
		t.Fatal("Expected user is active")
	}
	if user.PasswordExpires != true {
		t.Fatal("Expected user password to expire")
	}
	if user.PasswordLastSet != 1000 {
		t.Fatal("Unexpected password last set value")
	}
}

func TestParseWindowsUsersBad(t *testing.T) {
	// empty string
	bs := []byte("")
	users := parseWindowsUsers(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// bad string
	bs = []byte("csv")
	users = parseWindowsUsers(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// incorrect number
	bs = []byte("1,2,3,4,5\r\n1,2,3,4,5")
	users = parseWindowsUsers(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// incorrect number
	bs = []byte("1,2,3,4,5,6,7\r\n1,2,3,4,5,6,7")
	users = parseWindowsUsers(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// just header
	bs = []byte("1,2,3,4,5,6")
	users = parseWindowsUsers(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// mismatch between header and row
	bs = []byte("1,2,3,4,5,6\r\n1,2,3,4,5")
	users = parseWindowsUsers(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// mismatch between header and later row
	bs = []byte("1,2,3,4,5,6\r\n1,2,3,4,5,6\r\n1,2,3,4,5")
	users = parseWindowsUsers(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}
}

func TestParseWindowsUsers(t *testing.T) {
	// missing name
	bs := []byte("1,2,3,4,5,6\r\n,2,3,4,5,6")
	users := parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].Name != "" {
		t.Fatal("Expected no user name")
	}
	// given name
	bs = []byte("1,2,3,4,5,6\r\nname,2,3,4,5,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].Name != "name" {
		t.Fatal("Unexpected user name;", users[0].Name)
	}

	// present should always be true
	bs = []byte("1,2,3,4,5,6\r\n1,2,3,4,5,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].AccountPresent != true {
		t.Fatal("Expected account present to be true")
	}

	// missing id
	bs = []byte("1,2,3,4,5,6\r\nname,,3,4,5,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].ID != "" {
		t.Fatal("Expected no user ID")
	}

	// given id
	bs = []byte("1,2,3,4,5,6\r\nname,id,3,4,5,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].ID != "id" {
		t.Fatal("Unexpected user ID;", users[0].ID)
	}

	// missing account active
	bs = []byte("1,2,3,4,5,6\r\nname,id,,4,5,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].AccountActive != false {
		t.Fatal("Expected account active to be false when value is empty")
	}

	// account active true
	bs = []byte("1,2,3,4,5,6\r\nname,id,True,4,5,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].AccountActive != true {
		t.Fatal("Expected account active to be true")
	}

	// account active false
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,4,5,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].AccountActive != false {
		t.Fatal("Expected account active to be false")
	}

	// account doesn't expire
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,,5,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].AccountExpires != false {
		t.Fatal("Expected account expire to be false")
	}

	// account doesn't expire
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,expire,5,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].AccountExpires != true {
		t.Fatal("Expected account expire to be true")
	}

	// missing password last set
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,expire,,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].PasswordLastSet != math.MaxInt64 {
		t.Fatal("Expected password last set to be max int64 value")
	}

	// cannot parse password last set
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,expire,today,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].PasswordLastSet != 0 {
		t.Fatal("Expected password last set to be 0")
	}

	// password last set
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,expire,1/1/2000 12:34:56 AM,6")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	// this is timezone dependent...
	timezone, _ := time.Now().Zone()
	timestamp, _ := time.Parse("1/2/2006 3:04:05 PM MST", "1/1/2000 12:34:56 AM"+" "+timezone)
	if users[0].PasswordLastSet != timestamp.Unix() {
		t.Fatal("Expected password last set to be " + string(timestamp.Unix()))
	}

	// missing password expires
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,expire,1/1/2000 12:34:56 AM,")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].PasswordExpires != false {
		t.Fatal("Expected password expire to be false")
	}

	// password expires
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,expire,1/1/2000 12:34:56 AM,tomorrow")
	users = parseWindowsUsers(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].PasswordExpires != true {
		t.Fatal("Expected password expire to be true")
	}

	// multiple users
	bs = []byte("1,2,3,4,5,6\r\nname1,id,False,expire,1/1/2000 12:34:56 AM,tomorrow\r\nname2,id,False,expire,1/1/2000 12:34:56 AM,tomorrow")
	users = parseWindowsUsers(bs)
	if len(users) != 2 {
		t.Fatal("Expected 2 users")
	}
	if users[0].Name != "name1" {
		t.Fatal("Expected user name does not match")
	}
	if users[1].Name != "name2" {
		t.Fatal("Expected user name does not match")
	}
}

func TestParseWindowsTCPNetConnsBad(t *testing.T) {
	// empty string
	bs := []byte("")
	conns := parseWindowsTCPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// bad string
	bs = []byte("not")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// incorrect number
	bs = []byte("1,2,3,4,5\r\n1,2,3,4,5")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// incorrect number
	bs = []byte("1,2,3,4,5,6,7\r\n1,2,3,4,5,6,7")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// just header
	bs = []byte("1,2,3,4,5,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// mismatch between header and row
	bs = []byte("1,2,3,4,5,6\r\n1,2,3,4,5")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// mismatch between header and later row
	bs = []byte("1,2,3,4,5,6\r\n1,2,3,4,5,6\r\n1,2,3,4,5")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}
}

func TestParseWindowsTCPNetConns(t *testing.T) {
	// protocol should always be set
	bs := []byte("1,2,3,4,5,6\r\n1,2,3,4,5,6")
	conns := parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].Protocol != "TCP" {
		t.Fatal("Expected TCP protocol")
	}

	// empty PID
	bs = []byte("1,2,3,4,5,6\r\n,2,3,4,5,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].PID != 0 {
		t.Fatal("Expected default PID of 0")
	}

	// given PID
	bs = []byte("1,2,3,4,5,6\r\n726,2,3,4,5,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].PID != 726 {
		t.Fatal("Unexpected PID")
	}

	// missing state
	bs = []byte("1,2,3,4,5,6\r\n726,,3,4,5,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].State != model.NetworkConnectionUnknown {
		t.Fatal("Expected unknown state")
	}

	// invalid state
	bs = []byte("1,2,3,4,5,6\r\n726,invalid,3,4,5,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].State != model.NetworkConnectionUnknown {
		t.Fatal("Expected unknown state")
	}

	// given state
	bs = []byte("1,2,3,4,5,6\r\n726,established,3,4,5,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].State != model.NetworkConnectionEstablished {
		t.Fatal("Expected unknown state")
	}

	// missing local address
	bs = []byte("1,2,3,4,5,6\r\n726,established,,4,5,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if len(conns[0].LocalAddress) != 0 {
		t.Fatal("Expected no local address")
	}

	// given local address
	bs = []byte("1,2,3,4,5,6\r\n726,established,127.0.0.1,4,5,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].LocalAddress != "127.0.0.1" {
		t.Fatal("Expected unknown local address")
	}

	// missing local port
	bs = []byte("1,2,3,4,5,6\r\n726,established,127.0.0.1,,5,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if len(conns[0].LocalPort) != 0 {
		t.Fatal("Expected no local port")
	}

	// given local port
	bs = []byte("1,2,3,4,5,6\r\n726,established,127.0.0.1,49124,5,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].LocalPort != "49124" {
		t.Fatal("Expected unknown local port")
	}

	// missing remote address
	bs = []byte("1,2,3,4,5,6\r\n726,established,127.0.0.1,49124,,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if len(conns[0].RemoteAddress) != 0 {
		t.Fatal("Expected no remote address")
	}

	// given remote address
	bs = []byte("1,2,3,4,5,6\r\n726,established,127.0.0.1,49124,192.168.1.109,6")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].RemoteAddress != "192.168.1.109" {
		t.Fatal("Expected unknown remote address")
	}

	// missing remote port
	bs = []byte("1,2,3,4,5,6\r\n726,established,127.0.0.1,49124,192.168.1.109,")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if len(conns[0].RemotePort) != 0 {
		t.Fatal("Expected no remote port")
	}

	// given remote port
	bs = []byte("1,2,3,4,5,6\r\n726,established,127.0.0.1,49124,192.168.1.109,443")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].RemotePort != "443" {
		t.Fatal("Expected unknown remote port")
	}

	// multiple connections
	bs = []byte("1,2,3,4,5,6\r\n726,established,127.0.0.1,49124,192.168.1.109,443\r\n1313,established,10.2.62.124,57321,192.168.1.108,8443")
	conns = parseWindowsTCPNetConns(bs)
	if len(conns) != 2 {
		t.Fatal("Expected 2 connections")
	}
}

func TestParseWindowsUDPNetConnsBad(t *testing.T) {
	// empty string
	bs := []byte("")
	conns := parseWindowsUDPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// bad string
	bs = []byte("not")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// incorrect number
	bs = []byte("1,2\r\n1,2")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// incorrect number
	bs = []byte("1,2,3,4\r\n1,2,3,4")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// just header
	bs = []byte("1,2,3")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// mismatch between header and row
	bs = []byte("1,2,3\r\n1,2")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}

	// mismatch between header and later row
	bs = []byte("1,2,3\r\n1,2,3\r\n1,2")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 0 {
		t.Fatal("Expected 0 connections")
	}
}

func TestParseWindowsUDPNetConns(t *testing.T) {
	// protocol should always be set
	bs := []byte("1,2,3\r\n1,2,3")
	conns := parseWindowsUDPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connections")
	}
	if conns[0].Protocol != "UDP" {
		t.Fatal("Expected UDP protocol")
	}

	// empty PID
	bs = []byte("1,2,3\r\n,2,3")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].PID != 0 {
		t.Fatal("Expected default PID of 0")
	}

	// given PID
	bs = []byte("1,2,3\r\n726,2,3")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].PID != 726 {
		t.Fatal("Unexpected PID")
	}

	// state should always be unknown
	bs = []byte("1,2,3\r\n1,2,3")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connection")
	}
	if conns[0].State != model.NetworkConnectionUnknown {
		t.Fatal("Expected unknown state")
	}

	// missing local address
	bs = []byte("1,2,3r\n726,,3")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connections")
	}
	if len(conns[0].LocalAddress) != 0 {
		t.Fatal("Expected no local address")
	}

	// given local address
	bs = []byte("1,2,3\r\n726,127.0.0.1,3")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connections")
	}
	if conns[0].LocalAddress != "127.0.0.1" {
		t.Fatal("Expected unknown local address")
	}

	// missing local port
	bs = []byte("1,2,3\r\n726,127.0.0.1,")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connections")
	}
	if len(conns[0].LocalPort) != 0 {
		t.Fatal("Expected no local port")
	}

	// given local port
	bs = []byte("1,2,3\r\n726,127.0.0.1,49124")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 1 {
		t.Fatal("Expected 1 connections")
	}
	if conns[0].LocalPort != "49124" {
		t.Fatal("Expected unknown local port")
	}

	// two connections
	bs = []byte("1,2,3\r\n726,127.0.0.1,49124\r\n3621,127.0.0.1,80")
	conns = parseWindowsUDPNetConns(bs)
	if len(conns) != 2 {
		t.Fatal("Expected 2 connections")
	}
}

func TestParseWindowsProcessesBad(t *testing.T) {
	// empty string
	bs := []byte("")
	processes := parseWindowsProcesses(bs)
	if len(processes) != 0 {
		t.Fatal("Expected 0 processes")
	}

	// bad string
	bs = []byte("not")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 0 {
		t.Fatal("Expected 0 processes")
	}

	// incorrect number
	bs = []byte("1,2\r\n1,2")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 0 {
		t.Fatal("Expected 0 processes")
	}

	// incorrect number
	bs = []byte("1,2,3,4\r\n1,2,3,4")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 0 {
		t.Fatal("Expected 0 processes")
	}

	// just header
	bs = []byte("1,2,3")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 0 {
		t.Fatal("Expected 0 processes")
	}

	// mismatch between header and row
	bs = []byte("1,2,3\r\n1,2")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 0 {
		t.Fatal("Expected 0 processes")
	}

	// mismatch between header and later row
	bs = []byte("1,2,3\r\n1,2,3\r\n1,2")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 0 {
		t.Fatal("Expected 0 processes")
	}
}

func TestParseWindowsProcesses(t *testing.T) {
	// missing ID
	bs := []byte("1,2,3\r\n,2,3")
	processes := parseWindowsProcesses(bs)
	if len(processes) != 1 {
		t.Fatal("Expected 1 process")
	}
	if processes[0].PID != 0 {
		t.Fatal("Expected PID 0")
	}

	// given ID
	bs = []byte("1,2,3\r\n63,2,3")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 1 {
		t.Fatal("Expected 1 process")
	}
	if processes[0].PID != 63 {
		t.Fatal("Expected PID 63")
	}

	// missing user
	bs = []byte("1,2,3\r\n63,,3")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 1 {
		t.Fatal("Expected 1 process")
	}
	if len(processes[0].User) != 0 {
		t.Fatal("Expected empty user")
	}

	// given user
	bs = []byte("1,2,3\r\n63,user,3")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 1 {
		t.Fatal("Expected 1 process")
	}
	if processes[0].User != "user" {
		t.Fatal("Expected user name user")
	}

	// given user with hostname
	hostname, _ := os.Hostname()
	bs = []byte("1,2,3\r\n63," + hostname + "\\user,3")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 1 {
		t.Fatal("Expected 1 process")
	}
	if processes[0].User != "user" {
		t.Fatal("Expected user name user")
	}

	// missing command line
	bs = []byte("1,2,3\r\n63,user,")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 1 {
		t.Fatal("Expected 1 process")
	}
	if len(processes[0].CommandLine) != 0 {
		t.Fatal("Expected empty command line")
	}

	// given command line
	bs = []byte("1,2,3\r\n63,user,notepad.exe")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 1 {
		t.Fatal("Expected 1 process")
	}
	if processes[0].CommandLine != "notepad.exe" {
		t.Fatal("Expected command line")
	}

	// multiple processes
	bs = []byte("1,2,3\r\n63,user,notepad.exe\r\n72,user,mspaint.exe")
	processes = parseWindowsProcesses(bs)
	if len(processes) != 2 {
		t.Fatal("Expected 2 processes")
	}
	if processes[0].CommandLine != "notepad.exe" {
		t.Fatal("Unexpected command line")
	}
	if processes[1].CommandLine != "mspaint.exe" {
		t.Fatal("Unexpected command line")
	}
}

func TestParseWindowsGroupsBad(t *testing.T) {
	// empty string
	bs := []byte("")
	groups := parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// bad string
	bs = []byte("not")
	groups = parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// incorrect number
	bs = []byte("1\r\n1")
	groups = parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// incorrect number
	bs = []byte("1,2,3\r\n1,2,3")
	groups = parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// just header
	bs = []byte("1,2")
	groups = parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// mismatch between header and row
	bs = []byte("1,2\r\n1")
	groups = parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// mismatch between header and later row
	bs = []byte("1,2\r\n1,2\r\n1")
	groups = parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}
}

func TestParseWindowsGroups(t *testing.T) {
	// missing group
	bs := []byte("1,2\r\n,2")
	groups := parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// given group not expected format
	bs = []byte("1,2\r\ngroup,2")
	groups = parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// given group empty
	bs = []byte("1,2\r\n\"extra,Name=\"\"\"\"\",user")
	groups = parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// missing user
	bs = []byte("1,2\r\ngroup,")
	groups = parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// given user not expected format
	bs = []byte("1,2\r\n\"extra,Name=\"\"group\"\"\",user")
	groups = parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// given user empty
	bs = []byte("1,2\r\n\"extra,Name=\"\"group\"\"\",\"extra,Name=\"\"\"\"\"")
	groups = parseWindowsGroups(bs)
	if len(groups) != 0 {
		t.Fatal("Expected 0 groups")
	}

	// given user and given user
	bs = []byte("1,2\r\n\"extra,Name=\"\"group\"\"\",\"extra,Name=\"\"user\"\"\"")
	groups = parseWindowsGroups(bs)
	if len(groups) != 1 {
		t.Fatal("Expected 1 group")
	}
	users, present := groups["group"]
	if !present {
		t.Fatal("Expected group not found")
	}
	if len(users) != 1 {
		t.Fatal("Did not find expected number of users in group")
	}
	if users[0] != "user" {
		t.Fatal("Did not find expected user in group")
	}

	// 2 users same group
	bs = []byte("1,2\r\n\"extra,Name=\"\"group\"\"\",\"extra,Name=\"\"user1\"\"\"\r\n\"extra,Name=\"\"group\"\"\",\"extra,Name=\"\"user2\"\"\"")
	groups = parseWindowsGroups(bs)
	if len(groups) != 1 {
		t.Fatal("Expected 1 group")
	}
	users, present = groups["group"]
	if !present {
		t.Fatal("Expected group not found")
	}
	if len(users) != 2 {
		t.Fatal("Did not find expected number of users in group")
	}
	if users[0] != "user1" {
		t.Fatal("Did not find expected user in group")
	}
	if users[1] != "user2" {
		t.Fatal("Did not find expected user in group")
	}

	// 2 users different groups
	bs = []byte("1,2\r\n\"extra,Name=\"\"group1\"\"\",\"extra,Name=\"\"user1\"\"\"\r\n\"extra,Name=\"\"group2\"\"\",\"extra,Name=\"\"user2\"\"\"")
	groups = parseWindowsGroups(bs)
	if len(groups) != 2 {
		t.Fatal("Expected 2 groups")
	}
	users, present = groups["group1"]
	if !present {
		t.Fatal("Expected group not found")
	}
	if len(users) != 1 {
		t.Fatal("Did not find expected number of users in group")
	}
	if users[0] != "user1" {
		t.Fatal("Did not find expected user in group")
	}
	users, present = groups["group2"]
	if !present {
		t.Fatal("Expected group not found")
	}
	if len(users) != 1 {
		t.Fatal("Did not find expected number of users in group")
	}
	if users[0] != "user2" {
		t.Fatal("Did not find expected user in group")
	}

	// 2 users, 1 in one group, other in two groups
	bs = []byte("1,2\r\n\"extra,Name=\"\"group1\"\"\",\"extra,Name=\"\"user1\"\"\"\r\n\"extra,Name=\"\"group1\"\"\",\"extra,Name=\"\"user2\"\"\"\r\n\"extra,Name=\"\"group2\"\"\",\"extra,Name=\"\"user1\"\"\"")
	groups = parseWindowsGroups(bs)
	if len(groups) != 2 {
		t.Fatal("Expected 2 groups")
	}
	users, present = groups["group1"]
	if !present {
		t.Fatal("Expected group not found")
	}
	if len(users) != 2 {
		t.Fatal("Did not find expected number of users in group")
	}
	if users[0] != "user1" {
		t.Fatal("Did not find expected user in group")
	}
	if users[1] != "user2" {
		t.Fatal("Did not find expected user in group")
	}
	users, present = groups["group2"]
	if !present {
		t.Fatal("Expected group not found")
	}
	if len(users) != 1 {
		t.Fatal("Did not find expected number of users in group")
	}
	if users[0] != "user1" {
		t.Fatal("Did not find expected user in group")
	}

	// 3 users, 2 same group, 1 other group
	bs = []byte("1,2\r\n\"extra,Name=\"\"group1\"\"\",\"extra,Name=\"\"user1\"\"\"\r\n\"extra,Name=\"\"group2\"\"\",\"extra,Name=\"\"user2\"\"\"\r\n\"extra,Name=\"\"group1\"\"\",\"extra,Name=\"\"user3\"\"\"")
	groups = parseWindowsGroups(bs)
	if len(groups) != 2 {
		t.Fatal("Expected 2 groups")
	}
	users, present = groups["group1"]
	if !present {
		t.Fatal("Expected group not found")
	}
	if len(users) != 2 {
		t.Fatal("Did not find expected number of users in group")
	}
	if users[0] != "user1" {
		t.Fatal("Did not find expected user in group")
	}
	if users[1] != "user3" {
		t.Fatal("Did not find expected user in group")
	}
	users, present = groups["group2"]
	if !present {
		t.Fatal("Expected group not found")
	}
	if len(users) != 1 {
		t.Fatal("Did not find expected number of users in group")
	}
	if users[0] != "user2" {
		t.Fatal("Did not find expected user in group")
	}
}

func TestParseWindowsSoftwareBad(t *testing.T) {
	// empty string
	bs := []byte("")
	software := parseWindowsSoftware(bs)
	if len(software) != 0 {
		t.Fatal("Expected 0 software")
	}

	// bad string
	bs = []byte("not")
	software = parseWindowsSoftware(bs)
	if len(software) != 0 {
		t.Fatal("Expected 0 software")
	}

	// incorrect number
	bs = []byte("1\r\n1")
	software = parseWindowsSoftware(bs)
	if len(software) != 0 {
		t.Fatal("Expected 0 software")
	}

	// incorrect number
	bs = []byte("1,2,3\r\n1,2,3")
	software = parseWindowsSoftware(bs)
	if len(software) != 0 {
		t.Fatal("Expected 0 software")
	}

	// just header
	bs = []byte("1,2")
	software = parseWindowsSoftware(bs)
	if len(software) != 0 {
		t.Fatal("Expected 0 software")
	}

	// mismatch between header and row
	bs = []byte("1,2\r\n1")
	software = parseWindowsSoftware(bs)
	if len(software) != 0 {
		t.Fatal("Expected 0 software")
	}

	// mismatch between header and later row
	bs = []byte("1,2\r\n1,2\r\n1")
	software = parseWindowsSoftware(bs)
	if len(software) != 0 {
		t.Fatal("Expected 0 software")
	}
}

func TestParseWindowsSoftware(t *testing.T) {
	// missing name
	bs := []byte("1,2\r\n,2")
	software := parseWindowsSoftware(bs)
	if len(software) != 1 {
		t.Fatal("Expected 1 software")
	}
	if len(software[0].Name) != 0 {
		t.Fatal("Expected software name to be empty")
	}

	// missing version
	bs = []byte("1,2\r\ncp-scoring,")
	software = parseWindowsSoftware(bs)
	if len(software) != 1 {
		t.Fatal("Expected 1 software")
	}
	if len(software[0].Version) != 0 {
		t.Fatal("Expected software version to be empty")
	}

	// given name and version
	bs = []byte("1,2\r\ncp-scoring,0.1.0")
	software = parseWindowsSoftware(bs)
	if len(software) != 1 {
		t.Fatal("Expected 1 software")
	}
	if software[0].Name != "cp-scoring" {
		t.Fatal("Unexpected software name")
	}
	if software[0].Version != "0.1.0" {
		t.Fatal("Unexpected software version")
	}

	// multiple software
	bs = []byte("1,2\r\ncp-scoring,0.1.0\r\nanother,2")
	software = parseWindowsSoftware(bs)
	if len(software) != 2 {
		t.Fatal("Expected 2 software")
	}
	if software[0].Name != "cp-scoring" {
		t.Fatal("Unexpected software name")
	}
	if software[0].Version != "0.1.0" {
		t.Fatal("Unexpected software version")
	}
	if software[1].Name != "another" {
		t.Fatal("Unexpected software name")
	}
	if software[1].Version != "2" {
		t.Fatal("Unexpected software version")
	}
}

func TestFromHexStringPort(t *testing.T) {
	// empty string
	s, err := fromHexStringPort("")
	if err == nil {
		t.Fatal("Parsed port out of empty string")
	}

	// bad string
	s, err = fromHexStringPort(" ")
	if err == nil {
		t.Fatal("Parsed port out of space string")
	}

	// bad string
	s, err = fromHexStringPort("asdf!")
	if err == nil {
		t.Fatal("Parsed port out of bad string;", err)
	}

	// short string
	s, err = fromHexStringPort("bad")
	if err == nil {
		t.Fatal("Parsed port out of short string;", err)
	}

	// long string
	s, err = fromHexStringPort("baddd")
	if err == nil {
		t.Fatal("Parsed port out of long string;", err)
	}

	// acceptable string
	s, err = fromHexStringPort("0bad")
	if s != "2989" {
		t.Fatal("Unexpected parsed port")
	}

	s, err = fromHexStringPort("0000")
	if s != "0" {
		t.Fatal("Unexpected parsed port")
	}

	s, err = fromHexStringPort("000A")
	if s != "10" {
		t.Fatal("Unexpected parsed port")
	}

	s, err = fromHexStringPort("01BB")
	if s != "443" {
		t.Fatal("Unexpected parsed port")
	}

	s, err = fromHexStringPort("FFFF")
	if s != "65535" {
		t.Fatal("Unexpected parsed port")
	}
}

func TestFromHexStringIP(t *testing.T) {
	// empty string
	s, err := fromHexStringIP("")
	if err == nil {
		t.Fatal("Parsed IP out of empty string")
	}

	// bad string
	s, err = fromHexStringIP(" ")
	if err == nil {
		t.Fatal("Parsed IP out of space string")
	}

	// bad string
	s, err = fromHexStringIP("asdf!")
	if err == nil {
		t.Fatal("Parsed IP out of bad string;", err)
	}

	// bad string
	s, err = fromHexStringIP("bad")
	if err == nil {
		t.Fatal("Parsed IP out of bad string;", err)
	}

	// too short
	s, err = fromHexStringIP("0000000")
	if err == nil {
		t.Fatal("Parsed IP out of short string;", err)
	}

	// too long
	s, err = fromHexStringIP("000000000")
	if err == nil {
		t.Fatal("Parsed IP out of long string;", err)
	}

	// acceptable string
	s, err = fromHexStringIP("00000000")
	if s != "0.0.0.0" {
		t.Fatal("Unexpected parsed IP")
	}

	s, err = fromHexStringIP("FFFFFFFF")
	if s != "255.255.255.255" {
		t.Fatal("Unexpected parsed IP")
	}

	s, err = fromHexStringIP("7F000001")
	if s != "127.0.0.1" {
		t.Fatal("Unexpected parsed IP")
	}
}

func TestEmptyParseProcNet(t *testing.T) {
	bs := []byte("")
	conns := parseProcNet("TCP", bs)
	if len(conns) != 0 {
		t.Fatal("Parsed tcp conn out of empty string")
	}
}

func TestBadParseProcNet(t *testing.T) {
	bs := []byte("bad")
	conns := parseProcNet("TCP", bs)
	if len(conns) != 0 {
		t.Fatal("Parsed tcp conn out of bad string")
	}

	// cut off
	bs = []byte("sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n0: 00000000:0")
	conns = parseProcNet("TCP", bs)
	if len(conns) != 0 {
		t.Fatal("Parsed tcp conn out of incomplete string")
	}
}

func TestParseProcNet(t *testing.T) {
	bs := []byte("  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode                                                     \n0: 7F000001:0386 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 23479 1 ffff9e697826e080 100 0 0 10 0")
	conns := parseProcNet("TCP", bs)
	if len(conns) != 1 {
		t.Fatal("Did not parse expected tcp conn")
	}
	conn := conns[0]
	if conn.Protocol != "TCP" {
		t.Fatal("Unexpected protocol")
	}
	if conn.State != model.NetworkConnectionListen {
		t.Fatal("Unexpected state")
	}
	if conn.LocalAddress != "127.0.0.1" {
		t.Fatal("Unexpected local address")
	}
	if conn.LocalPort != "902" {
		t.Fatal("Unexpected local port")
	}
	if conn.RemoteAddress != "0.0.0.0" {
		t.Fatal("Unexpected remote address")
	}
	if conn.RemotePort != "0" {
		t.Fatal("Unexpected remote port")
	}
	if conn.PID != 0 {
		t.Fatal("Unexpected PID")
	}
}

func TestParseEtcGroup(t *testing.T) {
	// empty string
	bs := []byte("")
	groups := parseEtcGroup(bs)
	if len(groups) != 0 {
		t.Fatal("Parsed groups out of empty string")
	}

	// bad string
	bs = []byte("bad")
	groups = parseEtcGroup(bs)
	if len(groups) != 0 {
		t.Fatal("Parsed groups out of bad string")
	}

	// cut off
	bs = []byte("root:x:")
	groups = parseEtcGroup(bs)
	if len(groups) != 0 {
		t.Fatal("Parsed groups out of cut off string")
	}

	// one group
	bs = []byte("root:x:0:")
	groups = parseEtcGroup(bs)
	if len(groups) != 1 {
		t.Fatal("Did not parse expected group")
	}
	groupMembers1, present := groups["root"]
	if !present {
		t.Fatal("Did not find group root")
	}
	if len(groupMembers1) != 0 {
		t.Fatal("Unexpected group members")
	}

	// two groups
	bs = []byte("root:x:0:\nusers:x:100:user1,user2")
	groups = parseEtcGroup(bs)
	if len(groups) != 2 {
		t.Fatal("Did not parse expected group2")
	}
	groupMembers1, present = groups["root"]
	if !present {
		t.Fatal("Did not find group root")
	}
	if len(groupMembers1) != 0 {
		t.Fatal("Unexpected group members")
	}
	groupMembers2, present := groups["users"]
	if !present {
		t.Fatal("Did not find group users")
	}
	if len(groupMembers2) != 2 {
		t.Fatal("Unexpected group members")
	}
	if groupMembers2[0] != "user1" {
		t.Fatal("Unexpected user")
	}
	if groupMembers2[1] != "user2" {
		t.Fatal("Unexpected user")
	}
}
