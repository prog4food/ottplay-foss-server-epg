package providers_downloader

import (
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"

	xxhash32 "github.com/OneOfOne/xxhash"

	"ottplay-foss-server-epg/libs/config_epg"
	"ottplay-foss-server-epg/libs/helpers"
	ms "ottplay-foss-server-epg/libs/mem_storage"
)

var HttpClient *fasthttp.Client

func processProvider(ps *config_epg.ProviderSource, b []byte) {
  _body_hash := xxhash32.Checksum32(b)
  if ps.BodySum == _body_hash {
    return    
  }
  _json_prov := ProviderData{}
  err := _json_prov.UnmarshalJSON(b)
  if err != nil {
    log.Err(err).Msgf("%s: cannot read json", ps.IdName)
    return
  }

  // Новая запись провайдера
  p := &ms.ProviderElement{
    Id: &_json_prov.Meta.Id,
    IdHashes: append(_json_prov.Meta.Urls, xxhash32.ChecksumString32(_json_prov.Meta.Id)),
  }
  // Пустое хранилище: pEpg
  pEpg := &ms.ProviderEpgData{
    ById:     make(ms.HashEpgStruct, len(_json_prov.Data)),
    ByName:   make(ms.HashEpgStruct, len(_json_prov.Data)),
    UrlBase:  &ps.UrlBase,
    UrlBaseS: &ps.UrlBaseS,
    Order:    ps.OrderEpg,
    Prov:     p,
  }
  // Пустое хранилище: pIco
  pIco := &ms.ProviderIcoData{
    ById:     make(ms.HashIcoStruct, len(_json_prov.Data)),
    ByName:   make(ms.HashIcoStruct, len(_json_prov.Data)),
    Order:    ps.OrderIco,
    Prov:     p,
  }
  p.Epg = pEpg
  p.Ico = pIco

  // Обработка каналов
  var (
    last_epg    int64              // Временная переменная: последняя передача
    epg_data   *ms.EpgChannelData  // Временная переменная: структура Epg
    ico_data   *string             // Временная переменная: dedup ico
    has_epg     bool               // Временная переменная: канал имеет EPG
    has_picon   bool )             // Временная переменная: канал имеет ICO
  for _id_hash, _payload := range _json_prov.Data {
    _ch_meta := strings.Split(_payload[0], "¦")
    // Чтение структуры [0] (метаданные)
    if len(_ch_meta) != 3 {
      log.Error().Msgf("%s: Cannot parse [0]: %s", ps.IdName, _payload[0])
      continue
    }

    last_epg, err = strconv.ParseInt(_ch_meta[1], 10, 64); if err != nil {
      log.Error().Msgf("%s: Cannot parse LastEpg: %s", ps.IdName, _payload[0])
      continue
    }
    has_epg   = (last_epg > 0)
    has_picon = (_ch_meta[2] != "")
    // Заполнение pEpg/pIco ( ById )
    if has_epg {
      epg_data = &ms.EpgChannelData{
        IdHash:   ms.DedupByte(ms.Ids, helpers.AppendUint(nil, _id_hash)),
        // // TODO: пока нет смысла сохранять имена и Id в виде строк
        // Id:       ms.DedupStr(ms.Str, &_ch_meta[0]),
        // Names:    make([]*string, 0, len(_payload)-1),
        LastEpg:  last_epg,
        ExpiredEpg: (last_epg < _now_unix),
      }
      if epg_data.ExpiredEpg { pEpg.Outdated++ }
      pEpg.ById[_id_hash] = epg_data
    }
    if has_picon {
      ico_data = ms.DedupStr(ms.Url, &_ch_meta[2])
      pIco.ById[_id_hash] = ico_data
    }
    // Заполнение pEpg/pIco ( ByName )
    for g := 1; g < len(_payload); g++ {
      _name_hash := xxhash32.ChecksumString32(strings.ToLower(_payload[g]))
      if has_epg {
        pEpg.ByName[_name_hash] = epg_data
        // TODO: пока нет смысла сохранять имена
        //epg_data.Names = append(epg_data.Names, ms.DedupStr(ms.Str, &_payload[g])) 
      }
      if has_picon {
        pIco.ByName[_name_hash] = ico_data
      }
    }
  }
  // "Финалим" объект в общем списке
  ms.Lock.Lock()
    ms.P[ps.IdHash] = p
    ps.BodySum = _body_hash
    // Если это не запуск при инициализации, то делаем пересортировку (тк объект в хранилище поменял адрес)
    if !sched_first_run { ms.PO.Sort() }
  ms.Lock.Unlock()
  log.Info().Msgf("%s updated, names=%d, epg=%d outdated=%d, icon=%d", *p.Id, len(pEpg.ByName), len(pEpg.ById), pEpg.Outdated, len(pIco.ById))
}

func DownloadProvider(p *config_epg.ProviderSource) {
  var _url string
  if p.UrlBase != "" { _url = p.UrlBase
  } else if p.UrlBaseS != "" { _url = p.UrlBaseS
  } else {
    log.Error().Msgf("cannot find provider url [%s]", p.IdName)
    return
  }
  req  := fasthttp.AcquireRequest()
  req.SetRequestURI(_url + "channels.json")
  req.Header.SetMethod(fasthttp.MethodGet)
  resp := fasthttp.AcquireResponse()
  err  := HttpClient.DoRedirects(req, resp, 3)
  fasthttp.ReleaseRequest(req)
  if (err == nil) && (resp.StatusCode() == fasthttp.StatusOK) {
    processProvider(p, resp.Body())
  } else {
    log.Error().Err(err).Msgf("cannot download channels for '%s' [status: %d, url: %s]", p.IdName, resp.StatusCode(), _url)
  }
  fasthttp.ReleaseResponse(resp)
}