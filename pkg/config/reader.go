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

	// Параметры providers_default по-умолчанию
	if conf.ProvidersDef.ChTTL == 0 { conf.ProvidersDef.ChTTL = 6 }
	if conf.ProvidersDef.XmltvTTL == 0 { conf.ProvidersDef.XmltvTTL = 6 }

	// Чтение провайдеров
	var provIdHash uint32
	for i := 0; i < len(conf.Providers); i++ {
		p = &conf.Providers[i]
		// xxhash32 для Id провайдера
		provIdHash = xxhash32.ChecksumString32(p.IdName)

		// Хеширование всех M3uUrls + имя провайдера
		len_M3uUrls = len(p.Xmltv)
		p.XmltvHashes = make([]uint32, len_M3uUrls+1)
		var t []byte  // Временная строка для нормализации url-tvg
		var l int     // Временная переменная для длины строки t
		for g := 0; g < len_M3uUrls; g++ {
			t = helpers.S2b(p.Xmltv[g])
			if l = len(t); l > 10 {   // http://a.co
				t = helpers.CutURLb_gz(t, l) // первое, тк необходима l
				t = helpers.CutHTTPb(t)
			}

			p.XmltvHashes[g] = xxhash32.Checksum32(t)
		}
		p.XmltvHashes[len_M3uUrls] = provIdHash

		// Время автообновления channels.json
		if p.ChTTL == 0 { p.ChTTL = conf.ProvidersDef.ChTTL }
	}
	return conf
}
