package abi

import (
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/sha3"
)

type JSONFieldType string

const (
	JSONTypeFunc        JSONFieldType = "function"
	JSONTypeConstructor JSONFieldType = "constructor"
	JSONTypeReceive     JSONFieldType = "receive"
	JSONTypeFallback    JSONFieldType = "fallback"
	JSONTypeEvent       JSONFieldType = "event"
	JSONTypeError       JSONFieldType = "error"
)

type StateMutability string

const (
	StateMutabilityPure       StateMutability = "pure"
	StateMutabilityView       StateMutability = "view"
	StateMutabilityNonpayable StateMutability = "nonpayable"
	StateMutabilityPayable    StateMutability = "payable"
)

type JSONField struct {
	Type            JSONFieldType    `json:"type"`
	Name            string           `json:"name"`
	Inputs          []*JSONParam     `json:"inputs"`
	Outputs         []*JSONParam     `json:"outputs"`
	StateMutability *StateMutability `json:"stateMutability"`
	Anonymous       *bool            `json:"anonymous"`
}

type JSONParam struct {
	Name       string       `json:"name"`
	Type       string       `json:"type"`
	Components []*JSONParam `json:"components"`
	Indexed    *bool        `json:"indexed"`
}

// Selector generate function selector string
func Selector(abi string) string {
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
