package secrets

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecureKey(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
