package abi

import (
	"bytes"
	"math/big"
	"reflect"

	"github.com/libs4go/errors"
	"github.com/libs4go/fixed"
	"golang.org/x/crypto/sha3"
)

var bigIntType = reflect.TypeOf((*big.Int)(nil)).Elem()
var bigIntPtrType = reflect.TypeOf((*big.Int)(nil))
var fixedPtrType = reflect.TypeOf((*fixed.Number)(nil))
var bytesPtrType = reflect.TypeOf((*[]byte)(nil)).Elem()

func paddingLeft(src []byte) []byte {
	paddingZero := len(src) % 32

	if paddingZero == 0 {
		return src
	}

	return append(bytes.Repeat([]byte{0}, 32-paddingZero), src...)

}

func paddingRight(src []byte) []byte {
	paddingZero := len(src) % 32

	if paddingZero == 0 {
		return src
	}

	return append(src, bytes.Repeat([]byte{0}, 32-paddingZero)...)
}

func selector(abi string) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(abi))
	data := hasher.Sum(nil)

	return data[0:4]
}

func bitsCheck(maxBits uint, bits uint) error {
	if bits%8 != 0 {
		return errors.Wrap(ErrBits, "bits % 8 != 0")
	}

	if bits == 0 || bits > maxBits {
		return errors.Wrap(ErrBits, "bits out of range: 0 < bits <= %d", maxBits)
	}

	return nil
}

