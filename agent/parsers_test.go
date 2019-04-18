package agent

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
	user := users[0]
	if user.Name != "root" {
		t.Fatal("Unexpected user name")
	}
	if user.ID != "0" {
		t.Fatal("Unexpected user ID")
	}

	// two users
	bs = []byte("root:x:0:0:root:/root:/bin/bash\nuser:x:1000:1000:user:/home/user:/bin/bash")
	users = parseEtcPasswd(bs)
	if len(users) != 2 {
		t.Fatal("Did not parse expected users")
	}
	user1 := users[0]
	if user1.Name != "root" {
		t.Fatal("Unexpected user name")
	}
	if user1.ID != "0" {
		t.Fatal("Unexpected user ID")
	}
	user2 := users[1]
	if user2.Name != "user" {
		t.Fatal("Unexpected user name")
	}
	if user2.ID != "1000" {
		t.Fatal("Unexpected user ID")
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
	user := users[0]
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
	user1 := users[0]
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
	user2 := users[1]
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
	user3 := users[2]
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
	user4 := users[3]
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

func TestMergeLinuxUsers(t *testing.T) {
	// empty maps
	usersPasswd := make([]model.User, 0)
	usersShadow := make([]model.User, 0)

	result := mergeLinuxUsers(usersPasswd, usersShadow)
	if len(result) != 0 {
		t.Fatal("Expected no users")
	}

	// expected to have same users in both /etc/passwd and /etc/shadow
	user1Part1 := model.User{}
	user1Part1.Name = "bob"
	user1Part1.ID = "175"
	usersPasswd = append(usersPasswd, user1Part1)
	user2Part1 := model.User{}
	user2Part1.Name = "alice"
	user2Part1.ID = "176"
	usersPasswd = append(usersPasswd, user2Part1)

	user1Part2 := model.User{}
	user1Part2.Name = user1Part1.Name
	user1Part2.AccountActive = true
	user1Part2.AccountExpires = false
	user1Part2.PasswordExpires = true
	user1Part2.PasswordLastSet = 1000
	usersShadow = append(usersShadow, user1Part2)
	user2Part2 := model.User{}
	user2Part2.Name = user2Part1.Name
	user2Part2.AccountActive = true
	user2Part2.AccountExpires = false
	user2Part2.PasswordExpires = true
	user2Part2.PasswordLastSet = 2000
	usersShadow = append(usersShadow, user2Part2)

	result = mergeLinuxUsers(usersPasswd, usersShadow)
	if len(result) != 2 {
		t.Fatal("Expected to have 2 user entries")
	}
	if result[0].Name != "bob" {
		t.Fatal("Expected user name is bob")
	}
	if result[0].ID != "175" {
		t.Fatal("Expected user ID is 175")
	}
	if result[0].AccountExpires != false {
		t.Fatal("Expected user account not to expire")
	}
	if result[0].AccountActive != true {
		t.Fatal("Expected user is active")
	}
	if result[0].PasswordExpires != true {
		t.Fatal("Expected user password to expire")
	}
	if result[0].PasswordLastSet != 1000 {
		t.Fatal("Unexpected password last set value")
	}
	if result[1].Name != "alice" {
		t.Fatal("Expected user name is bob")
	}
	if result[1].ID != "176" {
		t.Fatal("Expected user ID is 175")
	}
	if result[1].AccountExpires != false {
		t.Fatal("Expected user account not to expire")
	}
	if result[1].AccountActive != true {
		t.Fatal("Expected user is active")
	}
	if result[1].PasswordExpires != true {
		t.Fatal("Expected user password to expire")
	}
	if result[1].PasswordLastSet != 2000 {
		t.Fatal("Unexpected password last set value")
	}
}

func TestParseWindowsUsersGetLocalUserBad(t *testing.T) {
	// empty string
	bs := []byte("")
	users := parseWindowsUsersGetLocalUser(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// bad string
	bs = []byte("csv")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// incorrect number
	bs = []byte("1,2,3,4,5\r\n1,2,3,4,5")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// incorrect number
	bs = []byte("1,2,3,4,5,6,7\r\n1,2,3,4,5,6,7")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// just header
	bs = []byte("1,2,3,4,5,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// mismatch between header and row
	bs = []byte("1,2,3,4,5,6\r\n1,2,3,4,5")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// mismatch between header and later row
	bs = []byte("1,2,3,4,5,6\r\n1,2,3,4,5,6\r\n1,2,3,4,5")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}
}

func TestParseWindowsUsersGetLocalUser(t *testing.T) {
	// missing name
	bs := []byte("1,2,3,4,5,6\r\n,2,3,4,5,6")
	users := parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].Name != "" {
		t.Fatal("Expected no user name")
	}
	// given name
	bs = []byte("1,2,3,4,5,6\r\nname,2,3,4,5,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].Name != "name" {
		t.Fatal("Unexpected user name;", users[0].Name)
	}

	// present should always be true
	bs = []byte("1,2,3,4,5,6\r\n1,2,3,4,5,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}

	// missing id
	bs = []byte("1,2,3,4,5,6\r\nname,,3,4,5,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].ID != "" {
		t.Fatal("Expected no user ID")
	}

	// given id
	bs = []byte("1,2,3,4,5,6\r\nname,id,3,4,5,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].ID != "id" {
		t.Fatal("Unexpected user ID;", users[0].ID)
	}

	// missing account active
	bs = []byte("1,2,3,4,5,6\r\nname,id,,4,5,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].AccountActive != false {
		t.Fatal("Expected account active to be false when value is empty")
	}

	// account active true
	bs = []byte("1,2,3,4,5,6\r\nname,id,True,4,5,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].AccountActive != true {
		t.Fatal("Expected account active to be true")
	}

	// account active false
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,4,5,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].AccountActive != false {
		t.Fatal("Expected account active to be false")
	}

	// account doesn't expire
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,,5,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].AccountExpires != false {
		t.Fatal("Expected account expire to be false")
	}

	// account doesn't expire
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,expire,5,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].AccountExpires != true {
		t.Fatal("Expected account expire to be true")
	}

	// missing password last set
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,expire,,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].PasswordLastSet != math.MaxInt64 {
		t.Fatal("Expected password last set to be max int64 value")
	}

	// cannot parse password last set
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,expire,today,6")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].PasswordLastSet != 0 {
		t.Fatal("Expected password last set to be 0")
	}

	// password last set
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,expire,1/1/2000 12:34:56 AM,6")
	users = parseWindowsUsersGetLocalUser(bs)
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
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].PasswordExpires != false {
		t.Fatal("Expected password expire to be false")
	}

	// password expires
	bs = []byte("1,2,3,4,5,6\r\nname,id,False,expire,1/1/2000 12:34:56 AM,tomorrow")
	users = parseWindowsUsersGetLocalUser(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 users")
	}
	if users[0].PasswordExpires != true {
		t.Fatal("Expected password expire to be true")
	}

	// multiple users
	bs = []byte("1,2,3,4,5,6\r\nname1,id,False,expire,1/1/2000 12:34:56 AM,tomorrow\r\nname2,id,False,expire,1/1/2000 12:34:56 AM,tomorrow")
	users = parseWindowsUsersGetLocalUser(bs)
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

