package abi

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"reflect"
	"strings"

	"github.com/libs4go/errors"
	"github.com/libs4go/fixed"
)

func paddingLeft(src []byte) string {
	paddingZero := len(src) % 32

	right := hex.EncodeToString(src)

	if paddingZero == 0 {
		return right
	}

	return strings.Repeat("0", (32-paddingZero)*2) + right
}

func paddingRight(src []byte) string {
	paddingZero := len(src) % 32

	right := hex.EncodeToString(src)

	if paddingZero == 0 {
		return right
	}

	return right + strings.Repeat("0", (32-paddingZero)*2)
}

type ABIType interface {
	Dynamic() bool
	Encode(value interface{}) (string, error)
}

// IntegerType int/uint type object
type IntegerType struct {
	LengthOfBits uint
	Signed       bool
}

func NewIntegerType(bits uint, signed bool) (ABIType, error) {

	if bits == 0 || bits > 256 || bits%8 != 0 {
		return nil, errors.Wrap(ErrBits, "invalid integer bits %d", bits)
	}

	return &IntegerType{
		LengthOfBits: bits,
		Signed:       signed,
	}, nil
}

func (t *IntegerType) Dynamic() bool {
	return false
}

func (t *IntegerType) Encode(value interface{}) (string, error) {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:

		i := reflect.ValueOf(v).Int()

		neg := false

		if i < 0 {
			neg = true
			i = -1
		}

		var buff [8]byte
		binary.BigEndian.PutUint64(buff[:], uint64(i))

		content := paddingLeft(buff[:])

		if neg {
			return "ff" + string([]byte(content)[2:]), nil
		}

	case uint, uint8, uint16, uint32, uint64:
		var buff [8]byte
		binary.BigEndian.PutUint64(buff[:], uint64(reflect.ValueOf(v).Uint()))

		return paddingLeft(buff[:]), nil

	case *big.Int:

		content := paddingLeft(big.NewInt(0).Abs(v).Bytes())

		if t.Signed && v.Sign() < 0 {
			return "ff" + string([]byte(content)[2:]), nil
		}

		return content, nil
	}

	return "", errors.Wrap(ErrValue, "IntegerType not support value type %v", reflect.TypeOf(value))
}

// NewAddressType create new address type
func NewAddressType() (ABIType, error) {
	return NewIntegerType(160, false)
}

// equivalent to uint8 restricted to the values 0 and 1
type BoolType struct{}

func NewBoolType() (ABIType, error) {
	return &BoolType{}, nil
}

func (t *BoolType) Dynamic() bool {
	return false
}

func (t *BoolType) Encode(value interface{}) (string, error) {
	b, ok := value.(bool)

	if !ok {
		return "", errors.Wrap(ErrValue, "BoolType not support value type %v", reflect.TypeOf(value))
	}
	var v *big.Int

	if b {
		v = big.NewInt(1)
	} else {
		v = big.NewInt(0)
	}

	return paddingLeft(v.Bytes()), nil
}

// ixed-point decimal number of bits
type FixedType struct {
	Decimals     uint
	LengthOfBits uint
	Signed       bool
}

func NewFixedType(bits uint, decimals uint, signed bool) (ABIType, error) {
	if bits == 0 || bits > 256 || bits%8 != 0 {
		return nil, errors.Wrap(ErrBits, "invalid integer bits %d", bits)
	}

	return &FixedType{
		Decimals:     decimals,
		LengthOfBits: bits,
		Signed:       signed,
	}, nil
}

func (t *FixedType) Dynamic() bool {
	return false
}

func (t *FixedType) Encode(value interface{}) (string, error) {
	var c *fixed.Number
	var err error
	switch v := value.(type) {
	case float32, float64:

		c, err = fixed.New(int(t.Decimals), fixed.Float(reflect.ValueOf(value).Float()))

	case *big.Float:
		c, err = fixed.New(int(t.Decimals), fixed.BigFloat(v))
	default:
		return "", errors.Wrap(ErrValue, "FixedType not support value type %v", reflect.TypeOf(value))
	}

	if err != nil {
		return "", err
	}

	content := paddingLeft(big.NewInt(0).Abs(c.RawValue).Bytes())

	if t.Signed && c.Sign() < 0 {
		return "ff" + string([]byte(content)[2:]), nil
	}

	return content, nil
}

// FixedBytesType .
type FixedBytesType struct {
	LengthOfBytes uint
}

func NewFixedBytesType(length uint) (ABIType, error) {
	return &FixedBytesType{
		LengthOfBytes: length,
	}, nil
}

func (t *FixedBytesType) Dynamic() bool {
	return false
}

func (t *FixedBytesType) Encode(value interface{}) (string, error) {
	c := reflect.ValueOf(value).Bytes()

	if len(c) != int(t.LengthOfBytes) {
		return "", errors.Wrap(ErrFixedBytes, "value length %d", len(c))
	}

	content := paddingRight(c)

	return content, nil
}

type FixedArrayType struct {
	Element ABIType
	Length  uint
}

func NewFixedArrayType(length uint) (ABIType, error) {
	return &FixedBytesType{
		LengthOfBytes: length,
	}, nil
}

func (t *FixedArrayType) Dynamic() bool {
	return false
}

func (t *FixedArrayType) Encode(value interface{}) (string, error) {
	if t.Element.Dynamic() {
		return t.EncodeDynamic(value)
	}

	return t.EncodeStatic(value)
}

func (t *FixedArrayType) EncodeDynamic(value interface{}) (string, error) {
	return "", nil
}

func (t *FixedArrayType) EncodeStatic(value interface{}) (string, error) {
	return "", nil
}

type BytesType struct {
}

type StringType struct {
}

type ArrayType struct {
	Element ABIType
}

type TupleType struct {
	Fields []ABIType
}
