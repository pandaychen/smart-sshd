package util

import (
	"encoding/base64"
	"math/rand"
)

const (
	DEFAULT_TOKEN = "abcdefghABCDFGH"
)

//warning: should init seed
func GetRandomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return DEFAULT_TOKEN
	}
	return base64.StdEncoding.EncodeToString(b)
}
