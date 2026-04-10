// femabras/backend/internal/utils/crypto.go
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

var cachedKey []byte

func init() {
	key := os.Getenv("ENCRYPTION_KEY")
	if key != "" {
		if len(key) != 32 {
			panic("CRITICAL: ENCRYPTION_KEY must be exactly 32 bytes")
		}
		cachedKey = []byte(key)
	}
}

func EncryptPayload(payload []byte) (string, error) {
	if cachedKey == nil {
		return "", errors.New("encryption key not set")
	}

	block, err := aes.NewCipher(cachedKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, payload, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptPayload(encryptedStr string) ([]byte, error) {
	if cachedKey == nil {
		return nil, errors.New("encryption key not set")
	}

	data, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(cachedKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}
