package config_epg

//type ProviderSource [2]string
type ProviderSource struct {
	IdName      string  `json:"id"`
	IdHash      uint32  `json:"-"`
	UrlBase     string  `json:"base_url"`
	UrlBaseS    string  `json:"base_url_https"`
	OrderEpg    uint8   `json:"order_epg"`
	OrderIco    uint8   `json:"order_ico"`
	UpdMins     uint16  `json:"update_mins"`
	NextUpd     int64   `json:"-"`
	BodySum     uint32  `json:"-"`
	EpgSortIdx  uint8   `json:"-"`
	IcoSortIdx  uint8   `json:"-"`
}

type ConfigFile struct {
	Providers []ProviderSource  `json:"providers"`
	ApiTokens []string          `json:"api_tokens"`
}
