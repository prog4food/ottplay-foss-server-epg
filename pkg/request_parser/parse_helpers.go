package request_parser

import (
	"bytes"
	"io"
	"net/http"

	"github.com/OneOfOne/xxhash"
	"github.com/rs/zerolog/log"

	"ottplay-foss-server-epg/pkg/helpers"
	"ottplay-foss-server-epg/pkg/mem_storage"
)

// Разделяет и валидирует запрос на 3 части
func SplitMatchBody(r *http.Request) [][]byte {
  // Валидация запроса
  if r.Method != "POST" {
    return nil
  }
  // Чтение всего тела
  body, err := io.ReadAll(r.Body)
  r.Body.Close()
  if err != nil && err != io.EOF {
    log.Error().Err(err).Msgf("/match-*/ - cannot read body!")
    return nil
  }
  out := bytes.SplitN(body, c_sep_block, 4)
  if len(out) != 3 {
    log.Error().Msgf("/match-*/ - cannot bad body blocks count!")
    return nil
  }
  return out
}

// Проверяет наличие элемента в Uint32 массиве
func SliceExistUint32(slice []uint32, val uint32) bool {
  for i := 0; i < len(slice); i++ {
    if slice[i] == val { return true }
  }
  return false
}

// Сортирует listGlobal так, что listUser будут первыми
func SortByUser[T *mem_storage.ProviderEpgData|*mem_storage.ProviderIcoData](listUser []T, listGlobal []T, lenGlobal int) []T {
  lenUser := len(listUser)
  if lenUser == 0 { return listGlobal }
  capUser := cap(listUser)

  // Проходимся по списку всех провайдеров и добавляем новые в конец listUser
  // исключая уже добавленные (вначале списка, удобно отделены по lenUser)
  // должен получиться список того же размера, что и список всех провайдеров
  var _exist bool
  for i := 0; i < lenGlobal; i++ {
    // Начинаем 
    _exist = false
    for g := 0; g < lenUser; g++ { // Проверяем уже добавленные
      if listUser[g] == listGlobal[i] { _exist = true; break }
    }
    if !_exist {
      lenUser++
      if lenUser > capUser {
        // Надеюсь никогда не увидеть это сообщение
        newLenUser := capUser + lenGlobal - i
        listUser = append(make([]T, 0, newLenUser), listUser...)
        log.Error().Msgf("SortByUser: listUser oversize! resize to %d", newLenUser)
      }
      listUser = append(listUser, listGlobal[i])
    }
  }
  return listUser
}


// Дженерик добавления только новых записей существующему хранилищу (проверка по Pointer)
func Slice_AppendUniq[T *mem_storage.ProviderEpgData|*mem_storage.ProviderIcoData](uniqStore []T, newElement T) []T {
  for i := 0; i < len(uniqStore); i++ {
    if (uniqStore[i] == newElement) { return uniqStore }
  }
  return append(uniqStore, newElement)
}


// Дженерик сканирует блок с ссылками на XMLTV и переводит их в
func PrioritizeUserProviders[T *mem_storage.ProviderEpgData|*mem_storage.ProviderIcoData](inData []byte, listGlobal []T, look_func func (uint32, []T)T) ([]T, uint8) {
  if len(inData) == 0 { return listGlobal, 0 }

  var hashUserXmltv uint32
  // Создаем с запасом по емкости (для последующей StoreSlice_PrioUser)
  var lenGlobal = len(listGlobal)
  var listUser  = make([]T, 0, lenGlobal)
  var provUser T
  var reqTvgUrls = bytes.Split(inData, c_sep_line)  // "\n"
  var t []byte  // Временная строка для нормализации url-tvg
  var l int     // Временная переменная для длины строки t
  for i := 0; i < len(reqTvgUrls); i++ {
    // Нормализация url-tvg
    t = reqTvgUrls[i]
    if l = len(t); l > 10 {   // http://a.co
      t = helpers.CutURLb_gz(t, l) // первое, тк необходима l
      t = helpers.CutHTTPb(t)
    }

    hashUserXmltv = xxhash.Checksum32(t)
    if hashUserXmltv == helpers.EmptyXXHash32 { continue }

    // Проверяем, знаем ли такого провайдера
    provUser = look_func(hashUserXmltv, listGlobal)
    if provUser != nil {
      // Знаем: Добавляем в priority список (с проверкой на уникальность)
      listUser = Slice_AppendUniq(listUser, provUser)
    } else {
      // НЕ Знаем: Логируем его хеш и ссылку (частых можно будет добавлять)
      log.Warn().Msgf("match.epg: unlisted provider - %d: %s", hashUserXmltv, reqTvgUrls[i])
    }
  }

  lenUser := uint8(len(listUser))
  if lenUser == 0 {
    // Ничего не нашли, ищем по стандартному порядку
    return listGlobal, 0
  }
  return SortByUser(listUser, listGlobal, lenGlobal), lenUser
}
