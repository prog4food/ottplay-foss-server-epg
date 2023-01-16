package internal

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	xxhash32 "github.com/OneOfOne/xxhash"

	"ottplay-foss-server-epg/pkg/helpers"
	"ottplay-foss-server-epg/pkg/mem_storage"
	"ottplay-foss-server-epg/pkg/providers_downloader"
)

const (
	HeaderIfNoneMatch = "If-None-Match"
	HeaderETag        = "ETag"

	flagNoEpg  = 0
	flagNoIco  = 1
	flagDirect = 2
)

var (
	HttpClient *http.Client = &http.Client{
		Timeout: 10 * time.Second,
	}
	ChannelRequester, _ = http.NewRequest("GET", "", nil)
)

func loadProvider(p *mem_storage.ProviderElement, b []byte) {
	// Чтение данных channels.json
	var channels = providers_downloader.ProviderData{}
	var err = channels.UnmarshalJSON(b)
	if err != nil {
		log.Err(err).Msgf("channels.json: %s - bad json", p.Id)
		return
	}
	var channelsCount = len(channels.Data) // Количество каналов в channels.json

	var pEpg *mem_storage.ProviderEpgData
	var pIco *mem_storage.ProviderIcoData
	var flagSkipEpg = helpers.Bit_has(p.Flags, flagNoEpg)
	var flagSkipIco = helpers.Bit_has(p.Flags, flagNoIco)
	if !flagSkipEpg { // Создаем хранилище: pEpg
		pEpg = &mem_storage.ProviderEpgData{
			ById:    make(mem_storage.HashEpgStruct, channelsCount),
			ByName:  make(mem_storage.HashEpgStruct, channelsCount),
			ProvUrl: Config.BaseUrl + p.Id + "/",
			Prov:    p,
		}
		if helpers.Bit_has(p.Flags, flagDirect) {
			pEpg.ProvUrl = p.ChUrl
		}
	}

	if !flagSkipIco { // Создаем хранилище: pIco
		pIco = &mem_storage.ProviderIcoData{
			ById:   make(mem_storage.HashIcoStruct, channelsCount),
			ByName: make(mem_storage.HashIcoStruct, channelsCount),
			Prov:   p,
		}
	}

	p.Epg = pEpg
	p.Ico = pIco

	// Обработка каналов
	var (
		last_epg  int64              // Временная переменная: последняя передача
		chEpg     *mem_storage.EpgChannelData // Временная переменная: структура Epg
		icoUrl    []byte             // Временная переменная: dedup ico
		chSkipEpg bool               // Временная переменная: канал имеет EPG
		chSkipIco bool               // Временная переменная: канал имеет ICO
		chExpired uint16             // Временная переменная: последняя передача
	)
	for _chIdHash, _chPayload := range channels.Data {
		/*
		   Чтение _chPayload[0] (метаданные)
		   _ch_meta[0] - Ид канала
		   _ch_meta[1] - Время последней передачи
		   _ch_meta[2] - Ссылка на значок
		*/
		_ch_meta := strings.Split(_chPayload[0], "¦")
		if len(_ch_meta) != 3 {
			log.Error().Msgf("channels.json: %s - bad channel / %d: %s", p.Id, _chIdHash, _chPayload[0])
			continue
		}

		last_epg, err = strconv.ParseInt(_ch_meta[1], 10, 64)
		if err != nil {
			log.Error().Msgf("channels.json: %s - bad last_epg / %d: %s", p.Id, _chIdHash, _ch_meta[1])
			continue
		}
		chSkipEpg = flagSkipEpg || last_epg == 0
		chSkipIco = flagSkipIco || _ch_meta[2] == ""

		// Чтение Id канала для pEpg/pIco
		if !chSkipEpg {
			chEpg = &mem_storage.EpgChannelData{
				// TODO: пока нет смысла сохранять имена и Id в виде строк
				// Id:         mem_storage.DedupStr(mem_storage.Str, &_ch_meta[0]),
				// Names:      make([]*string, 0, len(_payload)-1),
				IdHash:     mem_storage.DedupByte(mem_storage.Ids, helpers.Uint32ToBytes(nil, _chIdHash)),
				LastEpg:    last_epg,
				ExpiredEpg: (last_epg < _now_unix),
			}
			if chEpg.ExpiredEpg {
				chExpired++
			}
			pEpg.ById[_chIdHash] = chEpg
		}
		if !chSkipIco {
			icoUrl = mem_storage.DedupByteByS(mem_storage.Url, &_ch_meta[2])
			pIco.ById[_chIdHash] = icoUrl
		}

		// Чтение альтернативных имен из _chPayload[1..] для pEpg/pIco
		for g := 1; g < len(_chPayload); g++ {
			_name_hash := xxhash32.ChecksumString32(strings.ToLower(_chPayload[g]))
			if !chSkipEpg {
				// TODO: пока нет смысла сохранять имена и Id в виде строк
				// epg_data.Names = append(epg_data.Names, mem_storage.DedupStr(mem_storage.Str, &_payload[g]))
				pEpg.ByName[_name_hash] = chEpg
			}
			if !chSkipIco {
				pIco.ByName[_name_hash] = icoUrl
			}
		}
	}
	totalEpgName := 0
	totalEpgId   := 0
	totalIcoId   := 0
	if pIco != nil {
		totalIcoId = len(pIco.ById)
	}
	if pEpg != nil {
		totalEpgName = len(pEpg.ByName)
		totalEpgId = len(pEpg.ById)
	}
	log.Info().Msgf("channels.net: %s, names=%d, epg=%d outdated=%d, icon=%d", p.Id, totalEpgName, totalEpgId, chExpired, totalIcoId)
	return
}

