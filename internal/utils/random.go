package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
)

func RandomUUID() string {
	return uuid.NewString()
}

func RandomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
