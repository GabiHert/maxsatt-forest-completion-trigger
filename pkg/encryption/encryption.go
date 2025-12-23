package encryption

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

// GenerateSecret generates a base64-encoded secret
func GenerateSecret() (*string, error) {
	salt := make([]byte, 64)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	secret := base64.RawStdEncoding.EncodeToString(salt)

	return &secret, nil
}

// HashString hashes the input string using the provided secret with Argon2
func HashString(input string, secret string) (*string, error) {
	secretBytes, err := base64.RawStdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey([]byte(input), secretBytes, 1, 64*1024, 4, 32)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return &encodedHash, nil
}
