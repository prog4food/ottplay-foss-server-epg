package internal

import (
	"net/http"

	xxhash32 "github.com/OneOfOne/xxhash"
	"github.com/rs/zerolog/log"

	"ottplay-foss-server-epg/pkg/config"
	"ottplay-foss-server-epg/pkg/mem_storage"
	"ottplay-foss-server-epg/pkg/request_parser"
)

func ReLoadConfig() {
  var isFirstRun = (P == nil)
  if isFirstRun {
    // >>> Первый запуск, нет нужды перечитывать конфиг
    // блокируем хранилище в начале, на все время работы
    ChDbLock.Lock()
    defer ChDbLock.Unlock()
  } else {
    // >>> Повторный запуск, надо перечитать конфиг
    // и заблокировать хранилище в конце
    _config := config.Load()
    if _config == nil {
      // Не даст упасть при hot-reload
      log.Error().Msg("config: Bad config! Сonfig NOT reloaded!")
      return
    }
    Config = _config
  }

  // Загрузка списков order_list
  LoadOrderLists()

  // Создание новой channel-структуры из конфига
  var provConfig *config.ProviderConfig
  var prov *mem_storage.ProviderElement
  var provIdHash uint32
  var countProviders = len(Config.Providers)
  var p = make(mem_storage.ProvStore, countProviders)

  for i := 0; i < countProviders; i++ {
    provConfig = &Config.Providers[i]
    provIdHash = xxhash32.ChecksumString32(provConfig.IdName)
    prov = &mem_storage.ProviderElement{
      Id:          provConfig.IdName,
      IdHash:      provIdHash,   
      ChTTL:       provConfig.ChTTL,
      Flags:       (provConfig.Flags & mem_storage.StaticFlagMask),
      XmltvHashes: make([]uint32, len(provConfig.XmltvHashes)),
    }
    copy(prov.XmltvHashes, provConfig.XmltvHashes)
    p[provIdHash] = prov
  }

  // Загрузка ссылок из channels_url во временное хранилище
  config.ChUrls_Load(p)

  // Загрузка данных из channel.json во временное хранилище
  SchedulerCall(p)
  if !isFirstRun {
    // Запуск во время работы, блокируем в конце на время подмены/сортировки
    ChDbLock.Lock()
    defer ChDbLock.Unlock()
  }
  // Подмена основного хранилища временным (тут уже должен быть включен Lock)
  P = p
  ReOrder()
}


func EpgMatch(w http.ResponseWriter, r *http.Request) {
  req_parts := request_parser.SplitMatchBody(r)
  if req_parts == nil {
    http.Error(w, "Bad Request", 400)
    return
  }
  // Читаем список провайдеров из заголовка
  // TODO: обработка метаданных из req_parts[0]
  ChDbLock.RLock()
  defer ChDbLock.RUnlock()
  _m3u_providers, prov_user_len := request_parser.PrioritizeUserProviders(req_parts[1], PO.Epg, request_parser.Epg_GetProvByHash)
  if len(_m3u_providers) == 0 {
    http.Error(w, "Bad Request", 500)
    return
 }
  request_parser.ParseEpg_ReqChannels(req_parts[2], _m3u_providers, w, prov_user_len)
}


func IcoMatch(w http.ResponseWriter, r *http.Request) {
  req_parts := request_parser.SplitMatchBody(r)
  if req_parts == nil {
    http.Error(w, "Bad Request", 400)
    return
  }
  // Читаем список провайдеров из заголовка
  // TODO: обработка метаданных из req_parts[1]
  ChDbLock.RLock()
  defer ChDbLock.RUnlock()
  _m3u_providers, prov_user_len := request_parser.PrioritizeUserProviders(req_parts[1], PO.Ico, request_parser.Ico_GetProvByHash)
  if len(_m3u_providers) == 0 {
    http.Error(w, "Bad Request", 500)
    return
 }
  request_parser.ParseIco_ReqChannels(req_parts[2], _m3u_providers, w, prov_user_len)
}