func TestParseWindowsUsersWin32UserAccountBad(t *testing.T) {
	// empty string
	bs := []byte("")
	users := parseWindowsUsersWin32UserAccount(bs)
	if len(users) != 0 {
		t.Fatal("Parsed users from empty string")
	}

	// bad string
	bs = []byte("bad")
	users = parseWindowsUsersWin32UserAccount(bs)
	if len(users) != 0 {
		t.Fatal("Parsed users from bad string")
	}

	// incorrect number
	bs = []byte("Name\r\nuser")
	users = parseWindowsUsersWin32UserAccount(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// incorrect number
	bs = []byte("Name,SID,other\r\nuser,5,other")
	users = parseWindowsUsersWin32UserAccount(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// just header
	bs = []byte("Name,SID")
	users = parseWindowsUsersWin32UserAccount(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// mismatch between header and row
	bs = []byte("Name,SID\r\nuser")
	users = parseWindowsUsersWin32UserAccount(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}

	// mismatch between header and later row
	bs = []byte("Name,SID\r\nuser1,5\r\nuser2")
	users = parseWindowsUsersWin32UserAccount(bs)
	if len(users) != 0 {
		t.Fatal("Expected 0 users")
	}
}

func TestParseWindowsUsersWin32UserAccount(t *testing.T) {
	// empty name
	bs := []byte("Name,SID\r\n,5")
	users := parseWindowsUsersWin32UserAccount(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 user")
	}
	if users[0].Name != "" {
		t.Fatal("Expected empty name")
	}

	// given name
	bs = []byte("Name,SID\r\nuser,5")
	users = parseWindowsUsersWin32UserAccount(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 user")
	}
	if users[0].Name != "user" {
		t.Fatal("Unexpected name")
	}

	// empty SID
	bs = []byte("Name,SID\r\nuser,")
	users = parseWindowsUsersWin32UserAccount(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 user")
	}
	if users[0].ID != "" {
		t.Fatal("Expected empty ID")
	}

	// given SID
	bs = []byte("Name,SID\r\nuser,5")
	users = parseWindowsUsersWin32UserAccount(bs)
	if len(users) != 1 {
		t.Fatal("Expected 1 user")
	}
	if users[0].ID != "5" {
		t.Fatal("Unexpected ID")
	}
}

func TestParseWindowsNetUser(t *testing.T) {
	// account active
	bs := []byte("Account active               Yes")
	user := parseWindowsNetUser(bs)
	if !user.AccountActive {
		t.Fatal("Expected account active")
	}

	// account not active
	bs = []byte("Account active                No")
	user = parseWindowsNetUser(bs)
	if user.AccountActive {
		t.Fatal("Expected account not active")
	}

	// account expires
	bs = []byte("Account expires               Yes")
	user = parseWindowsNetUser(bs)
	if !user.AccountExpires {
		t.Fatal("Expected account to expire")
	}

	// account not expires
	bs = []byte("Account expires               Never")
	user = parseWindowsNetUser(bs)
	if user.AccountExpires {
		t.Fatal("Expected account to not expire")
	}

	// password expires
	bs = []byte("Password expires              Yes")
	user = parseWindowsNetUser(bs)
	if !user.PasswordExpires {
		t.Fatal("Expected password to expire")
	}

	// password not expires
	bs = []byte("Password expires              Never")
	user = parseWindowsNetUser(bs)
	if user.PasswordExpires {
		t.Fatal("Expected password to not expire")
	}

	// password last set
	bs = []byte("Password last set             1/1/2018 09:34:56 PM")
	user = parseWindowsNetUser(bs)
	// using local timezone, so need to test that password last set is expected
	expected := time.Unix(user.PasswordLastSet, 0).String()
	if expected[0:19] != "2018-01-01 21:34:56" {
		t.Fatal("Unexpected password set value")
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
	if users[0].Name != "user" {
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
	if users[0].Name != "user1" {
		t.Fatal("Did not find expected user in group")
	}
	if users[1].Name != "user2" {
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
	if users[0].Name != "user1" {
		t.Fatal("Did not find expected user in group")
	}
	users, present = groups["group2"]
	if !present {
		t.Fatal("Expected group not found")
	}
	if len(users) != 1 {
		t.Fatal("Did not find expected number of users in group")
	}
	if users[0].Name != "user2" {
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
	if users[0].Name != "user1" {
		t.Fatal("Did not find expected user in group")
	}
	if users[1].Name != "user2" {
		t.Fatal("Did not find expected user in group")
	}
	users, present = groups["group2"]
	if !present {
		t.Fatal("Expected group not found")
	}
	if len(users) != 1 {
		t.Fatal("Did not find expected number of users in group")
	}
	if users[0].Name != "user1" {
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
	if users[0].Name != "user1" {
		t.Fatal("Did not find expected user in group")
	}
	if users[1].Name != "user3" {
		t.Fatal("Did not find expected user in group")
	}
	users, present = groups["group2"]
	if !present {
		t.Fatal("Expected group not found")
	}
	if len(users) != 1 {
		t.Fatal("Did not find expected number of users in group")
	}
	if users[0].Name != "user2" {
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

	// empty name
	bs = []byte(",")
	software = parseWindowsSoftware(bs)
	if len(software) != 0 {
		t.Fatal("Expected 0 software")
	}
}

func TestParseWindowsSoftware(t *testing.T) {
	// missing name
	bs := []byte("1,2\r\n,2")
	software := parseWindowsSoftware(bs)
	if len(software) != 0 {
		t.Fatal("Expected 0 software")
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

func TestFromHexStringIPv4(t *testing.T) {
	// empty string
	s, err := fromHexStringIPv4("")
	if err == nil {
		t.Fatal("Parsed IP out of empty string")
	}

	// bad string
	s, err = fromHexStringIPv4(" ")
	if err == nil {
		t.Fatal("Parsed IP out of space string")
	}

	// bad string
	s, err = fromHexStringIPv4("asdf!")
	if err == nil {
		t.Fatal("Parsed IP out of bad string;", err)
	}

	// bad string
	s, err = fromHexStringIPv4("bad")
	if err == nil {
		t.Fatal("Parsed IP out of bad string;", err)
	}

	// too short
	s, err = fromHexStringIPv4("0000000")
	if err == nil {
		t.Fatal("Parsed IP out of short string;", err)
	}

	// too long
	s, err = fromHexStringIPv4("000000000")
	if err == nil {
		t.Fatal("Parsed IP out of long string;", err)
	}

	// non-hex string
	s, err = fromHexStringIPv4("0000000G")
	if err == nil {
		t.Fatal("Parsed IP out of non-hex string;", err)
	}

	// acceptable string
	s, err = fromHexStringIPv4("00000000")
	if s != "0.0.0.0" {
		t.Fatal("Unexpected parsed IP")
	}

	s, err = fromHexStringIPv4("FFFFFFFF")
	if s != "255.255.255.255" {
		t.Fatal("Unexpected parsed IP")
	}

	s, err = fromHexStringIPv4("0100007F")
	if s != "127.0.0.1" {
		t.Fatal("Unexpected parsed IP")
	}
}

func TestFromHexStringIPv6(t *testing.T) {
	// empty string
	s, err := fromHexStringIPv6("")
	if err == nil {
		t.Fatal("Parsed IP out of empty string")
	}

	// bad string
	s, err = fromHexStringIPv6(" ")
	if err == nil {
		t.Fatal("Parsed IP out of space string")
	}

	// bad string
	s, err = fromHexStringIPv6("asdf!")
	if err == nil {
		t.Fatal("Parsed IP out of bad string;", err)
	}

	// bad string
	s, err = fromHexStringIPv6("bad")
	if err == nil {
		t.Fatal("Parsed IP out of bad string;", err)
	}

	// too short
	s, err = fromHexStringIPv6("0000000000000000000000000000000")
	if err == nil {
		t.Fatal("Parsed IP out of short string;", err)
	}

	// too long
	s, err = fromHexStringIPv6("000000000000000000000000000000000")
	if err == nil {
		t.Fatal("Parsed IP out of long string;", err)
	}

	// non-hex string
	s, err = fromHexStringIPv6("0000000000000000000000000000000G")
	if err == nil {
		t.Fatal("Parsed IP out of non-hex string;", err)
	}

	// acceptable string
	s, err = fromHexStringIPv6("00000000000000000000000000000000")
	if s != "0000:0000:0000:0000:0000:0000:0000:0000" {
		t.Fatal("Unexpected parsed IP")
	}

	s, err = fromHexStringIPv6("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	if s != "FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF" {
		t.Fatal("Unexpected parsed IP")
	}

	s, err = fromHexStringIPv6("00000000000000000000000001000000")
	if s != "0000:0000:0000:0000:0000:0000:0000:0001" {
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
	bs = []byte("  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n   0: 00000000:0")
	conns = parseProcNet("TCP", bs)
	if len(conns) != 0 {
		t.Fatal("Parsed tcp conn out of incomplete string")
	}
}

func TestParseProcNet(t *testing.T) {
	// example 1
	bs := []byte("  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n   0: 0100007F:0386 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 23479 1 ffff9e697826e080 100 0 0 10 0")
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

	// example 2
	bs = []byte("  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n   7: 0201A8C0:A2B8 0D0C0B0A:01BB 01 00000000:00000000 00:00000000 00000000  1000        0 73261 1 ffff9256263d0000 37 4 9 10 -1")
	conns = parseProcNet("TCP", bs)
	if len(conns) != 1 {
		t.Fatal("Did not parse expected tcp conn")
	}
	conn = conns[0]
	if conn.Protocol != "TCP" {
		t.Fatal("Unexpected protocol")
	}
	if conn.State != model.NetworkConnectionEstablished {
		t.Fatal("Unexpected state")
	}
	if conn.LocalAddress != "192.168.1.2" {
		t.Fatal("Unexpected local address")
	}
	if conn.LocalPort != "41656" {
		t.Fatal("Unexpected local port")
	}
	if conn.RemoteAddress != "10.11.12.13" {
		t.Fatal("Unexpected remote address")
	}
	if conn.RemotePort != "443" {
		t.Fatal("Unexpected remote port")
	}
	if conn.PID != 0 {
		t.Fatal("Unexpected PID")
	}
}

func TestEmptyParseProcNet6(t *testing.T) {
	bs := []byte("")
	conns := parseProcNet6("TCP", bs)
	if len(conns) != 0 {
		t.Fatal("Parsed tcp conn out of empty string")
	}
}

func TestBadParseProcNet6(t *testing.T) {
	bs := []byte("bad")
	conns := parseProcNet6("TCP", bs)
	if len(conns) != 0 {
		t.Fatal("Parsed tcp conn out of bad string")
	}

	// cut off
	bs = []byte("  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n   0: 00000000000000000000000000000000:0386 00000")
	conns = parseProcNet6("TCP", bs)
	if len(conns) != 0 {
		t.Fatal("Parsed tcp conn out of incomplete string")
	}
}

func TestParseProcNet6(t *testing.T) {
	// example 1
	bs := []byte("  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n	0: 00000000000000000000000001000000:0386 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 27887 1 ffffa01c77fde800 100 0 0 10 0")
	conns := parseProcNet6("TCP", bs)
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
	if conn.LocalAddress != "0000:0000:0000:0000:0000:0000:0000:0001" {
		t.Fatal("Unexpected local address")
	}
	if conn.LocalPort != "902" {
		t.Fatal("Unexpected local port")
	}
	if conn.RemoteAddress != "0000:0000:0000:0000:0000:0000:0000:0000" {
		t.Fatal("Unexpected remote address")
	}
	if conn.RemotePort != "0" {
		t.Fatal("Unexpected remote port")
	}
	if conn.PID != 0 {
		t.Fatal("Unexpected PID")
	}

	// example 2
	bs = []byte("sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n	0: 00000000000000000000000001000000:A2B8 3000000020000000100000003400007F:01BB 01 00000000:00000000 00:00000000 00000000     0        0 27887 1 ffffa01c77fde800 100 0 0 10 0")
	conns = parseProcNet6("TCP", bs)
	if len(conns) != 1 {
		t.Fatal("Did not parse expected tcp conn")
	}
	conn = conns[0]
	if conn.Protocol != "TCP" {
		t.Fatal("Unexpected protocol")
	}
	if conn.State != model.NetworkConnectionEstablished {
		t.Fatal("Unexpected state")
	}
	if conn.LocalAddress != "0000:0000:0000:0000:0000:0000:0000:0001" {
		t.Fatal("Unexpected local address")
	}
	if conn.LocalPort != "41656" {
		t.Fatal("Unexpected local port")
	}
	if conn.RemoteAddress != "0000:0030:0000:0020:0000:0010:7F00:0034" {
		t.Fatal("Unexpected remote address")
	}
	if conn.RemotePort != "443" {
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
	if groupMembers2[0].Name != "user1" {
		t.Fatal("Unexpected user")
	}
	if groupMembers2[1].Name != "user2" {
		t.Fatal("Unexpected user")
	}
}

func TestParseBinPs(t *testing.T) {
	// empty string
	bs := []byte("")
	processes := parseBinPs(bs)
	if len(processes) != 0 {
		t.Fatal("Parsed processes out of empty string")
	}

	// bad string
	bs = []byte("bad")
	processes = parseBinPs(bs)
	if len(processes) != 0 {
		t.Fatal("Parsed processes out of bad string")
	}

	// cut off
	bs = []byte("PID USER COMMAND\n123 us")
	processes = parseBinPs(bs)
	if len(processes) != 0 {
		t.Fatal("Parsed processes out of cut off string")
	}

	// cut off
	bs = []byte("PID USER COMMAND\n123 user")
	processes = parseBinPs(bs)
	if len(processes) != 0 {
		t.Fatal("Parsed processes out of cut off string")
	}

	// empty pid
	bs = []byte("PID USER COMMAND\n    user cmd")
	processes = parseBinPs(bs)
	if len(processes) != 1 {
		t.Fatal("Did not find expected entry")
	}
	if processes[0].PID != -1 {
		t.Fatal("Unexpected default PID value")
	}

	// empty user
	bs = []byte("PID USER COMMAND\n123      cmd")
	processes = parseBinPs(bs)
	if len(processes) != 1 {
		t.Fatal("Did not find expected entry")
	}
	if len(processes[0].User) != 0 {
		t.Fatal("Did not parse empty string for empty user")
	}

	// empty command
	bs = []byte("PID USER COMMAND\n123 user ")
	processes = parseBinPs(bs)
	if len(processes) != 1 {
		t.Fatal("Did not find expected entry")
	}
	if len(processes[0].CommandLine) != 0 {
		t.Fatal("Did not parse empty string for empty command")
	}

	// one entry
	bs = []byte("PID USER COMMAND\n123 user command")
	processes = parseBinPs(bs)
	if len(processes) != 1 {
		t.Fatal("Did not find expected entry")
	}
	if processes[0].PID != 123 {
		t.Fatal("Did not parse expected PID")
	}
	if processes[0].User != "user" {
		t.Fatal("Did not parse expected user")
	}
	if processes[0].CommandLine != "command" {
		t.Fatal("Did not parse expected command")
	}

	// two entries
	// second has larger PID, longer username, longer command
	bs = []byte(" PID USER                             COMMAND\n 123 user                             command\n5678 abcdefghijklmnopqrstuvwxyz789012 this is a test command with spaces")
	processes = parseBinPs(bs)
	if len(processes) != 2 {
		t.Fatal("Did not find expected entries")
	}
	if processes[0].PID != 123 {
		t.Fatal("Did not parse expected PID")
	}
	if processes[0].User != "user" {
		t.Fatal("Did not parse expected user")
	}
	if processes[0].CommandLine != "command" {
		t.Fatal("Did not parse expected command")
	}
	if processes[1].PID != 5678 {
		t.Fatal("Did not parse expected PID")
	}
	if processes[1].User != "abcdefghijklmnopqrstuvwxyz789012" {
		t.Fatal("Did not parse expected user")
	}
	if processes[1].CommandLine != "this is a test command with spaces" {
		t.Fatal("Did not parse expected command")
	}
}

func TestParseAptListInstalled(t *testing.T) {
	// empty string
	bs := []byte("")
	software := parseAptListInstalled(bs)
	if len(software) != 0 {
		t.Fatal("Parsed software from empty string")
	}

	// bad string
	bs = []byte("bad")
	software = parseAptListInstalled(bs)
	if len(software) != 0 {
		t.Fatal("Parsed software from bad string")
	}

	// cut off
	bs = []byte("testpackage/stable,now 1.0.0")
	software = parseAptListInstalled(bs)
	if len(software) != 0 {
		t.Fatal("Parsed software from cut off string")
	}

	// one package
	bs = []byte("testpackage/stable,now 1.0.0 amd64 [installed,automatic]")
	software = parseAptListInstalled(bs)
	if len(software) != 1 {
		t.Fatal("Did not parse expected entry")
	}
	if software[0].Name != "testpackage" {
		t.Fatal("Did not parse expected software name")
	}
	if software[0].Version != "1.0.0" {
		t.Fatal("Did not parse expected software version")
	}

	// two packages
	bs = []byte("testpackage/stable,now 1.0.0 amd64 [installed,automatic]\nsecond/stable,now 0.1.0 amd64 [installed]")
	software = parseAptListInstalled(bs)
	if len(software) != 2 {
		t.Fatal("Did not parse expected entries")
	}
	if software[0].Name != "testpackage" {
		t.Fatal("Did not parse expected software name")
	}
	if software[0].Version != "1.0.0" {
		t.Fatal("Did not parse expected software version")
	}
	if software[1].Name != "second" {
		t.Fatal("Did not parse expected software name")
	}
	if software[1].Version != "0.1.0" {
		t.Fatal("Did not parse expected software version")
	}
}

func TestParsePowerShellVersion(t *testing.T) {
	// empty string
	bs := []byte("")
	version := parsePowerShellVersion(bs)
	if len(version) != 0 {
		t.Fatal("Parsed version from empty string")
	}

	// only header
	bs = []byte("VERSION")
	version = parsePowerShellVersion(bs)
	if len(version) != 0 {
		t.Fatal("Parsed version from only header")
	}

	// example version
	bs = []byte("VERSION\n5.1")
	version = parsePowerShellVersion(bs)
	if version != "5.1" {
		t.Fatal("Incorrectly parsed version")
	}
}
