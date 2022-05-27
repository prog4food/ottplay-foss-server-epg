package helpers

import (
	"unsafe"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func B2s(b []byte) string {
  /* #nosec G103 */
  return *(*string)(unsafe.Pointer(&b))
}

func ParseUint32(in []byte, err_note string) uint32 {
  _t, err := fasthttp.ParseUint(in)
  if err != nil {
    log.Err(err).Msg(err_note + ": cannot parse var %s" + B2s(in))
    return 0
  }
  return uint32(_t)
}

func AppendUint(dst []byte, n uint32) []byte {
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