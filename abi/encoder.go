package abi

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/libs4go/errors"
	"github.com/libs4go/fixed"
	"golang.org/x/crypto/sha3"
)

var bigIntType = reflect.TypeOf((*big.Int)(nil)).Elem()
var fixedType = reflect.TypeOf((*fixed.Number)(nil)).Elem()
var stringType = reflect.TypeOf((*string)(nil)).Elem()

func paddingLen(origin uint) uint {
	l := ((origin + 31) / 32) * 32

	if l == 0 {
		return 32
	}

	return l
}

func paddingLeft(src []byte) []byte {
	paddingZero := len(src) % 32

	if paddingZero == 0 && len(src) != 0 {
		return src
	}

	return append(bytes.Repeat([]byte{0}, 32-paddingZero), src...)

}

func paddingRight(src []byte) []byte {
	paddingZero := len(src) % 32

	if paddingZero == 0 && len(src) != 0 {
		return src
	}

	return append(src, bytes.Repeat([]byte{0}, 32-paddingZero)...)
}

func bitsCheck(maxBits uint, bits uint) error {
	if bits%8 != 0 {
		return errors.Wrap(ErrBits, "bits %% 8 != 0")
	}

	if bits == 0 || bits > maxBits {
		return errors.Wrap(ErrBits, "bits out of range: 0 < bits <= %d", maxBits)
	}

	return nil
}

// Selector function selector
func Selector(abi string) []byte {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(abi))
	data := hasher.Sum(nil)

	return data[0:4]
}

