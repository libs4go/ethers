package abi

import (
	"encoding/hex"
	"math/big"
	"reflect"
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

	require.Equal(t, l, uint(len(data)))

	require.Equal(t, val, b)
}

func IntegerCheck(t *testing.T, signed bool, bits uint, val interface{}, hexstring string) {
	encoder, err := Integer(signed, bits)

	require.NoError(t, err)

	data, err := encoder.Marshal(val)

	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(data), hexstring)

	switch val.(type) {
	case *big.Int:
		v := reflect.New(reflect.TypeOf(val).Elem())

		l, err := encoder.Unmarshal(data, v.Interface())

		require.NoError(t, err)

		require.Equal(t, l, uint(len(data)))

		require.Equal(t, val, v.Interface())
	default:
		v := reflect.New(reflect.TypeOf(val))

		l, err := encoder.Unmarshal(data, v.Interface())

		require.NoError(t, err)

		require.Equal(t, l, uint(len(data)))

		require.Equal(t, val, v.Elem().Interface())
	}

}

func BytesCheck(t *testing.T, encoder Encoder, val []byte, hexstring string) {
	data, err := encoder.Marshal(val)

	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(data), hexstring)

	v := reflect.New(reflect.TypeOf(val))

	l, err := encoder.Unmarshal(data, v.Interface())

	require.NoError(t, err)

	require.Equal(t, l, uint(len(data)))

	require.Equal(t, val, v.Elem().Interface())
}

func StringCheck(t *testing.T, encoder Encoder, val string, hexstring string) {
	data, err := encoder.Marshal(val)

	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(data), hexstring)

	var v string

	l, err := encoder.Unmarshal(data, &v)

	require.NoError(t, err)

	require.Equal(t, l, uint(len(data)))

	require.Equal(t, val, v)
}

func ArrayCheck(t *testing.T, encoder Encoder, val interface{}, hexstring string) {
	data, err := encoder.Marshal(val)

	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(data), hexstring)

	var v string

	l, err := encoder.Unmarshal(data, &v)

	require.NoError(t, err)

	require.Equal(t, l, uint(len(data)))

	require.Equal(t, val, v)
}

func TestBool(t *testing.T) {
	encoder, err := Bool()

	require.NoError(t, err)

	BoolCheck(t, encoder, true, "0000000000000000000000000000000000000000000000000000000000000001")

	BoolCheck(t, encoder, false, "0000000000000000000000000000000000000000000000000000000000000000")
}

func TestInteger(t *testing.T) {
	IntegerCheck(t, false, 8, uint8(2), "0000000000000000000000000000000000000000000000000000000000000002")

	IntegerCheck(t, false, 32, big.NewInt(1), "0000000000000000000000000000000000000000000000000000000000000001")

	IntegerCheck(t, false, 8, big.NewInt(257), "0000000000000000000000000000000000000000000000000000000000000101")

	IntegerCheck(t, true, 256, big.NewInt(-16), "ff00000000000000000000000000000000000000000000000000000000000010")
}

func TestBytes(t *testing.T) {
	bytesEncoder, err := Bytes()

	require.NoError(t, err)

	BytesCheck(t, bytesEncoder, []byte{0xf0, 0xf0, 0xf0}, "0000000000000000000000000000000000000000000000000000000000000003f0f0f00000000000000000000000000000000000000000000000000000000000")
}

func TestFixedBytes(t *testing.T) {
	bytesEncoder, err := FixedBytes(3)

	require.NoError(t, err)

	BytesCheck(t, bytesEncoder, []byte{0xf0, 0xf0, 0xf0}, "f0f0f00000000000000000000000000000000000000000000000000000000000")
}

func TestString(t *testing.T) {
	encoder, err := String()

	require.NoError(t, err)
	StringCheck(t, encoder, "foobar", "0000000000000000000000000000000000000000000000000000000000000006666f6f6261720000000000000000000000000000000000000000000000000000")
}

func TestArray(t *testing.T) {

	elemEncoder, err := String()

	require.NoError(t, err)

	encoder, err := Array(elemEncoder)

	require.NoError(t, err)

	v := []string{"hello", "foobar", "foobar"}

	buff, err := encoder.Marshal(v)

	require.NoError(t, err)

	packed := "0000000000000000000000000000000000000000000000000000000000000003" + // len(array) = 2
		"0000000000000000000000000000000000000000000000000000000000000060" + // offset 64 to i = 0
		"00000000000000000000000000000000000000000000000000000000000000a0" + // offset 128 to i = 1
		"00000000000000000000000000000000000000000000000000000000000000e0" + // offset 160 to i = 2
		"0000000000000000000000000000000000000000000000000000000000000005" + // len(str[0]) = 5
		"68656c6c6f000000000000000000000000000000000000000000000000000000" + // str[0]
		"0000000000000000000000000000000000000000000000000000000000000006" + // len(str[1]) = 6
		"666f6f6261720000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000006" + // len(str[1]) = 6
		"666f6f6261720000000000000000000000000000000000000000000000000000"

	require.Equal(t, packed, hex.EncodeToString(buff))
}

func TestTuple(t *testing.T) {
	packed := "0000000000000000000000000000000000000000000000000000000000000001" + // struct[a]
		"0000000000000000000000000000000000000000000000000000000000000001" + // struct[b]
		"ff00000000000000000000000000000000000000000000000000000000000001" + // struct[c]
		"0000000000000000000000000000000000000000000000000000000000000001" + // struct[d]
		"0100000000000000000000000000000000000000000000000000000000000000" + // struct[e] array[0][0]
		"0200000000000000000000000000000000000000000000000000000000000000" + // struct[e] array[0][1]
		"0300000000000000000000000000000000000000000000000000000000000000" + // struct[e] array[0][2]
		"0300000000000000000000000000000000000000000000000000000000000000" + // struct[e] array[1][0]
		"0400000000000000000000000000000000000000000000000000000000000000" + // struct[e] array[1][1]
		"0500000000000000000000000000000000000000000000000000000000000000" // struct[e] array[1][2]

	ienc, err := Integer(true, 64)

	require.NoError(t, err)

	benc, err := Bool()

	require.NoError(t, err)

	byte32, err := FixedBytes(32)

	require.NoError(t, err)

	array3, err := FixedArray(byte32, 3)

	require.NoError(t, err)

	array2, err := Array(array3)

	require.NoError(t, err)

	enc, err := Tuple(ienc, ienc, ienc, benc, array2)

	require.NoError(t, err)

	buff, err := enc.Marshal([]interface{}{1, big.NewInt(1), big.NewInt(-1), true, [][3][32]byte{{{1}, {2}, {3}}, {{3}, {4}, {5}}}})

	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(buff), packed)
}
