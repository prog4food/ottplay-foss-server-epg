package mem_storage

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

// ICO: Ищет провайдера по URL
func Ico_GetProvByHash(p uint32) *ProviderIcoData {
  for _, v := range P {
    if SliceExistUint32(v.IdHashes, p) { return v.Ico }
  }
  return nil
}


func ParseIco_ReqChannels(in []byte, prov_list []*ProviderIcoData, ctx *fasthttp.RequestCtx, prov_user_len uint8) {
  if len(in) == 0 { return }
  clist := bytes.Split(in, c_sep_line)  // "\n"


  var providers  []*ProviderIcoData
  var cdata      [][]byte
  var (         // Блок переменных для результатов поиска канала
    ch_key     []byte
    ch_data   *string )

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
      providers, prov_user_len = ParseEpg_Channel(cdata[1], prov_list, Ico_GetProvByHash)
    } else {
      // Используем общий список
      providers = prov_list
    }

    ch_key, ch_data, _ = ParseIco_LookupChannel(cdata[0], providers, prov_user_len)
    if ch_data != nil {
      // Пишем данные о канале
      mini_buf = mini_buf[:0] // Сброс буфера
      if i > 0 { mini_buf = append(mini_buf, c_sep_line...) }
      mini_buf = append(mini_buf, ch_key...)
      mini_buf = append(mini_buf, c_sep_array...)
      mini_buf = append(mini_buf, *ch_data...)
      out_buf.AppendBody(mini_buf)
    }
  }
}
