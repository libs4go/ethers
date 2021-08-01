package abi

import (
	"bytes"
	"math/big"
	"reflect"
	"sort"
	"strconv"

	"github.com/libs4go/errors"
	"github.com/libs4go/fixed"
	"golang.org/x/crypto/sha3"
)

var bigIntType = reflect.TypeOf((*big.Int)(nil)).Elem()
var fixedType = reflect.TypeOf((*fixed.Number)(nil)).Elem()
var bytesPtrType = reflect.TypeOf((*[]byte)(nil)).Elem()
var stringType = reflect.TypeOf((*string)(nil)).Elem()

func paddingLen(origin uint) uint {
	return ((origin + 31) / 32) * 32
}

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

type tupleOrder struct {
	order  uint64
	offset int
}

type byOrder []*tupleOrder

func (a byOrder) Len() int           { return len(a) }
func (a byOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byOrder) Less(i, j int) bool { return a[i].order < a[j].order }

func tupleOrderFields(tuple reflect.Type) ([]int, error) {
	if tuple.Kind() != reflect.Struct {
		return nil, errors.Wrap(ErrValue, "expect struct type")
	}

	var fields []*tupleOrder

	for i := 0; i < tuple.NumField(); i++ {
		field := tuple.Field(i)

		order, ok := field.Tag.Lookup("tuple")

		if !ok {
			continue
		}

		o, err := strconv.ParseUint(order, 10, 64)

		if err != nil {
			return nil, errors.Wrap(ErrTag, "filed %s tuple tag must be number", field.Name)
		}

		fields = append(fields, &tupleOrder{
			offset: i,
			order:  o,
		})
	}

	sort.Sort(byOrder(fields))

	var result []int

	for _, field := range fields {
		result = append(result, field.offset)
	}

	return result, nil
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
	Unmarshal(data []byte, v interface{}) (uint, error)
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

		t.Elem().Set(reflect.ValueOf(i.Uint64()))

		return 32, nil

	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		if !enc.sign {
			return 0, errors.Wrap(ErrValue, "expect int/uint types ptr")
		}

		t.Elem().Set(reflect.ValueOf(i.Int64()))

		return 32, nil

	default:
		if bigIntType == t.Elem().Type() {
			t.Elem().Set(reflect.ValueOf(i).Elem())
			return 32, nil
		}

		if reflect.Ptr == t.Elem().Kind() && t.Elem().Elem().Type() == bigIntType {
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

func (enc *fixedBytesEncoder) Unmarshal(data []byte, v interface{}) (uint, error) {

	pl := paddingLen(enc.len)

	if uint(len(data)) < pl {
		return 0, errors.Wrap(ErrValue, "abi data len error")
	}

	t := reflect.ValueOf(v)

	if t.Kind() != reflect.Ptr || t.IsNil() {
		return 0, errors.Wrap(ErrValue, "expect []byte ptr")
	}

	if t.Elem().Type() != bytesPtrType {
		return 0, errors.Wrap(ErrValue, "expect []byte ptr")
	}

	t.Elem().Set(reflect.ValueOf(data[:enc.len]))

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

func (enc *fixedArrayEncoder) Marshal(value interface{}) ([]byte, error) {
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Slice && enc.len != uint(v.Len()) {
		return nil, errors.Wrap(ErrValue, "expect [%d]content", enc.len)
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
		headerBuff, err := ie.Marshal(uint(headerLen + i*32))

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
		return 0, errors.Wrap(ErrValue, "expect slice ptr")
	}

	if t.Elem().Kind() != reflect.Slice {
		return 0, errors.Wrap(ErrValue, "expect slice ptr")
	}

	if enc.elem.Static() {
		return enc.staticUnmarshal(data, t.Elem())
	} else {
		return enc.dynamicUnmarshal(data, t.Elem())
	}

}

func (enc *fixedArrayEncoder) staticUnmarshal(data []byte, v reflect.Value) (uint, error) {

	// create slice
	slice := reflect.New(v.Type())

	offset := uint(0)

	for i := 0; i < int(enc.len); i++ {
		if len(data) < int(offset) {
			return 0, errors.Wrap(ErrLength, "abi data too short")
		}

		content := reflect.New(v.Elem().Type())

		len, err := enc.elem.Unmarshal(data[offset:], content.Addr().Interface())

		if err != nil {
			return 0, err
		}

		slice = reflect.Append(slice, content)

		offset += len
	}

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

	// create slice
	slice := reflect.New(v.Type())

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

		content := reflect.New(v.Elem().Type())

		len, err := enc.elem.Unmarshal(data[elemOffset:], content.Addr().Interface())

		if err != nil {
			return 0, err
		}

		offset += 32

		contentLen = elemOffset + len

		slice = reflect.Append(slice, content)
	}

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

func (enc *bytesEncoder) Marshal(value interface{}) ([]byte, error) {
	b, ok := value.([]byte)

	if !ok {
		return nil, errors.Wrap(ErrValue, "input value must be []byte")
	}

	encoder, err := FixedBytes(uint(len(b)))

	if err != nil {
		return nil, err
	}

	content, err := encoder.Marshal(b)

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

	_, err = encoder.Unmarshal(data[32:], v)

	if err != nil {
		return 0, err
	}

	return k + 32, nil
}

type stringEncoder struct {
}

func String() (Encoder, error) {
	return &stringEncoder{}, nil
}

func (enc *stringEncoder) Static() bool {
	return false
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

	content, err := encoder.Marshal(value)

	if err != nil {
		return nil, err
	}

	return append(header, content...), nil
}

func (enc *arrayEncoder) Unmarshal(data []byte, v interface{}) (uint, error) {
	if len(data) < 32 {
		return 0, errors.Wrap(ErrLength, "abi data length < 32")
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

	l, err = encoder.Unmarshal(data[32:], v)

	if err != nil {
		return 0, err
	}

	return l + 32, nil
}

type tupleEncoder struct {
	elems []Encoder
}

func Tuple(encoders ...Encoder) (Encoder, error) {
	return &tupleEncoder{
		elems: encoders,
	}, nil
}

func (enc *tupleEncoder) Static() bool {
	return false
}

func (enc *tupleEncoder) structToSlice(value interface{}) ([]interface{}, error) {
	vv := reflect.ValueOf(value)

	if vv.Kind() != reflect.Ptr || vv.IsNil() || vv.Elem().Kind() != reflect.Struct {
		return nil, errors.Wrap(ErrValue, "expect struct ptr")
	}

	structT := vv.Elem().Type()

	fields, err := tupleOrderFields(structT)

	if err != nil {
		return nil, err
	}
	var result []interface{}
	for _, i := range fields {
		result = append(result, vv.Elem().Field(i).Interface())
	}

	return result, nil
}

func (enc *tupleEncoder) Marshal(value interface{}) ([]byte, error) {

	v, err := enc.structToSlice(value)

	if err != nil {
		return nil, err
	}

	if len(enc.elems) != len(v) {
		return nil, errors.Wrap(ErrLength, "expect slice len %d", len(enc.elems))
	}

	// calc header length
	cached := make(map[int][]byte)

	offset := 0

	for i := 0; i < len(enc.elems); i++ {
		if enc.elems[i].Static() {
			buff, err := enc.elems[i].Marshal(v[i])

			if err != nil {
				return nil, err
			}

			cached[i] = buff

			offset += len(buff)
		} else {
			offset += 32
		}
	}

	ienc, err := Integer(false, 256)

	if err != nil {
		return nil, err
	}

	var header bytes.Buffer
	var content bytes.Buffer

	for i := 0; i < len(enc.elems); i++ {

		if enc.elems[i].Static() {
			header.Write(cached[i])
		} else {
			buff, err := ienc.Marshal(offset)

			if err != nil {
				return nil, err
			}

			_, err = header.Write(buff)

			if err != nil {
				return nil, err
			}

			buff, err = enc.elems[i].Marshal(v[i])

			if err != nil {
				return nil, err
			}

			_, err = content.Write(buff)

			if err != nil {
				return nil, err
			}
		}
	}

	return append(header.Bytes(), content.Bytes()...), nil
}

func (enc *tupleEncoder) structPtrToFields(value interface{}) ([]interface{}, error) {
	vv := reflect.ValueOf(value)

	if vv.Kind() != reflect.Ptr || vv.IsNil() {
		return nil, errors.Wrap(ErrValue, "expect struct ptr ptr")
	}

	vv = vv.Elem()

	if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Struct {
		return nil, errors.Wrap(ErrValue, "expect struct ptr ptr")
	}

	orders, err := tupleOrderFields(vv.Elem().Type())

	if err != nil {
		return nil, err
	}

	obj := reflect.New(vv.Elem().Type())

	vv.Set(obj)

	var result []interface{}

	for _, i := range orders {
		result = append(result, obj.Field(i).Addr().Interface())
	}

	return result, nil
}

func (enc *tupleEncoder) Unmarshal(data []byte, value interface{}) (uint, error) {

	vv, err := enc.structPtrToFields(value)

	if err != nil {
		return 0, err
	}

	if len(vv) != len(enc.elems) {
		return 0, errors.Wrap(ErrLength, "unexpect input type")
	}

	ienc, err := Integer(false, 256)

	if err != nil {
		return 0, err
	}

	headerOffset := uint(0)

	dataLen := len(data)

	contentLen := uint(0)

	for i, elem := range enc.elems {

		if dataLen < int(headerOffset) {
			return 0, errors.Wrap(ErrLength, "abi data length too short")
		}

		if elem.Static() {
			l, err := elem.Unmarshal(data[headerOffset:], vv[i])

			if err != nil {
				return 0, err
			}

			headerOffset += l
		} else {
			var offset uint

			_, err := ienc.Unmarshal(data[headerOffset:], &offset)

			if err != nil {
				return 0, err
			}

			headerOffset += 32

			if dataLen < int(offset) {
				return 0, errors.Wrap(ErrLength, "abi data length too short")
			}

			l, err := elem.Unmarshal(data[offset:], vv[i])

			if err != nil {
				return 0, err
			}

			contentLen = offset + l
		}
	}

	return headerOffset + contentLen, nil
}
