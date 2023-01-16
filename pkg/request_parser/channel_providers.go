package request_parser

import (
	"bytes"

	"github.com/rs/zerolog/log"

	"ottplay-foss-server-epg/pkg/helpers"
	"ottplay-foss-server-epg/pkg/mem_storage"
)

// ICO: Ищет провайдера по хешу M3uURL
// TODO: На данный момент этот свособ самый быстрый
func Ico_GetProvByHash(p uint32, l []*mem_storage.ProviderIcoData) *mem_storage.ProviderIcoData {
	var _prov *mem_storage.ProviderElement
	for i := 0; i < len(l); i++ {
		_prov = l[i].Prov
		if SliceExistUint32(_prov.XmltvHashes, p) {
			return l[i]
		}
	}
	return nil
}

// EPG: Ищет провайдера по хешу M3uURL
// TODO: На данный момент этот свособ самый быстрый
func Epg_GetProvByHash(p uint32, l []*mem_storage.ProviderEpgData) *mem_storage.ProviderEpgData {
	var _prov *mem_storage.ProviderElement
	for i := 0; i < len(l); i++ {
		_prov = l[i].Prov
		if SliceExistUint32(_prov.XmltvHashes, p) {
			return l[i]
		}
	}
	return nil
}

// Определяет приоритет поиска по провайдерам, на соснове провайдеров
// указанных для канала, которые передаются в виде хешей
func ReadChannelProviders[T *mem_storage.ProviderEpgData | *mem_storage.ProviderIcoData](inData []byte, listGlobal []T, look_func func(uint32, []T) T) ([]T, uint8) {
	if len(inData) == 0 {
		return listGlobal, 0
	}

	var hashUserXmltv uint32
	// Создаем с запасом по емкости (для последующей StoreSlice_PrioUser)
	var lenGlobal = len(listGlobal)
	var listUser = make([]T, 0, lenGlobal)
	var provUser T
	var err error
	listUserRaw := bytes.Split(inData, c_sep_vars) // "-"
	for i := 0; i < len(listUserRaw); i++ {
		hashUserXmltv, _, err = helpers.ParseUint32Buf(listUserRaw[i])
		if err != nil {
			log.Err(err).Msgf("match-channels.epg: cannot parse var %s", helpers.B2s(listUserRaw[i]))
			continue
		}

		// Проверяем, знаем ли такого провайдера
		provUser = look_func(hashUserXmltv, listGlobal)
		if provUser != nil {
			// Знаем: Добавляем в priority список (с проверкой на уникальность)
			listUser = Slice_AppendUniq(listUser, provUser)
		} else {
			// НЕ Знаем: Логируем его хеш и ссылку (частых можно будет добавлять)
			log.Warn().Msgf("channel-epg: unlisted provider - %d", hashUserXmltv)
		}
	}

	lenUser := uint8(len(listUser))
	if lenUser == 0 {
		// Ничего не нашли, ищем по стандартному порядку
		return listGlobal, 0
	}
	return SortByUser(listUser, listGlobal, lenGlobal), lenUser
}
