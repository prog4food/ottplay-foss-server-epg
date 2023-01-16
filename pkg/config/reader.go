package config

import (
	"os"

	xxhash32 "github.com/OneOfOne/xxhash"
	"github.com/hjson/hjson-go/v4"
	"github.com/rs/zerolog/log"

	"ottplay-foss-server-epg/pkg/helpers"
)

const (
	fileConfig   = "config.hjson"
	FileOrderEpg = "order_epg.hjson"
	FileOrderIco = "order_ico.hjson"
)


func Load() *ConfigFile {
	var len_M3uUrls int
	var p *ProviderConfig
	conf := &ConfigFile{}
	jsonData, err := os.ReadFile(fileConfig)
	if err == nil {
		err = hjson.Unmarshal(jsonData, conf)
	}
	if err != nil {
		log.Err(err).Msg(fileConfig + ": cannot read config!!!")
		return nil
	}

	// Параметры bind по-умолчанию
	if conf.Bind == "" {
		conf.Bind = "127.0.0.1:3001"
	}

	// Чтение провайдеров
	var provIdHash uint32
	for i := 0; i < len(conf.Providers); i++ {
		p = &conf.Providers[i]
		// xxhash32 для Id провайдера
		provIdHash = xxhash32.ChecksumString32(p.IdName)

		// Хеширование всех M3uUrls + имя провайдера
		len_M3uUrls = len(p.Xmltv)
		p.XmltvHashes = make([]uint32, len_M3uUrls+1)
		for g := 0; g < len_M3uUrls; g++ {
			p.XmltvHashes[g] = xxhash32.ChecksumString32(helpers.CutHTTP(p.Xmltv[g]))
		}
		p.XmltvHashes[len_M3uUrls] = provIdHash

		// Обновление списка каналов = 6 часов
		if p.ChTTL == 0 {
			p.ChTTL = 6
		}
	}
	return conf
}
