package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	randBytes := make([]byte, 32)
	if _, err := rand.Read(randBytes); err != nil {
		return "", fmt.Errorf("MakeRefreshToken: %v", err)
	}
	return hex.EncodeToString(randBytes), nil
}
