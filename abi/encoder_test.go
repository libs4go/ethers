package abi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func BoolCheck(t *testing.T, encoder Encoder, val bool, hexstring string) {
	data, err := encoder.Marshal(val)

	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(data), hexstring)

	var b bool

	l, err := encoder.Unmarshal(data, &b)

	require.NoError(t, err)

	require.Equal(t, l, uint(32))

	require.Equal(t, val, b)
}

func TestBool(t *testing.T) {
	encoder, err := Bool()

	require.NoError(t, err)

	BoolCheck(t, encoder, true, "0000000000000000000000000000000000000000000000000000000000000001")

	BoolCheck(t, encoder, false, "0000000000000000000000000000000000000000000000000000000000000000")
}

func TestInteger(t *testing.T) {

}
