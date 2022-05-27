package mem_storage

import (
	"bytes"
	"unsafe"

	"github.com/OneOfOne/xxhash"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func EpgSelector(ctx *fasthttp.RequestCtx, priority_prov []uint32) {
  len_priority_prov := len(priority_prov)
  in_priority_prov := func(v uint32) bool {
    for i := 0; i < len_priority_prov; i++ {
      if priority_prov[i] == v { return true }
    }
    return false
  }
  _ = in_priority_prov
  

}


// Проверяет наличие элемента в Uint32 массиве
func SliceExistUint32(slice []uint32, val uint32) bool {
  for i := 0; i < len(slice); i++ {
    if slice[i] == val { return true }
  }
  return false
}


func b2s(b []byte) string {
  /* #nosec G103 */
  return *(*string)(unsafe.Pointer(&b))
}


// Пересортирует список так, что пользовательские записи будут первыми
func Slice_Prioritize[T *ProviderEpgData|*ProviderIcoData](user_prio []T, global_list []T, global_len int) []T {
  len_prio := len(user_prio)
  if len_prio == 0 { return global_list }
  cap_prio := cap(user_prio)

  // Проходимся по списку всех провайдеров и добавляем новые в конец user_prio
  // исключая уже добавленные (вначале списка, удобно отделены по len_prio)
  // должен получиться список того же размера, что и список всех провайдеров
  var _exist bool
  for i := 0; i < global_len; i++ {
    // Начинаем 
    _exist = false
    for g := 0; g < len_prio; g++ { // Проверяем уже добавленные
      if user_prio[g] == global_list[i] { _exist = true; break }
    }
    if !_exist {
      len_prio++
      if len_prio > cap_prio {
        // Надеюсь никогда не увидеть это сообщение
        newlen := cap_prio + global_len - i
        user_prio = append(make([]T, 0, newlen), user_prio...)
        log.Error().Msgf("Order.Prio: user_prio oversize! resize to %d", newlen)
      }
      user_prio = append(user_prio, global_list[i])
    }
  }
  return user_prio
}


// Дженерик добавления только новых записей существующему хранилищу 
func Slice_AppendUniq[T *ProviderEpgData|*ProviderIcoData](uniq_store []T, new_el T) []T {
  for i := 0; i < len(uniq_store); i++ {
    if (uniq_store[i] == new_el) { return uniq_store }
  }
  return append(uniq_store, new_el)
}


// Дженерик сканирует блок с ссылками на XMLTV и переводит их в
func PrioritizeUserProviders[T *ProviderEpgData|*ProviderIcoData](in_h []byte, base_list []T, look_func func (p uint32)T) []T {
  var hash_epg uint32

  if len(in_h) == 0 { return base_list }
  pdata := bytes.Split(in_h, c_sep_line)  // "\n"

  // Создаем с запасом по емкости (для последующей StoreSlice_PrioUser)
  base_list_len := len(base_list)
  in_list := make([]T, 0, base_list_len) 
  for i := 0; i < len(pdata); i++ {
    if len(pdata) == 0 { continue }
    url_lower := bytes.ToLower(pdata[i])
    hash_epg = xxhash.Checksum32(url_lower)

    // Проверяем, знаем ли такого провайдера
    _prov := look_func(hash_epg)
    if _prov != nil {
      // Знаем: Добавляем в priority список (с проверкой на уникальность)
      in_list = Slice_AppendUniq(in_list, _prov)
    } else {
      // НЕ Знаем: Логируем его хеш и ссылку (частых можно будет добавлять)
      log.Warn().Msgf("head-epg: unlisted provider - %d: %s", hash_epg, pdata[i])
    }
  }

  if len(in_list) == 0 {
  // Ничего не нашли, ищем по стандартному порядку
  return base_list
  }
  return Slice_Prioritize(in_list, base_list, base_list_len)
}
