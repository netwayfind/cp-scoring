package main

import (
	"math/rand"
	"strings"

	"github.com/netwayfind/cp-scoring/model"
	"golang.org/x/crypto/bcrypt"
)

func checkPasswordHash(cleartext string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(cleartext))
	if err != nil {
		return false
	}
	return true
}

func hashPassword(cleartext string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(cleartext), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func randHexStr(length int) string {
	var output strings.Builder
	for i := 0; i < length; i++ {
		random := rand.Intn(len(model.KeyCharset))
		randomChar := model.KeyCharset[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}
