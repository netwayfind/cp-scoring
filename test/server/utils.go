package main

import (
	"math/rand"
	"strings"

	"github.com/netwayfind/cp-scoring/test/model"
)

func randHexStr(length int) string {
	var output strings.Builder
	for i := 0; i < length; i++ {
		random := rand.Intn(len(model.KeyCharset))
		randomChar := model.KeyCharset[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}
