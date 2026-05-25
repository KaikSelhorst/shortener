package service

import (
	"fmt"

	"github.com/sqids/sqids-go"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// ShortcodeService encodes and decodes short codes using a fixed alphabet and
// a custom base-10 digit mapping. The encoding logic (numberToDigits /
// digitsToNumber) is a stable contract with the database — changing it would
// invalidate all existing short_code values.
type ShortcodeService struct {
	s *sqids.Sqids
}

func NewShortcodeService() (*ShortcodeService, error) {
	s, err := sqids.New(sqids.Options{Alphabet: alphabet})
	if err != nil {
		return nil, fmt.Errorf("shortcode: failed to initialize sqids: %w", err)
	}
	return &ShortcodeService{s: s}, nil
}

func (svc *ShortcodeService) GenerateShortCode(id uint64) (string, error) {
	return svc.s.Encode(numberToDigits(id))
}

func (svc *ShortcodeService) DecodeShortCode(code string) uint64 {
	return digitsToNumber(svc.s.Decode(code))
}

func numberToDigits(id uint64) []uint64 {
	var digits []uint64
	for id > 0 {
		digits = append(digits, id%10)
		id /= 10
	}
	return digits
}

func digitsToNumber(digits []uint64) uint64 {
	var id uint64
	for i := len(digits) - 1; i >= 0; i-- {
		id = id*10 + digits[i]
	}
	return id
}
