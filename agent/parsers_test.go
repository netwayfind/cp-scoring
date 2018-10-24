package main

import (
	"testing"

	"github.com/sumwonyuno/cp-scoring/model"
)

func TestEmptyEtcPasswd(t *testing.T) {
	bs := []byte("")
	users := parseEtcPasswd(bs)
	if len(users) != 0 {
		t.Error("Parsed users out of empty string")
	}
}
func TestBadEtcPasswd(t *testing.T) {
	bs := []byte("bad")
	users := parseEtcPasswd(bs)
	if len(users) != 0 {
		t.Error("Parsed users out of bad string")
	}

	// cut off
	bs = []byte("root:x:0:0:roo")
	users = parseEtcPasswd(bs)
	if len(users) != 0 {
		t.Error("Parsed users out of incomplete string")
	}
}

func TestEtcPasswd(t *testing.T) {
	// one user
	bs := []byte("root:x:0:0:root:/root:/bin/bash")
	users := parseEtcPasswd(bs)
	if len(users) != 1 {
		t.Error("Did not parse expected user")
	}
	user := users["root"]
	if user.Name != "root" {
		t.Error("Unexpected user name")
	}
	if user.ID != "0" {
		t.Error("Unexpected user ID")
	}
	if user.AccountPresent != true {
		t.Error("Account should be present")
	}

	// two users
	bs = []byte("root:x:0:0:root:/root:/bin/bash\nuser:x:1000:1000:user:/home/user:/bin/bash")
	users = parseEtcPasswd(bs)
	if len(users) != 2 {
		t.Error("Did not parse expected users")
	}
	user1 := users["root"]
	if user1.Name != "root" {
		t.Error("Unexpected user name")
	}
	if user1.ID != "0" {
		t.Error("Unexpected user ID")
	}
	if user1.AccountPresent != true {
		t.Error("Account should be present")
	}
	user2 := users["user"]
	if user2.Name != "user" {
		t.Error("Unexpected user name")
	}
	if user2.ID != "1000" {
		t.Error("Unexpected user ID")
	}
	if user2.AccountPresent != true {
		t.Error("Account should be present")
	}
}

func TestEmptyEtcShadow(t *testing.T) {
	bs := []byte("")
	users := parseEtcShadow(bs)
	if len(users) != 0 {
		t.Error("Parsed users out of empty string")
	}
}

func TestBadEtcShadow(t *testing.T) {
	bs := []byte("no")
	users := parseEtcShadow(bs)
	if len(users) != 0 {
		t.Error("Parsed users out of bad string")
	}

	// cut off
	bs = []byte("root:*:1648")
	users = parseEtcShadow(bs)
	if len(users) != 0 {
		t.Error("Parsed users out of incomplete string")
	}
}

func TestEtcShadow(t *testing.T) {
	// one user
	bs := []byte("root:*:16482:0:99999:7:::")
	users := parseEtcShadow(bs)
	if len(users) != 1 {
		t.Error("Did not parse expected user")
	}
	user := users["root"]
	if user.AccountExpires != false {
		t.Error("Account is expected to expire")
	}
	if user.AccountActive != true {
		t.Error("Account is expected to be active")
	}
	if user.PasswordExpires != false {
		t.Error("Account password is expected to not expire")
	}
	// 16482 * 86400 (seconds)
	if user.PasswordLastSet != 1424044800 {
		t.Error("Account last password set not expected value")
	}

	// multiple users
	bs = []byte("user1:*:16000:0:99999:7:::\nuser2:$6$notahash:16482:0:123:7:::\nuser3:!$6$notahash:16482:0:99999:7:::\nuser4:$6$notahash:16482:0:99999:7::17000:")
	users = parseEtcShadow(bs)
	if len(users) != 4 {
		t.Error("Did not parse expected users")
	}
	user1 := users["user1"]
	if user1.AccountExpires != false {
		t.Error("Account is expected to not expire")
	}
	if user1.AccountActive != true {
		t.Error("Account is expected to be active")
	}
	if user1.PasswordExpires != false {
		t.Error("Account password is expected to not expire")
	}
	// 16000 * 86400 (seconds)
	if user1.PasswordLastSet != 1382400000 {
		t.Error("Account last password set not expected value")
	}
	user2 := users["user2"]
	if user2.AccountExpires != false {
		t.Error("Account is expected to not expire")
	}
	if user2.AccountActive != true {
		t.Error("Account is expected to be active")
	}
	if user2.PasswordExpires != true {
		t.Error("Account password is expected to expire")
	}
	// 16482 * 86400 (seconds)
	if user2.PasswordLastSet != 1424044800 {
		t.Error("Account last password set not expected value")
	}
	user3 := users["user3"]
	if user3.AccountExpires != false {
		t.Error("Account is expected to not expire")
	}
	if user3.AccountActive != false {
		t.Error("Account is expected to be active")
	}
	if user3.PasswordExpires != false {
		t.Error("Account password is expected to not expire")
	}
	// 16482 * 86400 (seconds)
	if user3.PasswordLastSet != 1424044800 {
		t.Error("Account last password set not expected value")
	}
	user4 := users["user4"]
	if user4.AccountExpires != true {
		t.Error("Account is expected to expire")
	}
	if user4.AccountActive != true {
		t.Error("Account is expected to be not active/locked")
	}
	if user4.PasswordExpires != false {
		t.Error("Account password is expected to not expire")
	}
	// 16482 * 86400 (seconds)
	if user4.PasswordLastSet != 1424044800 {
		t.Error("Account last password set not expected value")
	}
}

func TestMergeUserMaps(t *testing.T) {
	// empty maps
	usersPasswd := make(map[string]model.User)
	usersShadow := make(map[string]model.User)

	result := mergeUserMaps(usersPasswd, usersShadow)
	if len(result) != 0 {
		t.Error("Expected no users")
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
		t.Error("Expected to have 1 user entry")
	}
	user := result[0]
	if user.Name != "bob" {
		t.Error("Expected user name is bob")
	}
	if user.ID != "175" {
		t.Error("Expected user ID is 175")
	}
	if user.AccountPresent != true {
		t.Error("Expected user is present")
	}
	if user.AccountExpires != false {
		t.Error("Expected user account not to expire")
	}
	if user.AccountActive != true {
		t.Error("Expected user is active")
	}
	if user.PasswordExpires != true {
		t.Error("Expected user password to expire")
	}
	if user.PasswordLastSet != 1000 {
		t.Error("Unexpected password last set value")
	}
}
