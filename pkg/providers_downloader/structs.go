package providers_downloader

type ProvMeta struct {
	Id      string   `json:"id"`         // ID провайдера
	Urls    []uint32 `json:"url-hashes"` // Хеш список ссылок XMLTV
	LastUpd uint64   `json:"last-upd"`   // Время создания XMLTV
	LastEpg uint64   `json:"last-epg"`   // Время последней программы XMLTV
}

//easyjson:json
type ProviderData struct {
	Data map[uint32][]string `json:"data"`
	Meta ProvMeta            `json:"meta"`
}
