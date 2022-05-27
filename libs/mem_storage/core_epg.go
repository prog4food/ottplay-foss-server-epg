package mem_storage

import (
	"bytes"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

var (
	c_sep_block = []byte("\n\t\n")    // \n\t\n
  c_sep_line  = []byte{0x0A} // \n
  c_sep_array = []byte{'~'}  // ~
  c_sep_vars  = []byte{'-'}  // -

  json_quote = []byte{'"'}        // "
  json_end_line = []byte("\",\n") // ",\n
  json_end_file = []byte("]}\n")  // ]}\n
  ch_head_epg = []byte("{\"epg\":[\n")        // {"epg":[\n
  ch_head_prov = []byte("\n],\"provs\": [\n") // \n],"provs": [\n
)


// EPG: Ищет провайдера по URL
func Epg_GetProvByHash(p uint32) *ProviderEpgData {
  for _, v := range P {
    if SliceExistUint32(v.IdHashes, p) { return v.Epg }
  }
  return nil
}

// Определяет приоритет поиска по провайдерам, на соснове провайдеров
// указанных для канала, которые передаются в виде хешей
func ParseEpg_Channel[T *ProviderEpgData|*ProviderIcoData](in_h []byte, base_list []T, look_func func (p uint32)T) []T {
  var err error
  var _hash_epg int
  var hash_epg uint32

  if len(in_h) == 0 { return base_list }
  pdata := bytes.Split(in_h, c_sep_vars)  // "-"

  // Создаем с запасом по емкости (для последующей StoreSlice_PrioUser)
  base_list_len := len(base_list)
  in_list := make([]T, 0, base_list_len) 
  for i := 0; i < len(pdata); i++ {
    _hash_epg, err = fasthttp.ParseUint(pdata[i])
    if err != nil {
      log.Err(err).Msgf("match-channels.epg: cannot parse var %s", b2s(pdata[i]))
      continue
    }; hash_epg = uint32(_hash_epg)

    // Проверяем, знаем ли такого провайдера
    _prov := look_func(hash_epg)
    if _prov != nil {
      // Знаем: Добавляем в priority список (с проверкой на уникальность)
      in_list = Slice_AppendUniq(in_list, _prov)
    } else {
      // НЕ Знаем: Логируем его хеш и ссылку (частых можно будет добавлять)
      log.Warn().Msgf("channel-epg: unlisted provider - %d", hash_epg)
    }
  }

  if len(in_list) == 0 {
    // Ничего не нашли, ищем по стандартному порядку
    return base_list
  }
  return Slice_Prioritize(in_list, base_list, base_list_len)
}


func ParseEpg_ReqChannels(in []byte, prov_list []*ProviderEpgData, ctx *fasthttp.RequestCtx) {
  if len(in) == 0 { return }
  clist := bytes.Split(in, c_sep_line)  // "\n"


  var providers  []*ProviderEpgData
  var cdata      [][]byte
  var (         // Блок переменных для используемых провайдеров
    used_provs = make([]*ProviderEpgData, 0, len(prov_list))
    used_provs_len int )
  var (         // Блок переменных для результатов поиска канала
    ch_key     []byte
    ch_data   *EpgChannelData
    ch_prov   *ProviderEpgData )

  // Пишем json заголовок в файл
  out_buf  := &ctx.Response
  mini_buf := make([]byte, 0, 256)

  // Блок: МЕТА
  out_buf.AppendBodyString("{}")


  // Блок: КАНАЛЫ
  out_buf.AppendBody(c_sep_block)
  // Проход по каналам
  for i := 0; i < len(clist); i++ {
    if len(clist[i]) == 0 { continue }

    cdata = bytes.Split(clist[i], c_sep_array)  // "~" - делим данные канала и его epg
    if len(cdata) == 2 {
      // Приоритизируем пользовательские источники
      providers = ParseEpg_Channel(cdata[1], prov_list, Epg_GetProvByHash)
    } else {
      // Используем общий список
      providers = prov_list
    }

    ch_key, ch_data, ch_prov = ParseEpg_LookupChannel(cdata[0], providers)
    if ch_data != nil {
      // Пишем данные о канале
      mini_buf = mini_buf[:0] // Сброс буфера
      used_provs = Slice_AppendUniq(used_provs, ch_prov)
      if i > 0 { mini_buf = append(mini_buf, c_sep_line...) }
      mini_buf = append(mini_buf, ch_key...)
      mini_buf = append(mini_buf, c_sep_array...)
      mini_buf = append(mini_buf, *ch_prov.Prov.Id...)
      mini_buf = append(mini_buf, c_sep_array...)
      mini_buf = append(mini_buf, ch_data.IdHash...)
      out_buf.AppendBody(mini_buf)
    }
  }
  out_buf.AppendBody(c_sep_block)
  used_provs_len = len(used_provs)

  // Блок: ПРОВАЙДЕРЫ
  if used_provs_len > 0 {
    // Пишем данные об используемых провайдерах  
    for i := 0; i < used_provs_len; i++ {
      mini_buf = mini_buf[:0] // Сброс буфера
      if i > 0 { mini_buf = append(mini_buf, c_sep_line...) }
      mini_buf = append(mini_buf, *used_provs[i].Prov.Id...)
      mini_buf = append(mini_buf, c_sep_array...)
      mini_buf = append(mini_buf, *used_provs[i].UrlBase...)
      out_buf.AppendBody(mini_buf)
    }
  }
}

