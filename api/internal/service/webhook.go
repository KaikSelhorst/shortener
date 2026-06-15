package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
)

type WebhookService struct {
	repo      repository.WebhookRepo
	secretKey []byte
}

func NewWebhookService(repo repository.WebhookRepo, secretKey []byte) *WebhookService {
	return &WebhookService{repo: repo, secretKey: secretKey}
}

// EncryptSecret encrypts a plaintext secret using AES-256-GCM.
// The output is base64(nonce + ciphertext) stored in the database.
func (s *WebhookService) EncryptSecret(plaintext string) (string, error) {
	key := sha256.Sum256(s.secretKey) // normalise to 32 bytes
	block, err := aes.NewCipher(key[:])
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
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSecret reverses EncryptSecret.
func (s *WebhookService) DecryptSecret(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	key := sha256.Sum256(s.secretKey)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// Sign returns the HMAC-SHA256 signature of payload as "sha256=<hex>".
func Sign(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// GenerateWebhookSecret returns a random 32-byte hex secret.
func GenerateWebhookSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Dispatch inserts delivery rows for every enabled webhook in the project
// that subscribes to the given event. It only writes to the database —
// the worker goroutine handles the actual HTTP delivery.
func (s *WebhookService) Dispatch(ctx context.Context, projectID int64, event string, payload any) {
	webhooks, err := s.repo.ListByProjectAndEvent(ctx, projectID, event)
	if err != nil || len(webhooks) == 0 {
		return
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return
	}

	for _, w := range webhooks {
		d := &model.WebhookDelivery{
			WebhookID: w.ID,
			Event:     event,
			Payload:   raw,
		}
		_ = s.repo.CreateDelivery(ctx, d)
	}
}
