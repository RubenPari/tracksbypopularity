package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomString(len int) string {
	buff := make([]byte, len)
	_, _ = rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str
}