// Скачивает channels.json, если он изменился
func DownloadProvider(store mem_storage.ProvStore, p *mem_storage.ProviderElement) {
	var buf []byte
	var _body_hash uint32
	var err error
	var r *http.Response
	var wasActive = helpers.Bit_has(p.Flags, mem_storage.FlagReady)

	if p.ChURL == nil {
		goto errDownload
	}
	ChannelRequester.URL = p.ChURL
	// Добавляем ETag, если есть
	if p.ETag != "" {
		ChannelRequester.Header.Set(HeaderIfNoneMatch, p.ETag)
	} else {
		ChannelRequester.Header.Del(HeaderIfNoneMatch)
	}
	r, err = HttpClient.Do(ChannelRequester)
	if err != nil {
		goto errDownload
	}
	if p.ETag != "" && r.StatusCode == http.StatusNotModified {
		// Not changed (ETag)
		goto okNoChange
	}
	if r.StatusCode != http.StatusOK {
		goto errDownload
	}

	// Читаем тело ответа
	buf, err = io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		goto errDownload
	}
	_body_hash = xxhash32.Checksum32(buf)
	if p.ChHash == _body_hash {
		// Not changed (_body_hash)
		goto okNoChange
	}
	loadProvider(p, buf)
	// Сохраняем ETag и хэш ответа
	p.ChHash = _body_hash
	p.ETag = r.Header.Get(HeaderETag)
	// Меняем статус, если надо
	if !wasActive {
		p.Flags = helpers.Bit_set(p.Flags, mem_storage.FlagReady)
	}

	if helpers.HasSameDataPointer(store, P) {
		// Обновление основной структуры (Lock должен быть сейчас)
		ChDbLock.Lock()
		log.Printf("channels.net: %s - update with lock", p.Id)
		// TODO: Переделать на поиск и замену элемента без ReOrder
		// TODO: Учитывать, что Epg или Ico могут исчезнуть
		store[p.IdHash] = p
		ReOrder()
		ChDbLock.Unlock()
	} else {
		// Обновление не основной структуры (Lock и ReOrder будут дальше, при подмене)
		store[p.IdHash] = p
	}
	// GOTO_BLOCK: Выход без изменений
okNoChange:
	p.ErrorCount = 0
	return

	// GOTO_BLOCK: Ошибка загрузки
errDownload:
	if p.ErrorCount <= 3 { // Ошибки считаем до 4х
		p.ErrorCount++
		if p.ErrorCount == 3 && wasActive { // После 3х ошибок отключаем провайдер
			p.Flags = helpers.Bit_clear(p.Flags, mem_storage.FlagReady)
			log.Error().Msgf("channels.net: %s - was disabled after 3 fails", p.Id)
			if helpers.HasSameDataPointer(store, P) {
				ChDbLock.Lock()
				log.Printf("channels.net: %s - update with lock", p.Id)
				store[p.IdHash] = p
				ReOrder()
				ChDbLock.Unlock()
			}
		}
	}

	if p.ChURL == nil {
		log.Warn().Msgf("channels.net: %s - empty ChUrl", p.Id)
	} else {
		log.Error().Err(err).Msgf("channels.net: %s - download error [%d, %s]", p.Id, r.StatusCode, p.ChUrl)
	}
}
