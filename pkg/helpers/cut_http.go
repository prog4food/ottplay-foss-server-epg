package helpers

import "unsafe"


func HasHTTP(u string) bool {
  if (len(u) > 9) {
    if ( u[0:8] == "https://" ) || ( u[0:7] == "http://" ) {
      return true
    }
  }
  return false
}

func CutHTTP(u string) string {
  if (len(u) > 9) {
    if u[0:8] == "https://" { u = u[8:]
    } else if u[0:7] == "http://"  { u = u[7:] }
  }
  return u
}


/* Удаляет .gz из конца строки */
// inline cost = 26
func CutURLb_gz(t []byte,l int) []byte {
  const _gz = 0x7A67  // строка `gz`
  // Обрезаем .gz на конце
  if t[l-3] == '.' && *(*uint16)(unsafe.Pointer(&t[l-2])) == _gz {
    t = t[:l-3]
  }
  return t
}

/* Удаляет из byte строки протоколы http:// и https:// */
// inline cost = 39
func CutHTTPb(t []byte) []byte {
  const (
    _http       = 0x70747468  // строка `http`
    _protoid    = 0x2F2F3A73  // строка `s://`
    _protostart = 0x2F3A      // строка `:/`
  )
  // Обрезаем https:// и http://
  if t[6] == '/' && *(*uint32)(unsafe.Pointer(&t[0])) == _http {  // Ищем по маске `http??/`
    if *(*uint32)(unsafe.Pointer(&t[4])) == _protoid {            // проверка части после `http` `s://`
      t = t[8:]
    } else if *(*uint16)(unsafe.Pointer(&t[4])) == _protostart {  // проверка части после `http` `:/`
      t = t[7:]
    }
  }
  return t
}