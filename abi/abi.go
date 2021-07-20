package abi

import (
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Encode encode method string
func Encode(abi string) string {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(abi))
	data := hasher.Sum(nil)

	return hex.EncodeToString(data[0:4])
}

// PackNumeric .
func PackNumeric(value string, bytes int) string {
	return packNumeric(value, bytes)
}

func packNumeric(value string, bytes int) string {
	if value == "" {
		value = "0x0"
	}

	value = strings.TrimPrefix(value, "0x")

	chars := bytes * 2

	n := len(value)
	if n%chars == 0 {
		return value
	}
	return strings.Repeat("0", chars-n%chars) + value
}