// Encoder types encoder interface
type Encoder interface {
	Static() bool
	Marshal(value interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) (uint, error)
	fmt.Stringer
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

func (enc *integerEncoder) String() string {
	if enc.sign {
		return fmt.Sprintf("int%d", enc.bits)
	} else {
		return fmt.Sprintf("uint%d", enc.bits)
	}

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

func (enc *integerEncoder) Unmarshal(data []byte, v interface{}) (uint, error) {

	if len(data) < 32 {
		return 0, errors.Wrap(ErrLength, "unmarshal input data length < 32")
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
		return 0, errors.Wrap(ErrValue, "expect int/uint types ptr")
	}

	switch t.Elem().Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		if enc.sign {
			return 0, errors.Wrap(ErrValue, "expect int/uint types ptr")
		}

		t.Elem().Set(reflect.ValueOf(i.Uint64()).Convert(t.Elem().Type()))

		return 32, nil

	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		if !enc.sign {
			return 0, errors.Wrap(ErrValue, "expect int/uint types ptr")
		}

		t.Elem().Set(reflect.ValueOf(i.Int64()).Convert(t.Elem().Type()))

		return 32, nil

	default:
		if bigIntType == t.Elem().Type() {
			t.Elem().Set(reflect.ValueOf(i).Elem())
			return 32, nil
		}

		if reflect.Ptr == t.Type().Elem().Kind() && t.Type().Elem().Elem() == bigIntType {
			t.Elem().Set(reflect.ValueOf(i))
			return 32, nil
		}

		return 0, errors.Wrap(ErrValue, "expect int/uint types ptr")
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

func (enc *boolEncoder) Unmarshal(data []byte, v interface{}) (uint, error) {
	var i uint

	if _, err := enc.Encoder.Unmarshal(data, &i); err != nil {
		return 0, err
	}

	t := reflect.ValueOf(v)

	if t.Kind() != reflect.Ptr || t.IsNil() {
		return 0, errors.Wrap(ErrValue, "expect bool ptr")
	}

	if t.Elem().Kind() != reflect.Bool {
		return 0, errors.Wrap(ErrValue, "expect bool ptr")
	}

	if i == 0 {
		t.Elem().Set(reflect.ValueOf(false))
		return 32, nil
	} else if i == 1 {
		t.Elem().Set(reflect.ValueOf(true))
		return 32, nil
	} else {
		return 0, errors.Wrap(ErrValue, "unmarshal abi data error")
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

func (enc *fixedEncoder) Unmarshal(data []byte, v interface{}) (uint, error) {
	var i *big.Int

	len, err := enc.Encoder.Unmarshal(data, &i)

	if err != nil {
		return len, err
	}

	t := reflect.ValueOf(v)

	if t.Kind() != reflect.Ptr || t.IsNil() {
		return 0, errors.Wrap(ErrValue, "expect bool ptr")
	}

	number := &fixed.Number{
		RawValue: i,
		Decimals: int(enc.N),
	}

	switch t.Elem().Kind() {
	case reflect.Float32:
		f, _ := number.Float().Float32()
		t.Elem().Set(reflect.ValueOf(f))
		return len, nil
	case reflect.Float64:
		f, _ := number.Float().Float64()
		t.Elem().Set(reflect.ValueOf(f))
		return len, nil
	default:
		if reflect.PtrTo(fixedType) == t.Elem().Type() {
			t.Elem().Set(reflect.ValueOf(number))
			return len, nil
		}

		return 0, errors.Wrap(ErrValue, "expect fixed ptr")
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

func (enc *fixedBytesEncoder) String() string {
	return fmt.Sprintf("[%d]byte", enc.len)
}

func (enc *fixedBytesEncoder) Static() bool {
	return true
}

func (enc *fixedBytesEncoder) Marshal(value interface{}) ([]byte, error) {
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Array || reflect.TypeOf(value).Elem().Kind() != reflect.Uint8 {
		return nil, errors.Wrap(ErrValue, "expect %s, got %v", enc, v.Type())
	}

	vv := reflect.New(v.Type())

	vv.Elem().Set(v)

	return paddingRight(vv.Elem().Slice(0, vv.Elem().Len()).Bytes()), nil
}

func (enc *fixedBytesEncoder) Unmarshal(data []byte, v interface{}) (uint, error) {

	pl := paddingLen(enc.len)

	if uint(len(data)) < pl {
		return 0, errors.Wrap(ErrValue, "abi data len error")
	}

	t := reflect.ValueOf(v)

	if t.Kind() != reflect.Ptr || t.IsNil() {
		return 0, errors.Wrap(ErrValue, "expect %s ptr,got %v", enc, t.Type().Elem())
	}

	if t.Type().Elem().Kind() != reflect.Array {
		return 0, errors.Wrap(ErrValue, "expect %s ptr,got %v", enc, t.Type().Elem())
	}

	if t.Type().Elem().Elem().Kind() != reflect.Uint8 {
		return 0, errors.Wrap(ErrValue, "expect %s ptr,got %v", enc, t.Type().Elem())
	}

	d := data[:enc.len]

	t.Elem().Set(reflect.ValueOf(d).Convert(t.Type()).Elem())

	return pl, nil
}

type fixedArrayEncoder struct {
	len  uint
	elem Encoder
}

func FixedArray(elem Encoder, len uint) (Encoder, error) {
	return &fixedArrayEncoder{
		len:  len,
		elem: elem,
	}, nil
}

func (enc *fixedArrayEncoder) Static() bool {
	return true
}

func (enc *fixedArrayEncoder) String() string {
	return fmt.Sprintf("[%d]%s", enc.len, enc.elem)
}

func (enc *fixedArrayEncoder) Marshal(value interface{}) ([]byte, error) {
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Array || enc.len != uint(v.Len()) {
		return nil, errors.Wrap(ErrValue, "expect %s, got %v", enc, v.Type())
	}

	if enc.elem.Static() {
		return enc.staticMarshal(v)
	}

	return enc.dynamicMarshal(v)
}

func (enc *fixedArrayEncoder) staticMarshal(value reflect.Value) ([]byte, error) {
	var buff bytes.Buffer

	for i := 0; i < value.Len(); i++ {
		content, err := enc.elem.Marshal(value.Index(i).Interface())

		if err != nil {
			return nil, err
		}

		buff.Write(content)
	}

	return buff.Bytes(), nil
}

func (enc *fixedArrayEncoder) dynamicMarshal(value reflect.Value) ([]byte, error) {
	var header bytes.Buffer
	var content bytes.Buffer

	headerLen := 32 * value.Len()

	ie, err := Integer(false, 256)

	if err != nil {
		return nil, err
	}

	for i := 0; i < value.Len(); i++ {
		headerBuff, err := ie.Marshal(uint(headerLen + content.Len()))

		if err != nil {
			return nil, err
		}

		_, err = header.Write(headerBuff)

		if err != nil {
			return nil, err
		}

		contentBuff, err := enc.elem.Marshal(value.Index(i).Interface())

		if err != nil {
			return nil, err
		}

		_, err = content.Write(contentBuff)

		if err != nil {
			return nil, err
		}
	}

	_, err = header.Write(content.Bytes())

	return header.Bytes(), err
}

func (enc *fixedArrayEncoder) Unmarshal(data []byte, v interface{}) (uint, error) {

	t := reflect.ValueOf(v)

	if t.Kind() != reflect.Ptr || t.IsNil() {
		return 0, errors.Wrap(ErrValue, "expect %s ptr", enc)
	}

	if t.Elem().Kind() != reflect.Array {
		return 0, errors.Wrap(ErrValue, "expect %s,got %v", enc, t.Type())
	}

	if t.Elem().Len() != int(enc.len) {
		return 0, errors.Wrap(ErrValue, "expect %s,got %v", enc, reflect.Indirect(t).Type())
	}

	if enc.elem.Static() {
		return enc.staticUnmarshal(data, t)
	} else {
		return enc.dynamicUnmarshal(data, t)
	}

}

func (enc *fixedArrayEncoder) staticUnmarshal(data []byte, v reflect.Value) (uint, error) {

	// create slice
	array := reflect.New(v.Type().Elem())

	offset := uint(0)

	for i := 0; i < int(enc.len); i++ {
		if len(data) < int(offset) {
			return 0, errors.Wrap(ErrLength, "abi data too short")
		}

		content := reflect.New(array.Type().Elem().Elem())

		len, err := enc.elem.Unmarshal(data[offset:], content.Interface())

		if err != nil {
			return 0, err
		}

		array.Elem().Index(i).Set(content.Elem())

		offset += len
	}

	reflect.Indirect(v).Set(array.Elem())

	return offset, nil
}

func (enc *fixedArrayEncoder) dynamicUnmarshal(data []byte, v reflect.Value) (uint, error) {
	headerLen := enc.len * 32

	if len(data) < int(headerLen) {
		return 0, errors.Wrap(ErrLength, "input unmarshal data length error")
	}

	ienc, err := Integer(false, 256)

	if err != nil {
		return 0, err
	}

	offset := 0

	// create array ptr
	array := reflect.New(v.Type().Elem())

	contentLen := uint(0)

	for i := 0; i < int(enc.len); i++ {

		if len(data) < int(offset) {
			return 0, errors.Wrap(ErrLength, "abi data too short")
		}

		var elemOffset uint
		_, err := ienc.Unmarshal(data[offset:], &elemOffset)

		if err != nil {
			return 0, err
		}

		if len(data) < int(elemOffset) {
			return 0, errors.Wrap(ErrLength, "abi data too short")
		}

		content := reflect.New(v.Type().Elem().Elem())

		len, err := enc.elem.Unmarshal(data[elemOffset:], content.Interface())

		if err != nil {
			return 0, err
		}

		offset += 32

		contentLen = elemOffset + len

		array.Elem().Index(i).Set(content)
	}

	reflect.Indirect(v).Set(array.Elem())

	return headerLen + contentLen, nil
}

type bytesEncoder struct {
}

func Bytes() (Encoder, error) {
	return &bytesEncoder{}, nil
}

func (enc *bytesEncoder) Static() bool {
	return false
}

func (enc *bytesEncoder) String() string {
	return "[]byte"
}

func (enc *bytesEncoder) Marshal(value interface{}) ([]byte, error) {
	b, ok := value.([]byte)

	if !ok {
		return nil, errors.Wrap(ErrValue, "input value must be []byte")
	}

	encoder, err := FixedBytes(uint(len(b)))

	if err != nil {
		return nil, err
	}

	arrayT := reflect.PtrTo(reflect.ArrayOf(len(b), reflect.TypeOf(b).Elem()))

	content, err := encoder.Marshal(reflect.ValueOf(b).Convert(arrayT).Elem().Interface())

	if err != nil {
		return nil, err
	}

	ienc, err := Integer(false, 256)

	if err != nil {
		return nil, err
	}

	header, err := ienc.Marshal(len(b))

	if err != nil {
		return nil, err
	}

	return append(header, content...), nil
}

func (enc *bytesEncoder) Unmarshal(data []byte, v interface{}) (uint, error) {

	if len(data) < 32 {
		return 0, errors.Wrap(ErrLength, "abi data length < 32")
	}

	ienc, err := Integer(false, 256)

	if err != nil {
		return 0, err
	}

	var k uint

	_, err = ienc.Unmarshal(data, &k)

	if err != nil {
		return 0, err
	}

	if uint(len(data)) < k+32 {
		return 0, errors.Wrap(ErrLength, "abi data length < 32")
	}

	encoder, err := FixedBytes(k)

	if err != nil {
		return 0, err
	}

	vv := reflect.New(reflect.ArrayOf(int(k), reflect.TypeOf((*byte)(nil)).Elem()))

	contentLen, err := encoder.Unmarshal(data[32:], vv.Interface())

	if err != nil {
		return 0, err
	}

	reflect.ValueOf(v).Elem().Set(vv.Elem().Slice(0, vv.Elem().Len()))

	return contentLen + 32, nil
}

type stringEncoder struct {
}

func String() (Encoder, error) {
	return &stringEncoder{}, nil
}

func (enc *stringEncoder) Static() bool {
	return false
}

func (enc *stringEncoder) String() string {
	return "string"
}

func (enc *stringEncoder) Marshal(value interface{}) ([]byte, error) {
	s, ok := value.(string)

	if !ok {
		return nil, errors.Wrap(ErrValue, "input value must be []byte")
	}

	encoder, err := Bytes()

	if err != nil {
		return nil, err
	}

	return encoder.Marshal([]byte(s))
}

func (enc *stringEncoder) Unmarshal(data []byte, v interface{}) (uint, error) {

	vv := reflect.ValueOf(v)

	if vv.Kind() != reflect.Ptr || vv.Elem().Type() != stringType {
		return 0, errors.Wrap(ErrValue, "expect *string")
	}

	encoder, err := Bytes()

	if err != nil {
		return 0, err
	}

	var c []byte

	len, err := encoder.Unmarshal(data, &c)

	if err != nil {
		return 0, err
	}

	vv.Elem().Set(reflect.ValueOf(string(c)))

	return len, nil
}

type arrayEncoder struct {
	elem Encoder
}

func Array(elem Encoder) (Encoder, error) {
	return &arrayEncoder{
		elem: elem,
	}, nil
}

func (enc *arrayEncoder) Static() bool {
	return false
}

func (enc *arrayEncoder) String() string {
	return fmt.Sprintf("[]%s", enc.elem)
}

func (enc *arrayEncoder) Marshal(value interface{}) ([]byte, error) {
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Slice || v.IsNil() {
		return nil, errors.Wrap(ErrValue, "expect input slice not nil")
	}

	l := v.Len()

	ienc, err := Integer(false, 256)

	if err != nil {
		return nil, err
	}

	header, err := ienc.Marshal(uint(l))

	if err != nil {
		return nil, err
	}

	encoder, err := FixedArray(enc.elem, uint(l))

	if err != nil {
		return nil, err
	}

	arrayT := reflect.PtrTo(reflect.ArrayOf(l, v.Type().Elem()))

	content, err := encoder.Marshal(reflect.Indirect(v.Convert(arrayT)).Interface())

	if err != nil {
		return nil, err
	}

	return append(header, content...), nil
}

func (enc *arrayEncoder) Unmarshal(data []byte, value interface{}) (uint, error) {
	if len(data) < 32 {
		return 0, errors.Wrap(ErrLength, "abi data length < 32")
	}

	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Ptr || v.IsNil() || v.Type().Elem().Kind() != reflect.Slice {
		return 0, errors.Wrap(ErrValue, "expect input slice ptr got %v", v.Type())
	}

	ienc, err := Integer(false, 256)

	if err != nil {
		return 0, err
	}

	var l uint

	_, err = ienc.Unmarshal(data, &l)

	if err != nil {
		return 0, err
	}

	if uint(len(data)) < l+32 {
		return 0, errors.Wrap(ErrLength, "abi data length < %d", l+32)
	}

	encoder, err := FixedArray(enc.elem, l)

	if err != nil {
		return 0, err
	}

	arrayValue := reflect.New(reflect.ArrayOf(int(l), v.Type().Elem().Elem()))

	l, err = encoder.Unmarshal(data[32:], arrayValue.Interface())

	if err != nil {
		return 0, err
	}

	v.Elem().Set(arrayValue.Elem().Slice(0, arrayValue.Elem().Len()))

	return l + 32, nil
}

type tupleEncoder struct {
	elems []Encoder
}

func Tuple(elems ...Encoder) (Encoder, error) {

	return &tupleEncoder{
		elems: elems,
	}, nil
}

func (enc *tupleEncoder) Static() bool {
	return false
}

func (enc *tupleEncoder) String() string {
	var elems []string

	for _, el := range enc.elems {
		elems = append(elems, el.String())
	}

	return fmt.Sprintf("(%s)", strings.Join(elems, ","))
}

func (enc *tupleEncoder) Marshal(value interface{}) ([]byte, error) {

	slice, ok := value.([]interface{})

	if !ok {
		return nil, errors.Wrap(ErrValue, "Tuple: expect marshal value []interface{}")
	}

	if len(slice) != len(enc.elems) {
		return nil, errors.Wrap(ErrValue, "Tuple: marshal value []interface{} len must be %d", len(enc.elems))
	}

	iEncoder, err := Integer(false, 256)

	var header bytes.Buffer
	var body bytes.Buffer

	if err != nil {
		return nil, err

	}

	cached := make(map[int][]byte)
	headerLen := 0

	for i, elem := range enc.elems {
		if elem.Static() {

			buff, err := elem.Marshal(slice[i])

			if err != nil {
				return nil, err
			}

			cached[i] = buff

			headerLen += len(buff)

		} else {
			headerLen += 32
		}
	}

	for i, elem := range enc.elems {
		if elem.Static() {

			_, err = header.Write(cached[i])

			if err != nil {
				return nil, err
			}

		} else {

			buff, err := iEncoder.Marshal(body.Len() + headerLen)

			if err != nil {
				return nil, err
			}

			_, err = header.Write(buff)

			if err != nil {
				return nil, err
			}

			buff, err = elem.Marshal(slice[i])

			if err != nil {
				return nil, err
			}

			_, err = body.Write(buff)

			if err != nil {
				return nil, err
			}
		}
	}

	_, err = header.Write(body.Bytes())

	if err != nil {
		return nil, err
	}

	return header.Bytes(), nil
}

func (enc *tupleEncoder) Unmarshal(data []byte, value interface{}) (uint, error) {
	slice, ok := value.([]interface{})

	if !ok {
		return 0, errors.Wrap(ErrValue, "Tuple: expect umarshal value []interface{}")
	}

	if len(slice) != len(enc.elems) {
		return 0, errors.Wrap(ErrValue, "Tuple: marshal value []interface{} len must be %d", len(enc.elems))
	}

	offset := uint(0)
	maxLen := uint(len(data))
	contentLen := uint(0)

	iEncoder, err := Integer(false, 256)

	if err != nil {
		return 0, err
	}

	for i, elem := range enc.elems {
		if offset > maxLen {
			return 0, errors.Wrap(ErrLength, "Tuple: offset out of range")
		}

		if elem.Static() {
			len, err := elem.Unmarshal(data[offset:], slice[i])

			if err != nil {
				return 0, err
			}

			offset += len
		} else {
			var contentOffset uint

			len, err := iEncoder.Unmarshal(data[offset:], &contentOffset)

			if err != nil {
				return 0, err
			}

			offset += len

			if contentOffset > maxLen {
				return 0, errors.Wrap(ErrLength, "Tuple: content (%d,%s) offset out of range", i, elem)
			}

			len, err = elem.Unmarshal(data[contentOffset:], slice[i])

			if err != nil {
				return 0, err
			}

			contentLen += len
		}
	}

	return offset + contentLen, nil
}
