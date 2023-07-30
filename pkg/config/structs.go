package config

type ProviderConfig struct {
	IdName      string   `json:"id"`           // ID провайдера
	Xmltv       []string `json:"xmltv"`        // Список ссылок XMLTV
	XmltvHashes []uint32 `json:"-"`            // xxHash32: Хеш списка ссылок XMLTV
	ChTTL       uint16   `json:"channels_ttl"` // Кол-во часов автоматического обновления channels.json
	XmltvTTL    uint16   `json:"xmltv_ttl"`    // TODO: Кол-во часов минимального времени жизни XMLTV
	Flags       uint8    `json:"flags"`        // Флаги провайдера в конфиге
}

type _ProvidersDefault struct {
	ChTTL       uint16   `json:"channels_ttl"` // Кол-во часов автоматического обновления channels.json
	XmltvTTL    uint16   `json:"xmltv_ttl"`    // TODO: Кол-во часов минимального времени жизни XMLTV
}


type ConfigFile struct {
	Bind         string            `json:"bind"`               // Bind для сервера
	BaseUrl      string            `json:"base_url"`           // Базовая ссылка для телепрограммы
	Providers    []ProviderConfig  `json:"providers"`          // Список провайдеров
	AdminTokens  []string          `json:"admin_tokens"`       // Список токенов для управления сервером
	EpgTokens    []string          `json:"worker_tokens"`      // Список токенов для управления источниками EPG
	ProvidersDef _ProvidersDefault `json:"providers_default"`  // Блок default значений для списка провайдеров
}
