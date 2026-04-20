package service

import "github.com/sqids/sqids-go"

var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var s, _ = sqids.New(sqids.Options{Alphabet: alphabet})

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

func GenerateShortCode(id uint64) (string, error) {
	a := numberToDigits(id)
	return s.Encode(a)
}

func DecodeShortCode(code string) uint64 {
	idSlice := s.Decode(code)
	return digitsToNumber(idSlice)
}
