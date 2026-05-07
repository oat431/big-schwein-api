package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateVerifyToken creates a cryptographically random hex-encoded verification token.
func GenerateVerifyToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("generate verify token: %w", err)
	}
	return hex.EncodeToString(tokenBytes), nil
}
