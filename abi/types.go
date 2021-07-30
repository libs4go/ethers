package abi

import (
	"bytes"
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

func NewFixedArrayType(Element ABIType, length uint) (ABIType, error) {
	return &FixedArrayType{
		Element: Element,
		Length:  length,
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

	intType, err := NewIntegerType(256, false)

	if err != nil {
		return "", err
	}

	v := reflect.ValueOf(value)

	if v.Type().Kind() != reflect.Slice {
		return "", errors.Wrap(ErrValue, "expect slice got %v", v.Type())
	}

	if v.Len() != int(t.Length) {
		return "", errors.Wrap(ErrLength, "fixed array length err %d", t.Length)
	}

	var headerWriter bytes.Buffer
	var contentWriter bytes.Buffer

	var offset = 32 * t.Length

	for i := 0; i < v.Len(); i++ {

		header, err := intType.Encode(offset)

		if err != nil {
			return "", errors.Wrap(err, "write header error")
		}

		_, err = headerWriter.WriteString(header)

		if err != nil {
			return "", errors.Wrap(err, "write header error")
		}

		content, err := intType.Encode(offset)

		if err != nil {
			return "", errors.Wrap(err, "write content error")
		}

		_, err = contentWriter.WriteString(content)

		if err != nil {
			return "", errors.Wrap(err, "write content error")
		}

		offset += uint(len(content))
	}

	return headerWriter.String() + contentWriter.String(), nil

}

func (t *FixedArrayType) EncodeStatic(value interface{}) (string, error) {
	v := reflect.ValueOf(value)

	if v.Type().Kind() != reflect.Slice {
		return "", errors.Wrap(ErrValue, "expect slice got %v", v.Type())
	}

	if v.Len() != int(t.Length) {
		return "", errors.Wrap(ErrLength, "fixed array length err %d", t.Length)
	}

	var writer strings.Builder

	for i := 0; i < v.Len(); i++ {
		data, err := t.Element.Encode(v.Index(i).Interface())

		if err != nil {
			return "", errors.Wrap(err, "encode array element %v error", t.Element)
		}

		_, err = writer.WriteString(data)

		if err != nil {
			return "", errors.Wrap(err, "write string error")
		}
	}

	return writer.String(), nil
}

type BytesType struct {
}

func NewBytesType() (ABIType, error) {
	return &BytesType{}, nil
}

func (t *BytesType) Dynamic() bool {
	return true
}

func (t *BytesType) Encode(value interface{}) (string, error) {

	v := reflect.ValueOf(value)

	if v.Type().Kind() != reflect.Slice {
		return "", errors.Wrap(ErrValue, "expect slice got %v", v.Type())
	}

	if v.Type().Elem().Kind() != reflect.Uint8 {
		return "", errors.Wrap(ErrValue, "expect slice got %v", v.Type())
	}

	intType, err := NewIntegerType(256, false)

	if err != nil {
		return "", err
	}

	len, err := intType.Encode(len(v.Bytes()))

	if err != nil {
		return "", errors.Wrap(err, "enc bytes len error")
	}

	return len + paddingRight(v.Bytes()), nil
}

type StringType struct {
}

func NewStringType() (ABIType, error) {
	return &StringType{}, nil
}

func (t *StringType) Dynamic() bool {
	return true
}

func (t *StringType) Encode(value interface{}) (string, error) {

	s, ok := value.(string)

	if !ok {
		return "", errors.Wrap(ErrValue, "expect string got %v", value)
	}

	bytesType, err := NewBytesType()

	if err != nil {
		return "", err
	}

	return bytesType.Encode([]byte(s))
}

type ArrayType struct {
	Element ABIType
}

func NewArrayType() (ABIType, error) {
	return &ArrayType{}, nil
}

func (t *ArrayType) Dynamic() bool {
	return true
}

func (t *ArrayType) Encode(value interface{}) (string, error) {

	v := reflect.ValueOf(value)

	if v.Type().Kind() != reflect.Slice {
		return "", errors.Wrap(ErrValue, "expect slice got %v", v.Type())
	}

	arrayType, err := NewFixedArrayType(t.Element, uint(v.Len()))

	if err != nil {
		return "", err
	}

	content, err := arrayType.Encode(value)

	if err != nil {
		return "", err
	}

	intType, err := NewIntegerType(256, false)

	if err != nil {
		return "", err
	}

	len, err := intType.Encode(uint(v.Len()))

	if err != nil {
		return "", errors.Wrap(err, "encode array len error")
	}

	return len + content, nil
}

type TupleType struct {
	Fields []ABIType
}

func NewTupleType(fields []ABIType) (ABIType, error) {
	return &TupleType{
		Fields: fields,
	}, nil
}

func (t *TupleType) Dynamic() bool {
	return true
}

func (t *TupleType) Encode(value interface{}) (string, error) {
	var headerWriter bytes.Buffer
	var contentWriter bytes.Buffer

}
