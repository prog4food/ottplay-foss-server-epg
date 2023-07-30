package mem_storage

import (
	"net/url"

	"github.com/rs/zerolog/log"
)

// Хранилище провайдера Epg
type EpgChannelData struct {
	IdHash []byte
	// TODO: пока нет смысла сохранять имена и Id в виде строк
	// Id         *string
	// Names    []*string
	LastEpg    int64
}
type HashEpgStruct map[uint32]*EpgChannelData
type ProviderEpgData struct {
	ById    HashEpgStruct
	ByName  HashEpgStruct
	Prov    *ProviderElement
	ProvUrl string
}

// Хранилище провайдера Ico
type HashIcoStruct map[uint32][]byte
type ProviderIcoData struct {
	ById   HashIcoStruct
	ByName HashIcoStruct
	Prov   *ProviderElement
}

// Единица хранилища провайдера
type ProviderElement struct {
	Id          string
	ETag        string   // ETag последнего channels.json
	ChUrl       string   // Ссылка на channels.json
	ChURL       *url.URL // Подготовленная ссылка на channels.json
	Epg         *ProviderEpgData
	Ico         *ProviderIcoData
	XmltvHashes []uint32 // Копия блока  XmltvHashes из Config
	ChTTR       int64    // UnixTime: Время следующего обновления channels.json
	IdHash      uint32   // xxHash32: Хэш Id
	ChHash      uint32   // xxHash32: Хэш последнего channels.json
	ChTTL       uint16   // Кол-во часов автоматического обновления channels.json
	/*
		Дополнительные статусные флаги провайдера
		 7 - Провайдер готов к использованию
		 6 - Зарезервировано (FlagReserved1)
	*/
	Flags      uint8
	ErrorCount uint8 // Кол-во ошибок при обработке channels.json
}

type OrderedStorage struct {
	Epg []*ProviderEpgData
	Ico []*ProviderIcoData
}

type ProvStore map[uint32]*ProviderElement

const (
	StaticFlagMask = 0b00111111 // Маска флагов, которые надо забирать из конфига
	// Динамические флаги
	FlagReady     = 7 //Провайдер готов к использованию
	FlagReserved1 = 6 // Зарезервировано (FlagReserved1)

	TypeEpg = 1
	TypeIco = 2
)

// Устанавливает новый URL для channels.json
func (c *ProviderElement) SetUrl(p_url string) {
	c.ChUrl = p_url
	if p_url != "" {
		urlObject, err := url.Parse(p_url + "channels.json")
		if err == nil {
			c.ChURL = urlObject
			return
		}
		log.Err(err).Msgf("cannot parse url %s", p_url)
	}
	c.ChURL = nil
}
