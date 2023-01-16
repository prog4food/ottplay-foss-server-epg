package helpers

import (
	"errors"
	"reflect"
	"unsafe"
)

const EmptyXXHash32 = 46947589

func B2s(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

func S2b(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

func S2bP(s *string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

// Based on fasthttp.ParseUint (original has bug: read only 31 bit)
var (
	ErrEmptyInt               = errors.New("empty integer")
	ErrUnexpectedFirstChar    = errors.New("unexpected first char found. Expecting 0-9")
	ErrUnexpectedTrailingChar = errors.New("unexpected trailing char found. Expecting 0-9")
	ErrTooLongInt             = errors.New("too long int")
)

func ParseUint32Buf(b []byte) (uint32, int, error) {
	n := len(b)
	if n == 0 {
		return 0, 0, ErrEmptyInt
	}
	var (
		v, vNew uint32
		k       byte
	)
	for i := 0; i < n; i++ {
		k = b[i] - '0'
		if k > 9 {
			if i == 0 {
				return 0, i, ErrUnexpectedFirstChar
			} else {
				return v, i, ErrUnexpectedTrailingChar
			}
		}
		vNew = 10*v + uint32(k)
		// Test for overflow.
		if vNew < v {
			return 0, i, ErrTooLongInt
		}
		v = vNew
	}
	return v, n, nil
}

// Based on Uint32ToBytes from fasthttp
func Uint32ToBytes(dst []byte, n uint32) []byte {
	var b [20]byte
	buf := b[:]
	i := len(buf)
	var q uint32
	for n >= 10 {
		i--
		q = n / 10
		buf[i] = '0' + byte(n-q*10)
		n = q
	}
	i--
	buf[i] = '0' + byte(n)

	dst = append(dst, buf[i:]...)
	return dst
}

// Получает адрес блока данных объекта []byte
func getSliceDataPtr(s *string) uintptr {
	return (*reflect.SliceHeader)(unsafe.Pointer(s)).Data
}

// Получает адрес блока данных объекта String
func getStringDataPtr(s *string) uintptr {
	return (*reflect.StringHeader)(unsafe.Pointer(s)).Data
}

// Проверяет, что структуры ссылаются на одну и туже область памяти
// полезно при сравнении map
func HasSameDataPointer(a, b any) bool {
	return (reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer())
}