// Encoder types encoder interface
type Encoder interface {
	Static() bool
	Marshal(value interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type integerEncoder struct {
	bits uint
	sign bool
}

// Integer create Integer encoder
func Integer(sign bool, bits uint) (Encoder, error) {

	err := bitsCheck(256, bits)

	if err != nil {
		return nil, err
	}

	return &integerEncoder{
		bits: bits,
		sign: sign,
	}, nil
}

func (enc *integerEncoder) Static() bool {
	return true
}

func (enc *integerEncoder) marshalBigInt(v *big.Int) ([]byte, error) {

	content := paddingLeft(new(big.Int).Abs(v).Bytes())

	if enc.sign {
		if v.Sign() < 0 {
			content[0] = 0xff
		}
	}

	return content, nil
}

func (enc *integerEncoder) Marshal(value interface{}) ([]byte, error) {

	switch vv := value.(type) {

	case uint, uint8, uint16, uint32, uint64:
		v := reflect.ValueOf(value)
		return enc.marshalBigInt(big.NewInt(0).SetUint64(v.Uint()))
	case int, int8, int16, int32, int64:
		v := reflect.ValueOf(value)
		return enc.marshalBigInt(big.NewInt(v.Int()))
	case *big.Int:
		return enc.marshalBigInt(vv)
	default:
		return nil, errors.Wrap(ErrValue, "invalid value type %v", reflect.TypeOf(value))
	}
}

func (enc *integerEncoder) Unmarshal(data []byte, v interface{}) error {

	if len(data) < 32 {
		return errors.Wrap(ErrLength, "unmarshal input data length < 32")
	}

	negative := false

	var buff [32]byte

	copy(buff[:], data)

	if enc.sign && data[0] == 0xff {
		negative = true
		buff[0] = 0x00
	}

	i := new(big.Int).SetBytes(buff[:])

	if negative {
		i = new(big.Int).Neg(i)
	}

	t := reflect.ValueOf(v)

	if t.Kind() != reflect.Ptr || t.IsNil() {
		return errors.Wrap(ErrValue, "expect int/uint types ptr")
	}

	switch t.Elem().Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		if enc.sign {
			return errors.Wrap(ErrValue, "expect int/uint types ptr")
		}

		t.Elem().Set(reflect.ValueOf(i.Uint64()))

		return nil

	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		if !enc.sign {
			return errors.Wrap(ErrValue, "expect int/uint types ptr")
		}

		t.Elem().Set(reflect.ValueOf(i.Int64()))

		return nil

	default:
		if bigIntType == t.Elem().Type() {
			t.Elem().Set(reflect.ValueOf(i).Elem())
		}

		if bigIntPtrType == t.Elem().Type() {
			t.Elem().Set(reflect.ValueOf(i))
		}

		return errors.Wrap(ErrValue, "expect int/uint types ptr")
	}
}

type boolEncoder struct {
	Encoder
}

func Bool() (Encoder, error) {
	i, err := Integer(false, 8)

	if err != nil {
		return nil, err
	}

	return &boolEncoder{
		Encoder: i,
	}, nil
}

func (enc *boolEncoder) Marshal(value interface{}) ([]byte, error) {
	b, ok := value.(bool)

	if !ok {
		return nil, errors.Wrap(ErrValue, "expect bool")
	}

	if b {
		return enc.Encoder.Marshal(1)
	} else {
		return enc.Encoder.Marshal(0)
	}
}

func (enc *boolEncoder) Unmarshal(data []byte, v interface{}) error {
	var i uint

	if err := enc.Encoder.Unmarshal(data, &i); err != nil {
		return err
	}

	t := reflect.ValueOf(v)

	if t.Kind() != reflect.Ptr || t.IsNil() {
		return errors.Wrap(ErrValue, "expect bool ptr")
	}

	if t.Elem().Kind() != reflect.Bool {
		return errors.Wrap(ErrValue, "expect bool ptr")
	}

	if i == 0 {
		t.Elem().Set(reflect.ValueOf(false))
		return nil
	} else if i == 1 {
		t.Elem().Set(reflect.ValueOf(true))
		return nil
	} else {
		return errors.Wrap(ErrValue, "unmarshal abi data error")
	}
}

type fixedEncoder struct {
	Encoder
	M uint
	N uint
}

func Fixed(sign bool, M, N uint) (Encoder, error) {
	encoder, err := Integer(sign, 256)

	if err != nil {
		return nil, err
	}

	return &fixedEncoder{
		Encoder: encoder,
		M:       M,
		N:       N,
	}, nil
}

func (enc *fixedEncoder) Marshal(value interface{}) ([]byte, error) {
	var err error
	var number *fixed.Number

	switch v := value.(type) {
	case float32, float64:
		number, err = fixed.New(int(enc.N), fixed.Float(reflect.ValueOf(v).Float()))

		if err != nil {
			return nil, errors.Wrap(err, "create fixed error")
		}

	case *big.Float:
		number, err = fixed.New(int(enc.N), fixed.BigFloat(v))

		if err != nil {
			return nil, errors.Wrap(err, "create fixed error")
		}
	case *fixed.Number:
		if v.Decimals != int(enc.N) {
			return nil, errors.Wrap(ErrValue, "input fixed decimals error")
		}

		number = v
	}

	return enc.Encoder.Marshal(number.RawValue)
}

func (enc *fixedEncoder) Unmarshal(data []byte, v interface{}) error {
	var i *big.Int

	err := enc.Encoder.Unmarshal(data, &i)

	if err != nil {
		return err
	}

	t := reflect.ValueOf(v)

	if t.Kind() != reflect.Ptr || t.IsNil() {
		return errors.Wrap(ErrValue, "expect bool ptr")
	}

	number := &fixed.Number{
		RawValue: i,
		Decimals: int(enc.N),
	}

	switch t.Elem().Kind() {
	case reflect.Float32:
		f, _ := number.Float().Float32()
		t.Elem().Set(reflect.ValueOf(f))
		return nil
	case reflect.Float64:
		f, _ := number.Float().Float64()
		t.Elem().Set(reflect.ValueOf(f))
		return nil
	default:
		if t.Elem().Type() == fixedPtrType {
			t.Elem().Set(reflect.ValueOf(number))
			return nil
		}

		return errors.Wrap(ErrValue, "expect fixed ptr")
	}
}

type fixedBytesEncoder struct {
	len uint
}

func FixedBytes(len uint) (Encoder, error) {
	return &fixedBytesEncoder{
		len: len,
	}, nil
}

func (enc *fixedBytesEncoder) Static() bool {
	return true
}

func (enc *fixedBytesEncoder) Marshal(value interface{}) ([]byte, error) {
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Slice && v.Elem().Kind() != reflect.Uint8 {
		return nil, errors.Wrap(ErrValue, "expect [%d]byte", enc.len)
	}

	return paddingRight(v.Bytes()), nil
}

func (enc *fixedBytesEncoder) Unmarshal(data []byte, v interface{}) error {

	data = bytes.TrimSuffix(data, []byte{0})

	if len(data) != int(enc.len) {
		return errors.Wrap(ErrValue, "abi data len error")
	}

	t := reflect.ValueOf(v)

	if t.Kind() != reflect.Ptr || t.IsNil() {
		return errors.Wrap(ErrValue, "expect []byte ptr")
	}

	if t.Elem().Type() != bytesPtrType {
		return errors.Wrap(ErrValue, "expect []byte ptr")
	}

	t.Elem().Set(reflect.ValueOf(data[:enc.len]))

	return nil
}
