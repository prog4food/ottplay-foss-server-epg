package config

import (
	"encoding/json"
	"os"
	"time"

	xxhash32 "github.com/OneOfOne/xxhash"
	"github.com/rs/zerolog/log"

	ms "ottplay-foss-server-epg/pkg/mem_storage"
)

type ChannelsUrls map[string]string

const fileUrls = "channels_url.json"

var chUrls *ChannelsUrls


func getConfigProvidersLen(p ms.ProvStore) int {
	providers_count := len(p)
	if providers_count == 0 {
		log.Error().Msg("channels_url: providers config is empty!")
	}
	return providers_count
}


// Загрузка списка
func ChUrls_Load(p ms.ProvStore) {
	providersCount := getConfigProvidersLen(p)
	if providersCount == 0 {
		return
	}

	jsonData, err := os.ReadFile(fileUrls)
	if err != nil {
		log.Warn().Msg(fileUrls + ": cannot read")
		return
	}

	_chUrls := ChannelsUrls{}
	chUrls = &_chUrls
	err = json.Unmarshal(jsonData, &_chUrls)
	if err != nil {
		log.Error().Err(err).Msg(fileUrls + ": parse error!")
	}

	// Сбрасываем флаг сортировки
	var provider *ms.ProviderElement
	for _, provider = range p {
		url, ok := _chUrls[provider.Id]
		if ok {
			provider.SetUrl(url)
		} else {
			provider.SetUrl("")
			log.Warn().Msgf("channels.config: %s has no channels_url", provider.Id)
		}
	}
}

// Обновление ссылки и сохранение в файл
func ChUrls_Update(p ms.ProvStore, key, newUrl string) {
	providersCount := getConfigProvidersLen(p)
	if providersCount == 0 {
		return
	}

	var provider *ms.ProviderElement
	// Поиск провайдера в конфиге
	var providerHash = xxhash32.ChecksumString32(key)
	provider, ok := p[providerHash]
	if !ok {
		// Конфиг не найден
		log.Error().Msgf("channels_url: provider [%s]: not found in config!", key)
		return
	}

	// Проверка, что существующий url отличается
	if provider.ChUrl == newUrl {
		return
	}

	// Обновление ссылки у провайдера
	provider.SetUrl(newUrl)
	provider.ChTTR = time.Now().Unix()
	provider.ETag  = ""
	(*chUrls)[provider.Id] = newUrl
	ChUrls_Save(p, false)
	log.Info().Msgf("[%s]: new channel url", key)
}


// Сохранение списка в файл
func ChUrls_Save(p ms.ProvStore, onlyActual bool) {
	lenProviders := getConfigProvidersLen(p)
	if lenProviders == 0 {
		return
	}

	// Сохранить только провайдеры из конфига
	if onlyActual {
		var _urls = ChannelsUrls{}
		var _provider *ms.ProviderElement
		for _, _provider = range p {
			if _provider.ChUrl != "" {
				_urls[_provider.Id] = _provider.ChUrl
			}
		}
		chUrls = &_urls
	}

	// Пишем в файл
	jsonData, err := json.MarshalIndent(chUrls, "", "")
	if err == nil {
		err = os.WriteFile(fileUrls, jsonData, 0664)
		if err == nil {
			return
		}
	}
	log.Err(err).Msg(fileUrls + ": cannot save")
}

