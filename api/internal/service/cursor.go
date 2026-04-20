package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

func deriveKey(secret string) []byte {
	h := sha256.Sum256([]byte(secret))
	return h[:]
}

func EncodeCursor(direction string, id string, secret string) (string, error) {
	block, err := aes.NewCipher(deriveKey(secret))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	plaintext := fmt.Sprintf("%s:%s", direction, id)
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(sealed), nil
}

func DecodeCursor(cursor string, secret string) (direction string, id string, err error) {
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return "", "", err
	}

	block, err := aes.NewCipher(deriveKey(secret))
	if err != nil {
		return "", "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", "", errors.New("invalid cursor")
	}

	plaintext, err := gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return "", "", errors.New("invalid cursor")
	}

	parts := strings.SplitN(string(plaintext), ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid cursor format")
	}

	return parts[0], parts[1], nil
}
