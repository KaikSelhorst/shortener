package service

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const base = uint64(len(alphabet)) // 62

const phi64 = uint64(0x9e3779b97f4a7c15)

func modInverse64(a uint64) uint64 {
	x := a
	x *= 2 - a*x
	x *= 2 - a*x
	x *= 2 - a*x
	x *= 2 - a*x
	x *= 2 - a*x
	return x
}

type ShortcodeService struct {
	seed     uint64
	phi64Inv uint64
}

func NewShortcodeService(secret []byte) (*ShortcodeService, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("shortcode: secret must not be empty")
	}
	h := sha256.Sum256(secret)
	seed := binary.BigEndian.Uint64(h[:8])
	return &ShortcodeService{seed: seed, phi64Inv: modInverse64(phi64)}, nil
}

func (svc *ShortcodeService) scramble(id uint64) uint64  { return id*phi64 ^ svc.seed }
func (svc *ShortcodeService) unscramble(v uint64) uint64 { return (v ^ svc.seed) * svc.phi64Inv }

func (svc *ShortcodeService) GenerateShortCode(id uint64) (string, error) {
	if id == 0 {
		return "", fmt.Errorf("shortcode: id must be greater than zero")
	}
	v := svc.scramble(id)
	if v == 0 {
		return "", fmt.Errorf("shortcode: id %d produces zero after scramble; rotate SHORTCODE_SECRET", id)
	}
	var buf [11]byte
	n := 0
	for v > 0 {
		buf[n] = alphabet[v%base]
		v /= base
		n++
	}
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf[:n]), nil
}

func (svc *ShortcodeService) DecodeShortCode(code string) uint64 {
	if code == "" || len(code) > 11 {
		return 0
	}
	var v uint64
	for _, c := range code {
		idx := strings.IndexRune(alphabet, c)
		if idx < 0 {
			return 0
		}
		v = v*base + uint64(idx)
	}
	return svc.unscramble(v)
}
