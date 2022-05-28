package config_epg

//type ProviderSource [2]string
type ProviderSource struct {
  IdHash      uint32  `json:"-"`
  IdName      string  `json:"id"`
  UrlBase     string  `json:"base_url"`
  UrlBaseS    string  `json:"base_url_https"`
  OrderEpg    uint16  `json:"order_epg"`
  OrderIco    uint16  `json:"order_ico"`
  UpdMins     uint16  `json:"update_mins"`
  BodySum     uint32  `json:"-"`
  NextUpd     int64   `json:"-"`
}

type ConfigFile struct {
  Bind        string          `json:"bind"`
  Providers []ProviderSource  `json:"providers"`
  ApiTokens []string          `json:"api_tokens"`
}
